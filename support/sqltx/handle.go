package sqltx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrTxFinished is returned by PrepareCached after the finalizer has committed or rolled back.
// It WRAPS sql.ErrTxDone so connectors may use errors.Is(err, sql.ErrTxDone) while the
// user-facing text names the subflow instead of the opaque "sql: transaction has already been
// committed or rolled back".
var ErrTxFinished = fmt.Errorf(
	"SUBFLOW-TX-005: the enclosing transactional subflow has already ended; this statement was not executed: %w",
	sql.ErrTxDone)

// maxCachedStmts bounds the per-transaction statement memo.
//
// The connectors' EvaluateQuery expands an IN-clause into a placeholder list sized to the input
// array, so a loop iterating with differently-sized arrays produces a DISTINCT SQL text per
// iteration.
//
// Be clear about what this fixes. It bounds the Go-side map only. Statements prepared on a
// *sql.Tx are appended to tx.stmts and are closed by Commit/Rollback, so N distinct SQL texts
// inside one transaction pin N server-side cursors REGARDLESS of memoisation - the ORA-01000
// class of exhaustion is a property of the transaction's length, not of this cache. The cap
// stops the map growing without bound and records the fact so it is diagnosable; dynamic SQL in
// a loop inside a transactional subflow remains a documented anti-pattern.
const maxCachedStmts = 128

// Handle is the ambient transaction for exactly one transactional subflow invocation.
//
// THREE MUTEXES, ONE OF WHICH - mu - IS THE ONLY ONE A CONNECTOR EVER TOUCHES. A single mutex
// self-deadlocks: the connector takes the operation lock at the top of PreparedQuery and then
// calls getStatement, which re-enters the handle, and sync.Mutex is not reentrant.
//
//	mu      the D3 SERIALISATION lock. A connector takes it at the TOP of a whole Prepared*
//	        operation and releases it only after *sql.Rows has been fully drained, because
//	        database/sql pins the transaction's single connection until the Rows is closed.
//	        Exposed as Lock/Unlock/TryLockFor.
//	stmtsMu guards the statement memo, the done flag and the warn set ONLY.
//	opMu    guards the in-flight statement-cancel registrations ONLY.
//
// The lock order is acyclic because nothing inside this package ever takes mu, and opMu and
// stmtsMu are never held simultaneously. That is what makes PrepareCached, MarkDone, IsDone,
// WarnOnce and OpContext safe to call with mu held, which is the normal case.
//
// Every method is nil-receiver safe: connectors call them on values that are nil on the
// non-transactional path.
type Handle struct {
	connID string
	db     *sql.DB
	tx     *sql.Tx

	// baseCtx is the context the transaction was BEGUN on. Cancelling it makes database/sql roll
	// the transaction back asynchronously, via the watchdog goroutine it starts in BeginTx. The
	// ENGINE therefore never cancels it: the subflow activity roots it at context.Background(),
	// and only the finalizer cancels it, AFTER Commit or Rollback.
	baseCtx context.Context

	mu sync.Mutex // D3 operation lock

	stmtsMu        sync.Mutex
	stmts          map[string]*sql.Stmt // memoised prepares, keyed by SQL TEXT
	done           bool
	memoFull       bool
	warned         map[string]bool // D13 dedup
	memoFullWarned bool

	opMu      sync.Mutex
	opNextTok uint64
	opCancels map[uint64]context.CancelFunc
}

// NewHandle builds a handle for a transaction that has already been begun.
func NewHandle(connID string, db *sql.DB, tx *sql.Tx, baseCtx context.Context) *Handle {
	if baseCtx == nil {
		baseCtx = context.Background()
	}

	return &Handle{
		connID:  connID,
		db:      db,
		tx:      tx,
		baseCtx: baseCtx,
		stmts:   make(map[string]*sql.Stmt),
		warned:  make(map[string]bool),
	}
}

// ConnID is the connection id this transaction was opened on.
func (h *Handle) ConnID() string {
	if h == nil {
		return ""
	}
	return h.connID
}

// DB is the pool the transaction was taken from.
func (h *Handle) DB() *sql.DB {
	if h == nil {
		return nil
	}
	return h.db
}

// Tx is the underlying transaction.
func (h *Handle) Tx() *sql.Tx {
	if h == nil {
		return nil
	}
	return h.tx
}

// Context returns the transaction's own context, for read-only use such as a parent. To EXECUTE
// a statement use OpContext instead, so the finalizer can interrupt it.
func (h *Handle) Context() context.Context {
	if h == nil {
		return context.Background()
	}
	return h.baseCtx
}

// Lock takes the D3 operation lock. It must wrap a WHOLE Prepared* operation, including full row
// drainage.
func (h *Handle) Lock() {
	if h != nil {
		h.mu.Lock()
	}
}

// Unlock releases the D3 operation lock.
func (h *Handle) Unlock() {
	if h != nil {
		h.mu.Unlock()
	}
}

// TryLockFor polls Mutex.TryLock until d elapses. It exists so the engine finalizer never blocks
// forever behind an abandoned execTimeout Eval goroutine. Acquiring the lock is BEST EFFORT on
// the rollback path: the finalizer rolls back either way, after CancelInFlight.
func (h *Handle) TryLockFor(d time.Duration) bool {
	if h == nil {
		return false
	}
	if h.mu.TryLock() {
		return true
	}

	deadline := time.Now().Add(d)
	for {
		time.Sleep(2 * time.Millisecond)
		if h.mu.TryLock() {
			return true
		}
		if !time.Now().Before(deadline) {
			return false
		}
	}
}

// NoRelease is the release func for a statement the caller does not own: the process-wide cache
// owns it on the non-enlisted path, and the transaction owns it on the memoised path. Callers
// always `defer release()`, so having a shared no-op keeps every call site the same shape.
var NoRelease = func() {}

// PrepareCached prepares sqlText ON THE TRANSACTION and memoises it for the transaction's
// lifetime.
//
// On the enlisted path we must NEVER touch the pool. sql.DB.Prepare acquires a SECOND pooled
// connection while the transaction pins one; at maxOpenConnection:1 - an ordinary user setting -
// that is a hard deadlock with no deadline to break it. Tx.PrepareContext uses the transaction's
// already-pinned connection and appends to tx.stmts so Commit/Rollback closes it.
//
// The returned release func MUST be deferred by the caller. It is a no-op for a memoised
// statement, which the transaction owns and closes; it closes the statement only when the memo
// was at its cap and this statement was prepared for this call alone. Never Close the returned
// statement directly - use release.
//
// Safe to call with or without the operation lock held.
func (h *Handle) PrepareCached(sqlText string) (*sql.Stmt, func(), error) {
	if h == nil {
		return nil, NoRelease, errors.New("no ambient transaction")
	}
	// Pass the transaction's own context explicitly rather than nil. PrepareCachedContext still
	// tolerates a nil context for callers that have none, but routing through it here would be a
	// staticcheck SA1012 violation for no benefit.
	return h.PrepareCachedContext(h.baseCtx, sqlText)
}

// PrepareCachedContext is PrepareCached with an explicit context for the PREPARE itself.
func (h *Handle) PrepareCachedContext(ctx context.Context, sqlText string) (*sql.Stmt, func(), error) {
	if h == nil || h.tx == nil {
		return nil, NoRelease, errors.New("no ambient transaction")
	}
	if ctx == nil {
		ctx = h.baseCtx
	}

	h.stmtsMu.Lock()
	if h.done {
		h.stmtsMu.Unlock()
		return nil, NoRelease, ErrTxFinished
	}
	if st, ok := h.stmts[sqlText]; ok {
		h.stmtsMu.Unlock()
		return st, NoRelease, nil
	}
	h.stmtsMu.Unlock()

	// PREPARE with no lock of ours held: mu may be held by our own caller, and MarkDone must not
	// block behind network I/O.
	st, err := h.tx.PrepareContext(ctx, sqlText) // on the TX - never on the pool
	if err != nil {
		return nil, NoRelease, err
	}

	h.stmtsMu.Lock()
	if h.done { // finalised while we were preparing
		h.stmtsMu.Unlock()
		_ = st.Close()
		return nil, NoRelease, ErrTxFinished
	}
	if existing, ok := h.stmts[sqlText]; ok { // lost a benign race
		h.stmtsMu.Unlock()
		_ = st.Close()
		return existing, NoRelease, nil
	}
	if len(h.stmts) >= maxCachedStmts {
		// At the cap. Hand the statement over to the caller instead of memoising it, so a loop
		// generating distinct SQL - an IN-clause expanded per batch size, say - closes each one
		// as it goes rather than pinning a server-side cursor per iteration until the
		// transaction ends. That exhaustion (ORA-01000 and friends) is a property of the
		// transaction's length, which memoisation alone cannot bound.
		h.memoFull = true
		h.stmtsMu.Unlock()
		return st, func() { _ = st.Close() }, nil
	}
	h.stmts[sqlText] = st
	h.stmtsMu.Unlock()

	return st, NoRelease, nil
}

// MemoSaturated reports whether the statement memo hit its cap, so the connector can emit one
// diagnostic. See maxCachedStmts for what this does and does not tell you.
func (h *Handle) MemoSaturated() bool {
	if h == nil {
		return false
	}
	h.stmtsMu.Lock()
	defer h.stmtsMu.Unlock()
	return h.memoFull
}

// MemoSaturatedOnce is MemoSaturated, but true at most once per handle, for a one-shot log.
func (h *Handle) MemoSaturatedOnce() bool {
	if h == nil {
		return false
	}
	h.stmtsMu.Lock()
	defer h.stmtsMu.Unlock()
	if !h.memoFull || h.memoFullWarned {
		return false
	}
	h.memoFullWarned = true
	return true
}

// OpContext returns the context a connector must use to EXECUTE a statement inside the
// transaction, plus its release func.
//
// The parent is baseCtx, NOT context.Background(): a per-statement
// context.WithTimeout(context.Background(), queryTimeout) firing mid-statement cancels the query
// on the connection the transaction is pinned to and typically poisons it. Cancelling baseCtx
// itself is not an option either - database/sql's watchdog would roll the transaction back.
//
// Registrations are TOKEN-KEYED. A single cancel slot loses B's registration when two operations
// overlap, and then A's release clears B's, leaving CancelInFlight unable to unwedge anything.
// Overlap should not happen under the D3 discipline, but the handle must not depend on
// discipline it cannot enforce.
//
// A timeout of zero or less means no deadline.
func (h *Handle) OpContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	if h == nil {
		return context.Background(), func() {}
	}

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(h.baseCtx, timeout)
	} else {
		ctx, cancel = context.WithCancel(h.baseCtx)
	}

	h.opMu.Lock()
	h.opNextTok++
	tok := h.opNextTok
	if h.opCancels == nil {
		h.opCancels = make(map[uint64]context.CancelFunc, 2)
	}
	h.opCancels[tok] = cancel
	h.opMu.Unlock()

	return ctx, func() {
		cancel()
		h.opMu.Lock()
		delete(h.opCancels, tok)
		h.opMu.Unlock()
	}
}

// CancelInFlight cancels EVERY statement context OpContext has handed out and not yet released.
// The finalizer calls it on the ROLLBACK path only.
func (h *Handle) CancelInFlight() {
	if h == nil {
		return
	}

	h.opMu.Lock()
	pending := make([]context.CancelFunc, 0, len(h.opCancels))
	for tok, c := range h.opCancels {
		pending = append(pending, c)
		delete(h.opCancels, tok)
	}
	h.opMu.Unlock()

	for _, c := range pending {
		c()
	}
}

// MarkDone is called by the finalizer BEFORE Commit or Rollback. Afterwards PrepareCached returns
// ErrTxFinished rather than a statement the commit is about to close.
func (h *Handle) MarkDone() {
	if h == nil {
		return
	}
	h.stmtsMu.Lock()
	h.done = true
	h.stmtsMu.Unlock()
}

// IsDone reports whether the transaction has been finalised. A nil handle counts as done.
func (h *Handle) IsDone() bool {
	if h == nil {
		return true
	}
	h.stmtsMu.Lock()
	defer h.stmtsMu.Unlock()
	return h.done
}

// WarnOnce reports whether the caller should emit a warning for key on this handle. True exactly
// once per handle per key - that is, once per transactional subflow invocation per foreign
// connection, never process-wide.
func (h *Handle) WarnOnce(key string) bool {
	if h == nil {
		return false
	}
	h.stmtsMu.Lock()
	defer h.stmtsMu.Unlock()
	if h.warned[key] {
		return false
	}
	h.warned[key] = true
	return true
}

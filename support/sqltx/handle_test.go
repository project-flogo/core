package sqltx

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"
)

// beginFake opens a pool with the given MaxOpenConns, begins a transaction and returns a handle
// over it, plus the driver's counters and a cleanup func.
func beginFake(t *testing.T, maxOpen int) (*Handle, *fakeDriver, func()) {
	t.Helper()

	db, drv := newFakeDB(maxOpen)
	base := context.Background()

	tx, err := db.BeginTx(base, nil)
	if err != nil {
		t.Fatalf("BeginTx: %v", err)
	}

	h := NewHandle("conn-a", db, tx, base)

	return h, drv, func() {
		_ = tx.Rollback()
		_ = db.Close()
	}
}

// TestX2_LockThenPrepareCachedDoesNotDeadlock is the regression test for blocker X2.
//
// With a single mutex serving both the D3 operation lock and the statement memo, this sequence
// self-deadlocks on the first enlisted statement: the connector takes the operation lock at the
// top of PreparedQuery, then getStatement re-enters the handle, and sync.Mutex is not reentrant.
// The failure mode is a silent hang, so the whole suite runs with -timeout.
func TestX2_LockThenPrepareCachedDoesNotDeadlock(t *testing.T) {
	h, _, cleanup := beginFake(t, 0)
	defer cleanup()

	done := make(chan struct{})

	go func() {
		defer close(done)
		h.Lock()
		defer h.Unlock()

		if _, _, err := h.PrepareCached("SELECT 1"); err != nil {
			t.Errorf("first PrepareCached with the operation lock held: %v", err)
		}
		if _, _, err := h.PrepareCached("SELECT 1"); err != nil {
			t.Errorf("second PrepareCached (memo hit) with the operation lock held: %v", err)
		}
		// The other lock-held-safe methods, same reasoning.
		h.IsDone()
		h.WarnOnce("k")
		h.MemoSaturated()
		relCtx, release := h.OpContext(0)
		_ = relCtx
		release()
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("deadlocked: PrepareCached re-entered the operation lock")
	}
}

// TestX3_PrepareCachedNeverTouchesThePool is the regression test for blocker X3.
//
// A cache-miss db.Prepare needs a SECOND pooled connection while the transaction pins one. At
// maxOpenConnection:1 - an ordinary user setting on every connector's connection tile - that is
// a hard deadlock on context.Background() with no deadline to break it. Preparing on the
// transaction uses the already-pinned connection instead.
func TestX3_PrepareCachedNeverTouchesThePool(t *testing.T) {
	h, drv, cleanup := beginFake(t, 1) // the pool has exactly one connection, and the tx holds it
	defer cleanup()

	done := make(chan error, 1)
	go func() {
		_, _, err := h.PrepareCached("INSERT INTO t VALUES (?)")
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("PrepareCached at MaxOpenConns=1: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("deadlocked at MaxOpenConns=1: the enlisted path acquired a pooled connection")
	}

	if opened, _, _, _, _ := drv.counts(); opened != 1 {
		t.Fatalf("driver opened %d connections, want exactly 1 (the transaction's)", opened)
	}
}

// TestX3_PoolPrepareWouldDeadlock documents the hazard X3 describes, so the reasoning behind
// PrepareCached is verifiable rather than asserted. It is the behaviour we must NOT have.
func TestX3_PoolPrepareWouldDeadlock(t *testing.T) {
	db, _ := newFakeDB(1)
	defer db.Close()

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("BeginTx: %v", err)
	}
	defer tx.Rollback()

	blocked := make(chan struct{})
	go func() {
		// The tx holds the only connection, so this waits for one that cannot be returned.
		stmt, err := db.Prepare("SELECT 1")
		if err == nil {
			_ = stmt.Close()
		}
		close(blocked)
	}()

	select {
	case <-blocked:
		t.Fatal("expected db.Prepare to block while the transaction pins the only connection; " +
			"if this now returns, re-examine whether PrepareCached still needs to avoid the pool")
	case <-time.After(300 * time.Millisecond):
		// Blocked, as expected.
	}
}

func TestPrepareCachedMemoisesBySQLText(t *testing.T) {
	h, drv, cleanup := beginFake(t, 0)
	defer cleanup()

	first, _, err := h.PrepareCached("SELECT 1")
	if err != nil {
		t.Fatalf("first: %v", err)
	}
	second, _, err := h.PrepareCached("SELECT 1")
	if err != nil {
		t.Fatalf("second: %v", err)
	}

	if first != second {
		t.Fatal("PrepareCached must return the memoised statement for identical SQL text")
	}
	if _, _, prepares, _, _ := drv.counts(); prepares != 1 {
		t.Fatalf("driver prepared %d times, want 1", prepares)
	}
}

// TestMemoDoesNotGrowWithoutBound covers the dynamic-SQL loop: a connector expanding an IN-clause
// produces a distinct SQL text per iteration.
func TestMemoDoesNotGrowWithoutBound(t *testing.T) {
	h, _, cleanup := beginFake(t, 0)
	defer cleanup()

	for i := 0; i < 200; i++ {
		if _, _, err := h.PrepareCached("SELECT " + strconv.Itoa(i)); err != nil {
			t.Fatalf("PrepareCached #%d: %v", i, err)
		}
	}

	h.stmtsMu.Lock()
	n := len(h.stmts)
	h.stmtsMu.Unlock()

	if n != maxCachedStmts {
		t.Fatalf("memo holds %d statements, want the cap of %d", n, maxCachedStmts)
	}
	if !h.MemoSaturated() {
		t.Fatal("MemoSaturated should report true once the cap is hit")
	}
	if !h.MemoSaturatedOnce() {
		t.Fatal("MemoSaturatedOnce should report true the first time")
	}
	if h.MemoSaturatedOnce() {
		t.Fatal("MemoSaturatedOnce should report false the second time")
	}
}

// TestReleaseIsNoOpWhenMemoisedButClosesOnOverflow pins the contract the connectors rely on.
//
// A memoised statement is owned by the transaction and must NOT be closed by the caller -
// Commit/Rollback closes it. A statement handed out after the memo hit its cap is owned by the
// CALLER, and its release must close it, so a loop generating distinct SQL does not pin a
// server-side cursor per iteration for the rest of the transaction.
func TestReleaseIsNoOpWhenMemoisedButClosesOnOverflow(t *testing.T) {
	h, _, cleanup := beginFake(t, 0)
	defer cleanup()

	// Memoised: release must be the shared no-op, and the statement must stay usable afterwards.
	st, release, err := h.PrepareCached("SELECT 1")
	if err != nil {
		t.Fatalf("memoised prepare: %v", err)
	}
	release()
	again, _, err := h.PrepareCached("SELECT 1")
	if err != nil {
		t.Fatalf("after release: %v", err)
	}
	if again != st {
		t.Fatal("release() must not have discarded the memoised statement")
	}

	// Fill the memo, then verify the overflow statement really is closed by its release.
	for i := 0; i < maxCachedStmts+5; i++ {
		_, rel, err := h.PrepareCached("SELECT over " + strconv.Itoa(i))
		if err != nil {
			t.Fatalf("overflow prepare #%d: %v", i, err)
		}
		rel()
	}
	if !h.MemoSaturated() {
		t.Fatal("the memo should be saturated by now")
	}

	overflowStmt, rel, err := h.PrepareCached("SELECT the-overflow-one")
	if err != nil {
		t.Fatalf("overflow prepare: %v", err)
	}
	rel()
	// Closing twice is the observable proof the release closed it: a second Close on an already
	// closed *sql.Stmt is a no-op returning nil, whereas a live memoised statement would still be
	// present in h.stmts. Assert on the memo instead, which is unambiguous.
	h.stmtsMu.Lock()
	_, memoised := h.stmts["SELECT the-overflow-one"]
	n := len(h.stmts)
	h.stmtsMu.Unlock()
	if memoised {
		t.Fatal("an overflow statement must not be memoised")
	}
	if n != maxCachedStmts {
		t.Fatalf("memo holds %d, want it pinned at the cap %d", n, maxCachedStmts)
	}
	if overflowStmt == nil {
		t.Fatal("the overflow statement should still have been usable before release")
	}
}

func TestPrepareCachedAfterMarkDoneReturnsErrTxFinished(t *testing.T) {
	h, _, cleanup := beginFake(t, 0)
	defer cleanup()

	h.MarkDone()

	_, _, err := h.PrepareCached("SELECT 1")
	if !errors.Is(err, ErrTxFinished) {
		t.Fatalf("err = %v, want ErrTxFinished", err)
	}
	// Connectors branch on the stdlib sentinel, so the wrap must survive.
	if !errors.Is(err, sql.ErrTxDone) {
		t.Fatalf("ErrTxFinished must wrap sql.ErrTxDone; got %v", err)
	}
	if !h.IsDone() {
		t.Fatal("IsDone should be true after MarkDone")
	}
}

func TestNilHandleIsSafeEverywhere(t *testing.T) {
	var h *Handle // the non-transactional path

	if h.ConnID() != "" || h.DB() != nil || h.Tx() != nil {
		t.Fatal("nil handle accessors should return zero values")
	}
	if !h.IsDone() {
		t.Fatal("a nil handle counts as done")
	}
	if h.WarnOnce("k") {
		t.Fatal("a nil handle must not ask the caller to warn")
	}
	if h.MemoSaturated() || h.MemoSaturatedOnce() {
		t.Fatal("a nil handle is never saturated")
	}
	if h.TryLockFor(time.Millisecond) {
		t.Fatal("a nil handle cannot be locked")
	}
	if h.Context() == nil {
		t.Fatal("a nil handle should still yield a usable context")
	}
	if _, _, err := h.PrepareCached("SELECT 1"); err == nil {
		t.Fatal("PrepareCached on a nil handle should error, not panic")
	}

	// Must not panic.
	h.Lock()
	h.Unlock()
	h.MarkDone()
	h.CancelInFlight()
	_, release := h.OpContext(0)
	release()
}

func TestWarnOnceIsPerHandlePerKey(t *testing.T) {
	h := NewHandle("conn-a", nil, nil, nil)

	if !h.WarnOnce("conn-b") {
		t.Fatal("first warning for conn-b should fire")
	}
	if h.WarnOnce("conn-b") {
		t.Fatal("second warning for conn-b should be suppressed")
	}
	if !h.WarnOnce("conn-c") {
		t.Fatal("a different connection should warn independently")
	}

	// Per handle, not process-wide: the next subflow invocation warns again.
	if !NewHandle("conn-a", nil, nil, nil).WarnOnce("conn-b") {
		t.Fatal("a fresh handle must warn again")
	}
}

// TestOpContextTokenKeyed covers the overlap case a single cancel slot gets wrong: B's
// registration would evict A's, and then A's release would clear B's, leaving CancelInFlight
// with nothing to cancel.
func TestOpContextTokenKeyed(t *testing.T) {
	h, _, cleanup := beginFake(t, 0)
	defer cleanup()

	ctxA, releaseA := h.OpContext(0)
	ctxB, releaseB := h.OpContext(0)
	defer releaseA()
	defer releaseB()

	h.opMu.Lock()
	n := len(h.opCancels)
	h.opMu.Unlock()
	if n != 2 {
		t.Fatalf("registered %d cancels, want 2", n)
	}

	h.CancelInFlight()

	for name, ctx := range map[string]context.Context{"A": ctxA, "B": ctxB} {
		select {
		case <-ctx.Done():
		case <-time.After(time.Second):
			t.Fatalf("CancelInFlight did not cancel statement context %s", name)
		}
	}
}

func TestOpContextReleaseDeregisters(t *testing.T) {
	h, _, cleanup := beginFake(t, 0)
	defer cleanup()

	_, release := h.OpContext(0)
	release()

	h.opMu.Lock()
	n := len(h.opCancels)
	h.opMu.Unlock()

	if n != 0 {
		t.Fatalf("%d cancels still registered after release, want 0", n)
	}
}

// TestOpContextParentIsBaseCtxNotBackground guards the reason OpContext exists: a per-statement
// deadline must be able to interrupt the statement, but must derive from the transaction's own
// context rather than context.Background().
func TestOpContextParentIsBaseCtxNotBackground(t *testing.T) {
	db, _ := newFakeDB(0)
	defer db.Close()

	base, cancelBase := context.WithCancel(context.Background())
	tx, err := db.BeginTx(base, nil)
	if err != nil {
		t.Fatalf("BeginTx: %v", err)
	}
	h := NewHandle("conn-a", db, tx, base)

	opCtx, release := h.OpContext(0)
	defer release()

	cancelBase()

	select {
	case <-opCtx.Done():
	case <-time.After(time.Second):
		t.Fatal("statement context should derive from the transaction's base context")
	}
}

func TestTryLockForTimesOutRatherThanBlockingForever(t *testing.T) {
	h, _, cleanup := beginFake(t, 0)
	defer cleanup()

	h.Lock() // simulate an abandoned Eval goroutine still holding the operation lock
	defer h.Unlock()

	start := time.Now()
	if h.TryLockFor(50 * time.Millisecond) {
		t.Fatal("TryLockFor should not have acquired a held lock")
	}
	if elapsed := time.Since(start); elapsed < 40*time.Millisecond {
		t.Fatalf("TryLockFor returned after %v, want it to wait out the timeout", elapsed)
	}
}

func TestTryLockForAcquiresWhenFree(t *testing.T) {
	h, _, cleanup := beginFake(t, 0)
	defer cleanup()

	if !h.TryLockFor(time.Second) {
		t.Fatal("TryLockFor should acquire a free lock")
	}
	h.Unlock()
}

// TestConcurrentUseIsRaceFree drives the D3 discipline the way concurrent transition branches
// would. It exists to be run under -race.
func TestConcurrentUseIsRaceFree(t *testing.T) {
	h, _, cleanup := beginFake(t, 1)
	defer cleanup()

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			h.Lock()
			defer h.Unlock()

			if _, _, err := h.PrepareCached("SELECT " + strconv.Itoa(i%3)); err != nil {
				t.Errorf("goroutine %d: %v", i, err)
			}
			h.IsDone()
			h.WarnOnce("conn-" + strconv.Itoa(i%2))
			_, release := h.OpContext(0)
			release()
		}(i)
	}

	waited := make(chan struct{})
	go func() { wg.Wait(); close(waited) }()

	select {
	case <-waited:
	case <-time.After(10 * time.Second):
		t.Fatal("concurrent handle use deadlocked")
	}
}

// Package sqltx carries a database/sql transaction from the flow engine to the connector
// activities that must enlist in it (FLOGO-19484).
//
// The registry is ONE context key holding an immutable map[connID]*Handle, replaced
// copy-on-write. Copy-on-write is mandatory: with FLOGO_FLOW_EXECUTE_BRANCHES_CONCURRENTLY=true
// sibling branches read the map simultaneously, so WithHandle must never mutate a map that is
// already visible to another goroutine.
//
// Keying by connection id is what keeps the transaction scoped: an activity only enlists when it
// uses the very connection the transactional subflow declared. Two different connections inside
// one subflow cannot cross-contaminate.
package sqltx

import "context"

type ctxKey struct{}

// registry is immutable once stored in a context. Never mutate a map reached through
// ctx.Value; always copy.
type registry map[string]*Handle

func fromCtx(ctx context.Context) registry {
	if ctx == nil {
		return nil
	}
	r, _ := ctx.Value(ctxKey{}).(registry)
	return r
}

// WithHandle returns a context carrying h for connID.
//
// It is a no-op when h is nil or connID is empty, so a caller can never store a typed nil that
// later reads back as a non-nil interface.
func WithHandle(ctx context.Context, connID string, h *Handle) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if h == nil || connID == "" {
		return ctx
	}

	old := fromCtx(ctx)
	next := make(registry, len(old)+1)
	for id, existing := range old {
		next[id] = existing
	}
	next[connID] = h

	return context.WithValue(ctx, ctxKey{}, next)
}

// FromContext returns the handle registered for connID, or nil. Safe on a nil context.
func FromContext(ctx context.Context, connID string) *Handle {
	if connID == "" {
		return nil
	}
	return fromCtx(ctx)[connID]
}

// HasAny reports whether ctx carries any handle at all. This is the fast path every connector
// takes before doing anything else: one ctx.Value walk on the non-transactional path.
func HasAny(ctx context.Context) bool {
	return len(fromCtx(ctx)) > 0
}

// ConnIDs returns the connection ids ctx carries, unordered. Used by the nested-transaction
// guard and by the D13 warning message.
func ConnIDs(ctx context.Context) []string {
	r := fromCtx(ctx)
	if len(r) == 0 {
		return nil
	}

	ids := make([]string, 0, len(r))
	for id := range r {
		ids = append(ids, id)
	}

	return ids
}

// Propagate copies src's registry onto dst, preserving dst's own cancellation and deadline.
//
// It returns dst UNCHANGED - identically, with no allocation - when src carries no registry.
// That identity matters: flow's TestGoContextEvalCtxOverride asserts GoContext() returns exactly
// the context it was given, and every non-transactional flow goes through this path.
func Propagate(src, dst context.Context) context.Context {
	r := fromCtx(src)
	if len(r) == 0 {
		return dst
	}
	if dst == nil {
		dst = context.Background()
	}

	return context.WithValue(dst, ctxKey{}, r)
}

// Without strips the registry. Used for detached subflows, which outlive the transaction and
// must never enlist in it.
func Without(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	if !HasAny(ctx) {
		return ctx
	}

	return context.WithValue(ctx, ctxKey{}, registry(nil))
}

package sqltx

import (
	"context"
	"sort"
	"testing"
)

func TestWithHandleFromContextRoundTrip(t *testing.T) {
	h := NewHandle("conn-a", nil, nil, nil)

	ctx := WithHandle(context.Background(), "conn-a", h)

	if got := FromContext(ctx, "conn-a"); got != h {
		t.Fatalf("FromContext(conn-a) = %v, want the handle we stored", got)
	}
	if !HasAny(ctx) {
		t.Fatal("HasAny = false, want true")
	}
}

func TestFromContextWrongConnIDReturnsNil(t *testing.T) {
	ctx := WithHandle(context.Background(), "conn-a", NewHandle("conn-a", nil, nil, nil))

	// This is the silent-no-op guard: if a connector derives the wrong id, it must get nil
	// rather than somebody else's transaction.
	if got := FromContext(ctx, "conn-b"); got != nil {
		t.Fatalf("FromContext(conn-b) = %v, want nil", got)
	}
	if got := FromContext(ctx, ""); got != nil {
		t.Fatalf("FromContext(\"\") = %v, want nil", got)
	}
}

func TestNilContextIsSafe(t *testing.T) {
	// core/support/test.TestActivityContext.GoContext() and flow's LegacyCtx both return nil,
	// and connectors call straight through. None of these may panic.
	if FromContext(nil, "conn-a") != nil {
		t.Fatal("FromContext(nil, ...) should be nil")
	}
	if HasAny(nil) {
		t.Fatal("HasAny(nil) should be false")
	}
	if ConnIDs(nil) != nil {
		t.Fatal("ConnIDs(nil) should be nil")
	}
	if got := WithHandle(nil, "conn-a", NewHandle("conn-a", nil, nil, nil)); got == nil {
		t.Fatal("WithHandle(nil, ...) should return a usable context")
	}
	if got := Without(nil); got == nil {
		t.Fatal("Without(nil) should return a usable context")
	}
}

func TestWithHandleNilHandleOrEmptyIDIsNoOp(t *testing.T) {
	base := context.Background()

	if got := WithHandle(base, "conn-a", nil); HasAny(got) {
		t.Fatal("storing a nil handle must not create a registry entry")
	}
	if got := WithHandle(base, "", NewHandle("", nil, nil, nil)); HasAny(got) {
		t.Fatal("storing under an empty connID must not create a registry entry")
	}
}

func TestWithHandleIsCopyOnWrite(t *testing.T) {
	a := NewHandle("conn-a", nil, nil, nil)
	b := NewHandle("conn-b", nil, nil, nil)

	ctxA := WithHandle(context.Background(), "conn-a", a)
	ctxAB := WithHandle(ctxA, "conn-b", b)

	// Concurrent sibling branches read the map simultaneously, so adding to a derived context
	// must never mutate the map the parent context still points at.
	if FromContext(ctxA, "conn-b") != nil {
		t.Fatal("WithHandle mutated the parent context's registry")
	}
	if FromContext(ctxAB, "conn-a") != a || FromContext(ctxAB, "conn-b") != b {
		t.Fatal("derived context lost an entry")
	}
}

func TestConnIDs(t *testing.T) {
	ctx := WithHandle(context.Background(), "conn-a", NewHandle("conn-a", nil, nil, nil))
	ctx = WithHandle(ctx, "conn-b", NewHandle("conn-b", nil, nil, nil))

	ids := ConnIDs(ctx)
	sort.Strings(ids)

	if len(ids) != 2 || ids[0] != "conn-a" || ids[1] != "conn-b" {
		t.Fatalf("ConnIDs = %v, want [conn-a conn-b]", ids)
	}
}

// TestPropagateIsIdentityWhenSourceHasNoRegistry is load-bearing: flow's
// TestGoContextEvalCtxOverride asserts GoContext() returns *exactly* the context it was given,
// and every non-transactional flow goes through Propagate. If this allocates, that test breaks
// and every flow pays for a feature it does not use.
func TestPropagateIsIdentityWhenSourceHasNoRegistry(t *testing.T) {
	dst := context.WithValue(context.Background(), struct{ k string }{"unrelated"}, 1)

	if got := Propagate(context.Background(), dst); got != dst {
		t.Fatal("Propagate must return dst identically when src carries no registry")
	}
}

func TestPropagateCopiesRegistryAndKeepsDstCancellation(t *testing.T) {
	h := NewHandle("conn-a", nil, nil, nil)
	src := WithHandle(context.Background(), "conn-a", h)

	dst, cancel := context.WithCancel(context.Background())
	defer cancel()

	merged := Propagate(src, dst)

	if FromContext(merged, "conn-a") != h {
		t.Fatal("Propagate did not copy the registry")
	}

	cancel()
	select {
	case <-merged.Done():
	default:
		t.Fatal("Propagate must preserve dst's cancellation")
	}
}

func TestWithoutStripsTheRegistry(t *testing.T) {
	ctx := WithHandle(context.Background(), "conn-a", NewHandle("conn-a", nil, nil, nil))

	stripped := Without(ctx)

	if HasAny(stripped) {
		t.Fatal("Without must strip the registry")
	}
	if FromContext(stripped, "conn-a") != nil {
		t.Fatal("Without must make handles unreachable")
	}
	// A detached subflow must not inherit the transaction; the parent keeps it.
	if !HasAny(ctx) {
		t.Fatal("Without must not mutate the source context")
	}
}

func TestWithoutIsIdentityWhenNoRegistry(t *testing.T) {
	base := context.Background()
	if got := Without(base); got != base {
		t.Fatal("Without should return the context unchanged when there is no registry")
	}
}

package sqltx

import (
	"context"
	"sort"
	"strings"

	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
)

// FromManager returns the ambient transaction handle for m, or nil.
//
// It is the ONLY entry point a connector needs, and it is the single home of the
// foreign-connection warning, so all four connectors get identical wording for free.
//
// Safe on a nil context: core/support/test.TestActivityContext.GoContext() returns nil, and so
// does flow's LegacyCtx, and both are on paths connectors already take.
//
// On the non-transactional path this costs one ctx.Value walk, because HasAny short-circuits
// before the manager is inspected at all.
func FromManager(ctx context.Context, m connection.Manager, logger log.Logger) *Handle {
	if ctx == nil || m == nil || !HasAny(ctx) {
		return nil
	}

	id := connection.GetId(m)
	if id != "" {
		if h := FromContext(ctx, id); h != nil {
			if logger != nil && logger.DebugEnabled() {
				// The feature's canary. The subflow activity logs the id it enlisted; this logs
				// the id it derived. If those two strings ever differ, the whole feature
				// silently no-ops with no error anywhere.
				logger.Debugf("enlisting in the transaction of the enclosing transactional subflow on connection '%s'", id)
			}
			return h
		}
	}

	// Reached only inside a transactional subflow, on a connection that is not the enlisted one.
	// Not an error: writing an audit row deliberately outside the transaction is legitimate. But
	// it must never be silent.
	warnForeign(ctx, m, id, logger)

	return nil
}

// warnForeign emits the D13 warning at most once per handle per foreign connection.
func warnForeign(ctx context.Context, m connection.Manager, id string, logger log.Logger) {
	if logger == nil {
		return
	}

	enlisted := ConnIDs(ctx)
	if len(enlisted) == 0 {
		return
	}
	sort.Strings(enlisted) // stable message and stable dedup anchor

	// Dedup through any ambient handle: per-handle means per subflow invocation, per-key means
	// per foreign connection.
	anchor := FromContext(ctx, enlisted[0])

	key := id
	if key == "" {
		key = "inline:" + m.Type()
	}
	if !anchor.WarnOnce(key) {
		return
	}

	enlistedIDs := strings.Join(enlisted, ", ")
	if id != "" {
		logger.Warnf("connection '%s' is not the connection enlisted in the enclosing transactional subflow (%s); its statements commit independently", id, enlistedIDs)
		return
	}

	logger.Warnf("this activity uses an inline (non-shared) connection, which is not the connection enlisted in the enclosing transactional subflow (%s); its statements commit independently", enlistedIDs)
}

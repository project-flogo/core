package trigger

import (
	"context"
)

type key int

var handlerKey key

type HandlerInfo struct {
	Name string
}

// NewHandlerContext add the handler info to a new child context
func NewHandlerContext(parentCtx context.Context, config *HandlerConfig) context.Context {
	if config != nil && config.Name != "" {
		return context.WithValue(parentCtx, handlerKey, &HandlerInfo{Name: config.Name})
	}
	return parentCtx
}

// HandlerFromContext returns the handler info stored in the context, if any.
func HandlerFromContext(ctx context.Context) (*HandlerInfo, bool) {
	u, ok := ctx.Value(handlerKey).(*HandlerInfo)
	return u, ok
}

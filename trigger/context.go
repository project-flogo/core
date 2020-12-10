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

// This allows trigger developer to pass request parameters to handler w/o going through output mapper, e.g.,
//   ctx := trigger.NewContextWithValues(context.Background(), values)
//   results, err := handler.Handle(ctx, triggerData.ToMap())

type valueKey string

const contextValueKey valueKey = "RequestParams"

// NewContextWithValues returns a new Context that carries specified request parameter values
func NewContextWithValues(ctx context.Context, values map[string]interface{}) context.Context {
	return context.WithValue(ctx, contextValueKey, values)
}

// ValuesFromContext extracts request parameters from a context
func ValuesFromContext(ctx context.Context) (map[string]interface{}, bool) {
	values, ok := ctx.Value(contextValueKey).(map[string]interface{})
	return values, ok
}

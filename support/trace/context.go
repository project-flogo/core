package trace

import "context"

type key struct{}
var id = key{}

func AppendTracingContext(goCtx context.Context, tracingContext TracingContext)  context.Context {
	return context.WithValue(goCtx, id, tracingContext)
}

func ExtractTracingContext(goCtx context.Context) TracingContext {
	tctx, _ := goCtx.Value(id).(TracingContext)
	return tctx
}

type TracingContext struct {
	spanContext interface{}
}

func (tc TracingContext) WrappedContext() interface{} {
	return tc.spanContext
}
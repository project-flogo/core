package trace

import "context"

type key struct{}

var id = key{}

func AppendTracingContext(goCtx context.Context, tracingContext TracingContext) context.Context {
	return context.WithValue(goCtx, id, tracingContext)
}

func ExtractTracingContext(goCtx context.Context) TracingContext {
	tc, _ := goCtx.Value(id).(TracingContext)
	return tc
}

type TracingContext struct {
	tContext interface{}
}

func (tc TracingContext) WrappedContext() interface{} {
	return tc.tContext
}

type Config struct {
	Operation string
	Tags      map[string]interface{}
}

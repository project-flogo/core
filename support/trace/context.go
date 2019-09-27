package trace

import "context"

type key struct{}
var id = key{}

func AppendTracingContext(goCtx context.Context, tracingContext interface{})  context.Context {
	return context.WithValue(goCtx, id, tracingContext)
}

func ExtractTracingContext(goCtx context.Context) interface{} {
	return goCtx.Value(id)
}
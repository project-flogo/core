package trace

import "context"

type key struct{}

var id = key{}


type TracingContext interface {
	// TraceObject() returns underlying tracing implementation
	TraceObject() interface{}
	// SetTags() allows you to set one or more tags to tracing object
	SetTags(tags map[string]interface{}) bool
	// SetTags() allows you to add tag to tracing object
	SetTag(tagKey string, tagValue interface{}) bool
	// LogKV() allows you to log additional details about the entity being traced
	LogKV(kvs map[string]interface{}) bool
}

type Config struct {
	Operation string
	Tags      map[string]interface{}
}

func AppendTracingContext(goCtx context.Context, tracingContext TracingContext) context.Context {
	return context.WithValue(goCtx, id, tracingContext)
}

func ExtractTracingContext(goCtx context.Context) TracingContext {
	tc, _ := goCtx.Value(id).(TracingContext)
	return tc
}
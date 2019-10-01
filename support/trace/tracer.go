package trace

// CarrierFormat defines the format of the propogation that needs to injected or extracted
type CarrierFormat byte

const (
	// Binary represents trace context as opaque binary data.
	Binary CarrierFormat = iota
	// TextMap represents trace context as key:value string pairs.
	//
	// Unlike HTTPHeaders, the TextMap format does not restrict the key or
	// value character sets in any way.
	TextMap
	// HTTPHeaders represents trace context as HTTP header string pairs.
	//
	// Unlike TextMap, the HTTPHeaders format requires that the keys and values
	// be valid as HTTP headers as-is (i.e., character casing may be unstable
	// and special characters are disallowed in keys, values should be
	// URL-escaped, etc).
	HTTPHeaders
)

// FlogoTracer interface to configure individual tracers
type FlogoTracer interface {
	// Name() returns the name of the registered tracer
	Name() string

	// configure() sets up the tracer by gathering required configuration from the trace config
	// environment variable. The tracer must also register as an event listenerbefore tracing can
	// begin.
	configure() error

	// Inject() takes the `flowId` and `taskInstanceId` and injects the current span context for
	// propagation within `carrier`. The actual type of `carrier` depends on the value of `format`.
	//
	// FlogoTracer defines a common set of `format` values, and each has an expected `carrier` type.
	//
	// Example usage:
	//
	// tracer := trace.GetTracer()
	// ti := ctx.(*instance.TaskInst)
	// err = tracer.Inject(ctx.ActivityHost().ID(), ti.InstanceId(), trace.HttpHeaders, req)
	Inject(flowID string, taskInstanceID string, format CarrierFormat, carrier interface{}) error

	// Extract() returns a TraceContext instance given `format` and `carrier`.
	//
	// FlogoTracer defines a common set of `format` values, and each has an expected `carrier` type.
	//
	// After extracting the trace context from the incoming request, the trace context must be appended
	// to the go context to propogate it to the action handler.
	// See trace.AppendTracingContext() and trace.ExtractTracingContext().
	//
	// Example usage:
	//
	//	ft := trace.GetTracer()
	//	tctx, err := ft.Extract(trace.HttpHeaders, httpReq)
	//	if err != nil {
	//		log.Errorf("failed to extract tracing context due to error: %s", err.Error())
	//	}
	//	ctx := trace.AppendTracingContext(context.Background(), tctx)
	//	..
	//	results, err := handler.Handle(ctx, outputData)
	Extract(format CarrierFormat, carrier interface{}) (interface{}, error)

	// SetTag() adds a tag to the span defined by the `spanKey`.
	// The `spanKey` is the flowId at the flow level and `flowId + taskInstanceId` at the activity level.
	//
	// If there is a pre-existing tag set for `key`, it is overwritten.
	//
	// Returns true if the span for the `spanKey` was found and a tag was set. Else returns false.
	SetTag(spanKey string, TagKey string, TagValue interface{}) bool

	// LogKV() is a concise, readable way to record key:value logging data about a span.
	// Similar to `SetTag()`, LogKV() takes a `spanKey` to log data about the specific span.
	//
	// The keys must all be strings.
	//
	LogKV(spanKey string, alternatingKeyValues ...interface{}) bool
}

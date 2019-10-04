package trace

import "github.com/project-flogo/core/support/managed"

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

// Tracer interface to configure individual tracers
type Tracer interface {
    // Implement Start() and Stop() methods to manage tracer lifecycle. These methods will be called during engine startup and engine shutdown.
	managed.Managed

	// Name() returns the name of the registered tracer
	Name() string

	// Inject() takes the tracing context `tctx` and injects the current trace context for
	// propagation within `carrier`. The actual type of `carrier` depends on the value of `format`.
	//
	// trace.Tracer defines a common set of `format` values, and each has an expected `carrier` type.
	//
	// Example usage:
	//
	// tracer := trace.GetTracer()
	// err = tracer.Inject(ctx.TracingContext(), trace.HTTPHeaders, req)
	Inject(tCtx TracingContext, format CarrierFormat, carrier interface{}) error

	// Extract() returns a TracingContext given `format` and `carrier`.
	//
	// FlogoTracer defines a common set of `format` values, and each has an expected `carrier` type.
	//
	// After extracting the trace context from the incoming request, the trace context must be appended
	// to the go context to propogate it to the action handler.
	// See trace.AppendTracingContext() and trace.ExtractTracingContext().
	//
	// Example usage:
	//
	//	tr := trace.GetTracer()
	//	tctx, err := tr.Extract(trace.HTTPHeaders, httpReq)
	//	if err != nil {
	//		log.Errorf("failed to extract tracing context due to error: %s", err.Error())
	//	}
	//	ctx := trace.AppendTracingContext(context.Background(), tctx)
	//	..
	//	results, err := handler.Handle(ctx, outputData)
	Extract(format CarrierFormat, carrier interface{}) (TracingContext, error)

	// SetTag() adds a tag to the span defined by the `tctx`.
	//
	// If there is a pre-existing tag set for `key`, it is overwritten.
	//
	// Returns true if the span for the `tctx` was found and a tag was set. Else returns false.
	SetTag(tCtx TracingContext, TagKey string, TagValue interface{}) bool

	// LogKV() is a concise, readable way to record key:value logging data about a span.
	// Similar to `SetTag()`, LogKV() takes a `tctx` to log data about the specific span.
	//
	// The keys must all be strings.
	//
	LogKV(tCtx TracingContext, alternatingKeyValues ...interface{}) bool

	// StartSpan() returns a wrapped span created by the underlying tracing implementation.
	// Non nil parent indicates child span.
	StartSpan(config Config, parent TracingContext) (TracingContext, error)

	// FinishSpan() finishes a span wrapped span
	// Non nil error indicates failure
	FinishSpan(tContext TracingContext, err error) error
}

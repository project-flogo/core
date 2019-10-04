package trace

import "github.com/project-flogo/core/support/log"

var tracer Tracer

// RegisterTracer registers the configured tracer
func registerTracer(t Tracer) {
	if !isTracerRegistered() {
		log.RootLogger().Infof("Registering tracer: %s", t.Name())
		tracer = t
	}
}

func isTracerRegistered() bool {
	return tracer != nil
}

// GetTracer returns the instance of the registered tracer.
// If no tracer is registered, a noop tracer is returned
func GetTracer() Tracer {
	if tracer == nil {
		log.RootLogger().Warn("No tracing configuration found. registering noop-tracer")
		tracer = &nooptracer{}
	}
	return tracer
}

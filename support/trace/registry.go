package trace

import "github.com/project-flogo/core/support/log"

var tracer FlogoTracer

// RegisterTracer registers the configured tracer
func registerTracer(t FlogoTracer) {
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
func GetTracer() FlogoTracer {
	if tracer == nil {
		log.RootLogger().Debugf("No tracing configuration found. registering noop-tracer")
		tracer = &nooptracer{}
	}
	return tracer
}

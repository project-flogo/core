package trace

import (
	"github.com/project-flogo/core/support/log"
)

var tracer Tracer

// RegisterTracer registers the configured tracer
func RegisterTracer(t Tracer) error {
	if !isTracerRegistered() {
		log.RootLogger().Infof("Registering tracer: %s", t.Name())
		tracer = t
	}
	return nil
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

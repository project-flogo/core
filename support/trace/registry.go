package trace

import (
	"errors"

	"github.com/project-flogo/core/support/log"
)

var tracer Tracer

// RegisterTracer registers the configured tracer
func RegisterTracer(t Tracer) error {
	if tracer == nil {
		log.RootLogger().Infof("Registering tracer: %s", t.Name())
		tracer = t
	} else {
		log.RootLogger().Warnf("Tracer: %s already registered", tracer.Name())
		return errors.New("Tracer is already registered")
	}

	return nil
}

func Enabled() bool {
	return tracer != nil
}

// GetTracer returns the instance of the registered tracer.
// If no tracer is registered, a noop tracer is returned
func GetTracer() Tracer {
	return tracer
}

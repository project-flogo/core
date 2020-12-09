package trigger

import (
	"fmt"
	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/log"
)

var (
	triggerFactories = make(map[string]Factory)
	triggerLoggers   = make(map[string]log.Logger)
)

var rootLogger = log.RootLogger()

func Register(trigger Trigger, f Factory) error {

	if trigger == nil {
		return fmt.Errorf("'trigger' must be specified when registering")
	}

	if f == nil {
		return fmt.Errorf("cannot register trigger with 'nil' trigger factory")
	}

	ref := support.GetRef(trigger)

	if triggerFactories[ref] != nil {
		return fmt.Errorf("trigger already registered for ref %s", ref)
	}

	log.RootLogger().Debugf("Registering trigger: %s", ref)

	triggerFactories[ref] = f

	triggerLoggers[ref] = log.CreateLoggerFromRef(rootLogger, "trigger", ref)
	return nil
}

func GetFactory(ref string) Factory {
	return triggerFactories[ref]
}

func Factories() map[string]Factory {
	//todo return copy
	return triggerFactories
}

func GetLogger(ref string) log.Logger {
	return triggerLoggers[ref]
}

package trigger

import (
	"fmt"
	"github.com/project-flogo/core/support/log"
	"path"

	"github.com/project-flogo/core/support"
)

var (
	triggerFactories = make(map[string]Factory)
	triggerLoggers   = make(map[string]log.Logger)
)

var triggerLogger = log.ChildLogger(log.RootLogger(), "trigger")

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

	triggerName := path.Base(ref) //todo get this from the descriptor or register trigger with name as well
	triggerLoggers[ref] = log.ChildLogger(triggerLogger, triggerName)

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

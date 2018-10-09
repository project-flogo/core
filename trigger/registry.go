package trigger

import (
	"fmt"

	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/logger"
)

var (
	triggerFactories = make(map[string]Factory)
)

func Register(trigger Trigger, f Factory) error {

	if trigger == nil {
		return fmt.Errorf("'trigger' must be specified when registering")
	}

	if f == nil {
		return fmt.Errorf("cannot register trigger with 'nil' trigger factory")
	}

	ref := support.GetRef(trigger)

	if triggerFactories[ref] != nil {
		return fmt.Errorf("trigger already registered for ref '%s'", ref)
	}

	logger.Debugf("Registering trigger [ %s ]", ref)

	triggerFactories[ref] = f

	return nil
}

func GetFactory(ref string) Factory {
	return triggerFactories[ref]
}

func Factories() map[string]Factory {
	//todo return copy
	return triggerFactories
}

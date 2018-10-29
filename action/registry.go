package action

import (
	"fmt"
	"github.com/project-flogo/core/support/log"

	"github.com/project-flogo/core/support"
)

var (
	actionFactories = make(map[string]Factory)
)

func Register(action Action, f Factory) error {

	if action == nil {
		return fmt.Errorf("'action' must be specified when registering")
	}

	if f == nil {
		return fmt.Errorf("cannot register action with 'nil' action factory")
	}

	ref := support.GetRef(action)

	if _, dup := actionFactories[ref]; dup {
		return fmt.Errorf("action already registered: %s", ref)
	}

	log.RootLogger().Debugf("Registering action: %s", ref)

	actionFactories[ref] = f

	return nil
}

//DEPRECATED
func LegacyRegister(ref string, f Factory) error {

	if ref == "" {
		return fmt.Errorf("'action ref' must be specified when registering")
	}

	if f == nil {
		return fmt.Errorf("cannot register action with 'nil' action factory")
	}

	if _, dup := actionFactories[ref]; dup {
		return fmt.Errorf("action already registered: %s", ref)
	}

	log.RootLogger().Debugf("Registering legacy action: %s", ref)

	actionFactories[ref] = f

	return nil
}

func GetFactory(ref string) Factory {

	//temp hack
	if ref == "github.com/TIBCOSoftware/flogo-contrib/action/flow" {
		ref = "github.com/project-flogo/flow"
	}

	return actionFactories[ref]
}

func Factories() map[string]Factory {
	//todo return copy
	return actionFactories
}

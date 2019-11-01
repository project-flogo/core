package service

import (
	"fmt"

	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/log"
)

var (
	serviceFactories = make(map[string]Factory)
)

func RegisterFactory(factory Factory) error {

	if factory == nil {
		return fmt.Errorf("cannot register 'nil' service factory")
	}

	ref := support.GetRef(factory)

	if _, dup := serviceFactories[ref]; dup {
		return fmt.Errorf("service factory already registered: %s", ref)
	}

	log.RootLogger().Debugf("Registering service factory: %s", ref)

	serviceFactories[ref] = factory

	return nil
}

func GetFactory(ref string) Factory {
	return serviceFactories[ref]
}
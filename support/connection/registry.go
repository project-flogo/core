package connection

import (
	"fmt"

	"github.com/project-flogo/core/support/managed"
	"github.com/project-flogo/core/support/log"
)

var (
	managerFactories = make(map[string]ManagerFactory)
	managers = make(map[string]Manager)
)

func RegisterManagerFactory(connectionType string, factory ManagerFactory) error {

	if connectionType == "" {
		return fmt.Errorf("'connectionType' must be specified when registering")
	}

	if factory == nil {
		return fmt.Errorf("cannot register with 'nil' manager factory")
	}

	if _, dup := managerFactories[connectionType]; dup {
		return fmt.Errorf("manager factory already registered for type: %s", connectionType)
	}

	log.RootLogger().Debugf("Registering manager factory for type: %s", connectionType)

	managerFactories[connectionType] = factory

	return nil
}

func GetManagerFactory(id string) ManagerFactory {
	return managerFactories[id]
}

func RegisterManager(connectionId string, manager Manager) error {

	if connectionId == "" {
		return fmt.Errorf("'id' must be specified when registering")
	}

	if manager == nil {
		return fmt.Errorf("cannot register with 'nil' manager")
	}

	if _, dup := managers[connectionId]; dup {
		return fmt.Errorf("manager already registered: %s", connectionId)
	}

	log.RootLogger().Debugf("Registering manager: %s", connectionId)

	managers[connectionId] = manager

	return nil
}

func GetManager(id string) Manager {
	return managers[id]
}

func StartManagers() error {
	for id, manager := range managers {
		if m, ok:= manager.(managed.Managed); ok {
			err := m.Start()
			if err != nil {
				return fmt.Errorf("unable to start connection manager for '%s': %v", id, err)
			}
		}
	}

	return nil
}

func StopManagers() []error {
	var errors []error
	for id, manager := range managers {
		if m, ok:= manager.(managed.Managed); ok {
			err := m.Stop()
			if err != nil {
				errors = append(errors, fmt.Errorf("unable to start connection manager for '%s': %v", id, err))
			}
		}
	}

	return errors
}
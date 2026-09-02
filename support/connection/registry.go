package connection

import (
	"fmt"

	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/log"
)

var (
	managerFactories = make(map[string]ManagerFactory)
	managers         = make(map[string]Manager)
)

func RegisterManagerFactory(factory ManagerFactory) error {

	if factory == nil {
		return fmt.Errorf("cannot register with 'nil' connection manager factory")
	}

	ref := support.GetRef(factory)

	if _, dup := managers[ref]; dup {
		return fmt.Errorf("connection manager factory '%s' already registered", ref)
	}

	managerFactories[ref] = factory

	log.RootLogger().Debugf("Registering '%s' connection manager factory: %s", factory.Type(), ref)

	return nil
}

func ReplaceManagerFactory(ref string, factory ManagerFactory) error {

	if ref == "" {
		return fmt.Errorf("'ref' must be specified when registering")
	}

	if factory == nil {
		return fmt.Errorf("cannot register with 'nil' connection manager factory")
	}

	managerFactories[ref] = factory

	log.RootLogger().Debugf("Replacing '%s' connection manager factory: %s", factory.Type, ref)

	return nil
}

func GetManagerFactory(ref string) ManagerFactory {
	return managerFactories[ref]
}

func ManagerFactories() map[string]ManagerFactory {
	ret := make(map[string]ManagerFactory, len(managerFactories))
	for id, managerFactory := range managerFactories {
		ret[id] = managerFactory
	}

	return ret
}

func RegisterManager(connectionId string, manager Manager) error {

	if connectionId == "" {
		return fmt.Errorf("'id' must be specified when registering")
	}

	if manager == nil {
		return fmt.Errorf("cannot register with 'nil' manager")
	}

	if _, dup := managers[connectionId]; dup {
		return fmt.Errorf("connection manager already registered: %s", connectionId)
	}

	log.RootLogger().Debugf("Registering connection manager: %s", connectionId)

	managers[connectionId] = manager

	return nil
}

func GetManager(id string) Manager {
	return managers[id]
}

func Managers() map[string]Manager {
	ret := make(map[string]Manager, len(managers))
	for id, manager := range managers {
		ret[id] = manager
	}

	return ret
}

// GetId returns the shared-connection id `manager` was registered under, or "" when the
// manager is not a shared connection (for example an inline connection config). It mirrors
// IsShared, which already performs the same linear scan; `managers` holds one entry per app
// connection, so the scan is single-digit.
//
// A map[Manager]string reverse index was deliberately rejected: indexing by a manager panics
// with "hash of unhashable type" for a value-receiver manager over a struct holding a map or
// slice, which would be an unrecoverable panic at app startup for every app in the org,
// including apps that never use a transaction. The recover below closes the same pre-existing
// hazard IsShared already carries.
func GetId(manager Manager) (id string) {
	if manager == nil {
		return ""
	}

	defer func() {
		if r := recover(); r != nil {
			log.RootLogger().Debugf("connection.GetId: manager is not comparable: %v", r)
			id = ""
		}
	}()

	for cid, mgr := range managers {
		if manager == mgr {
			return cid
		}
	}

	return ""
}

package connection

import (
	"fmt"
	"runtime/debug"

	"github.com/project-flogo/core/support/log"
)

type Manager interface {
	Type() string

	GetConnection() interface{}

	ReleaseConnection(connection interface{})
}

// ReconfigurableConnection allows dynamic update for existing connection instance
type ReconfigurableConnection interface {
	Reconfigure(settings map[string]interface{}) error
}

type ManagerFactory interface {
	Type() string

	NewManager(settings map[string]interface{}) (Manager, error)
}

func NewManager(config *Config) (Manager, error) {

	//resolve settings
	err := ResolveConfig(config)
	if err != nil {
		return nil, err
	}

	f := GetManagerFactory(config.Ref)

	if f == nil {
		return nil, fmt.Errorf("connection factory '%s' not registered", config.Ref)
	}

	cm, err := f.NewManager(config.Settings)
	if err != nil {
		return nil, err
	}

	return cm, err
}

func NewSharedManager(id string, config *Config) (Manager, error) {

	cm, err := NewManager(config)
	if err != nil {
		return nil, err
	}

	err = RegisterManager(id, cm)
	if err != nil {
		return nil, err
	}

	return cm, err
}

func ReconfigureConnections(connections map[string]*Config) (err error) {
	defer func() {
		// Handle panic in implementation code
		if r := recover(); r != nil {
			log.RootLogger().Errorf("Unhandled error while reconfiguring connection: %v", r)
			if log.RootLogger().DebugEnabled() {
				log.RootLogger().Debugf("StackTrace: %s", debug.Stack())
			}
			err = fmt.Errorf("Unhandled error while reconfiguring connection: %v", r)
		}
	}()
	for id, config := range connections {
		manager := managers[id]
		if manager == nil {
			return fmt.Errorf("connection not found for id '%s'", id)
		}
		reconfigurableConn, ok := manager.(ReconfigurableConnection)
		if ok {
			// Resolve connection configuration
			err = ResolveConfig(config)
			if err != nil {
				return err
			}

			// Update existing connection instance
			err = reconfigurableConn.Reconfigure(config.Settings)
			if err != nil {
				return fmt.Errorf("Failed to reconfigure connection: [%s] due to error: %v : ", id, err)
			}
			log.RootLogger().Infof("Connection: %s successfully reconfigured", id)
		}
	}
	return err
}

func IsShared(manager Manager) bool {
	for _, mgr := range managers {
		if manager == mgr {
			return true
		}
	}
	return false
}

package connection

import (
	"fmt"
	
	"github.com/project-flogo/core/support/managed"
)

type Manager interface {
	Type() string
	
	GetConnection() interface{}
	
	ReleaseConnection(connection interface{})
}

// Connections that support dynamic updates
type ReconfigurableManager interface {
	ReconfigureConnection(settings map[string]interface{}) error
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

func ReconfigureManager(id string, config *Config) error {
	manager := managers[id]
	if manager == nil {
		return fmt.Errorf("connection not found for id '%s'", id)
	}
	var err error
	// Resolve connection configuration
	err = ResolveConfig(config)
	if err != nil {
		return err
	}
	
	reconfigurableConn, ok := manager.(ReconfigurableManager)
	if ok {
		// Update existing connection instance
		err = reconfigurableConn.ReconfigureConnection(config.Settings)
		if err != nil {
			return err
		}
	} else {
		// Create new connection instance
		m, isManaged := manager.(managed.Managed)
		if isManaged {
			// Stop previous connection instance
			err = m.Stop()
			if err != nil {
				return err
			}
			// Create new connection instance
			manager, err = NewManager(config)
			if err != nil {
				return err
			}
			// Start new connection instance
			err = manager.(managed.Managed).Start()
			if err != nil {
				return err
			}
			// Replace existing instance with new instance
			managers[id] = manager
		}
	}
	return nil
}

func IsShared(manager Manager) bool {
	for _, mgr := range managers {
		if manager == mgr {
			return true
		}
	}
	return false
}

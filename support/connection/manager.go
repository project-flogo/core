package connection

import (
	"fmt"

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

func Reconfigure(id string, config *Config) error {
	var err error
	manager := managers[id]
	if manager == nil {
		return fmt.Errorf("connection not found for id '%s'", id)
	}
	// Resolve connection configuration
	err = ResolveConfig(config)
	if err != nil {
		return err
	}
	reconfigurableConn, ok := manager.(ReconfigurableConnection)
	if ok {
		// Update existing connection instance
		err = reconfigurableConn.Reconfigure(config.Settings)
		if err != nil {
			return err
		}
		log.RootLogger().Infof("Connection: %s successfully reconfigured", id)
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

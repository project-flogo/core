package connection

import (
	"fmt"
)

type Manager interface {
	Type() string

	GetConnection() interface{}

	ReleaseConnection(connection interface{})
}

type ManagerFactory interface {
	Type() string

	NewManager(settings map[string]interface{}) (Manager, error)
}

func NewManager(config *Config) (Manager, error)  {

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

func NewSharedManager(id string, config *Config) (Manager, error)  {

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

func IsShared(manager Manager) bool{
	for _, mgr := range managers {
		if manager == mgr {
			return true
		}
	}
	return false
}


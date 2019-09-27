package property

import "github.com/project-flogo/core/data/property/tmp"

func init() {
	SetDefaultManager(NewManager(make(map[string]interface{})))
}

var defaultManager *Manager

func SetDefaultManager(manager *Manager) {
	defaultManager = manager

	tmp.SetDefaultManager(manager)
}

func DefaultManager() *Manager {
	return defaultManager
}

func NewManager(properties map[string]interface{}) *Manager {

	manager := &Manager{properties: properties}
	return manager
}

type Manager struct {
	properties map[string]interface{}
}

func (m *Manager) GetProperty(name string) (interface{}, bool) {
	val, exists := m.properties[name]
	return val, exists
}

func (m *Manager) Finalize(processors ...PostProcessor) error {

	for _, processor := range processors {
		_ = processor(m.properties)
	}

	return nil
}

type PostProcessor func(properties map[string]interface{}) error

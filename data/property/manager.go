package property

import (
	"github.com/project-flogo/core/support/log"
)

func init() {
	SetDefaultManager(NewManager(make(map[string]interface{})))
}

var defaultManager *Manager

func SetDefaultManager(manager *Manager) {
	defaultManager = manager
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

func (m *Manager) Finalize(useExternalResolvers bool, processors ...PostProcessor) error {

	logger := log.RootLogger()

	if useExternalResolvers {
		for name := range m.properties {
			newVal, found := ResolveExternally(name)

			if !found {
				logger.Warnf("Property '%s' could not be resolved using external resolver(s) '%s'. Using default value.", name)
			} else {
				m.properties[name] = newVal
			}
		}
	}

	for _, processor := range processors {
		processor(m.properties)
	}

	return nil
}

type PostProcessor func(properties map[string]interface{}) error

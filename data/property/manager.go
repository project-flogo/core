package property

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
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

func (m *Manager) GetProperties() map[string]interface{} {
	return m.properties
}

func (m *Manager) Finalize(processors ...PostProcessor) error {

	for _, processor := range processors {
		_ = processor(m.properties)
	}

	return nil
}

func (m *Manager) UpdateFromRequest(propFromReq map[string]interface{}) error {

	for name := range m.properties {
		newVal, found := propFromReq[name]
		if found {
			// Get datatype of old value
			dType, _ := data.GetType(m.properties[name])
			if dType != data.TypeUnknown {
				coercedVal, err := coerce.ToType(newVal, dType)
				if err == nil {
					m.properties[name] = coercedVal
					continue
				}
			}
			m.properties[name] = newVal
		}
	}

	return nil
}

type PostProcessor func(properties map[string]interface{}) error

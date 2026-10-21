package property

import (
	"strings"

	"github.com/project-flogo/core/support/log"
)

func init() {
	SetDefaultManager(NewManager(make(map[string]Property)))
}

var defaultManager *Manager

func SetDefaultManager(manager *Manager) {
	defaultManager = manager
}

func DefaultManager() *Manager {
	return defaultManager
}

type Property struct {
	Value       interface{}
	Overridable bool
}

func NewManager(properties map[string]Property) *Manager {
	manager := &Manager{properties: properties}

	var readOnly []string
	for name, prop := range properties {
		if !prop.Overridable {
			readOnly = append(readOnly, name)
		}
	}
	if len(readOnly) > 0 {
		log.RootLogger().Infof("Non-overridable (read-only) properties: [%s]", strings.Join(readOnly, ", "))
	}

	return manager
}

type Manager struct {
	properties map[string]Property
}

func (m *Manager) IsOverridable(name string) bool {
	prop, exists := m.properties[name]
	if !exists {
		return true
	}
	return prop.Overridable
}

func (m *Manager) GetProperty(name string) (interface{}, bool) {
	prop, exists := m.properties[name]
	return prop.Value, exists
}

func (m *Manager) GetProperties() map[string]interface{} {
	result := make(map[string]interface{}, len(m.properties))
	for k, v := range m.properties {
		result[k] = v.Value
	}
	return result
}

func (m *Manager) Finalize(processors ...PostProcessor) error {
	valueMap := make(map[string]interface{}, len(m.properties))
	for k, v := range m.properties {
		valueMap[k] = v.Value
	}

	for _, processor := range processors {
		_ = processor(valueMap)
	}

	for k, v := range valueMap {
		if prop, exists := m.properties[k]; exists {
			prop.Value = v
			m.properties[k] = prop
		}
	}

	return nil
}

type PostProcessor func(properties map[string]interface{}) error

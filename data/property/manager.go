package property

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
	overridables  map[string]bool
}

func (m *Manager) SetOverridables(overridables map[string]bool) {
	m.overridables = overridables
}

func (m *Manager) IsOverridable(name string) bool {
	if m.overridables == nil {
		return true
	}
	return m.overridables[name]
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

type PostProcessor func(properties map[string]interface{}) error

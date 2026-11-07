package property

import (
	"strings"
	"sync"
	"sync/atomic"

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

// NewManager returns a Manager that holds a copy of properties: later changes to the map do not affect it.
func NewManager(properties map[string]Property) *Manager {
	snapshot := make(map[string]Property, len(properties))
	var readOnly []string
	for name, prop := range properties {
		snapshot[name] = prop
		if !prop.Overridable {
			readOnly = append(readOnly, name)
		}
	}
	if len(readOnly) > 0 {
		log.RootLogger().Infof("Non-overridable (read-only) properties: [%s]", strings.Join(readOnly, ", "))
	}

	manager := &Manager{}
	manager.properties.Store(snapshot)
	return manager
}

// Manager is safe for concurrent use. Each read sees one whole snapshot, which is never modified once published;
// values are not deep-copied, so callers must not modify a map or slice value.
type Manager struct {
	mu         sync.Mutex   // serialises Finalize; readers never take it
	properties atomic.Value // the current snapshot, a map[string]Property
}

// snapshot returns the current properties, or nil for a zero Manager. The map must not be modified.
func (m *Manager) snapshot() map[string]Property {
	properties, _ := m.properties.Load().(map[string]Property)
	return properties
}

func (m *Manager) IsOverridable(name string) bool {
	prop, exists := m.snapshot()[name]
	if !exists {
		return true
	}
	return prop.Overridable
}

func (m *Manager) GetProperty(name string) (interface{}, bool) {
	prop, exists := m.snapshot()[name]
	return prop.Value, exists
}

// GetProperties returns a copy of the current property values. The caller owns the map.
func (m *Manager) GetProperties() map[string]interface{} {
	properties := m.snapshot()
	result := make(map[string]interface{}, len(properties))
	for k, v := range properties {
		result[k] = v.Value
	}
	return result
}

// Finalize runs the processors in order on a copy of the values, ignoring their errors and any key they add or
// remove, then publishes the result. Calls run one at a time.
func (m *Manager) Finalize(processors ...PostProcessor) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	current := m.snapshot()
	valueMap := make(map[string]interface{}, len(current))
	for k, v := range current {
		valueMap[k] = v.Value
	}

	for _, processor := range processors {
		_ = processor(valueMap)
	}

	next := make(map[string]Property, len(current))
	for k, prop := range current {
		if v, exists := valueMap[k]; exists {
			prop.Value = v
		}
		next[k] = prop
	}
	m.properties.Store(next)

	return nil
}

// PostProcessor may change the map it receives, a private copy, but not a map or slice value in it: the published
// snapshot shares those. It may read the Manager but must not call Finalize, which would deadlock.
type PostProcessor func(properties map[string]interface{}) error

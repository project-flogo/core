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

// Manager holds the app properties. It is safe for concurrent use.
//
// The properties live in a snapshot map that is never modified once published. Readers take no lock, and
// Finalize publishes a new snapshot in one step, so a reader never sees a half-applied refresh (two separate
// reads may still straddle one). Property values are not deep-copied: a map or slice value is shared by every
// snapshot and every GetProperties result, so callers must not modify it.
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

// Finalize runs the processors, in order, on a copy of the current values, then publishes the new values of the
// existing properties: keys that a processor adds are ignored, and keys that it removes keep their value.
// Processor errors are ignored. Calls run one at a time; a processor may read the Manager but must not call
// Finalize.
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

// PostProcessor transforms the property values during Finalize. It receives a private copy of the current values and
// may change that map in place; only keys that already exist are kept. It may read the Manager, but it must not call
// Finalize, which would deadlock, and it must not modify a value itself in place: a map or slice value is shared with
// the published snapshot.
type PostProcessor func(properties map[string]interface{}) error

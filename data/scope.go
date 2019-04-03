package data

import "sync"

// Scope is a set of attributes that are accessible
type Scope interface {
	// GetValue gets the specified value
	GetValue(name string) (value interface{}, exists bool)

	// SetValue sets the specified  value
	SetValue(name string, value interface{}) error
}

// SimpleScope is a basic implementation of a scope
type SimpleScope struct {
	parentScope Scope
	values      map[string]interface{}
}

// NewSimpleScope creates a new SimpleScope
func NewSimpleScope(values map[string]interface{}, parentScope Scope) Scope {

	scope := &SimpleScope{
		parentScope: parentScope,
		values:      make(map[string]interface{}),
	}

	for name, value := range values {
		scope.values[name] = value
	}

	return scope
}

// GetValue implements Scope.GetValue
func (s *SimpleScope) GetValue(name string) (value interface{}, exists bool) {
	value, found := s.values[name]

	if found {
		return value, true
	}

	if s.parentScope != nil {
		return s.parentScope.GetValue(name)
	}

	return nil, false
}

// SetValue implements Scope.SetValue
func (s *SimpleScope) SetValue(name string, value interface{}) error {
	s.values[name] = value
	return nil
}

// SimpleSyncScope is a basic implementation of a synchronized scope
type SimpleSyncScope struct {
	scope Scope
	mutex sync.RWMutex
}

// NewSimpleSyncScope creates a new SimpleSyncScope
func NewSimpleSyncScope(values map[string]interface{}, parentScope Scope) Scope {

	var syncScope SimpleSyncScope
	syncScope.scope = NewSimpleScope(values, parentScope)

	return &syncScope
}

// GetValue implements Scope.GetValue
func (s *SimpleSyncScope) GetValue(name string) (value interface{}, exists bool) {

	s.mutex.RLock()
	v,e := s.scope.GetValue(name)
	s.mutex.RUnlock()

	return v,e
}

// SetValue implements Scope.SetValue
func (s *SimpleSyncScope) SetValue(name string, value interface{}) error {

	s.mutex.Lock()
	err := s.scope.SetValue(name, value)
	s.mutex.Unlock()

	return err
}
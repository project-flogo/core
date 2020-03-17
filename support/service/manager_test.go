package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestService struct {
	name             string
	started, stopped bool
}

func (t *TestService) Start() error {
	t.started = true
	return nil
}

func (t *TestService) Stop() error {
	t.stopped = true
	return nil
}

func (t *TestService) Name() string {
	return t.name
}

func TestManager_FindService(t *testing.T) {

	sm := NewServiceManager()
	assert.NotNil(t, sm)
	sm.services = map[string]Service{"test1": &TestService{name: "test1"}, "test2": &TestService{name: "test2"}}

	s := sm.FindService(func(service Service) bool {
		return service.Name() == "test2"
	})

	assert.NotNil(t, s)
	assert.Equal(t, "test2", s.Name())

	s = sm.FindService(func(service Service) bool {
		return service.Name() == "test3"
	})
	assert.Nil(t, s)
}

func TestManager_GetService(t *testing.T) {

	sm := NewServiceManager()
	assert.NotNil(t, sm)
	ts := &TestService{}
	sm.services = map[string]Service{"test": ts}

	s := sm.GetService("test")
	assert.Equal(t, ts, s)
	s = sm.GetService("test2")
	assert.Nil(t, s)
}

func TestManager_RegisterService(t *testing.T) {

	sm := NewServiceManager()
	assert.NotNil(t, sm)
	ts := &TestService{name: "test"}

	err := sm.RegisterService(ts)
	assert.Nil(t, err)

	s := sm.services["test"]
	assert.Equal(t, ts, s)

	err = sm.RegisterService(ts)
	assert.NotNil(t, err)
}

func TestManager_Services(t *testing.T) {

	sm := NewServiceManager()
	assert.NotNil(t, sm)
	sm.services = map[string]Service{"test1": &TestService{name: "test1"}, "test2": &TestService{name: "test2"}}

	services := sm.Services()
	assert.Len(t, services, 2)
	sName := services[0].Name()
	assert.True(t, sName == "test1" || sName == "test2")
	sName = services[1].Name()
	assert.True(t, sName == "test1" || sName == "test2")
}

func TestManager_Start(t *testing.T) {

	sm := NewServiceManager()
	assert.NotNil(t, sm)
	sm.services = map[string]Service{"test1": &TestService{name: "test1"}, "test2": &TestService{name: "test2"}}

	err := sm.Start()
	assert.Nil(t, err)

	for _, s := range sm.services {
		assert.True(t, s.(*TestService).started)
	}
}

func TestManager_Stop(t *testing.T) {
	sm := NewServiceManager()
	assert.NotNil(t, sm)
	s1 := &TestService{name: "test1"}
	s2 := &TestService{name: "test2"}
	sm.services = map[string]Service{"test1": s1, "test2": s2}
	sm.started = []Service{ s1, s2}

	err := sm.Stop()
	assert.Nil(t, err)

	for _, s := range sm.services {
		assert.True(t, s.(*TestService).stopped)
	}
}

func TestNewServiceManager(t *testing.T) {

	sm := NewServiceManager()
	assert.NotNil(t, sm)
	assert.NotNil(t, sm.services)
}

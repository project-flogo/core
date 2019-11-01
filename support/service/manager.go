package service

import (
	"fmt"
	"sync"

	"github.com/project-flogo/core/support/managed"
)

// Manager is a simple service manager
type Manager struct {
	sync.RWMutex
	services   map[string]Service
	started    []Service
}

// NewServiceManager creates a new Manager
func NewServiceManager() *Manager {

	var manager Manager
	manager.services = make(map[string]Service)

	return &manager
}

// RegisterService registers the specified service
func (sm *Manager) RegisterService(service Service) error {
	sm.Lock()
	defer sm.Unlock()

	if service == nil {
		return fmt.Errorf("cannot register 'nil' service")
	}

	serviceName := service.Name()

	if _, dup := sm.services[serviceName]; dup {
		return fmt.Errorf("service already registered: %s",serviceName)
	}

	sm.services[serviceName] = service

	return nil
}

// Services gets all the registered Service Services
func (sm *Manager) Services() []Service {

	sm.RLock()
	defer sm.RUnlock()

	return sm.allServices()
}

// Services gets all the registered Service Services
func (sm *Manager) allServices() []Service {

	var curServices = sm.services

	list := make([]Service, 0, len(curServices))

	for _, value := range curServices {
		list = append(list, value)
	}

	return list
}

// GetService gets specified Service
func (sm *Manager) GetService(name string) Service {

	sm.RLock()
	defer sm.RUnlock()

	return sm.services[name]
}

// GetService gets specified Service
func (sm *Manager) FindService(f func (Service) bool) Service {

	sm.RLock()
	defer sm.RUnlock()

	for _, service := range sm.services {
		if f(service) {
			return service
		}
	}

	return nil
}


// Start implements util.Managed.Start()
func (sm *Manager) Start() error {

	sm.Lock()
	defer sm.Unlock()

	if len(sm.started) == 0 {
		services := sm.allServices()

		sm.started = make([]Service, 0, len(services))

		for _, service := range services {
			err := managed.Start(service.Name(), service)

			if err == nil {
				sm.started = append(sm.started, service)
			} else {
				return err
			}
		}
	}

	return nil
}

// Stop implements util.Managed.Stop()
func (sm *Manager) Stop() error {

	sm.Lock()
	defer sm.Unlock()

	var err error

	if len(sm.started) > 0 {

		var notStopped []Service

		for _, service := range sm.started {

			err = managed.Stop(service.Name(), service)

			if err != nil {
				notStopped = append(notStopped, service)
			}
		}

		sm.started = notStopped
	}

	return err
}

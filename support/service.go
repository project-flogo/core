package support

import (
	"errors"
	"sync"

	"github.com/project-flogo/core/support/managed"
)

// Service is an interface for defining/managing a service
type Service interface {
	managed.Managed

	Name() string
	Enabled() bool
}

// ServiceConfig is a simple service configuration object
type ServiceConfig struct {
	Name     string            `json:"name"`
	Enabled  bool              `json:"enabled"`
	Settings map[string]string `json:"settings,omitempty"`
}

// ServiceManager is a simple service manager
type ServiceManager struct {
	servicesMu sync.Mutex
	services   map[string]Service
	started    []Service
}

var defaultServiceManager *ServiceManager

func init() {
	defaultServiceManager = NewServiceManager()
}

func GetDefaultServiceManager() *ServiceManager {
	return defaultServiceManager
}

// NewServiceManager creates a new ServiceManager
func NewServiceManager() *ServiceManager {

	var manager ServiceManager
	manager.services = make(map[string]Service)

	return &manager
}

// RegisterService registers the specified service
func (sm *ServiceManager) RegisterService(service Service) error {
	sm.servicesMu.Lock()
	defer sm.servicesMu.Unlock()

	if service == nil {
		panic("ServiceManager.RegisterService: service is nil")
	}

	serviceName := service.Name()

	if _, dup := sm.services[serviceName]; dup {
		return errors.New("service already registered: " + serviceName)
	}

	sm.services[serviceName] = service

	return nil
}

// Services gets all the registered Service Services
func (sm *ServiceManager) Services() []Service {

	sm.servicesMu.Lock()
	defer sm.servicesMu.Unlock()

	return sm.allServices()
}

// Services gets all the registered Service Services
func (sm *ServiceManager) allServices() []Service {

	var curServices = sm.services

	list := make([]Service, 0, len(curServices))

	for _, value := range curServices {
		list = append(list, value)
	}

	return list
}

// GetService gets specified Service
func (sm *ServiceManager) GetService(name string) Service {

	sm.servicesMu.Lock()
	defer sm.servicesMu.Unlock()

	return sm.services[name]
}

// Start implements util.Managed.Start()
func (sm *ServiceManager) Start() error {

	sm.servicesMu.Lock()
	defer sm.servicesMu.Unlock()

	if len(sm.started) == 0 {
		services := sm.allServices()

		sm.started = make([]Service, 0, len(services))

		for _, service := range services {

			if service.Enabled() {
				err := managed.Start(service.Name(), service)

				if err == nil {
					sm.started = append(sm.started, service)
				} else {
					return err
				}
			}
		}
	}

	return nil
}

// Stop implements util.Managed.Stop()
func (sm *ServiceManager) Stop() error {

	sm.servicesMu.Lock()
	defer sm.servicesMu.Unlock()

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

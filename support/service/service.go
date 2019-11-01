package service

import (
	"github.com/project-flogo/core/support/managed"
)

// Service is an interface for defining/managing a service
type Service interface {
	managed.Managed
	Name() string
}

// Config is a simple service configuration object
type Config struct {
	Settings map[string]interface{} `json:"settings,omitempty"`
}

type Factory interface {
	NewService(config *Config) (Service, error)
}

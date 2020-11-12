package app

import (
	"os"
	"strconv"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/app/resource"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/schema"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/trigger"
)

const (
	EnvKeyDelayedAppStopInterval = "FLOGO_APP_DELAYED_STOP_INTERVAL"
	EnvKeyEnableFlowControl      = "FLOGO_APP_ENABLE_FLOW_CONTROL"
)

// Def is the configuration for the App
type Config struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Description string `json:"description"`
	AppModel    string `json:"appModel"`

	Imports     []string                      `json:"imports,omitempty"`
	Properties  []*data.Attribute             `json:"properties,omitempty"`
	Channels    []string                      `json:"channels,omitempty"`
	Triggers    []*trigger.Config             `json:"triggers"`
	Resources   []*resource.Config            `json:"resources,omitempty"`
	Actions     []*action.Config              `json:"actions,omitempty"`
	Schemas     map[string]*schema.Def        `json:"schemas,omitempty"`
	Connections map[string]*connection.Config `json:"connections,omitempty"`
}

func GetDelayedStopInterval() string {
	intervalEnv := os.Getenv(EnvKeyDelayedAppStopInterval)
	if len(intervalEnv) > 0 {
		return intervalEnv
	}
	return ""
}

func EnableFlowControl() bool {
	enable := os.Getenv(EnvKeyEnableFlowControl)
	if len(enable) > 0 {
		b, _ := strconv.ParseBool(enable)
		return b
	}
	return false
}

type LifecycleAware interface {
	OnStartup() error
	OnShutdown() error
}

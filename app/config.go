package app

import (
	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/app/resource"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/schema"
	"github.com/project-flogo/core/trigger"
)

// Def is the configuration for the App
type Config struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Description string `json:"description"`
	AppModel    string `json:"appModel"`

	Imports    []string               `json:"imports,omitempty"`
	Properties []*data.Attribute      `json:"properties,omitempty"`
	Channels   []string               `json:"channels,omitempty"`
	Triggers   []*trigger.Config      `json:"triggers"`
	Resources  []*resource.Config     `json:"resources"`
	Actions    []*action.Config       `json:"actions,omitempty"`
	Schemas    map[string]*schema.Def `json:"schemas,omitempty"`
}

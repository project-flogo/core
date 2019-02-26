package app

import (
	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/app/resource"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/schema"
	"github.com/project-flogo/core/trigger"
)

// Config is the configuration for the App
type Config struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Description string `json:"description"`

	Imports    []string           `json:"imports"`
	Properties []*data.Attribute  `json:"properties"`
	Channels   []string           `json:"channels"`
	Triggers   []*trigger.Config  `json:"triggers"`
	Resources  []*resource.Config `json:"resources"`
	Actions    []*action.Config   `json:"actions"`
	Schemas    []*SchemaConfig    `json:"schemas,omitempty"`
}

type SchemaConfig struct {
	*schema.Def
	Id string `json:"id"`
}
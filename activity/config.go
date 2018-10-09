package activity

import "github.com/project-flogo/core/data"

type Config struct {
	Ref      string                 `json:"ref"`
	Settings map[string]interface{} `json:"settings,omitempty"`
	Input    map[string]interface{} `json:"input,omitempty"`
	Output   map[string]interface{} `json:"output,omitempty"`
}

type ConfigMetadata struct {
	Input  map[string]data.TypedValue `json:"input,omitempty"`
	Output map[string]data.TypedValue `json:"output,omitempty"`
}

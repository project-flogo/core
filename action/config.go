package action

import "encoding/json"

// Config is the configuration for the Action
type Config struct {
	//inline action
	Ref      string                 `json:"ref,omitempty"`
	Type     string                 `json:"type,omitempty"` //an alias to the ref, can be used if imported
	Settings map[string]interface{} `json:"settings,omitempty"`

	//referenced action
	Id string `json:"id,omitempty"`

	//DEPRECATED
	Data json.RawMessage `json:"data,omitempty"`
}

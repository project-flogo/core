package action

import "encoding/json"

// Config is the configuration for the Action
type Config struct {
	//inline action
	Ref      string                 `json:"ref"`
	Type     string                 `json:"type"` //an alias to the ref, can be used if imported
	Settings map[string]interface{} `json:"settings"`

	//referenced action
	Id string `json:"id"`

	//DEPRECATED
	Data json.RawMessage `json:"data"`
}

package action

import "encoding/json"

// Config is the configuration for the Action
type Config struct {
	//inline action
	Ref      string                 `json:"ref"`
	Settings map[string]interface{} `json:"settings"`

	//referenced action
	Id string `json:"id"`

	//DEPRECATED
	Data json.RawMessage `json:"data"`
}

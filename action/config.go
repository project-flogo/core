package action

// Config is the configuration for the Action
type Config struct {
	//inline action
	Ref      string                 `json:"ref,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`

	//referenced action
	Id string `json:"id,omitempty"`

	//DEPRECATED
	Type string `json:"type,omitempty"`
}

package activity

type Config struct {
	Ref      string                 `json:"ref"`
	Settings map[string]interface{} `json:"settings,omitempty"`
	Input    map[string]interface{} `json:"input,omitempty"`
	Output   map[string]interface{} `json:"output,omitempty"`
}

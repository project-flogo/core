package connection

type Config struct {
	Ref      string                 `json:"ref,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
}

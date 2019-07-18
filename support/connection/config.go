package connection

type Config struct {
	Type     string            `json:"ref,omitempty"`
	Settings map[string]string `json:"settings,omitempty"`
}

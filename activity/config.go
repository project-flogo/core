package activity

type Config struct {
	Ref      string                 `json:"ref"`
	Type     string                 `json:"type"` //an alias to the ref, can be used if imported
	Settings map[string]interface{} `json:"settings,omitempty"`
	Input    map[string]interface{} `json:"input,omitempty"`
	Output   map[string]interface{} `json:"output,omitempty"`
	Schemas  *SchemaConfig          `json:"schemas,omitempty"`
}

type SchemaConfig struct {
	Input  map[string]interface{} `json:"input,omitempty"`
	Output map[string]interface{} `json:"output,omitempty"`
}

package resource

import "encoding/json"

type ResourcesConfig struct {
	Resources []*Config `json:"resources"`
}

type Config struct {
	ID   string          `json:"id"`
	Data json.RawMessage `json:"data"`
}

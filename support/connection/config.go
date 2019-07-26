package connection

import "fmt"

type Config struct {
	Ref      string                 `json:"ref,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
}

func ToConfig(config map[string]interface{}) (*Config,error) {

	if v,ok := config["ref"]; ok {
		if ref, ok := v.(string); ok {
			cfg := &Config{}
			cfg.Ref = ref
			if v,ok := config["settings"]; ok {
				if settings, ok := v.(map[string]interface{}); ok {
					cfg.Settings = settings

					return cfg, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("invalid connection config: %+v", config)
}
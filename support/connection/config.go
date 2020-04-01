package connection

import (
	"fmt"

	"github.com/project-flogo/core/app/resolve"
	"github.com/project-flogo/core/support"
)

type Config struct {
	Ref      string                 `json:"ref,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
	resolved bool
}

func ToConfig(config map[string]interface{}) (*Config, error) {

	if v, ok := config["ref"]; ok {
		if ref, ok := v.(string); ok {

			cfg := &Config{Settings: make(map[string]interface{})}
			cfg.Ref = ref

			err := resolveRef(cfg)
			if err != nil {
				return nil, err
			}
			if v, ok := config["settings"]; ok {
				if settings, ok := v.(map[string]interface{}); ok {
					// Resolve property/env value
					for name, value := range settings {
						strVal, ok := value.(string)
						if ok && len(strVal) > 0 && strVal[0] == '=' {
							var err error
							value, err = resolve.Resolve(strVal[1:], nil)
							if err != nil {
								return nil, err
							}
						}
						cfg.Settings[name] = value
					}
					return cfg, nil
				} else if settings, ok := v.([]interface{}); ok {
					//For backward compatible
					for _, v := range settings {
						val, ok := v.(map[string]interface{})
						if ok {
							name, _ := val["name"].(string)
							value := val["value"]
							strVal, ok := value.(string)
							if ok && len(strVal) > 0 && strVal[0] == '=' {
								var err error
								value, err = resolve.Resolve(strVal[1:], nil)
								if err != nil {
									return nil, err
								}
							}
							cfg.Settings[name] = value
						}

					}
					return cfg, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("invalid connection config: %+v", config)
}

func ResolveConfig(config *Config) error {

	if config.resolved {
		return nil
	}

	err := resolveRef(config)
	if err != nil {
		return err
	}

	for name, value := range config.Settings {

		if strVal, ok := value.(string); ok && len(strVal) > 0 && strVal[0] == '=' {
			var err error
			value, err = resolve.Resolve(strVal[1:], nil)
			if err != nil {
				return err
			}

			config.Settings[name] = value
		}
	}

	config.resolved = true
	return nil
}

func resolveRef(config *Config) error {
	if config.Ref[0] == '#' {
		var ok bool
		config.Ref, ok = support.GetAliasRef("connection", config.Ref)
		if !ok {
			return fmt.Errorf("connection '%s' not imported", config.Ref)
		}
	}
	return nil
}

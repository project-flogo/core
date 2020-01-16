package engine

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/project-flogo/core/support"
)

// Config is the configuration for the Engine, assumes all necessary imports have been add to go code
type Config struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`

	StopEngineOnError bool `json:"stopEngineOnError,omitempty"`
	RunnerType        string `json:"runnerType,omitempty"`

	Imports        []string                          `json:"imports,omitempty"`
	ActionSettings map[string]map[string]interface{} `json:"actionSettings,omitempty"`
	Services       []*ServiceConfig                  `json:"services,omitempty"`
}

// ServiceConfig is the configuration for Engine Services
type ServiceConfig struct {
	Ref      string
	Enabled  bool
	Settings map[string]interface{}
}

func LoadEngineConfig(engineJson string, compressed bool) (*Config, error) {

	var jsonBytes []byte

	if engineJson == "" {

		// a json string wasn't provided, so lets lookup the file in path
		configPath := GetFlogoEngineConfigPath()

		if _, err := os.Stat(configPath); err == nil {
			flogo, err := os.Open(configPath)
			if err != nil {
				return nil, err
			}

			jsonBytes, err = ioutil.ReadAll(flogo)
			if err != nil {
				return nil, err
			}
		}
	} else {

		if compressed {
			var err error
			jsonBytes, err = support.DecodeAndUnzip(engineJson)
			if err != nil {
				return nil, err
			}
		} else {
			jsonBytes = []byte(engineJson)
		}
	}

	cfg := &Config{}
	cfg.StopEngineOnError = StopEngineOnError()
	cfg.RunnerType = GetRunnerType()

	if jsonBytes != nil {
		err := json.Unmarshal(jsonBytes, &cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

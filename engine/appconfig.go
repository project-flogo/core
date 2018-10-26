package engine

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/project-flogo/core/app"
	"github.com/project-flogo/core/engine/secret"
	"github.com/project-flogo/core/support"
)

// AppConfigProvider interface to implement to provide the app configuration
type AppConfigProvider interface {
	GetAppConfig() (*app.Config, error)
}

var appName, appVersion string

// Returns name of the application
func GetAppName() string {
	return appName
}

// Returns version of the application
func GetAppVersion() string {
	return appVersion
}

func LoadAppConfig(flogoJson string, compressed bool) (*app.Config, error) {

	var jsonBytes []byte

	if flogoJson == "" {

		// a json string wasn't provided, so lets lookup the file in path
		configPath := GetFlogoConfigPath()

		flogo, err := os.Open(configPath)
		if err != nil {
			return nil, err
		}

		jsonBytes, err = ioutil.ReadAll(flogo)
		if err != nil {
			return nil, err
		}
	} else {

		if compressed {
			var err error
			jsonBytes, err = support.DecodeAndUnzip(flogoJson)
			if err != nil {
				return nil, err
			}
		} else {
			jsonBytes = []byte(flogoJson)
		}
	}

	updated, err := secret.PreProcessConfig(jsonBytes)
	if err != nil {
		return nil, err
	}

	appConfig := &app.Config{}
	err = json.Unmarshal(updated, &appConfig)
	if err != nil {
		return nil, err
	}

	appName = appConfig.Name
	appVersion = appConfig.Version

	return appConfig, nil
}

// DefaultAppConfigProvider returns the default App Config Provider
func DefaultAppConfigProvider() AppConfigProvider {
	return &defaultConfigProvider{}
}

// defaultConfigProvider implementation of AppConfigProvider
type defaultConfigProvider struct {
}

// GetApp returns the app configuration
func (d *defaultConfigProvider) GetAppConfig() (*app.Config, error) {
	return LoadAppConfig("", false)
}

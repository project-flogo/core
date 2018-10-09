package engine

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/project-flogo/core/app"
	"github.com/project-flogo/core/support"
)

func GetAppConfig(flogoJson string, compressed bool) (*app.Config, error) {

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

	updated, err := preProcessConfig(jsonBytes)
	if err != nil {
		return nil, err
	}

	appConfig := &app.Config{}
	err = json.Unmarshal(updated, &appConfig)
	if err != nil {
		return nil, err
	}
	return appConfig, nil
}

func preProcessConfig(appJson []byte) ([]byte, error) {

	re := regexp.MustCompile("SECRET:[^\\\\\"]*")
	for _, match := range re.FindAll(appJson, -1) {
		encodedValue := string(match[7:])
		decodedValue, err := GetSecretValueHandler().DecodeValue(encodedValue)
		if err != nil {
			return nil, err
		}
		appString := strings.Replace(string(appJson), string(match), decodedValue, -1)
		appJson = []byte(appString)
	}

	return appJson, nil
}

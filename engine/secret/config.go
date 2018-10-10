package secret

import (
	"regexp"
	"strings"
)

func PreProcessConfig(appJson []byte) ([]byte, error) {

	// For now decode secret values
	re := regexp.MustCompile("SECRET:[^\\\\\"]*")
	for _, match := range re.FindAll(appJson, -1) {
		decodedValue, err := resolveSecretValue(string(match))
		if err != nil {
			return nil, err
		}
		appstring := strings.Replace(string(appJson), string(match), decodedValue, -1)
		appJson = []byte(appstring)
	}

	return appJson, nil
}

func resolveSecretValue(encrypted string) (string, error) {
	encodedValue := string(encrypted[7:])
	decodedValue, err := GetSecretValueHandler().DecodeValue(encodedValue)
	if err != nil {
		return "", err
	}
	return decodedValue, nil
}

func PropertyProcessor(properties map[string]interface{}) error {

	for key, value := range properties {

		if strVal, ok := value.(string); ok && strings.HasPrefix(strVal, "SECRET:") {

			// Resolve secret value
			newVal, err := resolveSecretValue(strVal)
			if err != nil {
				return err
			}
			properties[key] = newVal
		}
	}

	return nil
}

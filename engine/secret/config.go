package secret

import (
	"strings"
)

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

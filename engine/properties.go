package engine

import (
	"fmt"
	"os"
	"strings"
)

func EnvPropertyProcessor(properties map[string]interface{}) error {

	for key, value := range properties {

		if strVal, ok := value.(string); ok && strings.HasPrefix(strVal, "$env[") {

			envVar := strVal[5 : len(strVal)-1]
			newVal, exists := os.LookupEnv(envVar)
			if !exists {
				err := fmt.Errorf("failed to resolve Environment Variable: '%s', ensure that variable is configured", envVar)
				return err
			}
			properties[key] = newVal
		}
	}

	return nil
}

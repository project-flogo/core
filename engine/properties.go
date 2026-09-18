package engine

import (
	"fmt"
	"os"
	"strings"

	"github.com/project-flogo/core/engine/secret"
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
			// Engine variables marked as 'password' are injected as "SECRET:"-prefixed
			// ciphertext; decrypt them transparently before storing the property value.
			decoded, err := secret.Decode(newVal)
			if err != nil {
				return fmt.Errorf("failed to decrypt Environment Variable '%s': %w", envVar, err)
			}
			properties[key] = decoded
		}
	}

	return nil
}

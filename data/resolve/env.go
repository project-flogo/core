package resolve

import (
	"fmt"
	"os"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/engine/secret"
)

var envResolverInfo = NewResolverInfo(true, true)

type EnvResolver struct {
}

func (*EnvResolver) GetResolverInfo() *ResolverInfo {
	return envResolverInfo
}

//EnvResolver Environment Resolver $env[item]
func (*EnvResolver) Resolve(scope data.Scope, item string, field string) (interface{}, error) {
	// Environment resolution
	value, exists := os.LookupEnv(item)
	if !exists {
		err := fmt.Errorf("failed to resolve Environment Variable: '%s', ensure that variable is configured", item)
		return "", err
	}

	// Engine variables marked as 'password' are injected as "SECRET:"-prefixed
	// ciphertext; decrypt them transparently so callers get the plaintext value.
	decoded, err := secret.Decode(value)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt Environment Variable '%s': %w", item, err)
	}

	return decoded, nil
}

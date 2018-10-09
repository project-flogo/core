package resolve

import (
	"fmt"
	"os"

	"github.com/project-flogo/core/data"
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

	return value, nil
}

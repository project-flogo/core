package resolve

import (
	"fmt"

	"github.com/project-flogo/core/data"
)

var scopeResolverInfo = NewResolverInfo(false, false)

type ScopeResolver struct {
}

func (*ScopeResolver) GetResolverInfo() *ResolverInfo {
	return scopeResolverInfo
}

//ScopeResolver Scope Resolver $.
func (*ScopeResolver) Resolve(scope data.Scope, item string, field string) (interface{}, error) {
	// Scope resolution
	value, exists := scope.GetValue(field)
	if !exists {
		err := fmt.Errorf("failed to resolve variable: '%s' in scope", field)
		return "", err
	}

	return value, nil
}

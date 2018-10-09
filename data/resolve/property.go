package resolve

import (
	"fmt"

	"github.com/project-flogo/core/data"
)

var propertyResolverInfo = NewResolverInfo(true, true)

type PropertyResolver struct {
}

func (*PropertyResolver) GetResolverInfo() *ResolverInfo {
	return propertyResolverInfo
}

//PropertyResolver Property Resolver $property[item]
func (*PropertyResolver) Resolve(scope data.Scope, item string, field string) (interface{}, error) {
	provider := data.GetPropertyProvider()
	value, exists := provider.GetProperty(item) //should we add the path and reset it to ""
	if !exists {
		return nil, fmt.Errorf("failed to resolve Property: '%s', ensure that property is configured in the application", item)
	}

	return value, nil
}

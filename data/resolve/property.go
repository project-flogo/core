package resolve

import (
	"fmt"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/property/tmp"
)


var propertyResolverInfo = NewResolverInfo(true, true)

//DEPRECATED
type PropertyResolver struct {
}

func (*PropertyResolver) GetResolverInfo() *ResolverInfo {
	return propertyResolverInfo
}

//Resolver Property Resolver $property[item]
func (*PropertyResolver) Resolve(scope data.Scope, item string, field string) (interface{}, error) {

	manager := tmp.DefaultManager()
	value, exists := manager.GetProperty(item) //should we add the path and reset it to ""
	if !exists {
		return nil, fmt.Errorf("failed to resolve Property: '%s', ensure that property is configured in the application", item)
	}

	return value, nil
}

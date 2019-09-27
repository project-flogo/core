package property

import (
	"fmt"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/resolve"
)

var propertyResolverInfo = resolve.NewResolverInfo(true, true)

type Resolver struct {
}

func (*Resolver) GetResolverInfo() *resolve.ResolverInfo {
	return propertyResolverInfo
}

//Resolver Property Resolver $property[item]
func (*Resolver) Resolve(scope data.Scope, item string, field string) (interface{}, error) {

	manager := DefaultManager()
	value, exists := manager.GetProperty(item) //should we add the path and reset it to ""
	if !exists {
		return nil, fmt.Errorf("failed to resolve Property: '%s', ensure that property is configured in the application", item)
	}

	return value, nil
}

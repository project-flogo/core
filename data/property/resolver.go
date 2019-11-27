package property

import (
	"fmt"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/resolve"
)

var propertyResolverInfo = resolve.NewResolverInfo(!IsPropertySnapshotEnabled(), true)

type Resolver struct {
}

func (*Resolver) GetResolverInfo() *resolve.ResolverInfo {
	return propertyResolverInfo
}

//Resolver Property Resolver $property[item]
func (*Resolver) Resolve(scope data.Scope, item string, field string) (interface{}, error) {

	// Resolve property value from the local copy of instance
	if !propertyResolverInfo.IsStatic() && scope != nil {
		props, found := scope.GetValue("_PROPERTIES")
		if found {
			pMap, ok := props.(map[string]interface{})
			if ok {
				pValue, exists := pMap[item]
				if exists {
					return pValue, nil
				}
			}
		}
	}

	manager := DefaultManager()
	value, exists := manager.GetProperty(item) //should we add the path and reset it to ""
	if !exists {
		return nil, fmt.Errorf("failed to resolve Property: '%s', ensure that property is configured in the application", item)
	}

	return value, nil
}

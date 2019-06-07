package resolve

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/path"
)

var loopResolverInfo = NewImplicitResolverInfo(false, true)

type LoopResolver struct {
}

func (*LoopResolver) GetResolverInfo() *ResolverInfo {
	return loopResolverInfo
}

//LoopResolver Loop Resolver $Loop[item]
func (*LoopResolver) Resolve(scope data.Scope, item string, field string) (interface{}, error) {
	if item == "" {
		v, exist := scope.GetValue(field)
		if !exist {
			return nil, fmt.Errorf("failed to resolve current Loop: '%s', ensure that Loop is configured in the application", field)
		}
		return v, nil
	} else {
		value, exists := scope.GetValue(item)
		if !exists {
			return nil, fmt.Errorf("failed to resolve Loop: '%s', ensure that Loop is configured in the application", item)
		}
		return path.GetValue(value, "."+field)
	}
}

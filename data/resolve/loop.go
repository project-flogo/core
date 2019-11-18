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
	var value interface{}
	var exist bool
	if item == "" {
		value, exist = scope.GetValue("_loop")
		if !exist {
			return nil, fmt.Errorf("failed to resolve current Loop: '%s', ensure that Loop is configured in the application", field)
		}
	} else {
		value, exist = scope.GetValue(item)
		if !exist {
			return nil, fmt.Errorf("failed to resolve Loop: '%s', ensure that Loop is configured in the application", item)
		}
	}
	if field != "" {
		return path.GetValue(value, "."+field)
	} else {
		return value, nil
	}
}

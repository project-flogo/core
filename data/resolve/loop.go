package resolve

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/path"
)

var loopResolverInfo = NewResolverInfo(false, true)

type LoopResolver struct {
}

func (*LoopResolver) GetResolverInfo() *ResolverInfo {
	return loopResolverInfo
}

//LoopResolver Loop Resolver $Loop[item]
func (*LoopResolver) Resolve(scope data.Scope, item string, field string) (interface{}, error) {
	value, exists := scope.GetValue(item)
	if !exists {
		return nil, fmt.Errorf("failed to resolve Loop: '%s', ensure that Loop is configured in the application", item)
	}
	return path.GetValue(value, "."+field)
}

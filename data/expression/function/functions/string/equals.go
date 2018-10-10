package string

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnEquals{})
}

type fnEquals struct {
}

func (fnEquals) Name() string {
	return "string.equals"
}

func (fnEquals) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString}, false
}

func (fnEquals) Eval(params ...interface{}) (interface{}, error) {

	s1 := params[0].(string)
	s2 := params[1].(string)
	return s1 == s2, nil
}

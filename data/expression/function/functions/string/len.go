package string

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnLen{})
}

type fnLen struct {
}

func (fnLen) Name() string {
	return "string.len"
}

func (fnLen) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

func (fnLen) Eval(params ...interface{}) (interface{}, error) {

	s := params[0].(string)

	return len(s), nil
}

package string

import (
	"strings"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

type fnEqualsIgnoreCase struct {
}

func init() {
	function.Register(&fnEqualsIgnoreCase{})
}

func (s *fnEqualsIgnoreCase) Name() string {
	return "string.equalsIgnoreCase"
}

func (fnEqualsIgnoreCase) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString}, false
}

func (fnEqualsIgnoreCase) Eval(params ...interface{}) (interface{}, error) {
	str1 := params[0].(string)
	str2 := params[1].(string)
	return strings.EqualFold(str1, str2), nil
}

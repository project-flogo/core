package builtin

import (
	"reflect"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	_ = function.Register(&fnLen{})
}

type fnLen struct {
}

func (*fnLen) Name() string {
	return "len"
}

func (*fnLen) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeAny}, false
}

func (*fnLen) Eval(params ...interface{}) (interface{}, error) {

	switch t := params[0].(type) {
	case string:
		return len(t), nil
	case nil:
		return 0, nil
	default:
		return reflect.ValueOf(t).Len(), nil

		//return 0, fmt.Errorf("Unable to coerce %#v to integer", val)
	}
}

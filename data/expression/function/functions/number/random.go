package string

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"math/rand"
	"time"
)

func init() {
	function.Register(&fnRandom{})
}

type fnRandom struct {
}

func (fnRandom) Name() string {
	return "number.random"
}

func (fnRandom) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeInt}, true
}

func (fnRandom) Eval(params ...interface{}) (interface{}, error) {

	limit := 10
	if len(params) > 0 {
		limit = params[0].(int)
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(limit), nil
}

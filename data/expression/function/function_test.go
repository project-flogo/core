package function

import (
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/stretchr/testify/assert"
)

type SimpleFunction struct {
}

func (s *SimpleFunction) Name() string {
	return "simple"
}

func (s *SimpleFunction) Sig() (params []data.Type, isVariadic bool) {
	return []data.Type{data.TypeInt}, false
}

func (s *SimpleFunction) Eval(params ...interface{}) (interface{}, error) {
	val := params[0].(int)
	return val, nil
}

func TestRegister(t *testing.T) {

	var err error
	//Register Nil
	err = Register(nil)
	assert.NotNil(t, err)

	err = Register(&SimpleFunction{})
	assert.Nil(t, err)

	ResolveAliases()

	//Registering duplicate
	err = Register(&SimpleFunction{})
	assert.NotNil(t, err)
}

func TestFunction(t *testing.T) {

	val, err := Eval(&SimpleFunction{}, 2)
	assert.Equal(t, 2, val.(int))
	assert.Nil(t, err)

	_, err = Eval(&SimpleFunction{})
	assert.NotNil(t, err)
}

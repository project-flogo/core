package script

import (
	"bytes"
	"fmt"
	"github.com/project-flogo/core/data/resolve"
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/stretchr/testify/assert"
)

func TestBuiltinFuncExpr(t *testing.T) {

	expr, err := factory.NewExpr(`builtin.len("test")`)
	assert.Nil(t, err)

	v, err := expr.Eval(nil)
	assert.Nil(t, err)

	assert.Equal(t, 4, v)
}

func TestFuncExprNoSpace(t *testing.T) {
	expr, err := factory.NewExpr(`script.concat("a","b")`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, "ab", v)
}

func TestFuncExprNested(t *testing.T) {

	expr, err := factory.NewExpr(`script.concat("This", "is",script.concat("my","first"),"gocc",script.concat("lexer","and","parser"),script.concat("go","program","!!!"))`)
	assert.Nil(t, err)

	v, err := expr.Eval(nil)
	assert.Nil(t, err)

	assert.Equal(t, "Thisismyfirstgocclexerandparsergoprogram!!!", v.(string))
}

func TestFuncExprNestedMultiSpace(t *testing.T) {

	expr, err := factory.NewExpr(`script.concat("This",   " is" , " Flogo")`)
	assert.Nil(t, err)

	v, err := expr.Eval(nil)
	assert.Nil(t, err)

	assert.Equal(t, "This is Flogo", v.(string))
}

func TestFunctionWithRef(t *testing.T) {

	scope := data.NewSimpleScope(map[string]interface{}{"queryParams": map[string]interface{}{"id": "flogo"}}, nil)
	factory := NewExprFactory(resolve.GetBasicResolver())
	testcases := make(map[string]interface{})
	testcases[`script.concat("This", " is ", $.queryParams.id)`] = "This is flogo"

	for k, v := range testcases {
		vv, err := factory.NewExpr(k)
		assert.Nil(t, err)
		result, err := vv.Eval(scope)
		assert.Nil(t, err)
		if !assert.ObjectsAreEqual(v, result) {
			assert.Fail(t, fmt.Sprintf("test function [%s] failed, expected [%+v] but actual [%+v]", k, v, result))
		}
	}

}

func init() {
	_ = function.Register(&fnConcat{})
}

type fnConcat struct {
}

func (fnConcat) Name() string {
	return "concat"
}

func (fnConcat) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, true
}

func (fnConcat) Eval(params ...interface{}) (interface{}, error) {
	if len(params) >= 2 {
		var buffer bytes.Buffer

		for _, v := range params {
			buffer.WriteString(v.(string))
		}
		return buffer.String(), nil
	}

	return "", fmt.Errorf("fnConcat function must have at least two arguments")
}

func TestFuncExprSingleQuote(t *testing.T) {
	expr, err := factory.NewExpr("script.concat('abc','def')")
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, "abcdef", v)
}

func init() {
	function.Register(&tLength{})
	function.ResolveAliases()

}

type tLength struct {
}

func (tLength) Name() string {
	return "length"
}

func (tLength) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

func (tLength) Eval(params ...interface{}) (interface{}, error) {
	p := params[0].(string)
	return len(p), nil
}

package script

import (
	"testing"

	_ "github.com/project-flogo/core/data/expression/script/function/functions/number"
	_ "github.com/project-flogo/core/data/expression/script/function/functions/string"
	"github.com/stretchr/testify/assert"
)

func TestBuiltinFuncExpr(t *testing.T) {

	expr, err := factory.NewExpr(`len("test")`)
	assert.Nil(t, err)

	v, err := expr.Eval(nil)
	assert.Nil(t, err)

	assert.Equal(t, 4, v)
}

func TestFuncExprNoSpace(t *testing.T) {
	expr, err := factory.NewExpr(`string.concat("a","b")`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, "ab", v)
}

func TestFuncExprNested(t *testing.T) {

	expr, err := factory.NewExpr(`string.concat("This", "is",string.concat("my","first"),"gocc",string.concat("lexer","and","parser"),string.concat("go","program","!!!"))`)
	assert.Nil(t, err)

	v, err := expr.Eval(nil)
	assert.Nil(t, err)

	assert.Equal(t, "Thisismyfirstgocclexerandparsergoprogram!!!", v.(string))
}

func TestFuncExprNestedMultiSpace(t *testing.T) {

	expr, err := factory.NewExpr(`string.concat("This",   " is" , " Flogo")`)
	assert.Nil(t, err)

	v, err := expr.Eval(nil)
	assert.Nil(t, err)

	assert.Equal(t, "This is Flogo", v.(string))
}

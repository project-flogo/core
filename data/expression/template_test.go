package expression

import (
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/stretchr/testify/assert"
)

type testExprFactory struct {
}

func (f *testExprFactory) NewExpr(exprStr string) (Expr, error) {
	return &testExpr{str: exprStr}, nil
}

//stringExpr simple wrapper for a string
type testExpr struct {
	str string
}

func (s *testExpr) Eval(scope data.Scope) (interface{}, error) {
	return s.str, nil
}

func TestParse(t *testing.T) {
	f := &testExprFactory{}

	s := `abcde{x{123}}fg{{h}}i`
	ss, err := parse(s, f)
	assert.Nil(t, err)
	assert.Len(t, ss, 3)
	v, _ := ss[0].Eval(nil)
	assert.Equal(t, `abcde{x{123}}fg`, v)
	v, _ = ss[1].Eval(nil)
	assert.Equal(t, `h`, v)
	v, _ = ss[2].Eval(nil)
	assert.Equal(t, `i`, v)

	s = `{{abcde{x{123}}`
	ss, err = parse(s, f)
	assert.Nil(t, err)
	assert.Len(t, ss, 1)
	v, _ = ss[0].Eval(nil)
	assert.Equal(t, `abcde{x{123`, v)

	s = `{{abcde{x{123}}fg{{h}}i`
	ss, err = parse(s, f)
	assert.Nil(t, err)
	assert.Len(t, ss, 4)
	v, _ = ss[0].Eval(nil)
	assert.Equal(t, `abcde{x{123`, v)
	v, _ = ss[1].Eval(nil)
	assert.Equal(t, `fg`, v)
	v, _ = ss[2].Eval(nil)
	assert.Equal(t, `h`, v)
	v, _ = ss[3].Eval(nil)
	assert.Equal(t, `i`, v)

	s = `{{abcde{x{123`
	ss, err = parse(s, f)
	assert.Nil(t, err)
	assert.Len(t, ss, 1)
	v, _ = ss[0].Eval(nil)
	assert.Equal(t, `{{abcde{x{123`, v)

	s = `abcde{x{123}}`
	ss, err = parse(s, f)
	assert.Nil(t, err)
	assert.Len(t, ss, 1)
	v, _ = ss[0].Eval(nil)
	assert.Equal(t, `abcde{x{123}}`, v)

	s = `abcde{{x{123}`
	ss, err = parse(s, f)
	assert.Nil(t, err)
	assert.Len(t, ss, 1)
	v, _ = ss[0].Eval(nil)
	assert.Equal(t, `abcde{{x{123}`, v)
}

//func TestMapObject(t *testing.T) {
//	str := "5"
//	isExpr := expression.IsExpression(str)
//
//	fmt.Printf("result: %v", isExpr)
//
//	_, err := strconv.Atoi("2.2")
//	fmt.Println(err)
//}

//='blah'

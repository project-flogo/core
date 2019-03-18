package script

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/resolve"
	"github.com/stretchr/testify/assert"
)

var resolver = resolve.NewCompositeResolver(map[string]resolve.Resolver{"static": &TestStaticResolver{}, ".": &TestResolver{}, "env": &resolve.EnvResolver{}})
var factory = NewExprFactory(resolver)

func TestLitExprInt(t *testing.T) {
	expr, err := factory.NewExpr(`123`)
	assert.Nil(t, err)
	assert.NotNil(t, expr)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 123, v)
}

func TestLitExprFloat(t *testing.T) {
	expr, err := factory.NewExpr(`123.5`)
	assert.Nil(t, err)
	assert.NotNil(t, expr)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 123.5, v)
}

func TestLitExprBool(t *testing.T) {
	expr, err := factory.NewExpr(`true`)
	assert.Nil(t, err)
	assert.NotNil(t, expr)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)
}

func TestLitExprStringSQ(t *testing.T) {
	expr, err := factory.NewExpr(`'foo bar'`)
	assert.Nil(t, err)
	assert.NotNil(t, expr)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, "foo bar", v)
}

func TestLitExprStringDQ(t *testing.T) {
	expr, err := factory.NewExpr(`"foo bar"`)
	assert.Nil(t, err)
	assert.NotNil(t, expr)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, "foo bar", v)
}

func TestLitExprNil(t *testing.T) {
	expr, err := factory.NewExpr(`nil`)
	assert.Nil(t, err)
	assert.NotNil(t, expr)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Nil(t, v)

	expr, err = factory.NewExpr(`null`)
	assert.Nil(t, err)
	assert.NotNil(t, expr)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Nil(t, v)
}

func TestLitExprRef(t *testing.T) {

	expr, err := factory.NewExpr(`$.foo`)
	assert.Nil(t, err)
	assert.NotNil(t, expr)

	scope := newScope(map[string]interface{}{"foo": "bar"})
	v, err := expr.Eval(scope)
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)
}

const testJsonData = `{
    "store": {
        "book": [
            {
                "category": "reference",
                "author": "Nigel Rees",
                "title": "Sayings of the Century",
                "price": 8.95
            },
            {
                "category": "fiction",
                "author": "Evelyn Waugh",
                "title": "Sword of Honour",
                "price": 12.99
            },
            {
                "category": "fiction",
                "author": "Herman Melville",
                "title": "Moby Dick",
                "isbn": "0-553-21311-3",
                "price": 8.99
            },
            {
                "category": "fiction",
                "author": "J. R. R. Tolkien",
                "title": "The Lord of the Rings",
                "isbn": "0-395-19395-8",
                "price": 22.99
            }
        ],
        "bicycle": {
            "color": "red",
            "price": 19.95
        }
    },
    "expensive": 10
}`

func TestJsonExpr(t *testing.T) {
	var data interface{}
	err := json.Unmarshal([]byte(testJsonData), &data)
	assert.Nil(t, err)

	scope := newScope(map[string]interface{}{"foo": data})

	expr, err := factory.NewExpr("$.foo.store.book[0].price")
	assert.Nil(t, err)
	assert.NotNil(t, expr)

	v, err := expr.Eval(scope)
	assert.Nil(t, err)
	assert.Equal(t, 8.95, v)
}

func TestLitExprStaticRef(t *testing.T) {

	expr, err := factory.NewExpr(`$static.foo`)
	assert.Nil(t, err)
	assert.NotNil(t, expr)

	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)
}

func TestEnvResolve(t *testing.T) {

	os.Setenv("FOO", "bar")
	expr, err := factory.NewExpr(`$env[FOO]`)
	assert.Nil(t, err)
	assert.NotNil(t, expr)

	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)
}

func TestCmpExprEq(t *testing.T) {
	expr, err := factory.NewExpr(`123==123`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`123==321`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`123==123.0`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`123==123.5`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`"foo"=="foo"`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`"foo"=='foo'`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`"foo"=="bar"`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)
}

func TestCmpExprNotEq(t *testing.T) {
	expr, err := factory.NewExpr(`123!=123`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`123!=321`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`123!=123.0`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`123!=123.5`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`"foo"!="foo"`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`"foo"!='foo'`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`"foo"!="bar"`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)
}

func TestCmpExprLt(t *testing.T) {
	expr, err := factory.NewExpr(`123<123`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`123<321`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`123<123.0`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`123<123.5`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`123.5<123`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`"ab"<"ac"`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)
}

func TestCmpExprLtEq(t *testing.T) {
	expr, err := factory.NewExpr(`123<=123`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`123<=321`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`123<=123.0`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`123<=123.5`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`123.5<=123`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`"ab"<="ac"`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)
}

func TestCmpExprGt(t *testing.T) {
	expr, err := factory.NewExpr(`123>123`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`123>321`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`123>123.0`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`123>123.5`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`123.5>123`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`"ab">"ac"`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)
}

func TestCmpExprGtEq(t *testing.T) {
	expr, err := factory.NewExpr(`123>=123`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`123>=321`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`123>=123.0`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`123>=123.5`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`123.5>=123`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`"ab">="ac"`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)
}

func TestArithExprAdd(t *testing.T) {
	expr, err := factory.NewExpr(`12+13`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 25, v)

	expr, err = factory.NewExpr(`12+13.5`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 25.5, v)

	expr, err = factory.NewExpr(`"foo"+'bar'`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, "foobar", v)
}

func TestArithExprSub(t *testing.T) {
	expr, err := factory.NewExpr(`13-12`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, v)

	expr, err = factory.NewExpr(`12-13`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, -1, v)

	expr, err = factory.NewExpr(`13.5-12`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 1.5, v)
}

func TestArithExprMul(t *testing.T) {
	expr, err := factory.NewExpr(`2*5`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 10, v)

	//expr, err = factory.NewExpr(`2*-5`)
	//assert.Nil(t, err)
	//v, err = expr.Eval(nil)
	//assert.Nil(t, err)
	//assert.Equal(t, -10, v)

	expr, err = factory.NewExpr(`2*.1`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, .2, v)
}

func TestArithExprDiv(t *testing.T) {
	expr, err := factory.NewExpr(`10/2`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 5, v)

	expr, err = factory.NewExpr(`2/10`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, v)

	expr, err = factory.NewExpr(`2/.5`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 4.0, v)
}

func TestArithExprMod(t *testing.T) {
	expr, err := factory.NewExpr(`10%2`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, v)

	expr, err = factory.NewExpr(`10%3`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, v)

	expr, err = factory.NewExpr(`10.5%2`) //todo should we throw an error in this case?
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, v)
}

func TestBoolExprOr(t *testing.T) {
	expr, err := factory.NewExpr(`true || false`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`false || false`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`1 || 0`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`0 || 0`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)
}

func TestBoolExprAnd(t *testing.T) {
	expr, err := factory.NewExpr(`true && true`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`true && false`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`1 && 1`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(`1 && 0`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)
}

func TestArithPrecedence(t *testing.T) {
	expr, err := factory.NewExpr(`1+5*2`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 11, v)

	expr, err = factory.NewExpr(`1+5/2.0`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 3.5, v)

	expr, err = factory.NewExpr(`6/2+1*2`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 5, v)

	expr, err = factory.NewExpr(`1+5%2`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 2, v)
}

func TestArithParen(t *testing.T) {
	expr, err := factory.NewExpr(`(1+5)*2`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 12, v)

	expr, err = factory.NewExpr(`10/(5-3)`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 5, v)

	expr, err = factory.NewExpr(`11/(5-3.0)`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 5.5, v)
}

func TestUnaryExpr(t *testing.T) {
	expr, err := factory.NewExpr(`-1`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, -1, v)

	expr, err = factory.NewExpr(`-1.5`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, -1.5, v)

	expr, err = factory.NewExpr(`!true`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`!false`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)
}

func TestUnaryExprComplex(t *testing.T) {
	expr, err := factory.NewExpr(`-(1+2)`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, -3, v)

	expr, err = factory.NewExpr(`-1*2`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, -2, v)

	expr, err = factory.NewExpr(`2*-1`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, -2, v)

	expr, err = factory.NewExpr(`!(false||true)`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(`!(false&&true)`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)
}

func TestTernaryExpr(t *testing.T) {
	expr, err := factory.NewExpr(` 1<2 ? 10 : 20`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 10, v)

	expr, err = factory.NewExpr(`4>3 ? 40 :30`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 40, v)
}

func TestExpression(t *testing.T) {

	scope := data.NewSimpleScope(map[string]interface{}{"queryParams": map[string]interface{}{"id": "helloworld"}}, nil)
	factory := NewExprFactory(resolve.GetBasicResolver())
	os.Setenv("name", "flogo")
	os.Setenv("address", "tibco")

	testcases := make(map[string]interface{})
	testcases[`1>2?tstring.concat("sss","ddddd"):"fff"`] = "fff"
	testcases[`1<2?"helloworld":"fff"`] = "helloworld"
	testcases["200>100?true:false"] = true
	testcases["1 + 2 * 3 + 2 * 6"] = 19
	testcases[`tstring.length($.queryParams.id) == 0 ? "Query Id cannot be null" : tstring.length($.queryParams.id)`] = 10
	testcases[`tstring.length("helloworld")>11?"helloworld":"fff"`] = "fff"
	testcases["123==456"] = false
	testcases["123==123"] = true
	testcases[`tstring.concat("123","456")=="123456"`] = true
	testcases[`tstring.concat("123","456") == tstring.concat("12","3456")`] = true
	testcases[`("dddddd" == "dddd3dd") && ("133" == "123")`] = false
	testcases[`tstring.length("helloworld") == 10`] = true
	testcases[`tstring.length("helloworld") > 10`] = false
	testcases[`tstring.length("helloworld") >= 10`] = true
	testcases[`tstring.length("helloworld") < 10`] = false
	testcases[`tstring.length("helloworld") >= 10`] = true
	testcases[`(tstring.length("sea") == 3) == true`] = true

	testcases[`(1&&1)==(1&&1)`] = true
	testcases[`(true && true) == false`] = false
	testcases[`nil==nil`] = true

	//Nested Ternary
	testcases[`(tstring.length("1234") == 4 ? true : false) ? (2 >1 ? (3>2?"Yes":"nono"):"No") : "false"`] = "Yes"
	testcases[`(4 == 4 ? true : false) ? "yes" : "no"`] = "yes"
	testcases[`(4 == 4 ? true : false) ? 4 < 3 ? "good" :"false" : "no"`] = "false"
	testcases[`4 > 3 ? 6<4 ?  "good2" : "false2" : "false"`] = "false2"
	testcases[`4 > 5 ? 6<4 ?  "good2" : "false2" : 3>2?"ok":"notok"`] = "ok"

	//Int vs float
	testcases[`1 == 1.23`] = false
	testcases[`1 < 1.23`] = true
	testcases[`1.23 == 1`] = false
	testcases[`1.23 > 1`] = true

	//Operator
	testcases[`1 + 2 * 3 + 2 * 6 / 2`] = 13
	testcases[` 1 + 4 * 5 + -6 `] = 15
	testcases[` 2 < 3 && 5 > 4 && 6 < 7 && 56 > 44`] = true
	testcases[` 2 < 3 && 5 > 4 ||  6 < 7 && 56 < 44`] = true
	testcases[`3-2`] = 1
	testcases[`3 - 2`] = 1
	testcases[`3+-2`] = 1
	testcases[`3- -2`] = 5

	//testcases[`tstring.length("helloworld")>11?$env[name]:$env[address]`] = "tibco"
	//testcases[`$env[name] != nil`] = true
	//testcases[`$env[name] == "flogo"`] = true

	for k, v := range testcases {
		vv, err := factory.NewExpr(k)
		assert.Nil(t, err)
		result, err := vv.Eval(scope)
		assert.Nil(t, err)
		if !assert.ObjectsAreEqual(v, result) {
			assert.Fail(t, fmt.Sprintf("test expr [%s] failed, expected [%+v] but actual [%+v]", k, v, result))
		}
	}
}

var result interface{}

func BenchmarkLit(b *testing.B) {
	var r interface{}

	expr, _ := factory.NewExpr(`123`)

	for n := 0; n < b.N; n++ {

		r, _ = expr.Eval(nil)
	}
	result = r
}

/////////////////////////
// Resolver Helpers

func newScope(values map[string]interface{}) data.Scope {
	return &TestScope{values: values}
}

type TestScope struct {
	values map[string]interface{}
}

func (s *TestScope) GetValue(name string) (value interface{}, exists bool) {
	value, exists = s.values[name]
	return
}

func (TestScope) SetValue(name string, value interface{}) error {
	//ignore
	return nil
}

type TestResolver struct {
}

func (*TestResolver) GetResolverInfo() *resolve.ResolverInfo {
	return resolve.NewResolverInfo(false, false)
}

func (*TestResolver) Resolve(scope data.Scope, item string, field string) (interface{}, error) {

	value, exists := scope.GetValue(field)
	if !exists {
		err := fmt.Errorf("failed to resolve variable: '%s', not in scope", field)
		return "", err
	}

	return value, nil
}

type TestStaticResolver struct {
}

func (*TestStaticResolver) GetResolverInfo() *resolve.ResolverInfo {
	return resolve.NewResolverInfo(true, false)
}

func (*TestStaticResolver) Resolve(scope data.Scope, item string, field string) (interface{}, error) {
	if field == "foo" {
		return "bar", nil
	}

	if field == "bar" {
		return "for", nil
	}

	return nil, fmt.Errorf("failed to resolve variable: '%s'", field)
}

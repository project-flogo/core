package script

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/project-flogo/core/data/property"
	"os"
	"testing"
	"time"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/resolve"
	"github.com/stretchr/testify/assert"
)

var resolver = resolve.NewCompositeResolver(map[string]resolve.Resolver{"static": &TestStaticResolver{}, ".": &TestResolver{}, "env": &resolve.EnvResolver{}, "property": &property.Resolver{}})
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
	var testData interface{}
	err := json.Unmarshal([]byte(testJsonData), &testData)
	assert.Nil(t, err)

	scope := newScope(map[string]interface{}{"foo": testData, "key": 2})

	expr, err := factory.NewExpr("$.foo.store.book[0].price")
	assert.Nil(t, err)
	assert.NotNil(t, expr)

	v, err := expr.Eval(scope)
	assert.Nil(t, err)
	assert.Equal(t, 8.95, v)
}

func TestArrayIndexExpr(t *testing.T) {
	var testData interface{}
	err := json.Unmarshal([]byte(testJsonData), &testData)
	assert.Nil(t, err)

	scope := newScope(map[string]interface{}{"foo": testData, "key": 2})

	expr, err := factory.NewExpr("$.foo.store.book[$.key].price")
	assert.Nil(t, err)
	assert.NotNil(t, expr)

	v, err := expr.Eval(scope)
	assert.Nil(t, err)
	assert.Equal(t, 8.99, v)

	expr, err = factory.NewExpr("$.foo.store.book[$.key > 2 ? 2:3].price")
	assert.Nil(t, err)
	assert.NotNil(t, expr)

	v, err = expr.Eval(scope)
	assert.Nil(t, err)
	assert.Equal(t, 22.99, v)

	expr, err = factory.NewExpr(`$.foo.store.book[$.key > 2 ? 2:"aa"].price`)
	assert.Nil(t, err)
	assert.NotNil(t, expr)

	v, err = expr.Eval(scope)
	assert.EqualError(t, err, "Invalid array index: aa")
}

func TestArrayIndexExprWithFunction(t *testing.T) {
	var testData interface{}
	err := json.Unmarshal([]byte(testJsonData), &testData)
	assert.Nil(t, err)

	scope := newScope(map[string]interface{}{"foo": testData, "key": 2})

	expr, err := factory.NewExpr("$.foo.store.book[script.length(\"23\")].price")
	assert.Nil(t, err)
	assert.NotNil(t, expr)

	v, err := expr.Eval(scope)
	assert.Nil(t, err)
	assert.Equal(t, 8.99, v)

	expr, err = factory.NewExpr("$.foo.store.book[2].price")
	assert.Nil(t, err)
	assert.NotNil(t, expr)

	v, err = expr.Eval(scope)
	assert.Nil(t, err)
	assert.Equal(t, 8.99, v)

}

func TestRefWithQuotes(t *testing.T) {
	var testData interface{}
	err := json.Unmarshal([]byte(testJsonData), &testData)
	assert.Nil(t, err)

	os.Setenv("index", "2")
	scope := newScope(map[string]interface{}{"foo": testData, "key": 2})

	expr, err := factory.NewExpr(`$.foo["store"].book[script.length("123")].price`)
	assert.Nil(t, err)
	assert.NotNil(t, expr)

	v, err := expr.Eval(scope)
	assert.Nil(t, err)
	assert.Equal(t, 22.99, v)

	expr, err = factory.NewExpr("$.foo['store'].book[2].price")
	assert.Nil(t, err)
	assert.NotNil(t, expr)

	v, err = expr.Eval(scope)
	assert.Nil(t, err)
	assert.Equal(t, 8.99, v)

	expr, err = factory.NewExpr("$.foo[`store`].book[0].price")
	assert.Nil(t, err)
	assert.NotNil(t, expr)

	v, err = expr.Eval(scope)
	assert.Nil(t, err)
	assert.Equal(t, 8.95, v)

	expr, err = factory.NewExpr("$.foo[`store`].book[$env[index]].price")
	assert.Nil(t, err)
	assert.NotNil(t, expr)

	v, err = expr.Eval(scope)
	assert.Nil(t, err)
	assert.Equal(t, 8.99, v)
	defer os.Unsetenv("index")

}
func TestLitExprStaticRef(t *testing.T) {

	expr, err := factory.NewExpr(`$static.foo`)
	assert.Nil(t, err)
	assert.NotNil(t, expr)

	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)

}

func TestPropertyName(t *testing.T) {
	property.SetDefaultManager(property.NewManager(map[string]interface{}{"Marketo.Connection.client_id.secret-id": "abc"}))
	expr, err := factory.NewExpr(`$property["Marketo.Connection.client_id.secret-id"]`)
	assert.Nil(t, err)
	assert.NotNil(t, expr)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, "abc", v)
}

func TestEnvResolve(t *testing.T) {
	_ = os.Setenv("FOO", "bar")
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

	scope := newScope(map[string]interface{}{"foo": "foo", "key": 2})
	expr, err = factory.NewExpr(` true || $.NOTEXIST > 200`)
	assert.Nil(t, err)
	v, err = expr.Eval(scope)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	expr, err = factory.NewExpr(` false || $.foo == "foo"`)
	assert.Nil(t, err)
	v, err = expr.Eval(scope)
	assert.Nil(t, err)
	assert.Equal(t, true, v)

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

	scope := newScope(map[string]interface{}{"foo": "foo", "key": 2})
	expr, err = factory.NewExpr(` false && $.NOTEXIST > 200`)
	assert.Nil(t, err)
	v, err = expr.Eval(scope)
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	expr, err = factory.NewExpr(` true && $.NOTEXIST > 200`)
	assert.Nil(t, err)
	v, err = expr.Eval(scope)
	assert.NotNil(t, err)

	expr, err = factory.NewExpr(` $.key == 2 && $.foo == "foo"`)
	assert.Nil(t, err)
	v, err = expr.Eval(scope)
	assert.Nil(t, err)
	assert.Equal(t, true, v)
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
	_ = os.Setenv("name", "flogo")
	_ = os.Setenv("address", "tibco")

	testcases := make(map[string]interface{})
	testcases[`1>2?script.concat("sss","ddddd"):"fff"`] = "fff"
	testcases[`1<2?"helloworld":"fff"`] = "helloworld"
	testcases["200>100?true:false"] = true
	testcases["1 + 2 * 3 + 2 * 6"] = 19
	testcases[`script.length($.queryParams.id) == 0 ? "Query Id cannot be null" : script.length($.queryParams.id)`] = 10
	testcases[`script.length("helloworld")>11?"helloworld":"fff"`] = "fff"
	testcases["123==456"] = false
	testcases["123==123"] = true
	testcases[`script.concat("123","456")=="123456"`] = true
	testcases[`script.concat("123","456") == script.concat("12","3456")`] = true
	testcases[`("dddddd" == "dddd3dd") && ("133" == "123")`] = false
	testcases[`script.length("helloworld") == 10`] = true
	testcases[`script.length("helloworld") > 10`] = false
	testcases[`script.length("helloworld") >= 10`] = true
	testcases[`script.length("helloworld") < 10`] = false
	testcases[`script.length("helloworld") >= 10`] = true
	testcases[`(script.length("sea") == 3) == true`] = true

	testcases[`(1&&1)==(1&&1)`] = true
	testcases[`(true && true) == false`] = false
	testcases[`nil==nil`] = true

	//Nested Ternary
	testcases[`(script.length("1234") == 4 ? true : false) ? (2 >1 ? (3>2?"Yes":"nono"):"No") : "false"`] = "Yes"
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

	//testcases[`script.length("helloworld")>11?$env[name]:$env[address]`] = "tibco"
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

func TestEscapedExpr(t *testing.T) {
	expr, err := factory.NewExpr(`script.concat("\"Hello\" ", '\'FLOGO\'')`)
	assert.Nil(t, err)
	v, err := expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, `"Hello" 'FLOGO'`, v)

	expr, err = factory.NewExpr("script.concat(`Hello `, `'FLOGO'`)")
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, `Hello 'FLOGO'`, v)

	expr, err = factory.NewExpr(`script.concat("Hello", "\world")`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, `Hello\world`, v)

	expr, err = factory.NewExpr(`script.concat("Hello", "wo\'rld")`)
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, `Hellowo\'rld`, v)

	expr, err = factory.NewExpr("script.concat('Hello', ` wo\nrld`)")
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	fmt.Println("=====", v)
	assert.Equal(t, `Hello wo
rld`, v)

	//Newline
	expr, err = factory.NewExpr("script.concat(\"Hello\n\", \"FLOGO\")")
	assert.Nil(t, err)
	v, err = expr.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, `Hello
FLOGO`, v)
}

func TestBuiltInFunction(t *testing.T) {

	var testData interface{}
	err := json.Unmarshal([]byte(testJsonData), &testData)
	assert.Nil(t, err)

	scope := newScope(map[string]interface{}{"foo": testData, "key": 2})

	tests := []struct {
		Expr         string
		ExpectResult interface{}
	}{
		{
			Expr:         "isDefined($.foo['store'].book[2].price)",
			ExpectResult: true,
		},
		{
			Expr:         "isDefined($.foo.store.exit)",
			ExpectResult: false,
		},
		{
			Expr:         "getValue($.foo.store.exit, \"flogo\")",
			ExpectResult: "flogo",
		},
		{
			Expr:         "getValue($.foo['store'].book[2].price, \"flogo\")",
			ExpectResult: 8.99,
		},
	}

	for i, tt := range tests {

		expr, err := factory.NewExpr(tt.Expr)
		assert.Nil(t, err)
		v, err := expr.Eval(scope)
		assert.Nil(t, err)
		if assert.NoError(t, err, "Unexpected error in case #%d.", i) {
			assert.Equal(
				t,
				tt.ExpectResult,
				v,
				"Unexpected expression output: expected to %v.", tt.ExpectResult,
			)
		}
	}
}

func TestDateTimeComparation(t *testing.T) {

	now := time.Now()
	scope := newScope(map[string]interface{}{"date1": now, "date2": now, "date3": "2020-03-19T15:02:03Z", "date4": "2050-03-19T15:02:03Z"})

	tests := []struct {
		Expr         string
		ExpectResult interface{}
	}{
		{
			Expr:         "$.date1 == $.date2",
			ExpectResult: true,
		},
		{
			Expr:         "$.date1 > $.date3",
			ExpectResult: true,
		},
		{
			Expr:         "$.date2 > $.date3",
			ExpectResult: true,
		},
		{
			Expr:         "$.date2 < $.date4",
			ExpectResult: true,
		},
		{
			Expr:         "$.date2 <= $.date4",
			ExpectResult: true,
		},
		{
			Expr:         "$.date4 >= $.date4",
			ExpectResult: true,
		},
		{
			Expr:         "$.date2 > $.date4",
			ExpectResult: false,
		},
	}

	for i, tt := range tests {

		expr, err := factory.NewExpr(tt.Expr)
		assert.Nil(t, err)
		v, err := expr.Eval(scope)
		assert.Nil(t, err)
		if assert.NoError(t, err, "Unexpected error in case #%d.", i) {
			assert.Equal(
				t,
				tt.ExpectResult,
				v,
				"Unexpected expression output: expected to %v.", tt.ExpectResult,
			)
		}
	}
}

func TestJsonNumberWithOperator(t *testing.T) {

	var jsonData = `
	{
"data": {
	"int":100,
	"float":45.33
}
    }
`
	var data map[string]interface{}

	d := json.NewDecoder(bytes.NewReader([]byte(jsonData)))
	d.UseNumber()
	err := d.Decode(&data)
	if err != nil {
		t.Fatal(err)
		return
	}

	scope := newScope(data)

	tests := []struct {
		Expr         string
		ExpectResult interface{}
	}{
		{
			Expr:         "$.data.int == 100",
			ExpectResult: true,
		},
		{
			Expr:         "$.data.float == 45.33",
			ExpectResult: true,
		},
		{
			Expr:         "$.data.float >= 46.33",
			ExpectResult: false,
		},
		{
			Expr:         "$.data.float <= 45.33",
			ExpectResult: true,
		},
		{
			Expr:         "$.data.float > 45.33",
			ExpectResult: false,
		},
		{
			Expr:         "$.data.float < 45.33",
			ExpectResult: false,
		},
		{
			Expr:         "$.data.float - $.data.int",
			ExpectResult: -54.67,
		},
		{
			Expr:         "$.data.float * $.data.int",
			ExpectResult: float64(4533),
		}, {
			Expr:         "$.data.int / 2",
			ExpectResult: int64(50),
		},
	}

	for i, tt := range tests {
		expr, err := factory.NewExpr(tt.Expr)
		assert.Nil(t, err)
		v, err := expr.Eval(scope)
		assert.Nil(t, err)
		if assert.NoError(t, err, "Unexpected error in case #%d.", i) {
			assert.Equal(
				t,
				tt.ExpectResult,
				v,
				"Unexpected expression output: expected to %v.", tt.ExpectResult,
			)
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

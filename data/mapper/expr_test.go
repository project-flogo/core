package mapper

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression/function"
	_ "github.com/project-flogo/core/data/expression/script"
	"github.com/project-flogo/core/data/resolve"
	"github.com/stretchr/testify/assert"
)

func TestLiteralMapper(t *testing.T) {

	mappings := map[string]interface{}{"One": "1", "Two": 2, "Three": 3, "Four": "4"}

	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	scope := data.NewSimpleScope(map[string]interface{}{}, nil)
	results, err := mapper.Apply(scope)
	assert.Nil(t, err)

	assert.Equal(t, "1", results["One"])
	assert.Equal(t, 2, results["Two"])
	assert.Equal(t, 3, results["Three"])
	assert.Equal(t, "4", results["Four"])

	//todo add util to do set to handle Obj.key, Params.paramKey etc
}

func TestAssignMapper(t *testing.T) {

	mappings := map[string]interface{}{"One": "=$.SimpleI", "Two": "=$.ObjI.key", "Three": "=$.ArrayI[2]", "Four": "=$.ParamsI.paramKey", "Five": "$.SimpleI"}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	inputValues := make(map[string]interface{})

	inputValues["SimpleI"] = 1

	objVal, _ := coerce.ToObject("{\"key\":2}")
	inputValues["ObjI"] = objVal

	arrVal, _ := coerce.ToArray("[1,2,3]")
	inputValues["ArrayI"] = arrVal

	paramVal, _ := coerce.ToParams("{\"paramKey\":\"val4\"}")
	inputValues["ParamsI"] = paramVal

	scope := data.NewSimpleScope(inputValues, nil)
	results, err := mapper.Apply(scope)
	assert.Nil(t, err)

	assert.Equal(t, 1, results["One"])
	assert.Equal(t, 2.0, results["Two"])
	assert.Equal(t, 3.0, results["Three"])
	assert.Equal(t, "val4", results["Four"])
	assert.Equal(t, "$.SimpleI", results["Five"])
}

func TestExpressionMapperFunction(t *testing.T) {

	mappings := map[string]interface{}{"SimpleO": `=tstring.concat("Hello ",$.SimpleI)`}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	inputValues := make(map[string]interface{})
	inputValues["SimpleI"] = "FLOGO"

	scope := data.NewSimpleScope(inputValues, nil)
	results, err := mapper.Apply(scope)
	assert.Nil(t, err)
	assert.Equal(t, "Hello FLOGO", results["SimpleO"])
}

func TestExpressionMapperConditionExpr(t *testing.T) {

	mappings := map[string]interface{}{"SimpleO": `=$.SimpleI == "FLOGO"`}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	inputValues := make(map[string]interface{})
	inputValues["SimpleI"] = "FLOGO"

	scope := data.NewSimpleScope(inputValues, nil)
	results, err := mapper.Apply(scope)
	assert.Nil(t, err)
	assert.Equal(t, true, results["SimpleO"])
}

func TestExpressionMapperTernaryExpr(t *testing.T) {

	mappings := map[string]interface{}{"SimpleO": `=$.SimpleI == "FLOGO" ? "Welcome FLOGO" : "Bye bye !"`}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	inputValues := make(map[string]interface{})
	inputValues["SimpleI"] = "FLOGO"

	scope := data.NewSimpleScope(inputValues, nil)
	results, err := mapper.Apply(scope)
	assert.Nil(t, err)
	assert.Equal(t, "Welcome FLOGO", results["SimpleO"])

	inputValues2 := make(map[string]interface{})
	inputValues2["SimpleI"] = "FLOGO2"

	scope = data.NewSimpleScope(inputValues2, nil)
	results, err = mapper.Apply(scope)
	assert.Nil(t, err)
	assert.Equal(t, "Bye bye !", results["SimpleO"])
}

func BenchmarkAssignMapper(b *testing.B) {

	mappings := map[string]interface{}{"SimpleO": "$.SimpleI"}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, _ := factory.NewMapper(mappings)

	inputValues := make(map[string]interface{})
	inputValues["SimpleI"] = 1

	scope := data.NewSimpleScope(inputValues, nil)

	for n := 0; n < b.N; n++ {

		results, err := mapper.Apply(scope)
		if err != nil {
			b.Error(err)
			b.Fail()
		}

		val, ok := results["SimpleO"]
		if ok {
			if val != 1 {
				b.Fail()
			}
		}
	}
}

func BenchmarkLiteralMapper(b *testing.B) {

	mappings := map[string]interface{}{"SimpleO": "testing"}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, _ := factory.NewMapper(mappings)

	scope := data.NewSimpleScope(map[string]interface{}{}, nil)

	for n := 0; n < b.N; n++ {

		results, err := mapper.Apply(scope)
		if err != nil {
			b.Error(err)
			b.Fail()
		}

		val, ok := results["SimpleO"]
		if ok {
			if val != "testing" {
				panic("Mapper error")
			}
		}
	}
}

func BenchmarkExpressionMapperFunction(b *testing.B) {

	mappings := map[string]interface{}{"SimpleO": `=tstring.concat("Hello ",$.SimpleI)`}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, _ := factory.NewMapper(mappings)

	inputValues := make(map[string]interface{})
	inputValues["SimpleI"] = "FLOGO"

	scope := data.NewSimpleScope(inputValues, nil)

	for n := 0; n < b.N; n++ {

		results, err := mapper.Apply(scope)
		if err != nil {
			b.Error(err)
			b.Fail()
		}

		val, ok := results["SimpleO"]
		if ok {
			if val != "Hello FLOGO" {
				b.Fail()
			}
		}
	}
}

func BenchmarkExpressionMapperConditionExpr(b *testing.B) {

	mappings := map[string]interface{}{"SimpleO": `=$.SimpleI == "FLOGO"`}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, _ := factory.NewMapper(mappings)

	inputValues := make(map[string]interface{})
	inputValues["SimpleI"] = "FLOGO"

	scope := data.NewSimpleScope(inputValues, nil)

	for n := 0; n < b.N; n++ {

		results, err := mapper.Apply(scope)
		if err != nil {
			b.Error(err)
			b.Fail()
		}

		val, ok := results["SimpleO"]
		if ok {
			if val != true {
				b.Fail()
			}
		}
	}
}

func BenchmarkExpressionMapperTernaryExpr(b *testing.B) {

	mappings := map[string]interface{}{"SimpleO": `=$.SimpleI == "FLOGO" ? "Welcome FLOGO" : "Bye bye !"`}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, _ := factory.NewMapper(mappings)

	inputValues := make(map[string]interface{})
	inputValues["SimpleI"] = "FLOGO"

	scope := data.NewSimpleScope(inputValues, nil)

	for n := 0; n < b.N; n++ {

		results, err := mapper.Apply(scope)
		if err != nil {
			b.Error(err)
			b.Fail()
		}

		val, ok := results["SimpleO"]
		if ok {
			if val != "Welcome FLOGO" {
				b.Fail()
			}
		}
	}
}

func init() {
	_ = function.Register(&fnConcat{})
	function.SetPackageAlias(reflect.ValueOf(fnConcat{}).Type().PkgPath(), "tstring")
	function.ResolveAliases()
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

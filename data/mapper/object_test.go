package mapper

import (
	"encoding/json"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/resolve"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectMappingWithFunction(t *testing.T) {
	mappingValue := `{
        "person2" : "person",
        "addresses": {
              "tostate"   : "=tstring.concat(\"State is \", \"tx\")",
               "tostreet": "3421 st",
              "addresses2": {
               			"tofield2": "=tstring.concat(\"field is \", \"ffff\")"
              }
        }
    }`

	arrayMapping := make(map[string]interface{})
	err := json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)
	mappings := map[string]interface{}{"addresses": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)

	results, err := mapper.Apply(nil)
	assert.Nil(t, err)

	arr := results["addresses"]
	assert.Equal(t, "person", arr.(map[string]interface{})["person2"])
	assert.Equal(t, "3421 st", arr.(map[string]interface{})["addresses"].(map[string]interface{})["tostreet"])
	assert.Equal(t, "State is tx", arr.(map[string]interface{})["addresses"].(map[string]interface{})["tostate"])
}

//1. array mapping with other activity output
func TestArrayMappingWithNest(t *testing.T) {
	mappingValue := `{
        "person2" : "person",
        "addresses": {
            "@foreach($.field.addresses, index)":
            {
              "tostate"   : "=$loop[index].state",
               "tostreet": "=$.street",
               "tozipcode":"=$.zipcode",
              "addresses2": {
                  "@foreach($.array)":{
                        "tofield1"  : "=$loop[index].street",
               			"tofield2": "=$.field2",
               			"tofield3":"=$.field3"
                  }
              }
            }
        }
    }`

	arrayData := `{
   "person": "name",
   "addresses": [
       {
           "street": "street",
           "zipcode": 77479,
           "state": "tx",
			"array":[
				{
					"field1":"field1value",
					"field2":"field2value",
					"field3":"field3value"
				},
				{
					"field1":"field1value2",
					"field2":"field2value2",
					"field3":"field3value2"
				}
			]
       },
 {
           "street": "street2",
           "zipcode": 774792,
           "state": "tx2",
			"array":[
				{
					"field1":"field1value2",
					"field2":"field2value2",
					"field3":"field3value2"
				},
				{
					"field1":"field1value22",
					"field2":"field2value22",
					"field3":"field3value22"
				}
			]
       }
   ]
}`

	arrayMapping := make(map[string]interface{})
	err := json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)
	mappings := map[string]interface{}{"addresses": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)

	attrs := map[string]interface{}{"field": arrayData}
	scope := data.NewSimpleScope(attrs, nil)
	results, err := mapper.Apply(scope)
	assert.Nil(t, err)

	arr := results["addresses"]
	assert.Equal(t, "person", arr.(map[string]interface{})["person2"])
	assert.Equal(t, float64(77479), arr.(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["tozipcode"])
	assert.Equal(t, "tx", arr.(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["tostate"])
}

//2. array mapping with function
func TestArrayMappingWithFunction(t *testing.T) {
	mappingValue := `{
        "person2" : "person",
        "addresses": {
            "@foreach($.field.addresses, index)":
            {
              "tostate"   : "=tstring.concat(\"State is \", $loop[index].state)",
               "tostreet": "=$.street",
               "tozipcode":"=$.zipcode",
              "addresses2": {
                  "@foreach($.array)":{
                        "tofield1"  : "=$loop[index].street",
               			"tofield2": "=tstring.concat(\"field is \", $.field2)",
               			"tofield3":"=$.field3"
                  }
              }
            }
        }
    }`

	arrayData := `{
   "person": "name",
   "addresses": [
       {
           "street": "street",
           "zipcode": 77479,
           "state": "tx",
			"array":[
				{
					"field1":"field1value",
					"field2":"field2value",
					"field3":"field3value"
				},
				{
					"field1":"field1value2",
					"field2":"field2value2",
					"field3":"field3value2"
				}
			]
       }
   ]
}`

	arrayMapping := make(map[string]interface{})
	err := json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)
	mappings := map[string]interface{}{"addresses": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)

	attrs := map[string]interface{}{"field": arrayData}
	scope := data.NewSimpleScope(attrs, nil)
	results, err := mapper.Apply(scope)
	assert.Nil(t, err)

	arr := results["addresses"]
	assert.Equal(t, "person", arr.(map[string]interface{})["person2"])
	assert.Equal(t, float64(77479), arr.(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["tozipcode"])
	assert.Equal(t, "State is tx", arr.(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["tostate"])
}

func TestGetSource(t *testing.T) {
	var s = "@foreach($activity[blah].out2)"
	foreach := newForeach(s, nil)
	assert.Equal(t, "$activity[blah].out2", foreach.sourceFrom)
	assert.Equal(t, "", foreach.index)

	s = "@foreach($activity[blah].out2, index)"
	foreach = newForeach(s, nil)
	assert.Equal(t, "$activity[blah].out2", foreach.sourceFrom)
	assert.Equal(t, "index", foreach.index)

}

package mapper

import (
	"encoding/json"
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/resolve"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectMappingWithFunction(t *testing.T) {
	mappingValue := `{"mapping": {
        "person2" : "person",
        "addresses": {
              "tostate"   : "=tstring.concat(\"State is \", \"tx\")",
               "tostreet": "3421 st",
              "addresses2": {
					"tofield2": "=tstring.concat(\"field is \", \"ffff\")"
              }
        }
    }
}`

	arrayMapping := make(map[string]interface{})
	err := json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)
	assert.False(t, IsLiteral(arrayMapping))
	mappings := map[string]interface{}{"addresses": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	results, err := mapper.Apply(nil)
	assert.Nil(t, err)

	arr := results["addresses"]
	assert.Equal(t, "person", arr.(map[string]interface{})["person2"])
	assert.Equal(t, "3421 st", arr.(map[string]interface{})["addresses"].(map[string]interface{})["tostreet"])
	assert.Equal(t, "State is tx", arr.(map[string]interface{})["addresses"].(map[string]interface{})["tostate"])
}

func TestObjectMappingWithArray(t *testing.T) {
	mappingValue := `{"mapping": {
  "person2": "person",
  "addresses": {
    "array": [
      {
        "ddd": "=tstring.concat(\"ddd is \", \"tx\")",
        "ccc": "=tstring.concat(\"ccc is \", \"tx\")"
      }
    ]
  }
}
}`

	arrayMapping := make(map[string]interface{})
	err := json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)
	assert.False(t, IsLiteral(arrayMapping))
	mappings := map[string]interface{}{"addresses": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	results, err := mapper.Apply(nil)
	assert.Nil(t, err)
	arr := results["addresses"]

	assert.Equal(t, "person", arr.(map[string]interface{})["person2"])
	assert.Equal(t, "ddd is tx", arr.(map[string]interface{})["addresses"].(map[string]interface{})["array"].([]interface{})[0].(map[string]interface{})["ddd"])
	assert.Equal(t, "ccc is tx", arr.(map[string]interface{})["addresses"].(map[string]interface{})["array"].([]interface{})[0].(map[string]interface{})["ccc"])

}

func TestRootObjectArray(t *testing.T) {
	mappingValue := `{"mapping": [
   {
      "id":"11111",
      "name":"nnnnn",
      "addresses": {
			"@foreach($.field.addresses, index)":{
				"id":"dddddd",
				"name":"=$.state"
			}
      }
   }
]
}`

	arrayData := `{
   "person": "name",
   "addresses": [
       {
           "street": "street",
           "zipcode": 77479,
           "state": "tx"
		}
   ]
}`
	var arrayValue interface{}
	err := json.Unmarshal([]byte(arrayData), &arrayValue)
	assert.Nil(t, err)
	attrs := map[string]interface{}{"field": arrayValue}
	scope := data.NewSimpleScope(attrs, nil)

	var arrayMapping interface{}
	err = json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)
	assert.False(t, IsLiteral(arrayMapping))
	mappings := map[string]interface{}{"target": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	results, err := mapper.Apply(scope)
	assert.Nil(t, err)
	arr := results["target"]

	assert.Equal(t, "11111", arr.([]interface{})[0].(map[string]interface{})["id"])
	assert.Equal(t, "nnnnn", arr.([]interface{})[0].(map[string]interface{})["name"])
	assert.Equal(t, "dddddd", arr.([]interface{})[0].(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["id"])
	assert.Equal(t, "tx", arr.([]interface{})[0].(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["name"])

}
func TestRootLiteralArray(t *testing.T) {
	mappingValue := `["1243","456"]`

	var arrayMapping interface{}
	err := json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)
	assert.True(t, IsLiteral(arrayMapping))
	mappings := map[string]interface{}{"target": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	results, err := mapper.Apply(nil)
	assert.Nil(t, err)
	arr := results["target"]

	assert.Equal(t, "1243", arr.([]interface{})[0])
	assert.Equal(t, "456", arr.([]interface{})[1])

}

func TestPrimitiveArray(t *testing.T) {
	mappingValue := `{"mapping": 
{
  "features": [
    {
      "name": "inputs",
      "data": [
        "=$.result.V1",
        "=$.result.V2",
        "=$.result.V3",
        "=$.result.V4"
      ]
    }
  ]
}
}`

	var arrayMapping interface{}
	err := json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)
	assert.False(t, IsLiteral(arrayMapping))
	mappings := map[string]interface{}{"target": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	attr := make(map[string]interface{})

	val := map[string]interface{}{"V1": "1111"}
	val["V2"] = "2222"
	val["V3"] = "3333"
	val["V4"] = "4444"
	attr["result"] = val

	scope := data.NewSimpleScope(attr, nil)

	results, err := mapper.Apply(scope)
	assert.Nil(t, err)
	arr := results["target"]

	assert.Equal(t, "inputs", arr.(map[string]interface{})["features"].([]interface{})[0].(map[string]interface{})["name"])
	assert.Equal(t, "2222", arr.(map[string]interface{})["features"].([]interface{})[0].(map[string]interface{})["data"].([]interface{})[1])

}

func TestRootLiteralArrayMapping(t *testing.T) {
	mappingValue := `{"mapping": ["=$.field.name", "=$.field.id"]}`
	arrayData := `{
   "name": "name",
	"id":"1001",
   "addresses": [
       {
           "street": "street",
           "zipcode": 77479,
           "state": "tx"
		}
   ]
}`
	var arrayValue interface{}
	err := json.Unmarshal([]byte(arrayData), &arrayValue)
	assert.Nil(t, err)

	attrs := map[string]interface{}{"field": arrayValue}
	scope := data.NewSimpleScope(attrs, nil)

	var arrayMapping interface{}
	err = json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)
	assert.False(t, IsLiteral(arrayMapping))
	mappings := map[string]interface{}{"target": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	results, err := mapper.Apply(scope)
	assert.Nil(t, err)
	arr := results["target"]

	assert.Equal(t, "name", arr.([]interface{})[0])
	assert.Equal(t, "1001", arr.([]interface{})[1])
}

func TestRootLiteralNumberArrayMapping(t *testing.T) {
	mappingValue := `{"mapping": ["=$.field.id2", "=$.field.id"]}`
	arrayData := `{
	"id2": 1002,
	"id":1001,
   "addresses": [
       {
           "street": "street",
           "zipcode": 77479,
           "state": "tx"
		}
   ]
}`
	var arrayValue interface{}
	err := json.Unmarshal([]byte(arrayData), &arrayValue)
	assert.Nil(t, err)

	attrs := map[string]interface{}{"field": arrayValue}
	scope := data.NewSimpleScope(attrs, nil)

	var arrayMapping interface{}
	err = json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)
	assert.False(t, IsLiteral(arrayMapping))

	mappings := map[string]interface{}{"target": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	results, err := mapper.Apply(scope)
	assert.Nil(t, err)
	arr := results["target"]

	assert.Equal(t, float64(1002), arr.([]interface{})[0])
	assert.Equal(t, float64(1001), arr.([]interface{})[1])
}

func TestRootArrayMapping(t *testing.T) {
	mappingValue := `{"mapping": {
			"@foreach($.field.addresses, index)":{
				"id":"dddddd",
				"name":"=$.state"
			}
   }}`

	arrayData := `{
   "person": "name",
   "addresses": [
       {
           "street": "street",
           "zipcode": 77479,
           "state": "tx"
		}
   ]
}
`
	var arrayValue interface{}
	err := json.Unmarshal([]byte(arrayData), &arrayValue)
	assert.Nil(t, err)

	attrs := map[string]interface{}{"field": arrayValue}
	scope := data.NewSimpleScope(attrs, nil)

	var arrayMapping interface{}
	err = json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)
	assert.False(t, IsLiteral(arrayMapping))

	mappings := map[string]interface{}{"target": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	results, err := mapper.Apply(scope)
	assert.Nil(t, err)
	arr := results["target"]

	assert.Equal(t, "dddddd", arr.([]interface{})[0].(map[string]interface{})["id"])
	assert.Equal(t, "tx", arr.([]interface{})[0].(map[string]interface{})["name"])

}

func TestStringStringMap(t *testing.T) {
	mappingValue := `
	{"id":"11111",
	"name":"nnnnn"}`

	arrayMapping := make(map[string]interface{})
	err := json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)
	assert.True(t, IsLiteral(arrayMapping))

	mappings := map[string]interface{}{"target": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	results, err := mapper.Apply(nil)
	assert.Nil(t, err)
	arr := results["target"]

	assert.Equal(t, "11111", arr.(map[string]interface{})["id"])
	assert.Equal(t, "nnnnn", arr.(map[string]interface{})["name"])

}

func TestArrayMappingWithNest(t *testing.T) {
	mappingValue := `{"mapping": {
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
    }}`

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
	assert.False(t, IsLiteral(arrayMapping))
	mappings := map[string]interface{}{"addresses": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	attrs := map[string]interface{}{"field": arrayData}
	scope := data.NewSimpleScope(attrs, nil)
	results, err := mapper.Apply(scope)
	assert.Nil(t, err)

	arr := results["addresses"]
	assert.Equal(t, "person", arr.(map[string]interface{})["person2"])
	assert.Equal(t, float64(77479), arr.(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["tozipcode"])
	assert.Equal(t, "tx", arr.(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["tostate"])
}

func TestArrayMappingWithFunction(t *testing.T) {
	mappingValue := `{"mapping": {
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
	assert.False(t, IsLiteral(arrayMapping))
	mappings := map[string]interface{}{"addresses": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	attrs := map[string]interface{}{"field": arrayData}
	scope := data.NewSimpleScope(attrs, nil)
	results, err := mapper.Apply(scope)
	assert.Nil(t, err)

	arr := results["addresses"]
	assert.Equal(t, "person", arr.(map[string]interface{})["person2"])
	assert.Equal(t, float64(77479), arr.(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["tozipcode"])
	assert.Equal(t, "State is tx", arr.(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["tostate"])
}

func TestArrayMappingWithFunction3Level(t *testing.T) {
	mappingValue := `{"mapping": {
   "person2":"person",
   "addresses":{
      "@foreach($.field.addresses, index)":{
         "tostate":"=tstring.concat(\"State is \", $loop[index].state)",
         "tostreet":"=$.street",
         "tozipcode":"=$.zipcode",
         "addresses2":{
            "@foreach($.array, index2)":{
               "tofield1":"=$loop[index].street",
               "tofield2":"=tstring.concat(\"field is \", $.field2)",
               "tofield3":"=$.field3",
               "addresses4":{
                  "@foreach($.level3)":{
                     "level3":"=$loop[index2].field1",
                     "level3-1":"=tstring.concat(\"field is \", $.field3)"
                  }
               }
            }
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
					"field3":"field3value",
					"level3": [
						{"field3":"ddddd"}
					]
				},
				{
					"field1":"field1value2",
					"field2":"field2value2",
					"field3":"field3value2",
					"level3": [
						{"field3":"ddddd2"}
					]
				}
			]
       }
   ]
}`
	var arrayValue interface{}
	err := json.Unmarshal([]byte(arrayData), &arrayValue)
	assert.Nil(t, err)

	arrayMapping := make(map[string]interface{})
	err = json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)
	assert.False(t, IsLiteral(arrayMapping))

	mappings := map[string]interface{}{"addresses": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	attrs := map[string]interface{}{"field": arrayValue}
	scope := data.NewSimpleScope(attrs, nil)
	results, err := mapper.Apply(scope)
	assert.Nil(t, err)

	arr := results["addresses"]

	assert.Equal(t, "person", arr.(map[string]interface{})["person2"])
	assert.Equal(t, float64(77479), arr.(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["tozipcode"])
	assert.Equal(t, "State is tx", arr.(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["tostate"])
}

func TestLiteral(t *testing.T) {
	mappingValue := `{"mapping": {
                      "Character": {
                        "appearsIn": [
                          "EMPIRE",
                          "JEDI"
                        ],
                        "friends": [
                          {
                            "appearsIn": [
                              "JEDI"
                            ],
                            "friends": [],
                            "id": "d123",
                            "name": "r2-d2",
                            "primaryFunction": "robot"
                          },
                          {
                            "appearsIn": [
                              "JEDI",
                              "NEWHOPE"
                            ],
                            "friends": [],
                            "homePlanet": "Mars",
                            "id": "h234",
                            "name": "Robert"
                          }
                        ],
                        "homePlanet": "Earth",
                        "id": "h123",
                        "name": "Luke"
                      }
                    }}`

	arrayMapping := make(map[string]interface{})
	err := json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)

	mappings := map[string]interface{}{"addresses": arrayMapping}
	factory := NewFactory(resolve.GetBasicResolver())
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	results, err := mapper.Apply(nil)
	assert.Nil(t, err)

	arr := results["addresses"]
	v, _ := json.Marshal(arr)
	fmt.Println(string(v))
	//assert.Equal(t, "person", arr.(map[string]interface{})["person2"])
	//assert.Equal(t, float64(77479), arr.(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["tozipcode"])
	//assert.Equal(t, "State is tx", arr.(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["tostate"])

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

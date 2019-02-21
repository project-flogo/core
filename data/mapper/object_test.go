package mapper

import (
	"encoding/json"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/resolve"
	"testing"

	"github.com/stretchr/testify/assert"
)

//
//func TestArrayMapping(t *testing.T) {
//	mappingValue := `
//{
//		"persion2":"$"
//       "@foreach[=$.field.addresses]":{
//              "street"  : "=$.street",
//              "zipcode" : "=$.zipcode",
//              "state"   : "=$.state"
//		}
//}
//
//`
//
//	arrayData := `{
//    "person": "name",
//    "addresses": [
//        {
//            "street": "street",
//            "zipcode": 77479,
//            "state": "tx"
//        }
//    ]
//}`
//
//	var val interface{}
//	err := json.Unmarshal([]byte(mappingValue), &val)
//	assert.Nil(t, err)
//	mappings := map[string]interface{}{"addresses": val}
//	factory := NewFactory(resolve.GetBasicResolver())
//	mapper, err := factory.NewMapper(mappings)
//
//	attrs := map[string]interface{}{"field": arrayData}
//	scope := data.NewSimpleScope(attrs, nil)
//	results, err := mapper.Apply(scope)
//	assert.Nil(t, err)
//	arr := results["addresses"]
//
//	assert.Equal(t, "street", arr.([]interface{})[0].(map[string]interface{})["street"])
//	assert.Equal(t, float64(77479), arr.([]interface{})[0].(map[string]interface{})["zipcode"])
//	assert.Equal(t, "tx", arr.([]interface{})[0].(map[string]interface{})["state"])
//
//}
//
////8. array mappping with static function and leaf field
//func TestArrayMappingWithFunction(t *testing.T) {
//	mappingValue := `{
//   "fields": [
//       {
//           "from": "=tstring.concat(\"this stree name: \", $.street)",
//           "to": "street",
//           "type": "primitive"
//       },
//       {
//           "from": "=tstring.concat(\"The zipcode is: \",$.zipcode)",
//           "to": "zipcode",
//           "type": "primitive"
//       },
//       {
//           "from": "=$.state",
//           "to": "state",
//           "type": "primitive"
//       }
//   ],
//   "from": "=$.field.addresses",
//   "to": "addresses",
//   "type": "foreach"
//}`
//
//	arrayData := `{
//   "person": "name",
//   "addresses": [
//       {
//           "street": "street",
//           "zipcode": 77479,
//           "state": "tx"
//       }
//   ]
//}`
//
//	mappings := map[string]interface{}{"addresses": mappingValue}
//	factory := NewFactory(resolve.GetBasicResolver())
//	mapper, err := factory.NewMapper(mappings)
//
//	attrs := map[string]interface{}{"field": arrayData}
//	scope := data.NewSimpleScope(attrs, nil)
//	results, err := mapper.Apply(scope)
//	assert.Nil(t, err)
//	arr := results["addresses"]
//	assert.Equal(t, "this stree name: street", arr.([]interface{})[0].(map[string]interface{})["street"])
//	assert.Equal(t, "The zipcode is: 77479", arr.([]interface{})[0].(map[string]interface{})["zipcode"])
//	assert.Equal(t, "tx", arr.([]interface{})[0].(map[string]interface{})["state"])
//
//}
//
////9. array mapping with other activity output
//func TestArrayMappingWithUpstreamingOutput(t *testing.T) {
//	mappingValue := `{
//   "fields": [
//       {
//           "from": "=tstring.concat(\"this stree name: \", $.field.person)",
//           "to": "street",
//           "type": "primitive"
//       },
//       {
//           "from": "=tstring.concat(\"The zipcode is: \",$.zipcode)",
//           "to": "zipcode",
//           "type": "primitive"
//       },
//       {
//           "from": "=$.state",
//           "to": "state",
//           "type": "primitive"
//       }
//   ],
//   "from": "=$.field.addresses",
//   "to": "addresses",
//   "type": "foreach"
//}`
//
//	arrayData := `{
//   "person": "name",
//   "addresses": [
//       {
//           "street": "street",
//           "zipcode": 77479,
//           "state": "tx"
//       }
//   ]
//}`
//	mappings := map[string]interface{}{"addresses": mappingValue}
//	factory := NewFactory(resolve.GetBasicResolver())
//	mapper, err := factory.NewMapper(mappings)
//
//	attrs := map[string]interface{}{"field": arrayData}
//	scope := data.NewSimpleScope(attrs, nil)
//	results, err := mapper.Apply(scope)
//	assert.Nil(t, err)
//	arr := results["addresses"]
//	assert.Equal(t, "this stree name: name", arr.([]interface{})[0].(map[string]interface{})["street"])
//	assert.Equal(t, "The zipcode is: 77479", arr.([]interface{})[0].(map[string]interface{})["zipcode"])
//	assert.Equal(t, "tx", arr.([]interface{})[0].(map[string]interface{})["state"])
//
//}

//9. array mapping with other activity output
func TestArrayMappingWithNest(t *testing.T) {
	mappingValue := `{
        "person2" : "person 2222",
        "addresses": {
            "@foreach($.field.addresses, index)":
            {
              "street"  : "=$.street",
              "zipcode" : "9999",
              "state"   : "$loop[index].state",
              "addresses2": {
                  "@foreach($.array)":{
                        "street2"  : "=$loop[index].street",
                        "zipcode2" : "=$.field1",
                        "state2"   : "=$.field3"
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
	assert.Equal(t, "this stree name: name", arr.([]interface{})[0].(map[string]interface{})["street"])
	assert.Equal(t, "The zipcode is: 77479", arr.([]interface{})[0].(map[string]interface{})["zipcode"])
	assert.Equal(t, "tx", arr.([]interface{})[0].(map[string]interface{})["state"])

	assert.Equal(t, "field1value", arr.([]interface{})[0].(map[string]interface{})["array"].([]interface{})[0].(map[string]interface{})["tofield1"])
	assert.Equal(t, "field2value", arr.([]interface{})[0].(map[string]interface{})["array"].([]interface{})[0].(map[string]interface{})["tofield2"])
	assert.Equal(t, "wangzai", arr.([]interface{})[0].(map[string]interface{})["array"].([]interface{})[0].(map[string]interface{})["tofield3"])

}

func TestGetSource(t *testing.T) {
	var s = "@foreach($activity[blah].out2)"
	foreach, err := getForeach(s)
	assert.Nil(t, err)
	assert.Equal(t, "$activity[blah].out2", foreach.sourceFrom)
	assert.Equal(t, "", foreach.index)

	s = "@foreach($activity[blah].out2, index)"
	foreach, err = getForeach(s)
	assert.Nil(t, err)
	assert.Equal(t, "$activity[blah].out2", foreach.sourceFrom)
	assert.Equal(t, "index", foreach.index)

}

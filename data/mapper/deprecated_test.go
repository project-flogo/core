package mapper

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayMapper(t *testing.T) {
	oldArray := `{
    "fields": [
        {
            "from": "tstring.concat(\"this street name: \", \"ddd\")",
            "to": "$.street",
            "type": "primitive"
        },
        {
            "from": "tstring.concat(\"The zipcode is: \",$.zipcode)",
            "to": "$.zipcode",
            "type": "primitive"
        },
        {
            "from": "$.state",
            "to": "$.state",
            "type": "primitive"
        },
		{
    		"from": "$.array",
    		"to": "$.array",
            "type": "foreach",
			"fields":[
				{
           			 "from": "$.field1",
           			 "to": "$.tofield1",
           			 "type": "assign"
        		},
				{
            		"from": "$.field2",
					"to": "$.tofield2",
            		"type": "assign"
        		},
				{
            		"from": "wangzai",
					"to": "$.tofield3",
            		"type": "assign"
        		}
			]

		}
    ],
    "from": "$activity[a1].field.addresses",
    "to": ".field.addresses",
    "type": "foreach"
}
`

	array, err := ParseArrayMapping(oldArray)
	assert.Nil(t, err)

	v, err := ToNewArray(array)
	assert.Nil(t, err)

	vv, _ := json.Marshal(v)
	fmt.Println(string(vv))
}

func TestNewArrayMapper(t *testing.T) {
	oldArray := `{
    "fields": [
        {
            "from": "tstring.concat(\"this street name: \", \"ddd\")",
            "to": "$.street",
            "type": "primitive"
        },
        {
            "from": "tstring.concat(\"The zipcode is: \",$.zipcode)",
            "to": "$.zipcode",
            "type": "primitive"
        },
        {
            "from": "$.state",
            "to": "$.state",
            "type": "primitive"
        },
		{
    		"from": "NEWARRAY",
    		"to": "$.array",
            "type": "foreach",
			"fields":[
				{
           			 "from": "$.field1",
           			 "to": "$.tofield1",
           			 "type": "assign"
        		},
				{
            		"from": "$.field2",
					"to": "$.tofield2",
            		"type": "assign"
        		},
				{
            		"from": "wangzai",
					"to": "$.tofield3",
            		"type": "assign"
        		}
			]

		}
    ],
    "from": "NEWARRAY",
    "to": ".field.addresses",
    "type": "foreach"
}
`

	array, err := ParseArrayMapping(oldArray)
	assert.Nil(t, err)

	v, err := ToNewArray(array)
	assert.Nil(t, err)

	vv, _ := json.Marshal(v)
	fmt.Println(string(vv))
}

func TestPathToObject(t *testing.T) {
	path := []string{"data", "field", "value"}
	obj := make(map[string]interface{})

	toObjectFromPath(path, "1234", obj)
	v, _ := json.Marshal(obj)
	fmt.Println(string(v))
	assert.Equal(t, "1234", obj["data"].(map[string]interface{})["field"].(map[string]interface{})["value"])
}

func TestPathToObjectArray(t *testing.T) {
	path := []string{"data[2]", "field[0]", "value"}
	obj := make(map[string]interface{})

	toObjectFromPath(path, "1234", obj)
	v, _ := json.Marshal(obj)
	fmt.Println(string(v))
	assert.Equal(t, "1234", obj["data"].([]interface{})[2].(map[string]interface{})["field"].([]interface{})[0].(map[string]interface{})["value"])
}

func TestMultiplePathToObject(t *testing.T) {
	path := []string{"data", "field", "value"}

	path2 := []string{"data", "field2", "value"}

	obj := make(map[string]interface{})

	toObjectFromPath(path, "1234", obj)

	toObjectFromPath(path2, "1234", obj)

	v, _ := json.Marshal(obj)
	fmt.Println(string(v))
	assert.Equal(t, "1234", obj["data"].(map[string]interface{})["field"].(map[string]interface{})["value"])
}

func TestMultiplePathToObjectArray(t *testing.T) {
	path := []string{"data[2]", "field[0]", "value"}
	obj := make(map[string]interface{})
	path2 := []string{"data[4]", "field[0]", "value"}

	toObjectFromPath(path, "1234", obj)

	toObjectFromPath(path2, "1234", obj)

	v, _ := json.Marshal(obj)
	fmt.Println(string(v))
	assert.Equal(t, "1234", obj["data"].([]interface{})[2].(map[string]interface{})["field"].([]interface{})[0].(map[string]interface{})["value"])
}

func TestMultiplePathToObjectArray2(t *testing.T) {
	path := []string{"data", "field", "value[0]"}
	obj := make(map[string]interface{})
	path2 := []string{"data", "field", "value[4]"}

	toObjectFromPath(path, 22, obj)

	toObjectFromPath(path2, 33, obj)

	v, _ := json.Marshal(obj)
	fmt.Println(string(v))
	assert.Equal(t, "1234", obj["data"].([]interface{})[2].(map[string]interface{})["field"].([]interface{})[0].(map[string]interface{})["value"])
}

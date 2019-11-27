package path

import (
	"testing"

	"github.com/project-flogo/core/data/coerce"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	Data map[string]interface{}
}

func TestGetValue(t *testing.T) {
	// Resolution of Old Trigger expression

	mapVal, _ := coerce.ToObject("{\"myParam\":5}")
	path := ".myParam"
	newVal, err := GetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 5.0, newVal)

	// Resolution of Old Trigger expression
	arrVal, _ := coerce.ToArray("[1,6,3]")
	path = "[1]"
	newVal, err = GetValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 6.0, newVal)

	mapVal, _ = coerce.ToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = ".myParam.nestedMap"
	newVal, err = GetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 1.0, newVal)

	mapVal, _ = coerce.ToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = `["myParam"].nestedMap`
	newVal, err = GetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 1.0, newVal)

	mapVal, _ = coerce.ToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = `['myParam'].nestedMap`
	newVal, err = GetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 1.0, newVal)

	mapVal, _ = coerce.ToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = `.myParam["nestedMap"]`
	newVal, err = GetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 1.0, newVal)

	arrVal, _ = coerce.ToArray("[{\"nestedMap1\":1},{\"nestedMap2\":2}]")
	path = "[1].nestedMap2"
	newVal, err = GetValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 2.0, newVal)

	mapVal, _ = coerce.ToObject("{\"myParam\":{\"nestedArray\":[7,8,9]}}")
	path = ".myParam.nestedArray[1]"
	newVal, err = GetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 8.0, newVal)

	arrVal, _ = coerce.ToArray("[{\"nestedMap1\":1},{\"nestedMap2\":{\"nestedArray\":[7,8,9]}}]")
	path = "[1].nestedMap2.nestedArray[2]"
	newVal, err = GetValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 9.0, newVal)

	mapVal, _ = coerce.ToObject("{\"myParam\":{\"nestedArray\":[7,8,9]}}")
	path = ".myParam.nestedArray"
	newVal, err = GetValue(mapVal, path)
	assert.Nil(t, err)
	//todo check if array

	arrVal, _ = coerce.ToArray("[{\"nestedMap1\":1},{\"nestedMap2\":{\"nestedArray\":[7,8,9]}}]")
	path = "[1].nestedMap2"
	newVal, err = GetValue(arrVal, path)
	assert.Nil(t, err)
	//todo check if map

	multiLevel := map[string]interface{}{
		"test": &Test{
			Data: map[string]interface{}{
				"foo": "bar",
			},
		},
	}
	path = ".test.Data.foo"
	newVal, err = GetValue(multiLevel, path)
	assert.Nil(t, err)
	assert.Equal(t, "bar", newVal)

	path = ".test.data.foo"
	newVal, err = GetValue(multiLevel, path)
	assert.Nil(t, err)
	assert.Equal(t, "bar", newVal)

	path = ".test.gah.foo"
	newVal, err = GetValue(multiLevel, path)
	assert.NotNil(t, err)
	assert.Nil(t, newVal)

	path = ".test.gah"
	newVal, err = GetValue(multiLevel, path)
	assert.NotNil(t, err)
	assert.Nil(t, newVal)
}

func TestSkipMissing(t *testing.T) {

	source := `
{
  "object": {
    "myParam": {
      "nestedMap": 1
    }
  },
  "array": [
    {
      "nestedMap1": 1
    },
    {
      "nestedMap2": {
        "nestedArray": [
          7,
          8,
          9
        ]
      }
    }
  ]
}
`

	skipMissing = true
	tests := []struct {
		Path  string
		Value interface{}
	}{
		{
			Path:  ".object.myParam.nestedMap",
			Value: float64(1),
		},
		{
			Path:  ".object.myParam.nestedMap2",
			Value: nil,
		},
		{
			Path:  ".object.myparam2",
			Value: nil,
		},
		{
			Path:  ".object.myparam2.nestedMap2",
			Value: nil,
		},
		{
			Path:  ".object.myparam[0].nestedMap2",
			Value: nil,
		},
		{
			Path:  ".object.myparam2[0]",
			Value: nil,
		},
		{
			Path:  ".object.myparam2[0].nestedMap2[0]",
			Value: nil,
		},
		//Arraty
		{
			Path:  ".array[0].nestedMap1",
			Value: float64(1),
		},
		{
			Path:  ".array[1].nestedMap2.nestedArray[0]",
			Value: float64(7),
		},
		{
			Path:  ".array[0].nestedMap2.nestedArray2[1]",
			Value: nil,
		},
		{
			Path:  ".array[0].nestedMap3[0]",
			Value: nil,
		},
		{
			Path:  ".array[0].nestedMap3",
			Value: nil,
		},
		{
			Path:  ".array[0].abc",
			Value: nil,
		},
		{
			Path:  ".array2[0].nestedMap2[0].nestedArray2[1]",
			Value: nil,
		},
		{
			Path:  ".array2",
			Value: nil,
		},
	}

	for _, test := range tests {
		newVal, err := GetValue(source, test.Path)
		assert.Nil(t, err)
		assert.Equal(t, test.Value, newVal)
	}

}

func TestSetValue(t *testing.T) {
	// Resolution of Old Trigger expression

	mapVal, _ := coerce.ToObject("{\"myParam\":5}")
	path := ".myParam"
	err := SetValue(mapVal, path, 6)
	assert.Nil(t, err)
	newVal, err := GetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 6, newVal)

	// Resolution of Old Trigger expression
	arrVal, _ := coerce.ToArray("[1,6,3]")
	path = "[1]"
	err = SetValue(arrVal, path, 4)
	assert.Nil(t, err)
	newVal, err = GetValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 4, newVal)

	mapVal, _ = coerce.ToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = ".myParam.nestedMap"
	assert.Nil(t, err)
	err = SetValue(mapVal, path, 7)
	newVal, err = GetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 7, newVal)

	mapVal, _ = coerce.ToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = `["myParam"].nestedMap`
	assert.Nil(t, err)
	err = SetValue(mapVal, path, 7)
	newVal, err = GetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 7, newVal)

	mapVal, _ = coerce.ToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = `['myParam'].nestedMap`
	assert.Nil(t, err)
	err = SetValue(mapVal, path, 7)
	newVal, err = GetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 7, newVal)

	mapVal, _ = coerce.ToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = `.myParam["nestedMap"]`
	assert.Nil(t, err)
	err = SetValue(mapVal, path, 7)
	newVal, err = GetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 7, newVal)

	arrVal, _ = coerce.ToArray("[{\"nestedMap1\":1},{\"nestedMap2\":2}]")
	path = "[1].nestedMap2"
	err = SetValue(arrVal, path, 3)
	assert.Nil(t, err)
	newVal, err = GetValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 3, newVal)

	mapVal, _ = coerce.ToObject("{\"myParam\":{\"nestedArray\":[7,8,9]}}")
	path = ".myParam.nestedArray[1]"
	err = SetValue(mapVal, path, 1)
	assert.Nil(t, err)
	newVal, err = GetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 1, newVal)

	arrVal, _ = coerce.ToArray("[{\"nestedMap1\":1},{\"nestedMap2\":{\"nestedArray\":[7,8,9]}}]")
	path = "[1].nestedMap2.nestedArray[2]"
	err = SetValue(arrVal, path, 5)
	assert.Nil(t, err)
	newVal, err = GetValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 5, newVal)

	//mapVal,_ = coerce.ToObject("{\"myParam\":{\"nestedArray\":[7,8,9]}}")
	//path = ".myParam.nestedArray"
	//err = SetValue(arrVal, path, 3)
	//assert.Nil(t, err)
	//newVal,err = GetValue(mapVal, path)
	//assert.Nil(t, err)
	////todo check if array
	//
	//arrVal,_ = coerce.ToArray("[{\"nestedMap1\":1},{\"nestedMap2\":{\"nestedArray\":[7,8,9]}}]")
	//path = "[1].nestedMap2"
	//assert.Nil(t, err)
	//err = SetValue(arrVal, path, 3)
	//newVal,err = GetValue(arrVal, path)
	//assert.Nil(t, err)
	//////todo check if map
}

func TestGetValueFromStruct(t *testing.T) {

	type object struct {
		ID     string  `json:"id"`
		Name   string  `json:"name"`
		TT     string  `json:"T!T"`
		Nested *object `json:"Nest.Object"`
	}

	source := struct {
		Object *object `json:"object"`
	}{
		Object: &object{
			ID:   "1001",
			Name: "flogo",
			TT:   "t!tTesing",
			Nested: &object{
				ID:     "1001-2",
				Name:   "flogo-2",
				TT:     "ddddd",
				Nested: nil,
			},
		},
	}

	skipMissing = true
	tests := []struct {
		Path  string
		Value interface{}
	}{
		{
			Path:  ".object.id",
			Value: "1001",
		},
		{
			Path:  ".object.name",
			Value: "flogo",
		},
		{
			Path:  `.object["T!T"]`,
			Value: "t!tTesing",
		},
		{
			Path:  `.object["Nest.Object"].id`,
			Value: "1001-2",
		},
		{
			Path:  `.object["Nest.Object"]["T!T"]`,
			Value: "ddddd",
		},
	}

	for _, test := range tests {
		newVal, err := GetValue(source, test.Path)
		assert.Nil(t, err)
		assert.Equal(t, test.Value, newVal)
	}

}

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

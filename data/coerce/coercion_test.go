package coerce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoerceToString(t *testing.T) {

	var valInt interface{} = 2
	cval, _ := ToString(valInt)
	assert.Equal(t, "2", cval, "not equal")

	var valBool interface{} = true
	cval, _ = ToString(valBool)
	assert.Equal(t, "true", cval, "not equal")

	var valStr interface{} = "12"
	cval, _ = ToString(valStr)
	assert.Equal(t, "12", cval, "not equal")

	var valFloat interface{} = 1.23
	cval, _ = ToString(valFloat)
	assert.Equal(t, "1.23", cval, "not equal")

	var valNil interface{} // = nil
	cval, _ = ToString(valNil)
	assert.Equal(t, "", cval, "not equal")
}

func TestCoerceToInt(t *testing.T) {

	var valInt interface{} = 2
	cval, _ := ToInt(valInt)
	assert.Equal(t, 2, cval, "not equal")

	var valBool interface{} = true
	cval, _ = ToInt(valBool)
	assert.Equal(t, 1, cval, "not equal")

	var valStr interface{} = "12"
	cval, _ = ToInt(valStr)
	assert.Equal(t, 12, cval, "not equal")

	var valFloat interface{} = 1.23
	cval, _ = ToInt(valFloat)
	assert.Equal(t, 1, cval, "not equal")

	var valNil interface{} //= nil
	cval, _ = ToInt(valNil)
	assert.Equal(t, 0, cval, "not equal")
}

func TestCoerceToBoolean(t *testing.T) {

	var valInt interface{} = 2
	cval, _ := ToBool(valInt)
	assert.Equal(t, true, cval, "not equal")

	var valBool interface{} = true
	cval, _ = ToBool(valBool)
	assert.Equal(t, true, cval, "not equal")

	var valStr interface{} = "false"
	cval, _ = ToBool(valStr)
	assert.Equal(t, false, cval, "not equal")

	var valFloat interface{} = 1.23
	cval, _ = ToBool(valFloat)
	assert.Equal(t, true, cval, "not equal")

	var valNil interface{} //= nil
	cval, _ = ToBool(valNil)
	assert.Equal(t, false, cval, "not equal")
}

func TestCoerceToObject(t *testing.T) {

	s := struct {
		Id   int
		Name string
	}{Id: 1001,
		Name: "flogo"}

	obj, err := ToObject(s)
	assert.Nil(t, err)
	assert.Equal(t, "flogo", obj["Name"])
}

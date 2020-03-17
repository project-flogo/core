package coerce

import (
	"reflect"
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

	var pointerString interface{} = StringPointer("flogo")
	cval, _ = ToString(pointerString)
	assert.Equal(t, "flogo", cval, "not equal")
}

func StringPointer(s string) *string {
	return &s
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

func TestCoerceToInt32(t *testing.T) {

	var valInt interface{} = 2
	cval, _ := ToInt32(valInt)
	assert.Equal(t, int32(2), cval)

	var valBool interface{} = true
	cval, _ = ToInt32(valBool)
	assert.Equal(t, int32(1), cval)

	var valStr interface{} = "12"
	cval, _ = ToInt32(valStr)
	assert.Equal(t, int32(12), cval)

	var valFloat interface{} = 1.23
	cval, _ = ToInt32(valFloat)
	assert.Equal(t, int32(1), cval)

	var valNil interface{} //= nil
	cval, _ = ToInt32(valNil)
	assert.Equal(t, int32(0), cval)
}

func TestCoerceToInt64(t *testing.T) {

	var valInt interface{} = 2
	cval, _ := ToInt64(valInt)
	assert.Equal(t, int64(2), cval)

	var valBool interface{} = true
	cval, _ = ToInt64(valBool)
	assert.Equal(t, int64(1), cval)

	var valStr interface{} = "12"
	cval, _ = ToInt64(valStr)
	assert.Equal(t, int64(12), cval)

	var valFloat interface{} = 1.23
	cval, _ = ToInt64(valFloat)
	assert.Equal(t, int64(1), cval)

	var valNil interface{} //= nil
	cval, _ = ToInt64(valNil)
	assert.Equal(t, int64(0), cval)
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

func TestArray(t *testing.T) {
	arr := []string{"a", "b", "c"}
	rArr, err := ToArray(arr)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(rArr))

	str := "a"
	rArr, err = ToArray(str)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(rArr))

}

func TestCoerceArray(t *testing.T) {

	var valInt interface{} = 2
	cval, _ := ToArray(valInt)
	assert.Equal(t, reflect.Slice, reflect.ValueOf(cval).Kind())

	valArr := []string{"a", "b"}
	cval, _ = ToArray(valArr)
	assert.Equal(t, reflect.Slice, reflect.ValueOf(cval).Kind())

	valInfArr := []interface{}{"a", "b", 2, 5}
	cval, _ = ToArray(valInfArr)
	assert.Equal(t, reflect.Slice, reflect.ValueOf(cval).Kind())

	var valNilArr interface{}
	cval, _ = ToArray(valNilArr)
	assert.Nil(t, cval)

}

func TestParams(t *testing.T) {
	var err error
	simpleParams := map[string]string{"params": "a"}
	cval, _ := ToParams(simpleParams)
	assert.NotNil(t, cval)

	simpleStringParams := `{"params": "a"}`
	cval, _ = ToParams(simpleStringParams)
	assert.NotNil(t, cval["params"])

	simpleStringParams = `{"a"}`
	_, err = ToParams(simpleStringParams)
	assert.NotNil(t, err)

	simpleDefault := 3
	_, err = ToParams(simpleDefault)
	assert.NotNil(t, err)
}

func TestConnection(t *testing.T) {

	var err error

	connString := "sampleConn"
	_, err = ToConnection(connString)
	assert.NotNil(t, err)

}

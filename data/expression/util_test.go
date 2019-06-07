package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLiteral(t *testing.T) {

	str := "1"
	val, isLiteral := GetLiteral(str)

	assert.True(t, isLiteral)
	assert.Equal(t, 1, val)

	str = "`abc`"
	val, isLiteral = GetLiteral(str)

	assert.True(t, isLiteral)
	assert.Equal(t, "abc", val)

	str = "1.1"
	val, isLiteral = GetLiteral(str)

	assert.True(t, isLiteral)
	assert.Equal(t, 1.1, val)

	str = "1"
	val, isLiteral = GetLiteral(str)

	assert.True(t, isLiteral)
	assert.Equal(t, 1, val)

	str = `{"a":1}`
	val, isLiteral = GetLiteral(str)

	assert.True(t, isLiteral)
	obj, ok := val.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, 1.0, obj["a"])

	str = `["a","b"]`
	val, isLiteral = GetLiteral(str)

	assert.True(t, isLiteral)
	arr, ok2 := val.([]interface{})
	assert.True(t, ok2)
	assert.Equal(t, 2, len(arr))

	str = "true"
	val, isLiteral = GetLiteral(str)
	assert.True(t, isLiteral)
	b, _ := val.(bool)
	assert.True(t, b)

	str = "false"
	val, isLiteral = GetLiteral(str)
	assert.True(t, isLiteral)
	b, ok2 = val.(bool)
	assert.True(t, ok2)
	assert.False(t, b)
}

func TestIsString(t *testing.T) {
	_, b := isQuotedString("lixingwang")
	assert.False(t, b)

	v, b := isQuotedString(`"ddddd"`)
	assert.True(t, b)
	assert.Equal(t, "ddddd", v)
	_, b = isQuotedString(`"ddddd" ==  "dddddd"`)
	assert.False(t, b)

	_, b = isQuotedString(`"ddddd" ==  'ddddd'`)
	assert.False(t, b)

	//Single

	v, b = isQuotedString(`'ddddd'`)
	assert.True(t, b)
	assert.Equal(t, "ddddd", v)
	_, b = isQuotedString(`'ddddd' ==  'dddd'`)
	assert.False(t, b)

	_, b = isQuotedString("`ddddd` ==  `ddddd`")
	assert.False(t, b)
}

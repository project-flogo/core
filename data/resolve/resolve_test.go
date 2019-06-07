package resolve

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResolveDirectiveDetails(t *testing.T) {

	a := "prop"
	details, err := GetResolveDirectiveDetails(a, false, false)
	assert.Nil(t, err)
	assert.Equal(t, "prop", details.ValueName)
	assert.Equal(t, "", details.ItemName)
	assert.Equal(t, "", details.Path)

	a = "[item]"
	details, err = GetResolveDirectiveDetails(a, true, false)
	assert.Nil(t, err)
	assert.Equal(t, "item", details.ItemName)
	assert.Equal(t, "", details.ValueName)
	assert.Equal(t, "", details.Path)

	a = "[foo.bar.item]"
	details, err = GetResolveDirectiveDetails(a, true, false)
	assert.Nil(t, err)
	assert.Equal(t, "foo.bar.item", details.ItemName)
	assert.Equal(t, "", details.ValueName)
	assert.Equal(t, "", details.Path)

	a = "[item].prop"
	details, err = GetResolveDirectiveDetails(a, true, false)
	assert.Nil(t, err)
	assert.Equal(t, "prop", details.ValueName)
	assert.Equal(t, "item", details.ItemName)
	assert.Equal(t, "", details.Path)

	a = "[mapitem]['foo']"
	details, err = GetResolveDirectiveDetails(a, true, false)
	assert.Nil(t, err)
	assert.Equal(t, "", details.ValueName)
	assert.Equal(t, "mapitem", details.ItemName)
	assert.Equal(t, "['foo']", details.Path)

	a = "[arritem][1]"
	details, err = GetResolveDirectiveDetails(a, true, false)
	assert.Nil(t, err)
	assert.Equal(t, "", details.ValueName)
	assert.Equal(t, "arritem", details.ItemName)
	assert.Equal(t, "[1]", details.Path)

	a = "[item].prop.path"
	details, err = GetResolveDirectiveDetails(a, true, false)
	assert.Nil(t, err)
	assert.Equal(t, "prop", details.ValueName)
	assert.Equal(t, "item", details.ItemName)
	assert.Equal(t, ".path", details.Path)

	a = "[item].prop['map']"
	details, err = GetResolveDirectiveDetails(a, true, false)
	assert.Nil(t, err)
	assert.Equal(t, "prop", details.ValueName)
	assert.Equal(t, "item", details.ItemName)
	assert.Equal(t, "['map']", details.Path)

	a = "[item].prop[0]"
	details, err = GetResolveDirectiveDetails(a, true, false)
	assert.Nil(t, err)
	assert.Equal(t, "prop", details.ValueName)
	assert.Equal(t, "item", details.ItemName)
	assert.Equal(t, "[0]", details.Path)

	a = ".value"
	details, err = GetResolveDirectiveDetails(a, false, false)
	assert.Nil(t, err)
	assert.Equal(t, "value", details.ValueName)
	assert.Equal(t, "", details.ItemName)
	assert.Equal(t, "", details.Path)
}

func TestIsResolveExpr(t *testing.T) {

	str := "1"
	isResolve := IsResolveExpr(str)
	assert.False(t, isResolve)

	str = "$1"
	isResolve = IsResolveExpr(str)
	assert.False(t, isResolve)

	str = "$."
	isResolve = IsResolveExpr(str)
	assert.False(t, isResolve)

	str = "$.1"
	isResolve = IsResolveExpr(str)
	assert.False(t, isResolve)

	str = "$a."
	isResolve = IsResolveExpr(str)
	assert.False(t, isResolve)

	str = "$env[blah]"
	isResolve = IsResolveExpr(str)
	assert.True(t, isResolve)

	str = "$property[blah].ooa"
	isResolve = IsResolveExpr(str)
	assert.True(t, isResolve)

	str = "$property[blah].1oa"
	isResolve = IsResolveExpr(str)
	assert.False(t, isResolve)

	str = "$test.test"
	isResolve = IsResolveExpr(str)
	assert.True(t, isResolve)

	str = "$.test.test"
	isResolve = IsResolveExpr(str)
	assert.True(t, isResolve)

	str = "$.test['blah blah'].test"
	isResolve = IsResolveExpr(str)
	assert.True(t, isResolve)

	str = "$.test[\"blah blah\"].test"
	isResolve = IsResolveExpr(str)
	assert.True(t, isResolve)

	str = "$.test[\"blah blah\"].test + 1"
	isResolve = IsResolveExpr(str)
	assert.False(t, isResolve)
}

func TestPropertyEncodeQuotes(t *testing.T) {
	a := `["foo.bar.item"]`
	details, err := GetResolveDirectiveDetails(a, true, false)
	assert.Nil(t, err)
	assert.Equal(t, "foo.bar.item", details.ItemName)
	assert.Equal(t, "", details.ValueName)
	assert.Equal(t, "", details.Path)

}

func TestIsResolveExprWithExpr(t *testing.T) {
	str := "$activity[REST].id"
	isResolve := IsResolveExpr(str)
	assert.True(t, isResolve)

	str = "$activity[REST].id==1"
	isResolve = IsResolveExpr(str)
	assert.False(t, isResolve)

	str = "$activity[REST][\"123\"]"
	isResolve = IsResolveExpr(str)
	assert.True(t, isResolve)

	str = "$activity[REST][\"123\"][$.index]"
	isResolve = IsResolveExpr(str)
	assert.False(t, isResolve)

	str = "$.body.addresses[$.index]"
	isResolve = IsResolveExpr(str)
	assert.False(t, isResolve)

	str = "$activity[REST].abc ? $.ok: $.false"
	isResolve = IsResolveExpr(str)
	assert.False(t, isResolve)

}

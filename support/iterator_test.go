package support

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_sliceIterator_HasNext(t *testing.T) {

	s := []interface{}{}
	itx := &sliceIterator{s, 0}
	assert.False(t, itx.HasNext())

	s = []interface{}{"foo1", "foo2"}
	itx = &sliceIterator{s, 0}

	assert.True(t, itx.HasNext())
}

func Test_sliceIterator_Next(t *testing.T) {

	s := []interface{}{"foo1", "foo2"}
	itx := &sliceIterator{s, 0}

	assert.Equal(t, "foo1", itx.Next())
	assert.Equal(t, "foo2", itx.Next())
}

func TestFixedDetails_Get(t *testing.T) {

	fd := FixedDetails{data: map[string]string{"v1": "foo1", "v2": "foo2"}}

	v := fd.Get("v")
	assert.Empty(t, v)
	v1 := fd.Get("v1")
	assert.Equal(t, "foo1", v1)
	v2 := fd.Get("v2")
	assert.Equal(t, "foo2", v2)
}

func TestFixedDetails_Iterate(t *testing.T) {
	fd := FixedDetails{data: map[string]string{"v1": "foo1", "v2": "foo2"}}

	var vals []string
	fd.Iterate(func(s string, s2 string) {
		vals = append(vals, s2+"bar")
	})

	assert.Len(t, fd.data, 2)

	for _, val := range vals {
		assert.True(t, strings.Contains(val, "bar"))
	}
}

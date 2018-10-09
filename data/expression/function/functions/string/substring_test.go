package string

import (
	"github.com/project-flogo/fscript/function"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFnSubstring_Eval(t *testing.T) {
	f := &fnSubstring{}
	v, err := function.Eval(f, "abc", 1, -1)
	assert.Nil(t, err)
	assert.Equal(t, "bc", v)

	v, err = function.Eval(f, "abc", 1, 1)
	assert.Nil(t, err)
	assert.Equal(t, "b", v)
}

package builtin

import (
	"github.com/project-flogo/fscript/function"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFnLen_Eval(t *testing.T) {
	f := &fnLen{}
	v, err := function.Eval(f, "abc")
	assert.Nil(t, err)
	assert.Equal(t, 3, v)

	v, err = function.Eval(f, []string{"a", "b"})
	assert.Nil(t, err)
	assert.Equal(t, 2, v)
}

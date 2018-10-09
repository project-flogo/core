package string

import (
	"github.com/project-flogo/fscript/function"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFnEquals_Eval(t *testing.T) {
	f := &fnEquals{}

	v, err := function.Eval(f, "foo", "bar")
	assert.Nil(t, err)
	assert.False(t, v.(bool))

	v, err = function.Eval(f, "foo", "foo")
	assert.Nil(t, err)
	assert.True(t, v.(bool))
}

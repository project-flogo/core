package string

import (
	"github.com/project-flogo/fscript/function"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFnConcat_Eval(t *testing.T) {
	f := &fnConcat{}
	v, err := function.Eval(f, "a", "b")

	assert.Nil(t, err)
	assert.Equal(t, "ab", v)
}

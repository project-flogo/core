package support

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGenerator(t *testing.T) {

	gen, err := NewGenerator()
	assert.Nil(t, err)
	assert.NotNil(t, gen)
	assert.Equal(t, uint64(0), gen.counter)
	assert.NotNil(t, gen.seed)
}

func TestGenerator_Next(t *testing.T) {

	gen, err := NewGenerator()
	assert.Nil(t, err)
	assert.NotNil(t, gen)

	uuid := gen.Next()
	assert.NotNil(t, uuid)
	assert.Equal(t, uint64(1), gen.counter)
}

func TestGenerator_NextAsString(t *testing.T) {

	gen, err := NewGenerator()
	assert.Nil(t, err)
	assert.NotNil(t, gen)

	uuid := gen.NextAsString()
	assert.NotEmpty(t, uuid)
	assert.Equal(t, uint64(1), gen.counter)
}

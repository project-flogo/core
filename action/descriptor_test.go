package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	inf := NewInfo(true, true)

	assert.NotNil(t, inf)

	assert.Equal(t, true, inf.Async())

	assert.Equal(t, true, inf.Passthru())

	assert.Equal(t, "", inf.Id())
}

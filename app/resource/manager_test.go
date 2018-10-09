package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTypeFromID(t *testing.T) {

	resType, err := GetTypeFromID("flow:myflow")
	assert.Nil(t, err)
	assert.Equal(t, "flow", resType)
}

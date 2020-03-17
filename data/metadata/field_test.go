package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestField(t *testing.T) {
	fieldDetails := NewFieldDetails("sample", "string", "required,allowed(GET)")

	assert.NotNil(t, fieldDetails)
	assert.Equal(t, 7, len(fieldDetails.AllowedToString()))

	assert.Nil(t, fieldDetails.Validate("GET"))
	assert.NotNil(t, fieldDetails.Validate(1))
}

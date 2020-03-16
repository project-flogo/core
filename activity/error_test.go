package activity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := NewError("sample error", "101", "sample error data")
	assert.NotNil(t, err)

	err = NewRetriableError("sample error", "102", "sample error data")
	assert.NotNil(t, err)

	assert.NotNil(t, err.Error())
	assert.NotNil(t, err.Code())
	assert.NotNil(t, err.Data())
	assert.NotNil(t, err.Retriable())
	assert.Equal(t, "", err.ActivityName())
}

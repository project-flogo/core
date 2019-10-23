package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIsMappingRelaxed(t *testing.T) {

	os.Setenv(EnvMappingIgnoreError, "true")

	assert.True(t, IsMappingIgnoreError())

	os.Unsetenv(EnvMappingIgnoreError)
	assert.False(t, IsMappingIgnoreError())
	defer func() {
		os.Unsetenv(EnvMappingIgnoreError)
	}()

}

package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIsMappingRelaxed(t *testing.T) {

	os.Setenv(EnvMappingIgnoreError, "true")

	assert.True(t, IsMappingIgnoreErrorsOn())

	os.Unsetenv(EnvMappingIgnoreError)
	assert.False(t, IsMappingIgnoreErrorsOn())
	defer func() {
		os.Unsetenv(EnvMappingIgnoreError)
	}()

}

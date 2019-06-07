package mapper

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIsMappingRelaxed(t *testing.T) {

	os.Setenv(EnvMappingRelexed, "true")

	assert.True(t, IsMappingRelaxed())

	os.Unsetenv(EnvMappingRelexed)
	assert.False(t, IsMappingRelaxed())
	defer func() {
		os.Unsetenv(EnvMappingRelexed)
	}()

}

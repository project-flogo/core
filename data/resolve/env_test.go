package resolve

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvResolver_GetResolverInfo(t *testing.T) {
	resolver := &EnvResolver{}
	assert.NotNil(t, resolver.GetResolverInfo())
	assert.True(t, resolver.GetResolverInfo().IsStatic())
	assert.True(t, resolver.GetResolverInfo().UsesItemFormat())
}

func TestEnvResolver_Resolve(t *testing.T) {
	resolver := &EnvResolver{}

	path, _ := os.LookupEnv("PATH")
	v, err := resolver.Resolve(nil, "PATH", "")
	assert.Nil(t, err)
	assert.Equal(t, path, v)

	env, _ := os.LookupEnv("NONEXISTANT_ENV_123")
	v, err = resolver.Resolve(nil, "NONEXISTANT_ENV_123", "")
	assert.NotNil(t, err)
	assert.Equal(t, env, v)
}

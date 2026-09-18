package resolve

import (
	"os"
	"testing"

	"github.com/project-flogo/core/engine/secret"
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

func TestEnvResolverDecryptsSecret(t *testing.T) {
	secret.SetSecretValueHandler(&secret.KeyBasedSecretValueHandler{Key: "mysecretkey"})
	defer secret.SetSecretValueHandler(nil)

	encoded, err := secret.GetSecretValueHandler().EncodeValue("s3cr3t")
	assert.Nil(t, err)

	_ = os.Setenv("FLOGO_TEST_SECRET", "SECRET:"+encoded)
	defer func() { _ = os.Unsetenv("FLOGO_TEST_SECRET") }()

	resolver := &EnvResolver{}
	val, err := resolver.Resolve(nil, "FLOGO_TEST_SECRET", "")
	assert.Nil(t, err)
	assert.Equal(t, "s3cr3t", val)
}

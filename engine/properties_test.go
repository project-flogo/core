package engine

import (
	"os"
	"testing"

	"github.com/project-flogo/core/engine/secret"
	"github.com/stretchr/testify/assert"
)

func TestEnvPropertyProcessorPlain(t *testing.T) {
	_ = os.Setenv("FLOGO_PROP_PLAIN", "plainvalue")
	defer func() { _ = os.Unsetenv("FLOGO_PROP_PLAIN") }()

	props := map[string]interface{}{"myProp": "$env[FLOGO_PROP_PLAIN]"}
	err := EnvPropertyProcessor(props)
	assert.Nil(t, err)
	assert.Equal(t, "plainvalue", props["myProp"])
}

func TestEnvPropertyProcessorDecryptsSecret(t *testing.T) {
	secret.SetSecretValueHandler(&secret.KeyBasedSecretValueHandler{Key: "mysecretkey"})
	defer secret.SetSecretValueHandler(nil)

	encoded, err := secret.GetSecretValueHandler().EncodeValue("s3cr3t")
	assert.Nil(t, err)

	_ = os.Setenv("FLOGO_PROP_SECRET", "SECRET:"+encoded)
	defer func() { _ = os.Unsetenv("FLOGO_PROP_SECRET") }()

	props := map[string]interface{}{"myProp": "$env[FLOGO_PROP_SECRET]"}
	err = EnvPropertyProcessor(props)
	assert.Nil(t, err)
	assert.Equal(t, "s3cr3t", props["myProp"])
}

func TestEnvPropertyProcessorMissingVar(t *testing.T) {
	props := map[string]interface{}{"myProp": "$env[FLOGO_PROP_MISSING]"}
	err := EnvPropertyProcessor(props)
	assert.NotNil(t, err)
}

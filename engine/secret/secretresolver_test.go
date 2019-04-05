package secret

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretKeyDefault(t *testing.T) {
	defer func() {
		SetSecretValueHandler(nil)
	}()
	handler := GetSecretValueHandler()
	encoded, err := handler.EncodeValue("mysecurepassword3")
	assert.Nil(t, err)
	decoded, err := handler.DecodeValue(encoded)
	assert.Nil(t, err)
	assert.Equal(t, "mysecurepassword3", decoded)
}

func TestSecretKeyEnv(t *testing.T) {
	_ = os.Setenv(EnvKeyDataSecretKey, "mysecretkey1")
	defer func() {
		_ = os.Unsetenv(EnvKeyDataSecretKey)
		SetSecretValueHandler(nil)
	}()

	handler := GetSecretValueHandler()
	encoded, err := handler.EncodeValue("mysecurepassword1")
	assert.Nil(t, err)
	decoded, err := handler.DecodeValue(encoded)
	assert.Nil(t, err)
	assert.Equal(t, "mysecurepassword1", decoded)
}

func TestSecretKey(t *testing.T) {
	defer func() {
		SetSecretValueHandler(nil)
	}()
	SetSecretValueHandler(&KeyBasedSecretValueHandler{Key: "mysecretkey2"})
	handler := GetSecretValueHandler()
	encoded, err := handler.EncodeValue("mysecurepassword1")
	assert.Nil(t, err)
	decoded, err := GetSecretValueHandler().DecodeValue(encoded)
	assert.Nil(t, err)
	assert.Equal(t, "mysecurepassword1", decoded)
}

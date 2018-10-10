package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	EnvKeyDataSecretKey  = "FLOGO_DATA_SECRET_KEY"
	defaultDataSecretKey = "flogo"
)

var secretValueHandler SecretValueHandler

// SecretValueDecoder defines method for decoding value
type SecretValueHandler interface {
	EncodeValue(value interface{}) (string, error)
	DecodeValue(value interface{}) (string, error)
}

// Set secret value decoder
func SetSecretValueHandler(pwdResolver SecretValueHandler) {
	secretValueHandler = pwdResolver
}

func GetDataSecretKey() string {
	key := os.Getenv(EnvKeyDataSecretKey)
	if len(key) > 0 {
		return key
	}
	return defaultDataSecretKey
}

// Get secret value handler. If not already set by SetSecretValueHandler(), will return default KeyBasedSecretValueDecoder
// where decoding key value is expected to be set through FLOGO_DATA_SECRET_KEY environment variable.
// If key is not set, a default key value(github.com/project-flogo/core/config.DATA_SECRET_KEY_DEFAULT) will be used.
func GetSecretValueHandler() SecretValueHandler {
	if secretValueHandler == nil {
		secretValueHandler = &KeyBasedSecretValueHandler{Key: GetDataSecretKey()}
	}
	return secretValueHandler
}

// A key based secret value decoder. Secret value encryption/decryption is based on SHA256
// and uses implementation from https://gist.github.com/willshiao/f4b03650e5a82561a460b4a15789cfa1
type KeyBasedSecretValueHandler struct {
	Key string
}

// Decode value based on a key
func (defaultResolver *KeyBasedSecretValueHandler) DecodeValue(value interface{}) (string, error) {
	if value != nil {
		if defaultResolver.Key != "" {
			kBytes := sha256.Sum256([]byte(defaultResolver.Key))
			return decryptValue(kBytes[:], value.(string))
		}
		return value.(string), nil
	}
	return "", nil
}

func (defaultResolver *KeyBasedSecretValueHandler) EncodeValue(value interface{}) (string, error) {
	if value != nil {
		if defaultResolver.Key != "" {
			kBytes := sha256.Sum256([]byte(defaultResolver.Key))
			return encryptValue(kBytes[:], value.(string))
		}
		return value.(string), nil
	}
	return "", nil
}

// encrypt string to base64 crypto using AES
func encryptValue(key []byte, text string) (string, error) {
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt from base64 to decrypted string
func decryptValue(key []byte, encryptedData string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return fmt.Sprintf("%s", ciphertext), nil
}

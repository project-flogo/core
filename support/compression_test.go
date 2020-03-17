package support

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeAndUnzip(t *testing.T) {

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write([]byte("this is a test"))
	assert.Nil(t, err)
	_ = zw.Close()
	if err != nil {
		return
	}
	cBytes := buf.Bytes()
	encoded := base64.StdEncoding.EncodeToString(cBytes)

	val, err := DecodeAndUnzip(encoded)
	assert.Nil(t, err)
	assert.Equal(t, []byte("this is a test"), val)
}

func TestUnzip(t *testing.T) {

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write([]byte("this is a test"))
	assert.Nil(t, err)
	_ = zw.Close()
	if err != nil {
		return
	}

	cBytes := buf.Bytes()
	val, err := Unzip(cBytes)
	assert.Nil(t, err)
	assert.Equal(t, []byte("this is a test"), val)
}

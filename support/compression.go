package support

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
)

func DecodeAndUnzip(encoded string) ([]byte, error) {

	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	return Unzip(decoded)
}

func Unzip(compressed []byte) ([]byte, error) {

	buf := bytes.NewBuffer(compressed)
	r, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	uncompressedBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return uncompressedBytes, nil
}

package log

import (
	"bytes"
	"testing"
)

var  sVal string
var val string

func BenchmarkConcat(b *testing.B) {

	val = "xyz"
	var str string
	for n := 0; n < b.N; n++ {
		str = "[" + val + "]"
	}
	sVal = str
}

func BenchmarkBuffer(b *testing.B) {

	var str string
	val = "xyz"
	for n := 0; n < b.N; n++ {
		var buffer bytes.Buffer
		buffer.Write([]byte("["))
		buffer.Write([]byte(val))
		buffer.Write([]byte("]"))
		str = buffer.String()
	}

	sVal = str
}

func BenchmarkBufferString(b *testing.B) {

	var str string
	val = "xyz"
	for n := 0; n < b.N; n++ {
		var buffer bytes.Buffer
		buffer.WriteString("[")
		buffer.WriteString(val)
		buffer.WriteString("]")
		str = buffer.String()
	}

	sVal = str
}
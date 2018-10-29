package log

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

type Field = interface{}

func Binary(key string, val []byte) Field {
	return zap.Binary(key,val)
}

func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

func Bools(key string, vals []bool) Field {
	return zap.Bools(key, vals)
}

func ByteString(key string, val []byte) Field {
	return zap.ByteString(key, val)
}

func ByteStrings(key string, vals [][]byte) Field {
	return zap.ByteStrings(key, vals)
}

func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

func Durations(key string, vals []time.Duration) Field {
	return zap.Durations(key, vals)
}

func Error(err error) Field {
	return zap.Error(err)
}

func Errors(key string, errs []error) Field {
	return zap.Errors(key, errs)
}

func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

func Float64s(key string, vals []float64) Field {
	return zap.Float64s(key, vals)
}

func Float32(key string, val float32) Field {
	return zap.Float32(key, val)
}

func Float32s(key string, vals []float32) Field {
	return zap.Float32s(key, vals)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Ints(key string, vals []int) Field {
	return zap.Ints(key, vals)
}

func Int32(key string, val int32) Field {
	return zap.Int32(key, val)
}

func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func Int64s(key string, vals []int64) Field {
	return zap.Int64s(key, vals)
}

func NamedError(key string, err error) Field {
	return zap.NamedError(key, err)
}

// Namespace see zap.Namespace
func Namespace(key string) Field {
	return zap.Namespace(key)
}

// Object encodes object using reflection
func Object(key string, val interface{}) Field {
	return zap.Reflect(key, val)
}

// Skip see zap.Skip
func Skip() Field {
	return zap.Skip()
}

func Stack(key string) Field {
	return zap.Stack(key)
}

func String(key string, val string) Field {
	return zap.String(key, val)
}

func Strings(key string, vals []string) Field {
	return zap.Strings(key, vals)
}

func Stringer(key string, val fmt.Stringer) Field {
	return zap.Stringer(key, val)
}

func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

func Times(key string, vals []time.Time) Field {
	return zap.Times(key, vals)
}

func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}

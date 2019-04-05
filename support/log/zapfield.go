package log

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

func FieldBinary(key string, val []byte) Field {
	return zap.Binary(key, val)
}

func FieldBool(key string, val bool) Field {
	return zap.Bool(key, val)
}

func FieldBools(key string, vals []bool) Field {
	return zap.Bools(key, vals)
}

func FieldByteString(key string, val []byte) Field {
	return zap.ByteString(key, val)
}

func FieldByteStrings(key string, vals [][]byte) Field {
	return zap.ByteStrings(key, vals)
}

func FieldDuration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

func FieldDurations(key string, vals []time.Duration) Field {
	return zap.Durations(key, vals)
}

func FieldError(err error) Field {
	return zap.Error(err)
}

func FieldErrors(key string, errs []error) Field {
	return zap.Errors(key, errs)
}

func FieldFloat64(key string, val float64) Field {
	return zap.Float64(key, val)
}

func FieldFloat64s(key string, vals []float64) Field {
	return zap.Float64s(key, vals)
}

func FieldFloat32(key string, val float32) Field {
	return zap.Float32(key, val)
}

func FieldFloat32s(key string, vals []float32) Field {
	return zap.Float32s(key, vals)
}

func FieldInt(key string, val int) Field {
	return zap.Int(key, val)
}

func FieldInts(key string, vals []int) Field {
	return zap.Ints(key, vals)
}

func FieldInt32(key string, val int32) Field {
	return zap.Int32(key, val)
}

func FieldInt64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func FieldInt64s(key string, vals []int64) Field {
	return zap.Int64s(key, vals)
}

func FieldNamedError(key string, err error) Field {
	return zap.NamedError(key, err)
}

// Namespace see zap.Namespace
func FieldNamespace(key string) Field {
	return zap.Namespace(key)
}

// Object encodes object using reflection
func FieldObject(key string, val interface{}) Field {
	return zap.Reflect(key, val)
}

// Skip see zap.Skip
func FieldSkip() Field {
	return zap.Skip()
}

func FieldStack(key string) Field {
	return zap.Stack(key)
}

func FieldString(key string, val string) Field {
	return zap.String(key, val)
}

func FieldStrings(key string, vals []string) Field {
	return zap.Strings(key, vals)
}

func FieldStringer(key string, val fmt.Stringer) Field {
	return zap.Stringer(key, val)
}

func FieldTime(key string, val time.Time) Field {
	return zap.Time(key, val)
}

func FieldTimes(key string, vals []time.Time) Field {
	return zap.Times(key, vals)
}

func FieldAny(key string, val interface{}) Field {
	return zap.Any(key, val)
}

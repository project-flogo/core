package coerce

import (
	"encoding/json"
	"strings"

	"github.com/project-flogo/core/data"
)

func init() {
	data.SetAttributeTypeConverter(ToType)
}

func NewTypedValue(dataType data.Type, value interface{}) (data.TypedValue, error) {
	newVal, err := ToType(value, dataType)
	if err != nil {
		return nil, err
	}
	return data.NewTypedValue(dataType, newVal), nil
}

// ToType coerce a value to the specified type
func ToType(value interface{}, dataType data.Type) (interface{}, error) {

	var coerced interface{}
	var err error

	switch dataType {
	case data.TypeAny:
		coerced, err = ToAny(value)
	case data.TypeString:
		coerced, err = ToString(value)
	case data.TypeInt:
		coerced, err = ToInt(value)
	case data.TypeInt32:
		coerced, err = ToInt32(value)
	case data.TypeInt64:
		coerced, err = ToInt64(value)
	case data.TypeFloat32:
		coerced, err = ToFloat32(value)
	case data.TypeFloat64:
		coerced, err = ToFloat64(value)
	case data.TypeBool:
		coerced, err = ToBool(value)
	case data.TypeBytes:
		coerced, err = ToBytes(value)
	case data.TypeParams:
		coerced, err = ToParams(value)
	case data.TypeObject:
		coerced, err = ToObject(value)
	case data.TypeArray:
		coerced, err = ToArrayIfNecessary(value)
	case data.TypeConnection:
		coerced, err = ToConnection(value)
	case data.TypeDateTime:
		coerced, err = ToDateTime(value)
	case data.TypeUnknown:
		coerced = value
	}
	if err != nil {
		return nil, err
	}

	return coerced, nil
}

// ToAny coerce a value to generic value
func ToAny(val interface{}) (interface{}, error) {

	switch t := val.(type) {

	case json.Number:
		if strings.Contains(t.String(), ".") {
			return t.Float64()
		} else {
			return t.Int64()
		}
	default:
		return val, nil
	}
}

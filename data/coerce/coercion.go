package coerce

import (
	"encoding/json"
	"fmt"
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
	case data.TypeComplexObject:
		coerced, err = CoerceToComplexObject(value)
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

//DEPRECATED
// CoerceToObject coerce a value to an complex object
func CoerceToComplexObject(val interface{}) (*data.ComplexObject, error) {
	//If the val is nil then just return empty struct
	var emptyComplexObject = &data.ComplexObject{Value: "{}"}
	if val == nil {
		return emptyComplexObject, nil
	}
	switch t := val.(type) {
	case string:
		if val == "" {
			return emptyComplexObject, nil
		} else {
			complexObject := &data.ComplexObject{}
			err := json.Unmarshal([]byte(t), complexObject)
			if err != nil {
				return nil, err

			}
			return handleComplex(complexObject), nil
		}
	case map[string]interface{}:
		v, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		complexObject := &data.ComplexObject{}
		err = json.Unmarshal(v, complexObject)
		if err != nil {
			return nil, err
		}
		return handleComplex(complexObject), nil
	case *data.ComplexObject:
		return handleComplex(val.(*data.ComplexObject)), nil
	default:
		return nil, fmt.Errorf("unable to coerce %#v to complex object", val)
	}
}

func handleComplex(complex *data.ComplexObject) *data.ComplexObject {
	if complex != nil {
		if complex.Value == "" {
			complex.Value = "{}"
		}
	}
	return complex
}

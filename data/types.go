package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Type denotes a data type
type Type int

//var intSize = strconv.IntSize

const (
	TypeUnknown Type = iota // interface{}

	TypeAny // interface{}

	// simple types
	TypeString   // string
	TypeInt      // int
	TypeInt32    // int32
	TypeInt64    // int64
	TypeFloat32  // float32
	TypeFloat64  // float64
	TypeBool     // bool
	TypeObject   // map[string]interface{}
	TypeBytes    // []byte
	TypeDateTime // time

	// compound types
	TypeParams // map[string]string
	TypeArray  // []interface{}
	TypeMap    // map[interface{}]interface{}

	//Special Type
	TypeConnection
)

var types = [...]string{
	"unknown",
	"any",
	"string",
	"int",
	"int32",
	"int64",
	"float32",
	"float64",
	"bool",
	"object",
	"bytes",
	"datetime",
	"params",
	"array",
	"map",
	"connection",
}

func (t Type) String() string {
	return types[t]
}

// IsSimple returns wither a type is simple
func (t Type) IsSimple() bool {
	return t > 1 && t < 8
}

var names = map[Type]string{
	TypeUnknown:    "TypeUnknown",
	TypeAny:        "TypeAny",
	TypeString:     "TypeString",
	TypeInt:        "TypeInt",
	TypeInt32:      "TypeInt32",
	TypeInt64:      "TypeInt64",
	TypeFloat32:    "TypeFloat32",
	TypeFloat64:    "TypeFloat64",
	TypeBool:       "TypeBool",
	TypeObject:     "TypeObject",
	TypeBytes:      "TypeBytes",
	TypeDateTime:   "TypeDateTime",
	TypeParams:     "TypeParams",
	TypeArray:      "TypeArray",
	TypeMap:        "TypeMap",
	TypeConnection: "TypeConnection",
}

// Name returns the name of the type
func (t Type) Name() string {
	return names[t]
}

// ToTypeEnum get the data type that corresponds to the specified name
func ToTypeEnum(typeStr string) (Type, error) {

	switch strings.ToLower(typeStr) {
	case "any":
		return TypeAny, nil
	case "string":
		return TypeString, nil
	case "int", "integer":
		return TypeInt, nil
	case "int32":
		return TypeInt32, nil
	case "int64", "long":
		return TypeInt64, nil
	case "float32", "float":
		return TypeFloat32, nil
	case "float64", "double":
		return TypeFloat64, nil
	case "bool", "boolean":
		return TypeBool, nil
	case "object":
		return TypeObject, nil
	case "bytes":
		return TypeBytes, nil
	case "datetime":
		return TypeDateTime, nil
	case "params":
		return TypeParams, nil
	case "array":
		return TypeArray, nil
	case "map":
		return TypeMap, nil
	case "connection":
		return TypeConnection, nil
	default:
		return TypeUnknown, errors.New("unknown type: " + typeStr)
	}
}

// GetType get the Type of the supplied value
func GetType(val interface{}) (Type, error) {

	switch t := val.(type) {
	case string:
		return TypeString, nil
	case int:
		if strconv.IntSize == 64 {
			return TypeInt64, nil
		}
		return TypeInt, nil
	case int32:
		return TypeInt32, nil
	case int64:
		return TypeInt64, nil
	case float32:
		return TypeFloat32, nil
	case float64:
		return TypeFloat64, nil
	case json.Number:
		if strings.Contains(t.String(), ".") {
			return TypeFloat64, nil
		} else {
			return TypeInt64, nil
		}
	case bool:
		return TypeBool, nil
	case map[string]interface{}:
		return TypeObject, nil
	case []interface{}:
		return TypeArray, nil
	case map[string]string:
		return TypeParams, nil
	case []byte:
		return TypeBytes, nil
	case time.Time:
		return TypeDateTime, nil
	default:
		return TypeUnknown, fmt.Errorf("unable to determine type of %#v", t)
	}
}

func ToTypeFromGoRep(strType string) Type {

	dt := TypeUnknown

	switch strType {
	case "interface{}", "interface {}":
		dt = TypeAny
	case "string":
		dt = TypeString
	case "int":
		dt = TypeInt
	case "int32":
		dt = TypeInt32
	case "int64":
		dt = TypeInt64
	case "float32":
		dt = TypeFloat32
	case "float64":
		dt = TypeFloat64
	case "bool":
		dt = TypeBool
	case "map[string]interface{}":
		dt = TypeObject
	case "[]byte":
		dt = TypeBytes
	case "time.Time":
		dt = TypeDateTime
	case "map[string]string":
		dt = TypeParams
	case "connection.Manager":
		dt = TypeConnection
	}

	return dt
}

func IsSimpleType(val interface{}) bool {
	switch val.(type) {
	case string, int, int32, int64, float32, float64, json.Number, bool:
		return true
	default:
		return false
	}
}

var typeConverter TypeConverter

func SetAttributeTypeConverter(converter TypeConverter) {
	typeConverter = converter
}

type TypeConverter func(value interface{}, toType Type) (interface{}, error)

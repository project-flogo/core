package coerce

import (
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"strconv"
	"time"
)

// ToString coerce a value to a string
func ToString(val interface{}) (string, error) {

	switch t := val.(type) {
	case string:
		return t, nil
	case *string:
		return *t, nil
	case int:
		return strconv.Itoa(t), nil
	case int64:
		return strconv.FormatInt(t, 10), nil
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 64), nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	case json.Number:
		return t.String(), nil
	case bool:
		return strconv.FormatBool(t), nil
	case nil:
		return "", nil
	case []byte:
		return string(t), nil
	case time.Time:
		b, err := t.MarshalText()
		if err != nil {
			return "", err
		}
		return string(b), nil
	default:
		b, err := json.Marshal(t)
		if err != nil {
			return "", fmt.Errorf("unable to coerce %#v to string", t)
		}
		return string(b), nil
	}
}

// ToInt coerce a value to an int
func ToInt(val interface{}) (int, error) {
	switch t := val.(type) {
	case int:
		return t, nil
	case uint:
		return int(t), nil
	case int8:
		return int(t), nil
	case uint8:
		return int(t), nil
	case int16:
		return int(t), nil
	case uint16:
		return int(t), nil
	case int32:
		return int(t), nil
	case uint32:
		return int(t), nil
	case int64:
		return int(t), nil
	case uint64:
		return int(t), nil
	case float32:
		return int(t), nil
	case float64:
		return int(t), nil
	case json.Number:
		i, err := t.Int64()
		return int(i), err
	case string:
		return strconv.Atoi(t)
	case bool:
		if t {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to coerce %#v to int", val)
	}
}

// ToInteger coerce a value to an integer
func ToInt32(val interface{}) (int32, error) {
	switch t := val.(type) {
	case int:
		return int32(t), nil
	case uint:
		return int32(t), nil
	case int8:
		return int32(t), nil
	case uint8:
		return int32(t), nil
	case int16:
		return int32(t), nil
	case uint16:
		return int32(t), nil
	case int32:
		return t, nil
	case uint32:
		return int32(t), nil
	case int64:
		return int32(t), nil
	case uint64:
		return int32(t), nil
	case float32:
		return int32(t), nil
	case float64:
		return int32(t), nil
	case json.Number:
		i, err := t.Int64()
		return int32(i), err
	case string:
		i, err := strconv.Atoi(t)
		return int32(i), err
	case bool:
		if t {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to coerce %#v to int32", val)
	}
}

// ToInteger coerce a value to an integer
func ToInt64(val interface{}) (int64, error) {
	switch t := val.(type) {
	case int:
		return int64(t), nil
	case uint:
		return int64(t), nil
	case int8:
		return int64(t), nil
	case uint8:
		return int64(t), nil
	case int16:
		return int64(t), nil
	case uint16:
		return int64(t), nil
	case int32:
		return int64(t), nil
	case uint32:
		return int64(t), nil
	case int64:
		return t, nil
	case uint64:
		return int64(t), nil
	case float32:
		return int64(t), nil
	case float64:
		return int64(t), nil
	case json.Number:
		return t.Int64()
	case string:
		return strconv.ParseInt(t, 10, 64)
	case bool:
		if t {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to coerce %#v to integer", val)
	}
}

// ToFloat32 coerce a value to a double/float32
func ToFloat32(val interface{}) (float32, error) {
	switch t := val.(type) {
	case float32:
		return t, nil
	case float64:
		return float32(t), nil
	case int:
		return float32(t), nil
	case int8:
		return float32(t), nil
	case int32:
		return float32(t), nil
	case int64:
		return float32(t), nil
	case uint:
		return float32(t), nil
	case uint8:
		return float32(t), nil
	case uint16:
		return float32(t), nil
	case uint32:
		return float32(t), nil
	case uint64:
		return float32(t), nil
	case json.Number:
		f, err := t.Float64()
		return float32(f), err
	case string:
		f, err := strconv.ParseFloat(t, 32)
		return float32(f), err
	case bool:
		if t {
			return 1.0, nil
		}
		return 0.0, nil
	case nil:
		return 0.0, nil
	default:
		return 0.0, fmt.Errorf("unable to coerce %#v to float32", val)
	}
}

// ToFloat64 coerce a value to a double/float64
func ToFloat64(val interface{}) (float64, error) {
	switch t := val.(type) {
	case float32:
		return float64(t), nil
	case float64:
		return t, nil
	case int:
		return float64(t), nil
	case int8:
		return float64(t), nil
	case int32:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case uint:
		return float64(t), nil
	case uint8:
		return float64(t), nil
	case uint16:
		return float64(t), nil
	case uint32:
		return float64(t), nil
	case uint64:
		return float64(t), nil
	case json.Number:
		return t.Float64()
	case string:
		return strconv.ParseFloat(t, 64)
	case bool:
		if t {
			return 1.0, nil
		}
		return 0.0, nil
	case nil:
		return 0.0, nil
	default:
		return 0.0, fmt.Errorf("unable to coerce %#v to float64", val)
	}
}

// ToBool coerce a value to a boolean
func ToBool(val interface{}) (bool, error) {
	switch t := val.(type) {
	case bool:
		return t, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return t != 0, nil
	case float64:
		return t != 0.0, nil
	case json.Number:
		i, err := t.Int64()
		return i != 0, err
	case string:
		return strconv.ParseBool(t)
	case nil:
		return false, nil
	default:
		str, err := ToString(val)
		if err != nil {
			return false, fmt.Errorf("unable to coerce %#v to bool", val)
		}
		return strconv.ParseBool(str)
	}
}

// ToBytes coerce a value to a byte array
func ToBytes(val interface{}) ([]byte, error) {

	switch t := val.(type) {
	case []byte:
		return t, nil
	case string:
		return []byte(t), nil
	case nil:
		return nil, nil
	default:
		//for now just convert everything to string then bytes
		s, err := ToString(val)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce %#v to bytes", t)
		}
		return []byte(s), nil
	}
}

func ToDateTime(val interface{}) (time.Time, error) {
	switch t := val.(type) {
	case time.Time:
		return t, nil
	case int64:
		return time.Unix(t, 0), nil
	case float64:
		return time.Unix(int64(t), 0), nil
	default:
		dateVal, err := ToString(val)
		if err != nil {
			return time.Time{}, nil
		}
		tm, err := dateparse.ParseAny(dateVal)
		if err != nil {
			return tm, fmt.Errorf("parse [%s] to time error: %s", dateVal, err.Error())
		}
		return tm, nil
	}
}

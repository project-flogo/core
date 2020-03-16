package coerce

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// ToParams coerce a value to params
func ToParams(val interface{}) (map[string]string, error) {

	switch t := val.(type) {
	case map[string]string:
		return t, nil
	case string:
		m := make(map[string]string)
		if t != "" {
			err := json.Unmarshal([]byte(t), &m)
			if err != nil {

				m, err = toParams(t)
				if err != nil {
					return nil, fmt.Errorf("unable to coerce %#v to params", val)
				}
			}
		}
		return m, nil
	case map[string]interface{}:

		var m = make(map[string]string, len(t))
		for k, v := range t {

			mVal, err := ToString(v)
			if err != nil {
				return nil, err
			}
			m[k] = mVal
		}
		return m, nil
	case map[interface{}]string:

		var m = make(map[string]string, len(t))
		for k, v := range t {

			mKey, err := ToString(k)
			if err != nil {
				return nil, err
			}
			m[mKey] = v
		}
		return m, nil
	case map[interface{}]interface{}:

		var m = make(map[string]string, len(t))
		for k, v := range t {

			mKey, err := ToString(k)
			if err != nil {
				return nil, err
			}

			mVal, err := ToString(v)
			if err != nil {
				return nil, err
			}
			m[mKey] = mVal
		}
		return m, nil
	case interface{}:
		s, err := ToString(t)
		if err != nil {
			return nil, err
		}
		return toParams(s)
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unable to coerce %#v to map[string]string", val)
	}
}

func toParams(values string) (map[string]string, error) {

	var params map[string]string

	result := strings.Split(values, ",")
	params = make(map[string]string)
	for _, pair := range result {
		nv := strings.Split(pair, "=")
		if len(nv) != 2 {
			return nil, fmt.Errorf("invalid params")
		}
		params[nv[0]] = nv[1]
	}

	return params, nil
}

// ToObject coerce a value to an object
func ToObject(val interface{}) (map[string]interface{}, error) {

	switch t := val.(type) {
	case map[string]interface{}:
		return t, nil
	case map[string]string:
		ret := make(map[string]interface{}, len(t))
		for key, value := range t {
			ret[key] = value
		}
		return ret, nil
	case string:
		m := make(map[string]interface{})
		if t != "" {
			err := json.Unmarshal([]byte(t), &m)
			if err != nil {
				return nil, fmt.Errorf("unable to coerce %#v to map[string]interface{}", val)
			}
		}
		return m, nil
	case interface{}:
		//Try to convert interface to an object.
		s, err := ToString(t)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce %#v to string", val)
		}
		return ToObject(s)
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unable to coerce %#v to map[string]interface{}", val)
	}
}

// ToArray coerce a value to an array of empty interface values
func ToArray(val interface{}) ([]interface{}, error) {

	switch t := val.(type) {
	case []interface{}:
		return t, nil

	case []map[string]interface{}:
		var a []interface{}
		for _, v := range t {
			a = append(a, v)
		}
		return a, nil
	case string:
		a := make([]interface{}, 0)
		if t != "" {
			err := json.Unmarshal([]byte(t), &a)
			if err != nil {
				a = append(a, t)
			}
		}
		return a, nil
	case nil:
		return nil, nil
	default:
		s := reflect.ValueOf(val)
		if s.Kind() == reflect.Slice {
			a := make([]interface{}, s.Len())

			for i := 0; i < s.Len(); i++ {
				a[i] = s.Index(i).Interface()
			}
			return a, nil
		}
		if s.IsValid() {
			a := make([]interface{}, 1)
			a[0] = val

			return a, nil
		}
		return nil, fmt.Errorf("unable to coerce %#v to []interface{}", val)
	}
}

// ToArrayIfNecessary coerce a value to an array if it isn't one already
func ToArrayIfNecessary(val interface{}) (interface{}, error) {

	if val == nil {
		return nil, nil
	}

	rt := reflect.TypeOf(val).Kind()

	if rt == reflect.Array || rt == reflect.Slice {
		return val, nil
	}

	switch t := val.(type) {
	case string:
		a := make([]interface{}, 0)
		if t != "" {
			err := json.Unmarshal([]byte(t), &a)
			if err != nil {
				return t, nil
			}
		}
		return a, nil
	default:
		return nil, fmt.Errorf("unable to coerce %#v to []interface{}", val)
	}
}

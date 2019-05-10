package path

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/project-flogo/core/data/coerce"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

//todo consolidate and optimize code

func GetValue(value interface{}, path string) (interface{}, error) {

	if path == "" {
		return value, nil
	}

	var newVal interface{}
	var err error
	var newPath string

	//To interface if it is an string
	if val, ok := value.(string); ok {
		var in interface{}
		err = json.Unmarshal([]byte(val), &in)
		if err != nil {
			return nil, err
		}
		value = in
	}

	if strings.HasPrefix(path, ".") {
		if objVal, ok := value.(map[string]interface{}); ok {
			newVal, newPath, err = getSetObjValue(objVal, path, nil, false)
		} else if paramsVal, ok := value.(map[string]string); ok {
			newVal, newPath, err = getSetParamsValue(paramsVal, path, nil, false)
		} else {
			fieldName, npIdx := getObjectKey(path[1:])
			newVal, err = getFieldValueByName(value, fieldName)
			if err != nil {
				return nil, err
			}
			newPath = path[npIdx:]
		}
	} else if hasMapKey(path) {
		if objVal, ok := value.(map[string]interface{}); ok {
			newVal, newPath, err = getSetMapValue(objVal, path, nil, false)
		} else if paramsVal, ok := value.(map[string]string); ok {
			newVal, newPath, err = getSetMapParamsValue(paramsVal, path, nil, false)
		} else {
			return nil, fmt.Errorf("unable to evaluate path: %s", path)
		}
	} else if strings.HasPrefix(path, "[") {
		newVal, newPath, err = getSetArrayValue(value, path, nil, false)
	} else {
		return nil, fmt.Errorf("unable to evaluate path: %s", path)
	}

	if err != nil {
		return nil, err
	}
	return GetValue(newVal, newPath)
}

func getFieldValueByName(object interface{}, name string) (interface{}, error) {
	val := reflect.ValueOf(object)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {

		field := val.FieldByName(NormalizeFieldName(name))
		if field.IsValid() {
			return field.Interface(), nil
		}

		typ := reflect.TypeOf(object)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		for i := 0; i < typ.NumField(); i++ {
			p := typ.Field(i)
			if !p.Anonymous {
				if p.Tag != "" && len(p.Tag) > 0 {
					if name == p.Tag.Get("json") {
						return val.FieldByName(typ.Field(i).Name).Interface(), nil
					}
				}
			}
		}

	} else if val.Kind() == reflect.Map {
		v := val.MapIndex(reflect.ValueOf(name))
		return v.Interface(), nil
	}
	return nil, fmt.Errorf("unable to evaluate path: %s", name)
}

func NormalizeFieldName(name string) string {
	symbols := []rune(name)
	symbols[0] = unicode.ToUpper(symbols[0])
	return string(symbols)
}

func SetValue(attrValue interface{}, path string, value interface{}) error {
	if path == "" || attrValue == nil {
		return nil
	}

	var newVal interface{}
	var err error
	var newPath string

	if strings.HasPrefix(path, ".") {

		if objVal, ok := attrValue.(map[string]interface{}); ok {
			newVal, newPath, err = getSetObjValue(objVal, path, value, true)
		} else if paramsVal, ok := attrValue.(map[string]string); ok {
			newVal, newPath, err = getSetParamsValue(paramsVal, path, value, true)
		} else {
			return fmt.Errorf("unable to evaluate path: %s", path)
		}
	} else if hasMapKey(path) {
		if objVal, ok := attrValue.(map[string]interface{}); ok {
			newVal, newPath, err = getSetMapValue(objVal, path, value, true)
		} else if paramsVal, ok := attrValue.(map[string]string); ok {
			newVal, newPath, err = getSetMapParamsValue(paramsVal, path, value, true)
		} else {
			return fmt.Errorf("unable to evaluate path: %s", path)
		}

	} else if strings.HasPrefix(path, "[") {
		newVal, newPath, err = getSetArrayValue(attrValue, path, value, true)
	} else {
		return fmt.Errorf("unable to evaluate path: %s", path)
	}

	if err != nil {
		return err
	}
	return SetValue(newVal, newPath, value)
}

func hasMapKey(path string) bool {
	return strings.HasPrefix(path, `["`) || strings.HasPrefix(path, `['`)
}

func equalMapKey(val string) bool {
	return val == `["` || val == `['`
}

func getObjectKey(s string) (string, int) {
	i := 0

	for i < len(s) {

		if s[i] == '.' || s[i] == '[' {
			return s[:i], i + 1
		}

		i += 1
	}

	return s, len(s) + 1
}

func getMapKey(s string) (string, int) {
	i := 0

	for i < len(s) {

		if s[i] == '"' || s[i] == '\'' {
			return s[:i], i + 4 // [" "]
		}

		i += 1
	}

	return s, len(s) + 1
}

func getSetArrayValue(obj interface{}, path string, value interface{}, set bool) (interface{}, string, error) {

	arrValue, valid := obj.([]interface{})
	if !valid {
		//Try to convert to a array in case it is a array string
		val, err := coerce.ToArray(obj)
		if err != nil {
			return nil, path, errors.New("'" + path + "' not an array")
		}
		arrValue = val
	}

	closeIdx := strings.Index(path, "]")

	if closeIdx == -1 {
		return nil, path, errors.New("'" + path + "' not an array")
	}

	arrayIdx, err := strconv.Atoi(path[1:closeIdx])
	if err != nil {
		return nil, path, errors.New("Invalid array index: " + path[1:closeIdx])
	}

	if arrayIdx >= len(arrValue) {
		return nil, path, errors.New("Array index '" + path + "' out of range.")
	}

	if set && closeIdx == len(path)-1 {
		arrValue[arrayIdx] = value
		return nil, "", nil
	}

	return arrValue[arrayIdx], path[closeIdx+1:], nil
}

func getSetObjValue(objValue map[string]interface{}, path string, value interface{}, set bool) (interface{}, string, error) {

	key, npIdx := getObjectKey(path[1:])
	if set && key == path[1:] {
		//end of path so set the value
		objValue[key] = value
		return nil, "", nil
	}

	val, found := objValue[key]
	if !found {
		if path == "."+key {
			return nil, "", nil
		}
		return nil, "", errors.New("Invalid path '" + path + "'. path not found.")
	}

	return val, path[npIdx:], nil
}

func getSetParamsValue(params map[string]string, path string, value interface{}, set bool) (interface{}, string, error) {

	key, _ := getObjectKey(path[1:])
	if set && key == path[1:] {
		//end of path so set the value
		paramVal, err := coerce.ToString(value)

		if err != nil {
			return nil, "", err
		}
		params[key] = paramVal
		return nil, "", nil
	}

	val, found := params[key]

	if !found {
		return "", "", nil
	}

	return val, "", nil
}

func getSetMapValue(objValue map[string]interface{}, path string, value interface{}, set bool) (interface{}, string, error) {

	key, npIdx := getMapKey(path[2:])

	if set && (key+`"]` == path[2:] || key+`']` == path[2:]) {
		//end of path so set the value
		objValue[key] = value
		return nil, "", nil
	}

	val, found := objValue[key]

	if !found {
		if path == "."+key {
			return nil, "", nil
		}
		return nil, "", errors.New("Invalid path '" + path + "'. path not found.")
	}

	return val, path[npIdx:], nil
}

func getSetMapParamsValue(params map[string]string, path string, value interface{}, set bool) (interface{}, string, error) {

	key, _ := getMapKey(path[2:])
	if set && key+`"]` == path[2:] {
		//end of path so set the value
		paramVal, err := coerce.ToString(value)

		if err != nil {
			return nil, "", err
		}
		params[key] = paramVal
		return nil, "", nil
	}

	val, found := params[key]

	if !found {
		return "", "", nil
	}

	return val, "", nil
}

func Deconstruct(fullPath string) (attrName string, path string, err error) {

	idx := strings.IndexFunc(fullPath, isSep)

	if idx == -1 {
		return fullPath, "", nil
	}

	return fullPath[:idx], fullPath[idx:], nil
}

func isSep(r rune) bool {
	return r == '.' || r == '['
}

package path

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/oliveagle/jsonpath"
	"github.com/project-flogo/core/data/coerce"
)

//todo consolidate and optimize code

func GetValue(value interface{}, path string) (interface{}, error) {

	if path == "" {
		return value, nil
	}

	var newVal interface{}
	var err error
	var newPath string

	if strings.HasPrefix(path, ".") {
		if objVal, ok := value.(map[string]interface{}); ok {
			newVal, newPath, err = getSetObjValue(objVal, path, nil, false)
		} else if paramsVal, ok := value.(map[string]string); ok {
			newVal, newPath, err = getSetParamsValue(paramsVal, path, nil, false)
			//} else if objVal, ok := value.(*ComplexObject); ok {
			//	return PathGetValue(objVal.Value, path)
		} else {

			val := reflect.ValueOf(value)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}

			if val.Kind() == reflect.Struct {
				fieldName, npIdx := getObjectKey(path[1:])
				newPath = path[npIdx:]
				f := val.FieldByName(fieldName)
				if f.IsValid() {
					return f.Interface(), nil
				}

				return nil, nil
			} else {
				return nil, fmt.Errorf("unable to evaluate path: %s", path)
			}
		}
	} else if strings.HasPrefix(path, `["`) {
		if objVal, ok := value.(map[string]interface{}); ok {
			newVal, newPath, err = getSetMapValue(objVal, path, nil, false)
		} else if paramsVal, ok := value.(map[string]string); ok {
			newVal, newPath, err = getSetMapParamsValue(paramsVal, path, nil, false)
			//} else if objVal, ok := value.(*ComplexObject); ok {
			//	return PathGetValue(objVal.Value, path)
		} else {
			return nil, fmt.Errorf("unable to evaluate path: %s", path)
		}
	} else if strings.HasPrefix(path, "[`") {
		jpath := strings.TrimSuffix(strings.TrimPrefix(path, "[`"), "`]")
		newVal, err = jsonpath.JsonPathLookup(value, jpath)
	} else if strings.HasPrefix(path, "[") {
		//if objVal, ok := value.(*ComplexObject); ok {
		//	newVal, newPath, err = getSetArrayValue(objVal.Value, path, nil, false)
		//} else {
		newVal, newPath, err = getSetArrayValue(value, path, nil, false)

		//}
	} else {
		return nil, fmt.Errorf("unable to evaluate path: %s", path)
	}

	if err != nil {
		return nil, err
	}
	return GetValue(newVal, newPath)
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
			//} else if objVal, ok := value.(*ComplexObject); ok {
			//	return PathSetValue(objVal.Value, path, value)
		} else {
			return fmt.Errorf("Unable to evaluate path: %s", path)
		}
	} else if strings.HasPrefix(path, `["`) {
		if objVal, ok := attrValue.(map[string]interface{}); ok {
			newVal, newPath, err = getSetMapValue(objVal, path, value, true)
		} else if paramsVal, ok := attrValue.(map[string]string); ok {
			newVal, newPath, err = getSetMapParamsValue(paramsVal, path, value, true)
			//} else if objVal, ok := value.(*ComplexObject); ok {
			//	return PathSetValue(objVal.Value, path, value)
		} else {
			return fmt.Errorf("unable to evaluate path: %s", path)
		}

	} else if strings.HasPrefix(path, "[") {
		//if objVal, ok := value.(*ComplexObject); ok {
		//	newVal, newPath, err = getSetArrayValue(attrValue, path, objVal.Value, true)
		//} else {
		newVal, newPath, err = getSetArrayValue(attrValue, path, value, true)
		//}
	} else {
		return fmt.Errorf("Unable to evaluate path: %s", path)
	}

	if err != nil {
		return err
	}
	return SetValue(newVal, newPath, value)
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

		if s[i] == '"' {
			return s[:i], i + 4 // [" "]
		}

		i += 1
	}

	return s, len(s) + 1
}

func getSetArrayValue(obj interface{}, path string, value interface{}, set bool) (interface{}, string, error) {

	arrValue, valid := obj.([]interface{})
	if !valid {
		//Try to convert to a array incase it is a array string
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
		return nil, "", errors.New("Invalid path '" + path + "'. path not found.")
	}

	return val, "", nil
}

func getSetMapValue(objValue map[string]interface{}, path string, value interface{}, set bool) (interface{}, string, error) {

	key, npIdx := getMapKey(path[2:])

	if set && key+`"]` == path[2:] {
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
		return nil, "", errors.New("Invalid path '" + path + "'. path not found.")
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

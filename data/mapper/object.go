package mapper

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/path"
	"github.com/project-flogo/core/support/log"
	"reflect"
	"runtime/debug"
	"strings"
)

var objectMapperLog = log.ChildLogger(log.RootLogger(), "object-mapper")

const (
	//To do an array mapping from upstreaming, use @foreach.
	/*
			"input": {

		    "val" : {
		        "a" : "=$activity[blah].out",
		        "addresses": {
		            "@foreach($activity[blah].out2, i)":
		            {
		              "street"  : "=$.street",
		              "zipcode" : "9999",
		              "state"   : "=$activity[test].state",
		              "addresses2": {
		                  "@foreach($.addresses2)":{
		                        "street2"  : "=$loop[i].street2",
		                        "zipcode2" : "=$.zipcode",
		                        "state2"   : "=$activity[test].state"
		                	  }
		               	}
		            	}
		        	}
		    	}
			}

	*/
	FOREACH = "@foreach"
)

type ObjectMapping struct {
	Mapping interface{} `json:"mapping"`
}

func GetObjectMapping(value interface{}) (interface{}, bool) {
	switch t := value.(type) {
	case *ObjectMapping:
		return t.Mapping, true
	case map[string]interface{}:

		if mapping, ok := t["mapping"]; ok {
			return mapping, true
		}
		return nil, false
	default:
		return nil, false
	}
}

type ObjectMapperFactory struct {
	exprFactory expression.Factory
}

func NewObjectMapperFactory(exprFactory expression.Factory) expression.Factory {
	return &ObjectMapperFactory{exprFactory: exprFactory}
}

func (am *ObjectMapperFactory) NewObjectMapper(value interface{}) (expression.Expr, error) {
	return &ObjectMapper{mappings: value, exprFactory: am.exprFactory}, nil
}

func (am *ObjectMapperFactory) NewExpr(value string) (expression.Expr, error) {
	v, err := coerce.ToObject(value)
	if err != nil {
		return nil, fmt.Errorf("unexpected object mapping error, %s", err.Error())
	}
	return am.NewObjectMapper(v)
}

type ObjectMapper struct {
	mappings    interface{}
	exprFactory expression.Factory
}

func (am *ObjectMapper) Eval(inputScope data.Scope) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			objectMapperLog.Error("%+v", r)
			objectMapperLog.Debugf("StackTrace: %s", debug.Stack())
		}
	}()
	objectMapperLog.Debugf("Handling object mapper %+v", am.mappings)
	return handleObjectMapping(am.mappings, am.exprFactory, inputScope)
}

type foreach struct {
	sourceFrom  string
	index       string
	filterExpr  expression.Expr
	exprFactory expression.Factory
}

func handleObjectMapping(objectMappings interface{}, exprF expression.Factory, inputScope data.Scope) (interface{}, error) {
	var err error
	switch t := objectMappings.(type) {
	case map[string]interface{}:
		objectVal := make(map[string]interface{})
		for mk, mv := range t {
			if strings.HasPrefix(mk, FOREACH) {
				foreach, err := newForeach(mk, exprF)
				if err != nil {
					return nil, err
				}
				return foreach.handle(mv.(map[string]interface{}), inputScope)
			}
			//Go second level to find possible @foreach node
			if nxtVal, ok := mv.(map[string]interface{}); ok {
				//second level object
				if hasForeach(nxtVal) {
					for foreachK, foreachV := range nxtVal {
						if strings.HasPrefix(foreachK, FOREACH) {
							foreach, err := newForeach(foreachK, exprF)
							if err != nil {
								return nil, err
							}
							objectVal[mk], err = foreach.handle(foreachV.(map[string]interface{}), inputScope)
							if err != nil {
								return nil, err
							}
						} else {
							switch t := foreachV.(type) {
							case []interface{}:
								arrayResult := make([]interface{}, len(t))
								for i, element := range t {
									var err error
									arrayResult[i], err = handlerArrayElement(element, exprF, inputScope)
									if err != nil {
										return nil, err
									}
								}
								objectVal[foreachK] = arrayResult
							case map[string]interface{}:
								if objectVal[mk] == nil {
									objectVal[mk] = make(map[string]interface{})
								}
								v, err := handleObjectMapping(t, exprF, inputScope)
								if err != nil {
									return nil, err
								}
								objectVal[mk].(map[string]interface{})[foreachK] = v
							default:
								if objectVal[mk] == nil {
									objectVal[mk] = make(map[string]interface{})
								}
								err := handleObject(objectVal[mk].(map[string]interface{}), foreachK, foreachV, exprF, inputScope)
								if err != nil {
									return nil, err
								}
							}
						}
					}
				} else {
					objectVal[mk], err = handleObjectMapping(mv, exprF, inputScope)
					if err != nil {
						return nil, err
					}
				}

			} else if arrayV, ok := mv.([]interface{}); ok {
				arrayResult := make([]interface{}, len(arrayV))
				for i, element := range arrayV {
					var err error
					arrayResult[i], err = handlerArrayElement(element, exprF, inputScope)
					if err != nil {
						return nil, err
					}
				}
				objectVal[mk] = arrayResult
			} else {
				err := handleObject(objectVal, mk, mv, exprF, inputScope)
				if err != nil {
					return nil, err
				}
			}
		}
		return objectVal, nil
	case []interface{}:
		//array with possible child object
		objArray := make([]interface{}, len(t))
		for i, element := range t {
			var err error
			objArray[i], err = handlerArrayElement(element, exprF, inputScope)
			if err != nil {
				return nil, err
			}
		}
		return objArray, nil
	default:
		return nil, fmt.Errorf("unsupport type [%s] for object mapper", reflect.TypeOf(objectMappings))
	}
}

func hasForeach(val map[string]interface{}) bool {
	for foreachK, _ := range val {
		if strings.HasPrefix(foreachK, FOREACH) {
			return true
		}
	}
	return false
}

func handlerArrayElement(element interface{}, exprF expression.Factory, inputScope data.Scope) (interface{}, error) {
	//Only handle object mapping
	switch element.(type) {
	case map[string]interface{}:
		return handleObjectMapping(element, exprF, inputScope)
	default:
		return getExpressionValue(element, exprF, inputScope)
	}
}

func (f *foreach) handle(arrayMappingFields map[string]interface{}, inputScope data.Scope) (interface{}, error) {
	fromValue, err := getExpressionValue("="+f.sourceFrom, f.exprFactory, inputScope)
	if err != nil {
		return nil, fmt.Errorf("foreach eval source array error, %s", err.Error())
	}
	newSourceArray, err := coerce.ToArray(fromValue)
	if err != nil {
		objectMapperLog.Errorf("foreach source [%+v] not an array, cast to array error %s", fromValue, err.Error())
		return nil, fmt.Errorf("foreach source [%+v] not an array", fromValue)
	}

	var targetValues []interface{}
	if hasArrayAssign(arrayMappingFields) {
		targetValues, err = f.handleArrayAssign(newSourceArray, arrayMappingFields, inputScope)
		if err != nil {
			return nil, fmt.Errorf("array assign error, %s", err.Error())
		}
	}

	arrayMappingFields = removeAssignFromArrayMappingFeild(arrayMappingFields)

	if len(arrayMappingFields) > 0 {
		requireUpdate := len(targetValues) > 0
		var skippedCount = 0
		for i, sourceValue := range newSourceArray {
			inputScope, err = newLoopScope(sourceValue, f.index, inputScope)
			if err != nil {
				return nil, err
			}
			passedFilter, err := f.Filter(inputScope)
			if err != nil {
				return nil, err
			}
			if passedFilter {
				item, err := handleObjectMapping(arrayMappingFields, f.exprFactory, inputScope)
				if err != nil {
					return nil, err
				}
				if requireUpdate {
					targetValueIndex := i - skippedCount
					if len(targetValues) <= 0 {
						targetValues = append(targetValues, item)
					} else if (len(targetValues) < targetValueIndex) || (targetValues[targetValueIndex] == nil) {
						// No target value, just append
						targetValues = append(targetValues, item)
					} else {
						//update value
						switch t := item.(type) {
						case map[string]interface{}:
							targetValue, err := ToObjectMap(targetValues[targetValueIndex])
							if err == nil {
								for k, v := range t {
									targetValue[k] = v
								}
							} else {
								return nil, fmt.Errorf("cannot assign map[string]interface to [%s]", reflect.TypeOf(targetValues[targetValueIndex]))
							}
						case []interface{}:
							targetValue, err := coerce.ToArray(targetValues[targetValueIndex])
							if err == nil {
								for k, v := range t {
									targetValue[k] = v
								}
							} else {
								return nil, fmt.Errorf("cannot assign []interface to [%s]", reflect.TypeOf(targetValues[targetValueIndex]))
							}
						}
					}
				} else {
					// No updated required, just go over the array mapping fields
					targetValues = append(targetValues, item)
				}
			} else {
				skippedCount++
			}
		}
		return targetValues, nil
	}

	return targetValues, nil
}

func (f *foreach) Filter(inputScope data.Scope) (bool, error) {
	if f.filterExpr != nil {

		v, err := f.filterExpr.Eval(inputScope)
		if err != nil {
			return false, fmt.Errorf("execute expression [%+v] error %s", f.filterExpr, err.Error())
		}
		return coerce.ToBool(v)
	}
	return true, nil
}

func removeAssignFromArrayMappingFeild(arrayMappingFields map[string]interface{}) map[string]interface{} {
	tmpArrayField := make(map[string]interface{})
	for k, v := range arrayMappingFields {
		if k != "=" {
			tmpArrayField[k] = v
		}
	}

	return tmpArrayField
}

func hasArrayAssign(arrayMappingFields map[string]interface{}) bool {
	field, ok := arrayMappingFields["="]
	if ok && field != nil {
		return true
	}
	return false
}

func (f *foreach) handleArrayAssign(sourceArray []interface{}, arrayMappingFields map[string]interface{}, inputScope data.Scope) ([]interface{}, error) {
	var targetValues []interface{}
	field, ok := arrayMappingFields["="]
	if ok && field != nil {
		if v, ok := field.(string); ok && v == "$loop" {
			for _, sourceValue := range sourceArray {
				var err error
				inputScope, err = newLoopScope(sourceValue, f.index, inputScope)
				if err != nil {
					return nil, err
				}
				passFilter, err := f.Filter(inputScope)
				if err != nil {
					return nil, err
				}
				if passFilter {
					targetValues = append(targetValues, sourceValue)
				}
			}
		} else {
			for _, sourceValue := range sourceArray {
				var err error
				inputScope, err = newLoopScope(sourceValue, f.index, inputScope)
				if err != nil {
					return nil, err
				}

				passFilter, err := f.Filter(inputScope)
				if err != nil {
					return nil, err
				}

				if passFilter {
					fromValue, err := getExpressionValue(field, f.exprFactory, inputScope)
					if err != nil {
						return nil, fmt.Errorf("eval expression failed %s", err.Error())
					}
					targetValues = append(targetValues, fromValue)
				}

			}
		}
	}
	return targetValues, nil
}

func handleObject(targetValue map[string]interface{}, targetName string, source interface{}, exprF expression.Factory, scope data.Scope) error {
	if isExpr(source) {
		fromValue, err := getExpressionValue(source, exprF, scope)
		if err != nil {
			return fmt.Errorf("eval expression failed %s", err.Error())
		}
		targetValue[targetName] = fromValue
	} else {
		targetValue[targetName] = source
	}
	return nil
}

func getExpressionValue(path interface{}, exprF expression.Factory, scope data.Scope) (interface{}, error) {
	var finalExpr expression.Expr
	if isExpr(path) {
		var err error
		finalExpr, err = exprF.NewExpr(path.(string)[1:])
		if err != nil {
			return nil, err
		}
	} else {
		finalExpr = expression.NewLiteralExpr(path)
	}
	return finalExpr.Eval(scope)
}

func newForeach(foreachpath string, exprF expression.Factory) (*foreach, error) {
	foreach := &foreach{exprFactory: exprF}
	foreachpath = strings.TrimSpace(foreachpath)
	if strings.HasPrefix(foreachpath, FOREACH) && strings.Contains(foreachpath, "(") && strings.Contains(foreachpath, ")") {
		paramsStr := foreachpath[strings.Index(foreachpath, "(")+1 : strings.LastIndex(foreachpath, ")")]
		sourceIdx := strings.Index(paramsStr, ",")
		if sourceIdx <= 0 {
			foreach.sourceFrom = strings.TrimSpace(paramsStr)
		} else {

			foreach.sourceFrom = strings.TrimSpace(paramsStr[:sourceIdx])
			if len(paramsStr) > sourceIdx+1 {
				//No more argument
				afterLoopNameParamStr := strings.TrimSpace(paramsStr[sourceIdx+1:])
				loopNameIdx := strings.Index(afterLoopNameParamStr, ",")
				if loopNameIdx >= 0 {
					foreach.index = strings.TrimSpace(afterLoopNameParamStr[:loopNameIdx])
				} else {
					foreach.index = afterLoopNameParamStr
					return foreach, nil
				}

				if len(afterLoopNameParamStr) > loopNameIdx+1 {
					filter := strings.TrimSpace(afterLoopNameParamStr[loopNameIdx+1:])
					if len(filter) > 0 {
						//create new filter expression
						filterExpr, err := exprF.NewExpr(filter)
						if err != nil {
							return nil, fmt.Errorf("create foreach filtering expression error: %s", err.Error())
						}
						foreach.filterExpr = filterExpr
					}
				}
			}
		}
	}
	return foreach, nil
}

func newLoopScope(arrayItem interface{}, indexName string, scope data.Scope) (data.Scope, error) {
	mapData, err := ToObjectMap(arrayItem)
	if err != nil {
		return nil, fmt.Errorf("convert array item data [%+v] to map failed, due to [%s]", arrayItem, err.Error())
	}

	loopData := make(map[string]interface{})
	loopData["_loop"] = mapData
	if len(indexName) > 0 {
		loopData[indexName] = mapData
	}

	return data.NewSimpleScope(loopData, scope), nil
}

func ToObjectMap(value interface{}) (map[string]interface{}, error) {
	switch t := value.(type) {
	case map[string]interface{}:
		return t, nil
	case map[string]string, string:
		return coerce.ToObject(value)
	default:
		out := make(map[string]interface{})
		v := reflect.ValueOf(t)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Map {
			for _, k := range v.MapKeys() {
				key, err := coerce.ToString(k.Interface())
				if err != nil {
					return nil, fmt.Errorf("unable to convert key [%+v] to string: %s", k.Interface(), err.Error())
				}
				out[key] = v.MapIndex(k).Interface()
			}
		} else if v.Kind() == reflect.Struct {
			typ := v.Type()
			for i := 0; i < v.NumField(); i++ {
				// gets us a StructField
				fi := typ.Field(i)
				if !fi.Anonymous {
					jsonTag := path.GetJsonTag(fi)
					if len(jsonTag) > 0 {
						out[jsonTag] = v.Field(i).Interface()
					} else {
						out[fi.Name] = v.Field(i).Interface()
					}
				}
			}
		} else {
			return coerce.ToObject(t)
		}

		return out, nil
	}
}

package mapper

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression"
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
	exprFactory expression.Factory
}

func handleObjectMapping(objectMappings interface{}, exprF expression.Factory, inputScope data.Scope) (interface{}, error) {
	var err error
	switch t := objectMappings.(type) {
	case map[string]interface{}:
		objectVal := make(map[string]interface{})
		for mk, mv := range t {
			if strings.HasPrefix(mk, FOREACH) {
				return newForeach(mk, exprF).handle(mv.(map[string]interface{}), inputScope)
			}
			//Go second level to find possible @foreach node
			if nxtVal, ok := mv.(map[string]interface{}); ok {
				//second level object
				if hasForeach(nxtVal) {
					for foreachK, foreachV := range nxtVal {
						if strings.HasPrefix(foreachK, FOREACH) {
							objectVal[mk], err = newForeach(foreachK, exprF).handle(foreachV.(map[string]interface{}), inputScope)
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

	targetValues := make([]interface{}, len(newSourceArray))
	if hasArrayAssign(arrayMappingFields) {
		targetValues, err = f.handleArrayAssign(newSourceArray, arrayMappingFields, inputScope)
		if err != nil {
			return nil, fmt.Errorf("array assign error, %s", err.Error())
		}
	}

	arrayMappingFields = removeAssignFromArrayMappingFeild(arrayMappingFields)

	if len(arrayMappingFields) > 0 {
		for i, sourceValue := range newSourceArray {
			inputScope = newLoopScope(sourceValue, f.index, inputScope)
			item, err := handleObjectMapping(arrayMappingFields, f.exprFactory, inputScope)
			if err != nil {
				return nil, err
			}
			if targetValues[i] == nil {
				targetValues[i] = item
			} else {
				//update value
				switch t := item.(type) {
				case map[string]interface{}:
					targetValue, ok := targetValues[i].(map[string]interface{})
					if ok {
						for k, v := range t {
							targetValue[k] = v
						}
					} else {
						return nil, fmt.Errorf("cannot assign map[string]interface to [%s]", reflect.TypeOf(targetValues[i]))
					}
				case []interface{}:
					targetValue, ok := targetValues[i].([]interface{})
					if ok {
						for k, v := range t {
							targetValue[k] = v
						}
					} else {
						return nil, fmt.Errorf("cannot assign []interface to [%s]", reflect.TypeOf(targetValues[i]))
					}
				}
			}
		}
		return targetValues, nil
	}

	return targetValues, nil
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
	targetValues := make([]interface{}, len(sourceArray))
	field, ok := arrayMappingFields["="]
	if ok && field != nil {
		if v, ok := field.(string); ok && v == "$loop" {
			targetValues = sourceArray
		} else {
			for i, sourceValue := range sourceArray {
				inputScope = newLoopScope(sourceValue, f.index, inputScope)
				fromValue, err := getExpressionValue(field, f.exprFactory, inputScope)
				if err != nil {
					return nil, fmt.Errorf("eval expression failed %s", err.Error())
				}
				targetValues[i] = fromValue
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

func newForeach(foreachpath string, exprF expression.Factory) *foreach {
	foreach := &foreach{exprFactory: exprF}
	foreachpath = strings.TrimSpace(foreachpath)
	if strings.HasPrefix(foreachpath, FOREACH) && strings.Contains(foreachpath, "(") && strings.Contains(foreachpath, ")") {
		paramsStr := foreachpath[strings.Index(foreachpath, "(")+1 : strings.Index(foreachpath, ")")]
		params := strings.Split(paramsStr, ",")
		if len(params) >= 2 {
			foreach.sourceFrom = strings.TrimSpace(params[0])
			foreach.index = strings.TrimSpace(params[1])
		} else {
			foreach.sourceFrom = strings.TrimSpace(params[0])
		}
	}
	return foreach
}

func newLoopScope(arrayItem interface{}, indexName string, scope data.Scope) data.Scope {
	if len(indexName) <= 0 {
		return data.NewSimpleScope(arrayItem.(map[string]interface{}), scope)
	} else {
		values := arrayItem.(map[string]interface{})
		values[indexName] = arrayItem
		return data.NewSimpleScope(values, scope)
	}
}

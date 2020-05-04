package mapper

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/path"
	"github.com/project-flogo/core/support/log"
	"reflect"
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
	return NewObjectMapper(value, am.exprFactory)
}

func (am *ObjectMapperFactory) NewExpr(value string) (expression.Expr, error) {
	v, err := coerce.ToObject(value)
	if err != nil {
		return nil, fmt.Errorf("unexpected object mapping error, %s", err.Error())
	}
	return am.NewObjectMapper(v)
}

type ObjectMapper struct {
	//Object
	objectFields map[string]expression.Expr
	//For Array mapping
	foreach *foreachExpr
	//For literal array mapping
	literalArray []expression.Expr
}

type foreachExpr struct {
	// source array
	sourceFrom expression.Expr
	// array scope name
	scopeName string
	// filter expression
	filterExpr expression.Expr
	// fields
	fields map[string]expression.Expr
	//Use to assign value
	assign expression.Expr
}

// assignAllExpr uses to indicate it is assign all.
type assignAllExpr struct {
}

func (a *assignAllExpr) Eval(scope data.Scope) (interface{}, error) {
	return nil, nil
}

func NewObjectMapper(mappings interface{}, exprF expression.Factory) (expr expression.Expr, err error) {
	switch t := mappings.(type) {
	case map[string]interface{}:
		objFields := make(map[string]expression.Expr)
		for mk, mv := range t {
			//Root Level foreach
			if strings.HasPrefix(mk, FOREACH) {
				foreach, err := newForeachExpr(mk, exprF)
				if err != nil {
					return nil, err
				}
				foreach.addFields(mv.(map[string]interface{}), exprF)
				return foreach, nil
			} else {
				objFields[mk], err = NewObjectMapper(mv, exprF)
			}

		}
		return &ObjectMapper{
			objectFields: objFields,
		}, nil
	case []interface{}:
		//array with possible child object
		objArray := make([]expression.Expr, len(t))
		for i, element := range t {
			var err error
			objArray[i], err = NewObjectMapper(element, exprF)
			if err != nil {
				return nil, err
			}
		}
		return &ObjectMapper{
			literalArray: objArray,
		}, nil
	case interface{}:
		return newExpr(t, exprF)
	default:
		return nil, fmt.Errorf("unsupport type [%s] for object mapper", reflect.TypeOf(t))
	}
}

func (f *foreachExpr) addFields(fields map[string]interface{}, exprF expression.Factory) (err error) {
	for key, value := range fields {
		if key == "=" {
			if value == "$loop" {
				f.assign = &assignAllExpr{}
			} else {
				f.assign, err = newExpr(value, exprF)
				if err != nil {
					return err
				}
			}
		} else {
			if f.fields == nil {
				f.fields = make(map[string]expression.Expr)
			}
			f.fields[key], err = NewObjectMapper(value, exprF)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func newExpr(path interface{}, exprF expression.Factory) (expression.Expr, error) {
	if isExpr(path) {
		return exprF.NewExpr(path.(string)[1:])
	} else {
		return expression.NewLiteralExpr(path), nil
	}
}

func newForeachExpr(foreachpath string, exprF expression.Factory) (*foreachExpr, error) {
	foreach := &foreachExpr{}
	foreachpath = strings.TrimSpace(foreachpath)
	if strings.HasPrefix(foreachpath, FOREACH) && strings.Contains(foreachpath, "(") && strings.Contains(foreachpath, ")") {
		paramsStr := foreachpath[strings.Index(foreachpath, "(")+1 : strings.LastIndex(foreachpath, ")")]
		sourceIdx := strings.Index(paramsStr, ",")
		if sourceIdx <= 0 {
			finalExpr, err := exprF.NewExpr(strings.TrimSpace(paramsStr))
			if err != nil {
				return nil, err
			}
			foreach.sourceFrom = finalExpr
		} else {

			finalExpr, err := exprF.NewExpr(strings.TrimSpace(paramsStr[:sourceIdx]))
			if err != nil {
				return nil, err
			}
			foreach.sourceFrom = finalExpr

			if len(paramsStr) > sourceIdx+1 {

				afterLoopNameParamStr := strings.TrimSpace(paramsStr[sourceIdx+1:])
				loopNameIdx := strings.Index(afterLoopNameParamStr, ",")
				if loopNameIdx >= 0 {
					foreach.scopeName = strings.TrimSpace(afterLoopNameParamStr[:loopNameIdx])
				} else {
					foreach.scopeName = afterLoopNameParamStr
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

func (obj *ObjectMapper) Eval(scope data.Scope) (value interface{}, err error) {
	if obj.foreach != nil {
		return obj.foreach.Eval(scope)
	} else if obj.literalArray != nil {
		var array []interface{}
		for _, v := range obj.literalArray {
			arrValue, err := v.Eval(scope)
			if err != nil {
				return err, nil
			}
			array = append(array, arrValue)
		}
		return array, nil
	} else {
		arrValue := make(map[string]interface{})
		for k, v := range obj.objectFields {
			arrValue[k], err = v.Eval(scope)
			if err != nil {
				return nil, err
			}
		}
		return arrValue, nil
	}
}

func (f *foreachExpr) Eval(scope data.Scope) (interface{}, error) {
	sourceAr, err := f.sourceFrom.Eval(scope)
	if err != nil {
		return nil, fmt.Errorf("foreach eval source array error, %s", err.Error())
	}

	newSourceArray, err := coerce.ToArray(sourceAr)
	if err != nil {
		return nil, fmt.Errorf("foreach source [%+v] not an array", f.sourceFrom)
	}

	var targetValues []interface{}
	if f.assign != nil {
		targetValues, err = f.handleAssign(newSourceArray, scope)
		if err != nil {
			return nil, fmt.Errorf("array assign error, %s", err.Error())
		}
	}

	if len(f.fields) > 0 {
		requireUpdate := len(targetValues) > 0
		var skippedCount = 0
		for i, sourceValue := range newSourceArray {
			scope, err = newLoopScope(sourceValue, f.scopeName, scope)
			if err != nil {
				return nil, err
			}
			passedFilter, err := f.Filter(scope)
			if err != nil {
				return nil, err
			}

			if passedFilter {
				item, err := f.HandleFields(scope)
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

func (f *foreachExpr) handleAssign(sourceArray []interface{}, inputScope data.Scope) ([]interface{}, error) {
	var targetValues []interface{}

	switch f.assign.(type) {
	case *assignAllExpr:
		for _, sourceValue := range sourceArray {
			if f.filterExpr != nil {
				var err error
				inputScope, err = newLoopScope(sourceValue, f.scopeName, inputScope)
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
			} else {
				targetValues = append(targetValues, sourceValue)
			}
		}
	default:
		for _, sourceValue := range sourceArray {
			var err error
			inputScope, err = newLoopScope(sourceValue, f.scopeName, inputScope)
			if err != nil {
				return nil, err
			}

			passFilter, err := f.Filter(inputScope)
			if err != nil {
				return nil, err
			}
			if passFilter {
				fromValue, err := f.assign.Eval(inputScope)
				if err != nil {
					return nil, fmt.Errorf("eval expression failed %s", err.Error())
				}
				targetValues = append(targetValues, fromValue)
			}

		}
	}
	return targetValues, nil
}

func (f *foreachExpr) Filter(inputScope data.Scope) (bool, error) {
	if f.filterExpr != nil {

		v, err := f.filterExpr.Eval(inputScope)
		if err != nil {
			return false, fmt.Errorf("execute expression [%+v] error %s", f.filterExpr, err.Error())
		}
		return coerce.ToBool(v)
	}
	return true, nil
}

func (f *foreachExpr) HandleFields(inputScope data.Scope) (interface{}, error) {
	vals := make(map[string]interface{})
	var err error
	for k, v := range f.fields {
		vals[k], err = v.Eval(inputScope)
		if err != nil {
			return nil, err
		}
	}
	return vals, nil
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

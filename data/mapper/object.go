package mapper

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/support/log"
	"runtime/debug"
	"strconv"
	"strings"
)

const (
	FOREACH  = "@foreach"
	NEWARRAY = "NEWARRAY"
)

type ObjectMapperFactory struct {
	exprFactory expression.Factory
}

func NewObjectMapperFactory(exprFactory expression.Factory) expression.Factory {
	return &ObjectMapperFactory{exprFactory: exprFactory}
}

func (am *ObjectMapperFactory) NewObjectMapper(value map[string]interface{}) (expression.Expr, error) {
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
	mappings    map[string]interface{}
	exprFactory expression.Factory
}

func (am *ObjectMapper) Eval(inputScope data.Scope) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			log.RootLogger().Error("%+v", r)
			log.RootLogger().Debugf("StackTrace: %s", debug.Stack())
		}
	}()

	return handleObjectMapping(am.mappings, am.exprFactory, inputScope)
}

func getForeach(foreachpath string) (*foreach, error) {
	foreachpath = strings.TrimSpace(foreachpath)
	if strings.HasPrefix(foreachpath, FOREACH) && strings.Contains(foreachpath, "(") && strings.Contains(foreachpath, ")") {
		paramsStr := foreachpath[strings.Index(foreachpath, "(")+1 : strings.Index(foreachpath, ")")]
		params := strings.Split(paramsStr, ",")
		if len(params) == 2 {
			return &foreach{sourceFrom: strings.TrimSpace(params[0]), index: strings.TrimSpace(params[1])}, nil
		} else {
			return &foreach{sourceFrom: strings.TrimSpace(params[0])}, nil
		}
	}
	return nil, fmt.Errorf("invalid foreach [%s], it must follow @foreach(sourceArray, optional)", foreachpath)
}

type foreach struct {
	sourceFrom string
	index      string
}

func handleObjectMapping(mappings map[string]interface{}, exprF expression.Factory, inputScope data.Scope) (interface{}, error) {
	var value interface{}
	if len(mappings) > 0 {
		for mk, mv := range mappings {
			//Array mapping should be only @foreach node
			if nxtVal, ok := mv.(map[string]interface{}); ok {
				for foreachK, foreachV := range nxtVal {
					if strings.HasPrefix(foreachK, FOREACH) {
						//Array mapping
						foreach, err := getForeach(foreachK)
						if err != nil {
							return nil, err
						}
						if foreach.sourceFrom == NEWARRAY {
							newArrayValues := make([]interface{}, 1)
							value = make([]interface{}, 1)
							err := handleArray(value.([]interface{}), newArrayValues, foreachV.(map[string]interface{}), exprF, newLoopScope(newArrayValues, foreach.index, inputScope))
							if err != nil {
								return nil, fmt.Errorf("new array mapping failed %s", err.Error())
							}

						} else {
							fromValue, err := getExpressionValue("="+foreach.sourceFrom, exprF, inputScope)
							if err != nil {
								return nil, fmt.Errorf("get value from source array error, %s", err.Error())
							}
							sourceValues, err := coerce.ToArray(fromValue)
							if err != nil {
								return nil, fmt.Errorf("source array not an array, error %s", err.Error())
							}
							value = make([]interface{}, len(sourceValues))
							err = handleArray(value.([]interface{}), sourceValues, foreachV.(map[string]interface{}), exprF, newLoopScope(sourceValues, foreach.index, inputScope))
							if err != nil {
								return nil, fmt.Errorf("array mapping failed %s", err.Error())
							}
						}

					} else {
						if value == nil {
							value = make(map[string]interface{})
						}
						err := handleObject(value, mk, mv, exprF, inputScope)
						if err != nil {
							return nil, err
						}
					}
				}
			} else {
				//Object mapping
				if value == nil {
					value = make(map[string]interface{})
				}
				err := handleObject(value, mk, mv, exprF, inputScope)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return value, nil
}

func handleObject(targetValue interface{}, targetName string, source interface{}, exprF expression.Factory, scope data.Scope) error {
	if objValue, ok := targetValue.(map[string]interface{}); ok {
		if isExpr(source) {
			fromValue, err := getExpressionValue(source, exprF, scope)
			if err != nil {
				return fmt.Errorf("get expression failed %s", err.Error())
			}
			objValue[targetName] = fromValue
		} else {
			objValue[targetName] = source
		}
	}
	return nil
}

func handleArray(targetValues []interface{}, sourceValues []interface{}, arrayMappingFields map[string]interface{}, exprF expression.Factory, inputScope data.Scope) error {
	for i, sourceValue := range sourceValues {
		for k, v := range arrayMappingFields {
			if strings.HasPrefix(k, FOREACH) {
				//Array mapping
				foreach, err := getForeach(k)
				if err != nil {
					return err
				}
				if foreach.sourceFrom == NEWARRAY {
					targetValues[i] = make([]interface{}, 1)
					err := handleArray(targetValues[i].([]interface{}), []interface{}{}, v.(map[string]interface{}), exprF, newLoopScope([]interface{}{}, foreach.index, inputScope))
					if err != nil {
						return fmt.Errorf("get expression failed %s", err.Error())
					}

				} else {
					fromValue, err := getExpressionValueWithData(sourceValue, "="+foreach.sourceFrom, exprF, inputScope)
					if err != nil {
						return fmt.Errorf("get value from source array error, %s", err.Error())
					}

					sourceValues, err := coerce.ToArray(fromValue)
					if err != nil {
						return fmt.Errorf("source array not an array, error %s", err.Error())
					}
					targetValues[i] = make([]interface{}, len(sourceValues))

					err = handleArray(targetValues[i].([]interface{}), sourceValues, v.(map[string]interface{}), exprF, newLoopScope(sourceValues, foreach.index, inputScope))
					if err != nil {
						return err
					}
				}

			} else {
				//Array fields
				if targetValues[i] == nil {
					targetValues[i] = make(map[string]interface{})
				}
				if objValue, ok := targetValues[i].(map[string]interface{}); ok {
					if isExpr(v) {
						fromValue, err := getExpressionValueWithData(sourceValue, v, exprF, inputScope)
						if err != nil {
							return fmt.Errorf("get expression failed %s", err.Error())
						}
						objValue[k] = fromValue
					} else {
						objValue[k] = v
					}
				}
			}
		}
	}
	return nil

}

func getExpressionValueWithData(value, path interface{}, exprF expression.Factory, scope data.Scope) (interface{}, error) {
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
	//Add current from value to scope
	if value != nil {
		scope = data.NewSimpleScope(value.(map[string]interface{}), scope)

	}
	return finalExpr.Eval(scope)
}

func getExpressionValue(path interface{}, exprF expression.Factory, scope data.Scope) (interface{}, error) {
	return getExpressionValueWithData(nil, path, exprF, scope)
}

func isExpr(value interface{}) bool {
	if strVal, ok := value.(string); ok && len(strVal) > 0 && (strVal[0] == '=' || strVal[0] == '$') {
		return true
	}
	return false
}

func newLoopScope(array []interface{}, indexName string, scope data.Scope) data.Scope {
	values := make(map[string]interface{})
	for i, v := range array {
		values["_L."+indexName+"."+strconv.Itoa(i)] = v
	}
	return data.NewSimpleScope(values, scope)
}

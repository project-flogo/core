package mapper

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/support/log"
	"runtime/debug"
	"strings"
)

const (
	//Foreach Handle array must start with @foreach.
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
	//NEWARRAY create one array element for target array
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

type foreach struct {
	sourceFrom  string
	index       string
	exprFactory expression.Factory
}

func handleObjectMapping(objectMappings map[string]interface{}, exprF expression.Factory, inputScope data.Scope) (interface{}, error) {
	objectVal := make(map[string]interface{})
	if len(objectMappings) > 0 {
		for mk, mv := range objectMappings {
			//Root level @foreach, it should be only itself
			if strings.HasPrefix(mk, FOREACH) {
				return newForeach(mk, exprF).handle(mv.(map[string]interface{}), inputScope)
			}

			if nxtVal, ok := mv.(map[string]interface{}); ok {
				//second level object
				for foreachK, foreachV := range nxtVal {
					if strings.HasPrefix(foreachK, FOREACH) {
						result, err := newForeach(foreachK, exprF).handle(foreachV.(map[string]interface{}), inputScope)
						if err != nil {
							return nil, err
						}
						objectVal[mk] = result
					} else {
						secondMap, ok := objectVal[mk]
						if !ok {
							secondMap = make(map[string]interface{})
							objectVal[mk] = secondMap
						}
						err := handleObject(secondMap, foreachK, foreachV, exprF, inputScope)
						if err != nil {
							return nil, err
						}
					}
				}
			} else {
				err := handleObject(objectVal, mk, mv, exprF, inputScope)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return objectVal, nil
}

func (f *foreach) handle(arrayMappingFields map[string]interface{}, inputScope data.Scope) (interface{}, error) {
	var newSourceArray []interface{}
	if f.sourceFrom == NEWARRAY {
		newSourceArray = make([]interface{}, 1)
	} else {
		fromValue, err := getExpressionValue("="+f.sourceFrom, f.exprFactory, inputScope)
		if err != nil {
			return nil, fmt.Errorf("get value from source array error, %s", err.Error())
		}
		newSourceArray, err = coerce.ToArray(fromValue)
		if err != nil {
			return nil, fmt.Errorf("source array not an array, error %s", err.Error())
		}
	}

	targetValues := make([]interface{}, len(newSourceArray))
	for i, sourceValue := range newSourceArray {
		inputScope = newLoopScope(sourceValue, f.index, inputScope)
		item, err := handleObjectMapping(arrayMappingFields, f.exprFactory, inputScope)
		if err != nil {
			return nil, err
		}
		targetValues[i] = item
	}
	return targetValues, nil

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

func isExpr(value interface{}) bool {
	if strVal, ok := value.(string); ok && len(strVal) > 0 && (strVal[0] == '=') {
		return true
	}
	return false
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

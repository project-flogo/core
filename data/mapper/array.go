package mapper

import (
	"encoding/json"
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/support/log"
	"runtime/debug"
	"strings"
)

const (
	PRIMITIVE = "primitive"
	FOREACH   = "foreach"
	NEWARRAY  = "NEWARRAY"
)

type ArrayMapperFactory struct {
	exprFactory expression.Factory
}

type ArrayMapping struct {
	From   interface{}     `json:"from"`
	To     string          `json:"to"`
	Type   string          `json:"type"`
	Fields []*ArrayMapping `json:"fields,omitempty"`
}

func ParseArrayMapping(arrayDatadata interface{}) (*ArrayMapping, error) {
	amapping := &ArrayMapping{}
	switch t := arrayDatadata.(type) {
	case string:
		err := json.Unmarshal([]byte(t), amapping)
		if err != nil {
			return nil, err
		}
	case interface{}:
		s, err := coerce.ToString(t)
		if err != nil {
			return nil, fmt.Errorf("Convert array mapping value to string error, due to [%s]", err.Error())
		}
		err = json.Unmarshal([]byte(s), amapping)
		if err != nil {
			return nil, err
		}
	}
	return amapping, nil
}

func IsArrayMapping(value interface{}) bool {
	a, err := ParseArrayMapping(value)
	if err != nil {
		return false
	}

	return a.Validate() == nil
}

func (a *ArrayMapping) Validate() error {
	//Validate root from/to field
	if a.From == nil {
		return fmt.Errorf("The array mapping validation failed for the mapping [%s]. Ensure valid array is mapped in the mapper. ", a.From)
	}

	if a.To == "" || len(a.To) <= 0 {
		return fmt.Errorf("The array mapping validation failed for the mapping [%s]. Ensure valid array is mapped in the mapper. ", a.To)
	}

	if a.Type == FOREACH {
		//Validate root from/to field
		if a.From == NEWARRAY {
			//Make sure no array ref fields exist
			for _, field := range a.Fields {
				if field.Type == FOREACH {
					return field.Validate()
				}
				stringVal, ok := field.From.(string)
				if ok && isArrayMapping(stringVal) {
					return fmt.Errorf("The array mapping validation failed, invalid new array mapping [%s]", stringVal)
				}

			}
		} else {
			for _, field := range a.Fields {
				if field.Type == FOREACH {
					return field.Validate()
				}
			}
		}
	}

	return nil
}

func (am *ArrayMapperFactory) NewArrayExpr(value interface{}) (expression.Expr, error) {
	aMapping, err := ParseArrayMapping(value)
	if err != nil {
		return nil, fmt.Errorf("parsing array mapping failed %s", err)
	}
	return &ArrayExpr{arrayMappings: aMapping, exprFactory: am.exprFactory}, nil
}

func isArrayMapping(ref string) bool {
	if ref != "" {
		return strings.HasPrefix(ref, "$.")
	}
	return false
}

type ArrayExpr struct {
	arrayMappings *ArrayMapping
	exprFactory   expression.Factory
}

func (am *ArrayExpr) Eval(inputScope data.Scope) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			log.RootLogger().Error("%+v", r)
			log.RootLogger().Debugf("StackTrace: %s", debug.Stack())
		}
	}()

	arrayValue, err := am.arrayMappings.Run(am.exprFactory, inputScope)
	if err != nil {
		return nil, err
	}
	return arrayValue, nil
}

func (a *ArrayMapping) Run(exprFactory expression.Factory, scope data.Scope) ([]interface{}, error) {
	//First level must be foreach
	if a.Type == FOREACH {
		//First Level
		var fromValue interface{}
		var err error

		stringVal, ok := a.From.(string)
		if !ok {
			return nil, fmt.Errorf("Invalid mapping root from %s", a.From)
		}

		if strings.EqualFold(stringVal, NEWARRAY) {
			//log.Debugf("Init a new array for field", a.To)
			fromValue = make([]interface{}, 1)
		} else {
			fromValue, err = getExpressionValue(nil, stringVal, exprFactory, scope)
			if err != nil {
				return nil, err
			}
			//fromValue, err = expr.Eval(scope)
		}

		//Loop array
		fromArrayvalues, err := coerce.ToArray(fromValue)
		if err != nil {
			return nil, fmt.Errorf("Failed to get array value from [%s], due to error- [%s] value not an array", a.From, a.From)
		}

		targetArray := make([]interface{}, len(fromArrayvalues))
		for i, _ := range targetArray {
			targetArray[i] = make(map[string]interface{})
		}

		//Check if fields is empty for primitive array mapping
		if a.Fields == nil || len(a.Fields) <= 0 {
			//Set value directlly to MapTo field
			targetArray = fromArrayvalues
			return targetArray, nil
		}

		for i, arrayV := range fromArrayvalues {
			err := a.Iterator(arrayV, targetArray[i].(map[string]interface{}), a.Fields, exprFactory, scope)
			if err != nil {
				log.Error(err)
				return nil, err
			}
		}
		return targetArray, nil
	}
	return nil, fmt.Errorf("array mapping root must be foreach")

}

func (a *ArrayMapping) Iterator(fromValue interface{}, targetValues map[string]interface{}, fields []*ArrayMapping, exprF expression.Factory, scope data.Scope) error {
	for _, arrayField := range fields {
		maptoKey := arrayField.To
		if strings.HasPrefix(arrayField.To, "$.") || strings.HasPrefix(arrayField.To, "$$") {
			maptoKey = arrayField.To[2:]
		}

		switch arrayField.Type {
		case FOREACH:
			var fromArrayvalues []interface{}
			if strings.EqualFold(arrayField.From.(string), NEWARRAY) {
				//log.Debugf("Init a new array for field", arrayField.To)
				fromArrayvalues = make([]interface{}, 1)
			} else {
				fValue, err := getExpressionValue(fromValue, arrayField.From, exprF, scope)
				if err != nil {
					return err
				}
				var ok bool
				fromArrayvalues, ok = fValue.([]interface{})
				if !ok {
					return fmt.Errorf("Failed to get array value from [%s], due to error- value not an array", fValue)
				}
			}

			objArray := make([]interface{}, len(fromArrayvalues))
			for i, _ := range objArray {
				objArray[i] = make(map[string]interface{})
			}

			targetValues[maptoKey] = objArray

			//Check if fields is empty for primitive array mapping
			if arrayField.Fields == nil || len(arrayField.Fields) <= 0 {
				objArray = fromArrayvalues
				continue
			}

			for i, arrayV := range fromArrayvalues {
				err := arrayField.Iterator(arrayV, objArray[i].(map[string]interface{}), arrayField.Fields, exprF, scope)
				if err != nil {
					log.Error(err)
					return err
				}
			}
		default:
			value, err := getExpressionValue(fromValue, arrayField.From, exprF, scope)
			if err != nil {
				return err
			}
			targetValues[maptoKey] = value
		}
	}
	return nil
}

func getExpressionValue(fromValue, fromPath interface{}, exprF expression.Factory, scope data.Scope) (interface{}, error) {
	var finalExpr expression.Expr
	if isExpr(fromPath) {
		var err error
		finalExpr, err = exprF.NewExpr(fromPath.(string)[1:])
		if err != nil {
			return nil, err
		}
	} else {
		finalExpr = expression.NewLiteralExpr(fromPath)
	}
	//Add current from value to scope
	if fromValue != nil {
		scope = data.NewSimpleScope(fromValue.(map[string]interface{}), scope)

	}
	return finalExpr.Eval(scope)
}

func isExpr(value interface{}) bool {
	if strVal, ok := value.(string); ok && len(strVal) > 0 && (strVal[0] == '=' || strVal[0] == '$') {
		return true
	}
	return false
}

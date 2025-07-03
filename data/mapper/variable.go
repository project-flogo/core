package mapper

import (
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression"
	"strings"
)

const (
	Variable = "@variables"
)

type VariableMapper struct {
	name       string
	vartype    interface{}
	expression expression.Expr
}

func IsVariableMapping(value interface{}) bool {
	switch t := value.(type) {
	case map[string]interface{}:
		for k := range t {
			if strings.HasPrefix(k, Variable) {
				return true
			}
		}
		return false
	default:
		obj, _ := coerce.ToObject(value)
		if obj != nil {
			for k, _ := range obj {
				if strings.HasPrefix(k, Variable) {
					return true
				}
			}
		}
		return false
	}

}

func createVariableMapper(vars []interface{}, exprF expression.Factory) ([]VariableMapper, error) {
	variables := make([]VariableMapper, 0, 5)
	for _, vv := range vars {
		variable := vv.(map[string]interface{})
		name := variable["name"].(string)
		expr, err := newExpr(variable["value"].(string), exprF)
		if err != nil {
			return nil, err
		}
		variables = append(variables, VariableMapper{
			name:       name,
			vartype:    variable["type"].(string),
			expression: expr,
		})
	}
	return variables, nil
}

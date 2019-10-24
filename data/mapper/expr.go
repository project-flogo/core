package mapper

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/mapper/config"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support/log"
	"strings"
)

type ExprMapperFactory struct {
	exprFactory   expression.Factory
	objectFactory expression.Factory
}

func NewFactory(resolver resolve.CompositeResolver) Factory {
	exprFactory := expression.NewFactory(resolver)
	objMapperFactory := NewObjectMapperFactory(exprFactory)
	return &ExprMapperFactory{exprFactory: exprFactory, objectFactory: objMapperFactory}
}

func (mf *ExprMapperFactory) NewMapper(mappings map[string]interface{}) (Mapper, error) {

	if len(mappings) == 0 {
		return nil, nil
	}

	exprMappings := make(map[string]expression.Expr)
	for key, value := range mappings {
		if value != nil {
			switch t := value.(type) {
			case string:
				if len(t) > 0 && t[0] == '=' {
					//it's an expression
					expr, err := mf.exprFactory.NewExpr(t[1:])
					if err != nil {
						return nil, err
					}
					exprMappings[key] = expr
				} else {
					exprMappings[key] = expression.NewLiteralExpr(value)
				}
			default:
				if mapping, ok := GetObjectMapping(t); ok {
					//Object mapping
					objectExpr, err := NewObjectMapperFactory(mf.exprFactory).(*ObjectMapperFactory).NewObjectMapper(mapping)
					if err != nil {
						return nil, err
					}
					exprMappings[key] = objectExpr
				} else {
					exprMappings[key] = expression.NewLiteralExpr(value)
				}
			}
		}
	}

	if len(exprMappings) == 0 {
		return nil, nil
	}

	return &ExprMapper{mappings: exprMappings}, nil
}

type ExprMapper struct {
	mappings map[string]expression.Expr
}

func (m *ExprMapper) Apply(inputScope data.Scope) (map[string]interface{}, error) {
	output := make(map[string]interface{}, len(m.mappings))
	for key, expr := range m.mappings {
		val, err := expr.Eval(inputScope)
		if err != nil {
			if config.IsMappingIgnoreErrorsOn() {
				log.RootLogger().Warnf("expresson eval error; %s", err.Error())
				//Skip value set.
				continue
			}
			//todo add some context to error (consider adding String() to exprImpl)
			return nil, err
		}
		output[key] = val
	}

	return output, nil
}

func IsLiteral(value interface{}) bool {
	if value != nil {
		switch t := value.(type) {
		case map[string]interface{}:
			for k, v := range t {
				if strings.HasPrefix(k, FOREACH) {
					return false
				}
				if !IsLiteral(v) {
					return false
				}
			}
		case []interface{}:
			for _, element := range t {
				if !IsLiteral(element) {
					return false
				}
			}
		default:
			return !isExpr(t)
		}
	}
	return true
}

func isExpr(value interface{}) bool {
	if strVal, ok := value.(string); ok && len(strVal) > 0 && (strVal[0] == '=') {
		return true
	}
	return false
}

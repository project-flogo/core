package mapper

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/resolve"
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
					expr, err := mf.exprFactory.NewExpr(t[1:])
					if err != nil {
						return nil, err
					}
					exprMappings[key] = expr
				} else {
					exprMappings[key] = expression.NewLiteralExpr(value)
				}
			//Object mapping
			case map[string]interface{}, []interface{}:
				objectExpr, err := NewObjectMapperFactory(mf.exprFactory).(*ObjectMapperFactory).NewObjectMapper(t)
				if err != nil {
					return nil, err
				}
				exprMappings[key] = objectExpr
			default:
				exprMappings[key] = expression.NewLiteralExpr(value)
			}
		}
	}
	if len(exprMappings) <= 0 {
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
			//todo add some context to error (consider adding String() to exprImpl)
			return nil, err
		}
		output[key] = val
	}

	return output, nil
}

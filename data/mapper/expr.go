package mapper

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/resolve"
)

type ExprMapperFactory struct {
	exprFactory expression.Factory
}

func NewFactory(resolver resolve.CompositeResolver) Factory {
	exprFactory := expression.NewFactory(resolver)
	return &ExprMapperFactory{exprFactory: exprFactory}
}

func (mf *ExprMapperFactory) NewMapper(mappings map[string]interface{}) (Mapper, error) {
	exprMappings := make(map[string]expression.Expr, len(mappings))

	for key, value := range mappings {

		if strVal, ok := value.(string); ok && len(strVal) > 0 && strVal[0] == '=' {
			expr, err := mf.exprFactory.NewExpr(strVal[1:])
			if err != nil {
				return nil, err
			}
			exprMappings[key] = expr
		} else {
			exprMappings[key] = expression.NewLiteralExpr(value)
		}
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

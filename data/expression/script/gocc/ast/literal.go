package ast

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/resolve"
	"strings"

	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression/script/gocc/token"
)

type literalExpr struct {
	val interface{}
}

func (*literalExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	return nil
}

func (e *literalExpr) Eval(scope data.Scope) (interface{}, error) {
	return e.val, nil
}

func NewLiteral(litType string, lit interface{}) (Expr, error) {

	litAsStr := strings.TrimSpace(string(lit.(*token.Token).Lit)) //todo is trim overkill

	switch litType {
	case "int":
		i, err := coerce.ToInt(litAsStr)
		return &literalExpr{val: i}, err
	case "float":
		f, err := coerce.ToFloat64(litAsStr)
		return &literalExpr{val: f}, err
	case "bool":
		b, err := coerce.ToBool(litAsStr)
		return &literalExpr{val: b}, err
	case "string":
		s := litAsStr[1 : len(litAsStr)-1] //remove quotes
		return &literalExpr{val: s}, nil
	case "nil":
		return &literalExpr{val: nil}, nil
	}

	return nil, fmt.Errorf("unsupported literal type '%s'", litType)
}

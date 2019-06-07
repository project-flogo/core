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
	typ string
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
		return &literalExpr{val: i, typ: "int"}, err
	case "float":
		f, err := coerce.ToFloat64(litAsStr)
		return &literalExpr{val: f, typ: "float"}, err
	case "bool":
		b, err := coerce.ToBool(litAsStr)
		return &literalExpr{val: b, typ: "bool"}, err
	case "string":
		return &literalExpr{val: removeQuotedAndEscaped(litAsStr), typ: "string"}, nil
	case "nil":
		return &literalExpr{val: nil, typ: "nil"}, nil
	}

	return nil, fmt.Errorf("unsupported literal type '%s'", litType)
}

func removeQuotedAndEscaped(str string) string {
	//Eascap string
	firstChar := str[0]
	switch firstChar {
	case '"':
		str = str[1 : len(str)-1]
		if strings.Contains(str, "\\\"") {
			str = strings.Replace(str, `\"`, `"`, -1)
		}
		str = handleNewline(str)
	case '\'':
		str = str[1 : len(str)-1]
		//Eascap string
		if strings.Contains(str, "\\'") {
			str = strings.Replace(str, "\\'", "'", -1)
		}
		str = handleNewline(str)
	default:
		str = str[1 : len(str)-1]
	}

	return str
}

func handleNewline(str string) string {
	//Handle \n,\r and \t
	str = strings.Replace(str, `\n`, string('\n'), -1)
	str = strings.Replace(str, `\r`, string('\r'), -1)
	str = strings.Replace(str, `\t`, string('\t'), -1)
	return str
}

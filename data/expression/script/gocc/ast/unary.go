package ast

import (
	"encoding/json"
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/resolve"
	"reflect"
	"strings"

	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression/script/gocc/token"
)

func NewUnaryExpr(operand, tok interface{}) (Expr, error) {
	expr := tok.(Expr)
	op := string(operand.(*token.Token).Lit)

	switch op {
	case "!":
		return &unaryNotExpr{expr: expr}, nil
	case "-":
		return &unaryNegExpr{expr: expr}, nil
	}

	return nil, fmt.Errorf("unsupported arithmetic operator '%s'", op)
}

type unaryNotExpr struct {
	expr Expr
}

func (e *unaryNotExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.expr.Init(resolver, false)
	return err
}

func (e *unaryNotExpr) Eval(scope data.Scope) (interface{}, error) {
	v, err := e.expr.Eval(scope)
	if err != nil {
		return nil, err
	}

	if v == nil {
		//todo validate
		return nil, fmt.Errorf("cannot not a nil")
	}

	switch ve := v.(type) {
	case bool:
		return !ve, nil
	}

	return false, fmt.Errorf("cannot not '%s'", reflect.TypeOf(v).String())
}

type unaryNegExpr struct {
	expr Expr
}

func (e *unaryNegExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.expr.Init(resolver, false)
	return err
}

func (e *unaryNegExpr) Eval(scope data.Scope) (interface{}, error) {
	v, err := e.expr.Eval(scope)
	if err != nil {
		return nil, err
	}

	if v == nil {
		//todo validate
		return nil, fmt.Errorf("cannot negate a nil")
	}

	switch ve := v.(type) {
	case int, int32, int64:
		vi, _ := coerce.ToInt(ve) //todo should this be Int64
		return -vi, nil
	case float32, float64:
		vf, _ := coerce.ToFloat64(ve)
		return -vf, nil
	case json.Number:
		if strings.Contains(ve.String(), ".") {
			vf, _ := ve.Float64()
			return -vf, nil
		} else {
			vf, _ := ve.Int64()
			return -vf, nil
		}
	}

	return false, fmt.Errorf("cannot not '%s'", reflect.TypeOf(v).String())
}

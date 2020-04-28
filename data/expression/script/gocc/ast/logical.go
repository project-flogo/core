package ast

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/resolve"

	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression/script/gocc/token"
)

func NewBoolExpr(left, operand, right interface{}) (Expr, error) {
	le := left.(Expr)
	re := right.(Expr)
	op := string(operand.(*token.Token).Lit)

	switch op {
	case "||":
		return &boolOrExpr{left: le, right: re}, nil
	case "&&":
		return &boolAndExpr{left: le, right: re}, nil
	}

	return nil, fmt.Errorf("unsupported boolean operator '%s'", op)
}

type boolOrExpr struct {
	left, right Expr
}

func (e *boolOrExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.left.Init(resolver, false)
	if err != nil {
		return err
	}
	err = e.right.Init(resolver, false)
	return err
}

func (e *boolOrExpr) Eval(scope data.Scope) (interface{}, error) {
	lv , err := e.left.Eval(scope)
	if err != nil {
		return nil, err
	}
	lb, err := coerce.ToBool(lv)
	if err != nil {
		return nil, err
	}

	// Return true if left side true
	if lb {
		return true, nil
	}

	rv, err := e.right.Eval(scope)
	if err != nil {
		return nil, err
	}

	rb, err := coerce.ToBool(rv)
	if err != nil {
		return nil, err
	}

	return lb || rb, nil
}

type boolAndExpr struct {
	left, right Expr
}

func (e *boolAndExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.left.Init(resolver, false)
	if err != nil {
		return err
	}
	err = e.right.Init(resolver, false)
	return err
}

func (e *boolAndExpr) Eval(scope data.Scope) (interface{}, error) {
	lv , err := e.left.Eval(scope)
	if err != nil {
		return nil, err
	}
	lb, err := coerce.ToBool(lv)
	if err != nil {
		return nil, err
	}

	// Return false if left side false
	if !lb {
		return false, nil
	}

	rv, err := e.right.Eval(scope)
	if err != nil {
		return nil, err
	}

	rb, err := coerce.ToBool(rv)
	if err != nil {
		return nil, err
	}
	return lb && rb, nil
}

package ast

import (
	"fmt"
	"strings"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/data/expression/script/gocc/token"
	"github.com/project-flogo/core/data/resolve"
)

func NewFuncExpr(name interface{}, args interface{}) (Expr, error) {
	funcName := string(name.(*token.Token).Lit)
	if strings.EqualFold(funcName, "isdefined") {
		return &IsDefinedExpr{refExpr: args.([]Expr)[0]}, nil
	} else if strings.EqualFold(funcName, "getValue") {
		return &GetValueExpr{refExpr: args.([]Expr)[0], valueExpr: args.([]Expr)[1]}, nil
	}

	f := function.Get(funcName)
	if f == nil {
		return nil, fmt.Errorf("unable to parse expression, function [%s] is not installed", funcName)
	}
	fe := &funcExpr{f: f}

	if args, ok := args.([]Expr); ok {
		if len(args) > 0 {
			exprParams := make([]Expr, len(args))
			for idx, arg := range args {
				exprParams[idx] = arg.(Expr)
			}

			fe.params = exprParams
		}
	}

	return fe, nil
}

type funcExpr struct {
	f      function.Function
	params []Expr
}

func (e *funcExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	for _, param := range e.params {
		err := param.Init(resolver, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *funcExpr) Eval(scope data.Scope) (interface{}, error) {

	vals := make([]interface{}, len(e.params))
	for idx, param := range e.params {
		v, err := param.Eval(scope)
		if err != nil {
			return nil, err
		}
		vals[idx] = v
	}

	return function.Eval(e.f, vals...)
}

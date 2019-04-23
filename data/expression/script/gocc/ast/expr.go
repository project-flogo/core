package ast

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/resolve"
	"strings"
)

type Expr interface {
	Init(resolver resolve.CompositeResolver, root bool) error //todo can use root to multi-thread eval of root node

	Eval(scope data.Scope) (interface{}, error)
}

func evalLR(left, right Expr, scope data.Scope) (lv, rv interface{}, err error) {
	lv, err = left.Eval(scope)
	if err != nil {
		return nil, nil, err
	}
	rv, err = right.Eval(scope)
	return lv, rv, err
}

func NewExprList(x interface{}) ([]Expr, error) {
	if x, ok := x.(Expr); ok {
		return []Expr{x}, nil
	}
	return nil, fmt.Errorf("invalid expression list expression type; expected ast.Expr, got %T", x)
}

func AppendToExprList(list, x interface{}) ([]Expr, error) {
	lst, ok := list.([]Expr)
	if !ok {
		return nil, fmt.Errorf("invalid expression list type; expected []ast.Expr, got %T", list)
	}
	if x, ok := x.(Expr); ok {
		return append(lst, x), nil
	}
	return nil, fmt.Errorf("invalid expression list expression type; expected ast.Expr, got %T", x)
}

func NewTernaryExpr(ifNode, thenNode, elseNode interface{}) (Expr, error) {

	ifExpr := ifNode.(Expr)
	thenExpr := thenNode.(Expr)
	elseExpr := elseNode.(Expr)

	return &exprTernary{ifExpr: ifExpr, thenExpr: thenExpr, elseExpr: elseExpr}, nil
}

func NewTernaryArgument(first interface{}) (Expr, error) {
	switch t := first.(type) {
	case Expr:
		return t, nil
	default:
		return nil, fmt.Errorf("unsupported ternary type %+v", first)
	}
}

type exprTernary struct {
	ifExpr, thenExpr, elseExpr Expr
}

func (e *exprTernary) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.ifExpr.Init(resolver, false)
	if err != nil {
		return err
	}
	err = e.thenExpr.Init(resolver, false)
	if err != nil {
		return err
	}
	err = e.elseExpr.Init(resolver, false)
	return err
}

func (e *exprTernary) Eval(scope data.Scope) (interface{}, error) {

	iv, err := e.ifExpr.Eval(scope)
	if err != nil {
		return nil, err
	}

	bv, err := coerce.ToBool(iv)
	if err != nil {
		return nil, err
	}

	if bv {
		tv, err := e.thenExpr.Eval(scope)
		if err != nil {
			return nil, err
		}
		return tv, nil
	} else {
		ev, err := e.elseExpr.Eval(scope)
		if err != nil {
			return nil, err
		}
		return ev, nil
	}
}

func NewRefExpr(refNode ...interface{}) (Expr, error) {
	refFields, err := Concat(refNode...)
	if err != nil {
		return nil, err
	}
	return &exprRef{fields: refFields}, nil
}

type exprRef struct {
	fields       []interface{}
	res          resolve.Resolution
	resolver     resolve.CompositeResolver
	hasIndexExpr bool
}

func (e *exprRef) Init(resolver resolve.CompositeResolver, root bool) error {
	e.resolver = resolver
	for _, v := range e.fields {
		switch t := v.(type) {
		case Expr:
			e.hasIndexExpr = true
			err := t.Init(resolver, root)
			if err != nil {
				return err
			}
		}
	}

	//Make resolution if no index expression
	if !e.hasIndexExpr {
		var err error
		t := make([]string, len(e.fields))
		for i, v := range e.fields {
			t[i] = v.(string)
		}
		e.res, err = resolver.GetResolution(strings.Join(t, ""))
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *exprRef) Eval(scope data.Scope) (data interface{}, err error) {
	if e.hasIndexExpr {
		//Get final resolved ref from index expression
		e.res, err = e.constructRealRef(scope)
		if err != nil {
			return nil, err
		}
	}
	return e.res.GetValue(scope)
}

func (e *exprRef) constructRealRef(scope data.Scope) (resolve.Resolution, error) {
	var ref = ""
	for _, v := range e.fields {
		switch t := v.(type) {
		case Expr:
			indexRef, err := t.Eval(scope)
			if err != nil {
				return nil, err
			}
			ref = ref + indexRef.(string)
		case string:
			ref = ref + t
		}
	}
	return e.resolver.GetResolution(ref)
}

type keyIndexExpr struct {
	expr Expr
}

func (e *keyIndexExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	return e.expr.Init(resolver, root)
}

func (e *keyIndexExpr) Eval(scope data.Scope) (interface{}, error) {
	v, err := e.expr.Eval(scope)
	if err != nil {
		return "", fmt.Errorf("eval array index expression error: %s", err.Error())
	}

	switch t := e.expr.(type) {
	case *literalExpr:
		if t.typ == "string" {
			//Add double quotes for string, no matter single qutoes or one tick
			v = `"` + v.(string) + `"`
		}
	}

	index, err := coerce.ToString(v)
	if err != nil {
		return nil, err
	}
	return "[" + index + "]", nil
}

package ast

import (
	"fmt"
	"reflect"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression/script/gocc/token"
	"github.com/project-flogo/core/data/resolve"
)

func NewArithExpr(left, operand, right interface{}) (Expr, error) {
	le := left.(Expr)
	re := right.(Expr)
	op := string(operand.(*token.Token).Lit)

	switch op {
	case "+":
		return &arithAddExpr{left: le, right: re}, nil
	case "-":
		return &arithSubExpr{left: le, right: re}, nil
	case "*":
		return &arithMulExpr{left: le, right: re}, nil
	case "/":
		return &arithDivExpr{left: le, right: re}, nil
	case "%":
		return &arithModExpr{left: le, right: re}, nil
	}

	return nil, fmt.Errorf("unsupported arithmetic operator '%s'", op)
}

type arithAddExpr struct {
	left, right Expr
}

func (e *arithAddExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.left.Init(resolver, false)
	if err != nil {
		return err
	}
	err = e.right.Init(resolver, false)
	return err
}

func (e *arithAddExpr) Eval(scope data.Scope) (interface{}, error) {
	lv, rv, err := evalLR(e.left, e.right, scope)
	if err != nil {
		return nil, err
	}

	if lv == nil || rv == nil {
		//todo validate
		return nil, fmt.Errorf("cannot add %v to %v", lv, rv)
	}

	rt := reflect.TypeOf(rv).Kind()
	switch le := lv.(type) {
	case int, int32, int64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 {
			li, _ := coerce.ToInt(lv) //todo should this be Int64
			ri, _ := coerce.ToInt(rv) //todo should this be Int64
			return li + ri, nil
		}
		if rt == reflect.Float32 || rt == reflect.Float64 {
			li, _ := coerce.ToFloat64(lv)
			ri, _ := coerce.ToFloat64(rv)
			return li + ri, nil
		}
	case float32, float64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 || rt == reflect.Float32 || rt == reflect.Float64 {
			lf, _ := coerce.ToFloat64(lv)
			rf, _ := coerce.ToFloat64(rv)
			return lf + rf, nil
		}
	case string:
		rs, _ := coerce.ToString(rv)
		return le + rs, nil
	}

	if rt == reflect.String {
		ls, _ := coerce.ToString(rv)
		rs, _ := coerce.ToString(rv)
		return ls + rs, nil
	}

	return false, fmt.Errorf("cannot add %s to %s", reflect.TypeOf(lv).String(), reflect.TypeOf(rv).String())
}

type arithSubExpr struct {
	left, right Expr
}

func (e *arithSubExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.left.Init(resolver, false)
	if err != nil {
		return err
	}
	err = e.right.Init(resolver, false)
	return err
}

func (e *arithSubExpr) Eval(scope data.Scope) (interface{}, error) {
	lv, rv, err := evalLR(e.left, e.right, scope)
	if err != nil {
		return nil, err
	}

	if lv == nil || rv == nil {
		//todo validate
		return nil, fmt.Errorf("cannot subtract %v from %v", rv, lv)
	}

	rt := reflect.TypeOf(rv).Kind()
	switch lv.(type) {
	case int, int32, int64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 {
			li, _ := coerce.ToInt(lv) //todo should this be Int64
			ri, _ := coerce.ToInt(rv) //todo should this be Int64
			return li - ri, nil
		}
		if rt == reflect.Float32 || rt == reflect.Float64 {
			li, _ := coerce.ToFloat64(lv)
			ri, _ := coerce.ToFloat64(rv)
			return li - ri, nil
		}
	case float32, float64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 || rt == reflect.Float32 || rt == reflect.Float64 {
			lf, _ := coerce.ToFloat64(lv)
			rf, _ := coerce.ToFloat64(rv)
			return lf - rf, nil
		}
	}

	return false, fmt.Errorf("cannot subtract %s from %s", reflect.TypeOf(rv).String(), reflect.TypeOf(lv).String())
}

type arithMulExpr struct {
	left, right Expr
}

func (e *arithMulExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.left.Init(resolver, false)
	if err != nil {
		return err
	}
	err = e.right.Init(resolver, false)
	return err
}

func (e *arithMulExpr) Eval(scope data.Scope) (interface{}, error) {
	lv, rv, err := evalLR(e.left, e.right, scope)
	if err != nil {
		return nil, err
	}

	if lv == nil || rv == nil {
		//todo validate
		return nil, fmt.Errorf("cannot multiply %v with %v", rv, lv)
	}

	rt := reflect.TypeOf(rv).Kind()
	switch lv.(type) {
	case int, int32, int64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 {
			li, _ := coerce.ToInt(lv) //todo should this be Int64
			ri, _ := coerce.ToInt(rv) //todo should this be Int64
			return li * ri, nil
		}
		if rt == reflect.Float32 || rt == reflect.Float64 {
			li, _ := coerce.ToFloat64(lv)
			ri, _ := coerce.ToFloat64(rv)
			return li * ri, nil
		}
	case float32, float64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 || rt == reflect.Float32 || rt == reflect.Float64 {
			lf, _ := coerce.ToFloat64(lv)
			rf, _ := coerce.ToFloat64(rv)
			return lf * rf, nil
		}
	}

	return false, fmt.Errorf("cannot multiply %s with %s", reflect.TypeOf(lv).String(), reflect.TypeOf(rv).String())
}

type arithDivExpr struct {
	left, right Expr
}

func (e *arithDivExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.left.Init(resolver, false)
	if err != nil {
		return err
	}
	err = e.right.Init(resolver, false)
	return err
}

func (e *arithDivExpr) Eval(scope data.Scope) (interface{}, error) {
	lv, rv, err := evalLR(e.left, e.right, scope)
	if err != nil {
		return nil, err
	}

	if lv == nil || rv == nil {
		//todo validate
		return nil, fmt.Errorf("cannot div %v with %v", lv, rv)
	}

	rt := reflect.TypeOf(rv).Kind()
	switch lv.(type) {
	case int, int32, int64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 {
			li, _ := coerce.ToInt(lv) //todo should this be Int64
			ri, _ := coerce.ToInt(rv) //todo should this be Int64
			return li / ri, nil
		}
		if rt == reflect.Float32 || rt == reflect.Float64 {
			li, _ := coerce.ToFloat64(lv)
			ri, _ := coerce.ToFloat64(rv)
			return li / ri, nil
		}
	case float32, float64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 || rt == reflect.Float32 || rt == reflect.Float64 {
			lf, _ := coerce.ToFloat64(lv)
			rf, _ := coerce.ToFloat64(rv)
			return lf / rf, nil
		}
	}

	return false, fmt.Errorf("cannot divide %s with %s", reflect.TypeOf(lv).String(), reflect.TypeOf(rv).String())
}

type arithModExpr struct {
	left, right Expr
}

func (e *arithModExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.left.Init(resolver, false)
	if err != nil {
		return err
	}
	err = e.right.Init(resolver, false)
	return err
}

func (e *arithModExpr) Eval(scope data.Scope) (interface{}, error) {
	lv, rv, err := evalLR(e.left, e.right, scope)
	if err != nil {
		return nil, err
	}

	if lv == nil || rv == nil {
		//todo validate
		return nil, fmt.Errorf("cannot mod %v with %v", rv, lv)
	}

	rt := reflect.TypeOf(rv).Kind()
	switch lv.(type) {
	case int, int32, int64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 {
			li, _ := coerce.ToInt(lv) //todo should this be Int64
			ri, _ := coerce.ToInt(rv) //todo should this be Int64
			return li % ri, nil
		}
		if rt == reflect.Float32 || rt == reflect.Float64 {
			li, _ := coerce.ToInt(lv) //todo should this be Int64
			ri, _ := coerce.ToInt(rv) //todo should this be Int64
			return li % ri, nil
		}
	case float32, float64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 || rt == reflect.Float32 || rt == reflect.Float64 {
			li, _ := coerce.ToInt(lv) //todo should this be Int64
			ri, _ := coerce.ToInt(rv) //todo should this be Int64
			return li % ri, nil
		}
	}

	return false, fmt.Errorf("cannot mod %s with %s", reflect.TypeOf(lv).String(), reflect.TypeOf(rv).String())
}

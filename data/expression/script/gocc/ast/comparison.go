package ast

import (
	"encoding/json"
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression/script/gocc/token"
	"github.com/project-flogo/core/data/resolve"
	"reflect"
	"strings"
	"time"
)

func NewCmpExpr(left, operand, right interface{}) (Expr, error) {
	le := left.(Expr)
	re := right.(Expr)
	op := string(operand.(*token.Token).Lit)

	switch op {
	case "==":
		return &cmpEqExpr{left: le, right: re}, nil
	case "!=":
		return &cmpNotEqExpr{left: le, right: re}, nil
	case "<":
		return &cmpLtExpr{left: le, right: re}, nil
	case "<=":
		return &cmpLtEqExpr{left: le, right: re}, nil
	case ">":
		return &cmpGtExpr{left: le, right: re}, nil
	case ">=":
		return &cmpGtEqExpr{left: le, right: re}, nil
	}

	return nil, fmt.Errorf("unsupported cmpative operator '%s'", op)
}

type cmpEqExpr struct {
	left, right Expr
}

func (e *cmpEqExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.left.Init(resolver, false)
	if err != nil {
		return err
	}
	err = e.right.Init(resolver, false)
	return err
}

func (e *cmpEqExpr) Eval(scope data.Scope) (interface{}, error) {
	lv, rv, err := evalLR(e.left, e.right, scope)
	if err != nil {
		return nil, err
	}

	if lv == nil || rv == nil {
		return rv == nil && lv == nil, nil
	}

	rt := reflect.TypeOf(rv).Kind()
	switch le := lv.(type) {
	case int, int32, int64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 {
			li, _ := coerce.ToInt64(lv)
			ri, _ := coerce.ToInt64(rv)
			return li == ri, nil
		}
		if rt == reflect.Float32 || rt == reflect.Float64 {
			li, _ := coerce.ToFloat64(lv)
			ri, _ := coerce.ToFloat64(rv)
			return li == ri, nil
		}

		if isJsonNumber(rv) {
			li, _ := coerce.ToInt64(lv)
			ri, _ := coerce.ToInt64(rv)
			return li == ri, nil
		}

	case float32, float64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 || rt == reflect.Float32 || rt == reflect.Float64 || isJsonNumber(rv) {
			lf, _ := coerce.ToFloat64(lv)
			rf, _ := coerce.ToFloat64(rv)
			return lf == rf, nil
		}
	case json.Number:
		if strings.Contains(le.String(), ".") {
			li, _ := le.Float64()
			ri, _ := coerce.ToFloat64(rv)
			return li == ri, nil
		} else {
			li, _ := le.Int64()
			ri, _ := coerce.ToInt64(rv)
			return li == ri, nil
		}
	case string:
		if rt == reflect.String {
			rs, _ := coerce.ToString(rv)
			return le == rs, nil
		}
	case bool:
		if rt == reflect.Bool {
			rb, _ := coerce.ToBool(rv)
			return le == rb, nil
		}
	case time.Time:
		switch re := rv.(type) {
		case time.Time:
			return le.Equal(re), nil
		default:
			t, err := coerce.ToDateTime(rv)
			if err != nil {
				return false, nil
			}
			return le.Equal(t), nil
		}
	}

	return false, nil
}

type cmpNotEqExpr struct {
	left, right Expr
}

func (e *cmpNotEqExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.left.Init(resolver, false)
	if err != nil {
		return err
	}
	err = e.right.Init(resolver, false)
	return err
}

func (e *cmpNotEqExpr) Eval(scope data.Scope) (interface{}, error) {

	lv, rv, err := evalLR(e.left, e.right, scope)
	if err != nil {
		return nil, err
	}

	if lv == nil || rv == nil {
		return !(rv == nil && lv == nil), nil
	}

	rt := reflect.TypeOf(rv).Kind()
	switch le := lv.(type) {
	case int, int32, int64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 {
			li, _ := coerce.ToInt64(lv)
			ri, _ := coerce.ToInt64(rv)
			return li != ri, nil
		}
		if rt == reflect.Float32 || rt == reflect.Float64 {
			li, _ := coerce.ToFloat64(lv)
			ri, _ := coerce.ToFloat64(rv)
			return li != ri, nil
		}

		if isJsonNumber(rv) {
			li, _ := coerce.ToInt64(lv)
			ri, _ := coerce.ToInt64(rv)
			return li != ri, nil
		}
	case float32, float64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 || rt == reflect.Float32 || rt == reflect.Float64 || isJsonNumber(rv) {
			lf, _ := coerce.ToFloat64(lv)
			rf, _ := coerce.ToFloat64(rv)
			return lf != rf, nil
		}
	case json.Number:
		if strings.Contains(le.String(), ".") {
			li, _ := le.Float64()
			ri, _ := coerce.ToFloat64(rv)
			return li != ri, nil
		} else {
			li, _ := le.Int64()
			ri, _ := coerce.ToInt64(rv)
			return li != ri, nil
		}
	case string:
		if rt == reflect.String {
			rs, _ := coerce.ToString(rv)
			return le != rs, nil
		}
	case bool:
		if rt == reflect.Bool {
			rb, _ := coerce.ToBool(rv)
			return le != rb, nil
		}
	case time.Time:
		switch re := rv.(type) {
		case time.Time:
			return !le.Equal(re), nil
		default:
			t, err := coerce.ToDateTime(rv)
			if err != nil {
				return true, nil
			}
			return !le.Equal(t), nil
		}
	}

	return true, nil
}

type cmpGtExpr struct {
	left, right Expr
}

func (e *cmpGtExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.left.Init(resolver, false)
	if err != nil {
		return err
	}
	err = e.right.Init(resolver, false)
	return err
}

func (e *cmpGtExpr) Eval(scope data.Scope) (interface{}, error) {
	lv, rv, err := evalLR(e.left, e.right, scope)
	if err != nil {
		return nil, err
	}

	if lv == nil || rv == nil {
		//todo validate this behavior
		return false, nil
	}

	rt := reflect.TypeOf(rv).Kind()
	switch le := lv.(type) {
	case int, int32, int64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 {
			li, _ := coerce.ToInt64(lv)
			ri, _ := coerce.ToInt64(rv)
			return li > ri, nil
		}
		if rt == reflect.Float32 || rt == reflect.Float64 {
			li, _ := coerce.ToFloat64(lv)
			ri, _ := coerce.ToFloat64(rv)
			return li > ri, nil
		}
		if isJsonNumber(rv) {
			li, _ := coerce.ToInt64(lv)
			ri, _ := coerce.ToInt64(rv)
			return li > ri, nil
		}
	case float32, float64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 || rt == reflect.Float32 || rt == reflect.Float64 || isJsonNumber(rv) {
			lf, _ := coerce.ToFloat64(lv)
			rf, _ := coerce.ToFloat64(rv)
			return lf > rf, nil
		}
	case json.Number:
		if strings.Contains(le.String(), ".") {
			li, _ := le.Float64()
			ri, _ := coerce.ToFloat64(rv)
			return li > ri, nil
		} else {
			li, _ := le.Int64()
			ri, _ := coerce.ToInt64(rv)
			return li > ri, nil
		}
	case string:
		if rt == reflect.String {
			rs, _ := coerce.ToString(rv)
			return le > rs, nil
		}
	case bool:
		return false, nil
	case time.Time:
		switch re := rv.(type) {
		case time.Time:
			return le.After(re), nil
		default:
			t, err := coerce.ToDateTime(rv)
			if err == nil {
				return le.After(t), nil
			}
		}
	}

	return false, fmt.Errorf("cannot compare %s with %s", reflect.TypeOf(lv).String(), reflect.TypeOf(rv).String())
}

type cmpGtEqExpr struct {
	left, right Expr
}

func (e *cmpGtEqExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.left.Init(resolver, false)
	if err != nil {
		return err
	}
	err = e.right.Init(resolver, false)
	return err
}

func (e *cmpGtEqExpr) Eval(scope data.Scope) (interface{}, error) {
	lv, rv, err := evalLR(e.left, e.right, scope)
	if err != nil {
		return nil, err
	}

	if lv == nil || rv == nil {
		return lv == nil && rv == nil, nil
	}

	rt := reflect.TypeOf(rv).Kind()
	switch le := lv.(type) {
	case int, int32, int64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 {
			li, _ := coerce.ToInt64(lv)
			ri, _ := coerce.ToInt64(rv)
			return li >= ri, nil
		}
		if rt == reflect.Float32 || rt == reflect.Float64 {
			li, _ := coerce.ToFloat64(lv)
			ri, _ := coerce.ToFloat64(rv)
			return li >= ri, nil
		}

		if isJsonNumber(rv) {
			li, _ := coerce.ToInt64(lv)
			ri, _ := coerce.ToInt64(rv)
			return li >= ri, nil
		}
	case float32, float64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 || rt == reflect.Float32 || rt == reflect.Float64 || isJsonNumber(rv) {
			lf, _ := coerce.ToFloat64(lv)
			rf, _ := coerce.ToFloat64(rv)
			return lf >= rf, nil
		}
	case json.Number:
		if strings.Contains(le.String(), ".") {
			li, _ := le.Float64()
			ri, _ := coerce.ToFloat64(rv)
			return li >= ri, nil
		} else {
			li, _ := le.Int64()
			ri, _ := coerce.ToInt64(rv)
			return li >= ri, nil
		}
	case string:
		if rt == reflect.String {
			rs, _ := coerce.ToString(rv)
			return le >= rs, nil
		}
	case bool:
		if rt == reflect.Bool {
			rb, _ := coerce.ToBool(rv)
			return le == rb, nil
		}
	case time.Time:
		switch re := rv.(type) {
		case time.Time:
			return le.After(re) || le.Equal(re), nil
		default:
			t, err := coerce.ToDateTime(rv)
			if err == nil {
				return le.After(t) || le.Equal(t), nil
			}
		}
	}

	return false, fmt.Errorf("cannot compare %s with %s", reflect.TypeOf(lv).String(), reflect.TypeOf(rv).String())
}

type cmpLtExpr struct {
	left, right Expr
}

func (e *cmpLtExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.left.Init(resolver, false)
	if err != nil {
		return err
	}
	err = e.right.Init(resolver, false)
	return err
}

func (e *cmpLtExpr) Eval(scope data.Scope) (interface{}, error) {
	lv, rv, err := evalLR(e.left, e.right, scope)
	if err != nil {
		return nil, err
	}

	if lv == nil || rv == nil {
		//todo validate this behavior
		return false, nil
	}

	rt := reflect.TypeOf(rv).Kind()
	switch le := lv.(type) {
	case int, int32, int64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 {
			li, _ := coerce.ToInt64(lv)
			ri, _ := coerce.ToInt64(rv)
			return li < ri, nil
		}
		if rt == reflect.Float32 || rt == reflect.Float64 {
			li, _ := coerce.ToFloat64(lv)
			ri, _ := coerce.ToFloat64(rv)
			return li < ri, nil
		}
		if isJsonNumber(rv) {
			li, _ := coerce.ToInt64(lv)
			ri, _ := coerce.ToInt64(rv)
			return li < ri, nil
		}
	case float32, float64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 || rt == reflect.Float32 || rt == reflect.Float64 || isJsonNumber(rv) {
			lf, _ := coerce.ToFloat64(lv)
			rf, _ := coerce.ToFloat64(rv)
			return lf < rf, nil
		}
	case json.Number:
		if strings.Contains(le.String(), ".") {
			li, _ := le.Float64()
			ri, _ := coerce.ToFloat64(rv)
			return li < ri, nil
		} else {
			li, _ := le.Int64()
			ri, _ := coerce.ToInt64(rv)
			return li < ri, nil
		}
	case string:
		if rt == reflect.String {
			rs, _ := coerce.ToString(rv)
			return le < rs, nil
		}
	case bool:
		if rt == reflect.Bool {
			return false, nil
		}
	case time.Time:
		switch re := rv.(type) {
		case time.Time:
			return le.Before(re), nil
		default:
			t, err := coerce.ToDateTime(rv)
			if err == nil {
				return le.Before(t), nil
			}
		}
	}

	return false, fmt.Errorf("cannot compare %s with %s", reflect.TypeOf(lv).String(), reflect.TypeOf(rv).String())
}

type cmpLtEqExpr struct {
	left, right Expr
}

func (e *cmpLtEqExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := e.left.Init(resolver, false)
	if err != nil {
		return err
	}
	err = e.right.Init(resolver, false)
	return err
}

func (e *cmpLtEqExpr) Eval(scope data.Scope) (interface{}, error) {
	lv, rv, err := evalLR(e.left, e.right, scope)
	if err != nil {
		return nil, err
	}

	if lv == nil || rv == nil {
		return lv == nil && rv == nil, nil
	}

	rt := reflect.TypeOf(rv).Kind()
	switch le := lv.(type) {
	case int, int32, int64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 {
			li, _ := coerce.ToInt64(lv)
			ri, _ := coerce.ToInt64(rv)
			return li <= ri, nil
		}
		if rt == reflect.Float32 || rt == reflect.Float64 {
			li, _ := coerce.ToFloat64(lv)
			ri, _ := coerce.ToFloat64(rv)
			return li <= ri, nil
		}
		if isJsonNumber(rv) {
			li, _ := coerce.ToInt64(lv)
			ri, _ := coerce.ToInt64(rv)
			return li <= ri, nil
		}
	case float32, float64:
		if rt == reflect.Int || rt == reflect.Int32 || rt == reflect.Int64 || rt == reflect.Float32 || rt == reflect.Float64 || isJsonNumber(rv) {
			lf, _ := coerce.ToFloat64(lv)
			rf, _ := coerce.ToFloat64(rv)
			return lf <= rf, nil
		}
	case json.Number:
		if strings.Contains(le.String(), ".") {
			li, _ := le.Float64()
			ri, _ := coerce.ToFloat64(rv)
			return li <= ri, nil
		} else {
			li, _ := le.Int64()
			ri, _ := coerce.ToInt64(rv)
			return li <= ri, nil
		}
	case string:
		if rt == reflect.String {
			rs, _ := coerce.ToString(rv)
			return le <= rs, nil
		}
	case bool:
		if rt == reflect.Bool {
			rb, _ := coerce.ToBool(rv)
			return le == rb, nil
		}
	case time.Time:
		switch re := rv.(type) {
		case time.Time:
			return le.Before(re) || le.Equal(re), nil
		default:
			t, err := coerce.ToDateTime(rv)
			if err == nil {
				return le.Before(t) || le.Equal(t), nil
			}
		}
	}

	return false, fmt.Errorf("cannot compare %s with %s", reflect.TypeOf(lv).String(), reflect.TypeOf(rv).String())
}

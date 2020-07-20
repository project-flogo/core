package ast

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/resolve"
	"strings"
)

// Those errors are indicate that resolver not able to found the attr or fields
var errorStrs = []string{
	"path not found",
	"unable to evaluate path",
	"failed to resolve variable",
	"failed to resolve Environment Variable",
	"failed to resolve Property",
	"failed to resolve Loop",
	"failed to resolve activity attr",
	"failed to resolve activity value",
	"not found in flow"}

type IsDefinedExpr struct {
	refExpr Expr
}

func (d *IsDefinedExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	return d.refExpr.Init(resolver, root)
}

func (d *IsDefinedExpr) Eval(scope data.Scope) (interface{}, error) {
	_, isDefine, err := isDefined(d.refExpr, scope)
	return isDefine, err
}

type GetValueExpr struct {
	refExpr   Expr
	valueExpr Expr
}

func (d *GetValueExpr) Init(resolver resolve.CompositeResolver, root bool) error {
	err := d.refExpr.Init(resolver, root)
	if err != nil {
		return err
	}
	return d.valueExpr.Init(resolver, root)

}

func (d *GetValueExpr) Eval(scope data.Scope) (interface{}, error) {
	v, isDefine, err := isDefined(d.refExpr, scope)
	if err != nil {
		return nil, err
	}
	if !isDefine {
		return d.valueExpr.Eval(scope)
	}
	return v, nil
}

func isDefined(expr Expr, scope data.Scope) (interface{}, bool, error) {
	v, err := expr.Eval(scope)
	if err != nil {
		msg := err.Error()
		if isNotFoundError(msg) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return v, v != nil, nil
}

func isNotFoundError(errStr string) bool {
	for _, s := range errorStrs {
		if strings.Contains(errStr, s) {
			return true
		}
	}
	return false
}

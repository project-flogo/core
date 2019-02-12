package expression

import (
	"fmt"

	"github.com/project-flogo/core/data/resolve"
)

var scriptFactoryCreator FactoryCreatorFunc

func SetScriptFactoryCreator(factory FactoryCreatorFunc) {
	if scriptFactoryCreator == nil {
		scriptFactoryCreator = factory
	}
}

func NewFactory(resolver resolve.CompositeResolver) Factory {
	var scriptFactory Factory
	if scriptFactoryCreator != nil {
		scriptFactory = scriptFactoryCreator(resolver)
	}

	return &factoryImpl{resolver: resolver, scriptExprFactory: scriptFactory}
}

type factoryImpl struct {
	resolver          resolve.CompositeResolver
	scriptExprFactory Factory
}

func (f *factoryImpl) NewExpr(exprStr string) (Expr, error) {

	if resolve.IsResolveExpr(exprStr) {

		resolution, err := f.resolver.GetResolution(exprStr)
		if err != nil {
			//todo if there is no resolver, should we do a simple scope resolution?
			return nil, err
		}
		if resolution.IsStatic() {
			val, _ := resolution.GetValue(nil)
			return &literalExpr{val: val}, nil
		} else {
			return &resolutionExpr{resolution: resolution}, nil
		}
	}

	// template expression
	if IsTemplateExpr(exprStr) {
		return NewTemplateExpr(exprStr, f)
	}

	l, isLiteral := GetLiteral(exprStr)

	if isLiteral {
		return &literalExpr{val: l}, nil
	}

	// script expression
	if f.scriptExprFactory != nil {
		expr, err := f.scriptExprFactory.NewExpr(exprStr)
		if err != nil {
			return nil, fmt.Errorf("unable to compile expression '%s': %s", exprStr, err.Error())
		}
		return expr, nil
	}

	return nil, fmt.Errorf("unable to compile expression '%s'", exprStr)
}

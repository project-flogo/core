package script

import (
	"fmt"

	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/expression/script/gocc/ast"
	"github.com/project-flogo/core/data/expression/script/gocc/errors"
	"github.com/project-flogo/core/data/expression/script/gocc/lexer"
	"github.com/project-flogo/core/data/expression/script/gocc/parser"
	"github.com/project-flogo/core/data/resolve"

	_ "github.com/project-flogo/core/data/expression/function/builtin"
)

func init() {
	expression.SetScriptFactoryCreator(NewExprFactory)
}

func NewExprFactory(resolver resolve.CompositeResolver) expression.Factory {
	return &factoryImpl{resolver: resolver}
}

type factoryImpl struct {
	resolver resolve.CompositeResolver
}

func (f *factoryImpl) NewExpr(exprStr string) (expression.Expr, error) {
	st, err := parse(exprStr)
	if err != nil {
		if gerr, ok := err.(*errors.Error); ok {

			if gerr.Err != nil {
				return nil, gerr.Err
			}
			//log details in debug
			return nil, fmt.Errorf("error parsing expression")
		}

		return nil, err
	}

	expr, ok := st.(ast.Expr)
	if ok {
		err := expr.Init(f.resolver, true)
		if err != nil {
			return nil, err
		}

		return expr, nil
	}
	return nil, fmt.Errorf("error parsing expression")
}

func parse(scriptExpr string) (interface{}, error) {
	lex := lexer.NewLexer([]byte(scriptExpr))
	p := parser.NewParser()
	st, err := p.Parse(lex)
	return st, err
}

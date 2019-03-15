package ast

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression/script/gocc/token"
	"github.com/project-flogo/core/data/resolve"
	"strconv"
)

func Concat(items ...interface{}) ([]interface{}, error) {
	var array []interface{}
	for _, item := range items {
		if item == nil {
			continue
		}
		switch t := item.(type) {
		case string:
			array = append(array, t)
		case *token.Token:
			array = append(array, string(t.Lit))
		case Expr:
			array = append(array, t)
		case []interface{}:
			array = append(array, t...)
		default:
			array = append(array, t)
		}
	}
	return array, nil
}

func Indexer(indexer interface{}) (interface{}, error) {
	switch t := indexer.(type) {
	case Expr:
		return &arrayIndexer{t}, nil
	}
	return nil, fmt.Errorf("invalid array indexer")
}

type arrayIndexer struct {
	expr Expr
}

func (e *arrayIndexer) ToRef(resolver resolve.CompositeResolver, root bool, scope data.Scope) (string, error) {
	if err := e.expr.Init(resolver, root); err != nil {
		return "", nil
	}

	v, err := e.expr.Eval(scope)
	if err != nil {
		return "", fmt.Errorf("eval array index expression error: %s", err.Error())
	}
	index, err := coerce.ToInt(v)
	if err != nil {
		return "", fmt.Errorf("array index [%s] must be int", v)
	}

	return "[" + strconv.Itoa(index) + "]", nil
}

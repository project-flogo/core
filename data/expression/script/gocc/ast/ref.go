package ast

import (
	"fmt"
	"github.com/project-flogo/core/data/expression/script/gocc/token"
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
		}
	}
	return array, nil
}

func Key(indexer interface{}) (interface{}, error) {
	switch t := indexer.(type) {
	case Expr:
		return &keyIndexExpr{t}, nil
	}
	return nil, fmt.Errorf("invalid array indexer")
}

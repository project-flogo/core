package ast

import (
	"fmt"
	"github.com/project-flogo/core/data/expression/script/gocc/token"
	"reflect"
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
		}
	}
	return array, nil
}

type indexserRef struct {
	expr Expr
}

func ConcatIndexer(fist interface{}, second interface{}, third interface{}) (interface{}, error) {
	fmt.Println(string(fist.(*token.Token).Lit))
	fmt.Println(reflect.TypeOf(second))

	fmt.Println(string(third.((*token.Token)).Lit))

	switch t := second.(type) {
	case Expr:
		return &indexserRef{t}, nil
	}

	return nil, fmt.Errorf("indexer not find expr")
}

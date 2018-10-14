package ast

import "github.com/project-flogo/core/data/expression/script/gocc/token"

func Concat(items ...interface{}) (*token.Token, error) {
	s := ""
	for _, item := range items {
		if item == nil {
			continue
		}
		switch t := item.(type) {
		case string:
			s += t
		case *token.Token:
			s += string(t.Lit)
		}
	}
	return &token.Token{Lit: []byte(s)}, nil
}

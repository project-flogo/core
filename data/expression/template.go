package expression

import (
	"bytes"
	"strings"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
)

func IsTemplateExpr(exprStr string) bool {
	//todo fix this
	return strings.Contains(exprStr, "{{")
}

// NewTemplateExpr creates a new Template Expr, the template only supports
// non-nested double braces as embedded expression tokens {{ }}
func NewTemplateExpr(exprStr string, factory Factory) (Expr, error) {

	exprs, err := parse(exprStr, factory)
	if err != nil {
		return nil, err
	}

	return &templateExpr{exprs: exprs}, nil
}

type templateExpr struct {
	exprs []Expr
}

func (s *templateExpr) Eval(scope data.Scope) (interface{}, error) {

	var buffer bytes.Buffer

	for _, expr := range s.exprs {
		v, err := expr.Eval(scope)
		if err != nil {
			return nil, err
		}
		str, err := coerce.ToString(v)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(str)
	}

	return buffer.String(), nil
}

//stringExpr simple wrapper for a string
type stringExpr struct {
	str string
}

func (s *stringExpr) Eval(scope data.Scope) (interface{}, error) {
	return s.str, nil
}

func parse(s string, f Factory) ([]Expr, error) {

	strLen := len(s)

	var ss []Expr
	start := 0
	for i := 0; i < strLen; i++ {
		if s[i] == '{' {
			if s[i+1] == '{' {
				j := i + 2
				newStart := j
				closed := false
				for j < len(s) {
					if s[j] == '}' {
						if j+1 < strLen && s[j+1] == '}' {
							if i > 0 {
								ss = append(ss, &stringExpr{str: s[start:i]})
							}
							e, err := f.NewExpr(s[newStart:j])
							if err != nil {
								return nil, err
							}
							ss = append(ss, e)
							i = j + 2
							start = i
							closed = true
							break
						}
					}
					j++
				}

				if !closed {
					//todo should we throw an error?
					ss = append(ss, &stringExpr{str: s[start:]})
					start = len(s)
				}
			}
		}
	}

	if start != len(s) {
		ss = append(ss, &stringExpr{str: s[start:]})
	}

	return ss, nil
}

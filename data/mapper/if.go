package mapper

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression"
	"strings"
)

const (
	If     = "@if"
	Else   = "@else"
	ElseIf = "@elseIf"
)

type IfElseMapper struct {
	IfExpr     *IfElseExpr
	ElseExpr   *IfElseExpr
	ElseIfExpr []*IfElseExpr
}

func (fs *IfElseMapper) Eval(scope data.Scope) (interface{}, error) {
	ifExpr := fs.IfExpr
	if ifExpr.Condition == nil {
		return nil, fmt.Errorf("if or elseif must have condition expression")
	}
	ok, err := ifExpr.EvalCondition(scope)
	if err != nil {
		return nil, err
	}
	if ok {
		// if is true
		if ifExpr.object != nil {
			return ifExpr.object.Eval(scope)
		}
		return nil, nil
	} else {
		//go to else if
		if len(fs.ElseIfExpr) > 0 {
			for _, elseIfExr := range fs.ElseIfExpr {
				ok, err := elseIfExr.EvalCondition(scope)
				if err != nil {
					return nil, err
				}
				if ok {
					if elseIfExr.object != nil {
						return elseIfExr.object.Eval(scope)
					}
					return nil, nil

				}
			}
		}
	}
	//go to else
	if fs.ElseExpr != nil && fs.ElseExpr.object != nil {
		return fs.ElseExpr.object.Eval(scope)
	}
	return nil, nil
}

type IfElseExpr struct {
	Condition expression.Expr
	// Object mapper
	object expression.Expr
}

func (f *IfElseExpr) EvalCondition(scope data.Scope) (bool, error) {
	if f.Condition != nil {
		ifCondition, err := f.Condition.Eval(scope)
		if err != nil {
			return false, err
		}
		ok, _ := coerce.ToBool(ifCondition)
		return ok, nil
	}
	return false, nil
}

func hasIfElse(value interface{}) bool {
	switch t := value.(type) {
	case map[string]interface{}:
		for k, _ := range t {
			if strings.HasPrefix(k, If) || strings.HasPrefix(k, Else) {
				return true
			}
		}
		return false
	default:
		obj, _ := coerce.ToObject(value)
		if obj != nil {
			for k, _ := range obj {
				if strings.HasPrefix(k, If) || strings.HasPrefix(k, Else) {
					return true
				}
			}
		}
		return false
	}
}

func newIfElseMapper(value interface{}, ef expression.Factory) (expression.Expr, error) {
	switch t := value.(type) {
	case map[string]interface{}:
		ifMapper := &IfElseMapper{}
		for k, v := range t {
			if strings.HasPrefix(k, If) || strings.HasPrefix(k, Else) {
				//if expr
				ifCondition, err := getIfCondition(k, ef)
				if err != nil {
					return nil, err
				}

				mapper, err := NewObjectMapper(v, ef)
				if err != nil {
					return nil, err
				}

				expr := &IfElseExpr{
					Condition: ifCondition,
					object:    mapper,
				}

				if strings.HasPrefix(k, If) {
					ifMapper.IfExpr = expr
				} else if strings.HasPrefix(k, ElseIf) {
					ifMapper.ElseIfExpr = append(ifMapper.ElseIfExpr, expr)
				} else if strings.HasPrefix(k, Else) {
					ifMapper.ElseExpr = expr
				}
			}
		}

		if ifMapper.IfExpr != nil {
			return ifMapper, nil
		} else {
			//Not an if else mapper
			return NewObjectMapper(value, ef)
		}
	default:
		return NewObjectMapper(value, ef)
	}
}

func getIfCondition(key string, ef expression.Factory) (expression.Expr, error) {
	start := strings.Index(key, "(")
	end := strings.LastIndex(key, ")")
	if start > 0 && end > 0 {
		return ef.NewExpr(key[start+1 : end])
	}
	return nil, nil
}

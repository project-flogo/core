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
	ElseIf = "@elseif"
)

type IfElseMapper struct {
	IfExpr     *ifElseExpr
	ElseExpr   *ifElseExpr
	ElseIfExpr []*ifElseExpr
}

func (fs *IfElseMapper) Eval(scope data.Scope) (interface{}, error) {
	ifExpr := fs.IfExpr
	if ifExpr.condition == nil {
		return nil, fmt.Errorf("if mapping must have condition expression")
	}
	ok, err := ifExpr.EvalCondition(scope)
	if err != nil {
		return nil, err
	}
	if ok {
		// if is true
		if ifExpr.body != nil {
			return ifExpr.body.Eval(scope)
		}
		return nil, nil
	} else {
		//go to else if
		if len(fs.ElseIfExpr) > 0 {
			for _, elseIfExr := range fs.ElseIfExpr {
				if elseIfExr.condition == nil {
					return nil, fmt.Errorf("elseif mapping must have condition expression")
				}
				ok, err := elseIfExr.EvalCondition(scope)
				if err != nil {
					return nil, err
				}
				if ok {
					if elseIfExr.body != nil {
						return elseIfExr.body.Eval(scope)
					}
					return nil, nil

				}
			}
		}
	}
	//go to else
	if fs.ElseExpr != nil && fs.ElseExpr.body != nil {
		return fs.ElseExpr.body.Eval(scope)
	}
	return nil, nil
}

type ifElseExpr struct {
	condition expression.Expr
	body      expression.Expr
}

// EvalCondition Execute the condition expression of if/elseif
func (f *ifElseExpr) EvalCondition(scope data.Scope) (bool, error) {
	if f.condition != nil {
		ifCondition, err := f.condition.Eval(scope)
		if err != nil {
			return false, fmt.Errorf("eval if/else condition [%s] error: %s", f.condition, err.Error())
		}
		ok, _ := coerce.ToBool(ifCondition)
		return ok, nil
	}
	return false, nil
}

// isIfElseMapping check to see if the mapping is an if/else mapping
func isIfElseMapping(value interface{}) bool {
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

// createIfElseMapper  create if/else mapper
func createIfElseMapper(value interface{}, ef expression.Factory) (expression.Expr, error) {
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
				expr := &ifElseExpr{
					condition: ifCondition,
					body:      mapper,
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
			//Not if/else mapper
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

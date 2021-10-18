package mapper

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression"
	"strings"
)

const (
	Conditional = "@conditional"
	Otherwise   = "@otherwise"
)

type ConditionalMapper struct {
	conditions []*conditionalExpr
	otherwise  expression.Expr
}

func (fs *ConditionalMapper) Eval(scope data.Scope) (interface{}, error) {
	if len(fs.conditions) > 0 {
		for _, condition := range fs.conditions {
			keyResult, err := condition.key.Eval(scope)
			if err != nil {
				return nil, err
			}
			ok, _ := coerce.ToBool(keyResult)
			if ok {
				return condition.mapper.Eval(scope)
			}
		}
		if fs.otherwise != nil {
			return fs.otherwise.Eval(scope)
		}
	}
	return nil, nil
}

type conditionalExpr struct {
	key    expression.Expr
	mapper expression.Expr
}

// IsConditionalMapping check to see if the mapping is an conditional mapping
func IsConditionalMapping(value interface{}) bool {
	switch t := value.(type) {
	case map[string]interface{}:
		for k, _ := range t {
			if strings.HasPrefix(k, Conditional) {
				return true
			}
		}
		return false
	default:
		obj, _ := coerce.ToObject(value)
		if obj != nil {
			for k, _ := range obj {
				if strings.HasPrefix(k, Conditional) {
					return true
				}
			}
		}
		return false
	}
}

// createConditionalMapper  create conditional mapper
func createConditionalMapper(value interface{}, ef expression.Factory) (expression.Expr, error) {
	switch t := value.(type) {
	case map[string]interface{}:
		ifMapper := &ConditionalMapper{}
		for k, v := range t {
			if strings.HasPrefix(k, Conditional) {
				var exprPrefix string
				conditionArg, exit := getConditionArgument(k)
				if exit {
					exprPrefix = conditionArg
				}
				conditionalArray, err := coerce.ToArray(v)
				if err != nil {
					return nil, err
				}

				for _, conditionElems := range conditionalArray {
					element, err := coerce.ToObject(conditionElems)
					if err != nil {
						return nil, fmt.Errorf("connditonal mapper element must be key-value pair: %s", err.Error())
					}

					if len(element) > 1 {
						return nil, fmt.Errorf("connditonal mapper element must be single key-value pair: %s", err.Error())
					}

					var keyExpr, valueExpr expression.Expr
					for exprK, exprV := range element {
						valueExpr, err = NewObjectMapper(exprV, ef)
						if err != nil {
							return nil, err
						}
						if Otherwise == exprK {
							ifMapper.otherwise = valueExpr
							continue
						} else {
							if len(exprK) == 0 {
								exprStr, _ := coerce.ToString(exprV)
								return nil, fmt.Errorf("no condition expression set for [%s]", exprStr)
							} else if len(exprPrefix) > 0 {
								keyExpr, err = ef.NewExpr(exprPrefix + " " + exprK)
								if err != nil {
									return nil, fmt.Errorf("creating conndtion expr error: %s", err.Error())
								}
							} else {
								keyExpr, err = ef.NewExpr(exprK)
								if err != nil {
									return nil, fmt.Errorf("creating conndtion expr error: %s", err.Error())
								}
							}

						}
						conditionExpr := &conditionalExpr{
							key:    keyExpr,
							mapper: valueExpr,
						}
						ifMapper.conditions = append(ifMapper.conditions, conditionExpr)
					}
				}
			}
		}
		if ifMapper.conditions != nil {
			return ifMapper, nil
		} else {
			//Not conditional mapper
			return NewObjectMapper(value, ef)
		}
	default:
		return NewObjectMapper(value, ef)
	}
}

func getConditionArgument(key string) (string, bool) {
	start := strings.Index(key, "(")
	end := strings.LastIndex(key, ")")
	if start > 0 && end > 0 {
		return key[start+1 : end], true
	}
	return "", false
}

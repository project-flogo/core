package trace

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression"
)

type TagDef struct {
	Name     string
	Value    string
	nameExpr expression.Expr
	expr     expression.Expr
}

// ParseTagDefs parses the nested tags configuration structure and compiles expression values.
// Expected JSON structure: {"mapping": {"tags": [{"name": "k1", "value": "v1"}, ...]}}
func ParseTagDefs(raw interface{}, ef expression.Factory) []*TagDef {
	if raw == nil {
		return nil
	}

	tagsMap, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}

	mapping, ok := tagsMap["mapping"].(map[string]interface{})
	if !ok {
		return nil
	}

	tagsList, ok := mapping["tags"].([]interface{})
	if !ok {
		return nil
	}

	var defs []*TagDef
	for _, tag := range tagsList {
		tagMap, ok := tag.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := tagMap["name"].(string)
		value, _ := tagMap["value"].(string)
		if name == "" {
			continue
		}

		td := &TagDef{Name: name, Value: value}

		if len(name) > 0 && name[0] == '=' && ef != nil {
			nameExpr, err := ef.NewExpr(name[1:])
			if err == nil {
				td.nameExpr = nameExpr
			}
		}

		if len(value) > 0 && value[0] == '=' && ef != nil {
			expr, err := ef.NewExpr(value[1:])
			if err == nil {
				td.expr = expr
			}
		}

		defs = append(defs, td)
	}

	return defs
}

// ResolveTagDefs resolves tag definitions at runtime, evaluating expressions against the provided scope.
func ResolveTagDefs(defs []*TagDef, scope data.Scope) map[string]interface{} {
	if len(defs) == 0 {
		return nil
	}

	result := make(map[string]interface{}, len(defs))
	for _, td := range defs {
		key := td.Name
		if td.nameExpr != nil && scope != nil {
			val, err := td.nameExpr.Eval(scope)
			if err == nil {
				if s, ok := val.(string); ok && s != "" {
					key = s
				}
			}
		}

		if td.expr != nil && scope != nil {
			val, err := td.expr.Eval(scope)
			if err == nil {
				result[key] = val
				continue
			}
		}
		result[key] = td.Value
	}

	return result
}

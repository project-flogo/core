package trace

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/mapper"
)

type TagDef struct {
	Name     string
	Value    string
	nameExpr expression.Expr
	expr     expression.Expr
}

type TagDefs struct {
	staticDefs  []*TagDef
	dynamicExpr expression.Expr
}

func (td *TagDefs) IsEmpty() bool {
	return td == nil || (len(td.staticDefs) == 0 && td.dynamicExpr == nil)
}

// ParseTagDefs parses the nested tags configuration structure and compiles expression values.
// Supports both static tags: {"mapping": {"tags": [{"name": "k1", "value": "v1"}, ...]}}
// and dynamic tags with @conditional/@foreach: {"mapping": {"tags": {"@conditional": [...]}}}
func ParseTagDefs(raw interface{}, ef expression.Factory) *TagDefs {
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

	tagsValue := mapping["tags"]
	if tagsValue == nil {
		return nil
	}

	if tagsList, ok := tagsValue.([]interface{}); ok {
		return &TagDefs{staticDefs: parseStaticTags(tagsList, ef)}
	}

	if tagsObj, ok := tagsValue.(map[string]interface{}); ok && ef != nil {
		expr, err := mapper.NewObjectMapper(tagsObj, ef)
		if err == nil && expr != nil {
			return &TagDefs{dynamicExpr: expr}
		}
	}

	return nil
}

func parseStaticTags(tagsList []interface{}, ef expression.Factory) []*TagDef {
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
func ResolveTagDefs(td *TagDefs, scope data.Scope) map[string]interface{} {
	if td.IsEmpty() {
		return nil
	}

	if td.dynamicExpr != nil {
		return resolveDynamicTags(td.dynamicExpr, scope)
	}

	return resolveStaticTags(td.staticDefs, scope)
}

func resolveStaticTags(defs []*TagDef, scope data.Scope) map[string]interface{} {
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

func resolveDynamicTags(expr expression.Expr, scope data.Scope) map[string]interface{} {
	val, err := expr.Eval(scope)
	if err != nil || val == nil {
		return nil
	}

	result := make(map[string]interface{})

	switch resolved := val.(type) {
	case []interface{}:
		for _, item := range resolved {
			if tagMap, ok := item.(map[string]interface{}); ok {
				name, _ := tagMap["name"].(string)
				value := tagMap["value"]
				if name != "" {
					result[name] = value
				}
			}
		}
	case map[string]interface{}:
		for k, v := range resolved {
			result[k] = v
		}
	}

	return result
}

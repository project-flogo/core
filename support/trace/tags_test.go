package trace

import (
	"fmt"
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression"
	"github.com/stretchr/testify/assert"
)

type testExprFactory struct{}

func (f *testExprFactory) NewExpr(exprStr string) (expression.Expr, error) {
	return &testExpr{str: exprStr}, nil
}

type failExprFactory struct{}

func (f *failExprFactory) NewExpr(exprStr string) (expression.Expr, error) {
	return nil, fmt.Errorf("failed to compile expression: %s", exprStr)
}

type testExpr struct {
	str string
}

func (e *testExpr) Eval(scope data.Scope) (interface{}, error) {
	val, ok := scope.GetValue(e.str)
	if !ok {
		return nil, fmt.Errorf("value not found: %s", e.str)
	}
	return val, nil
}

// --- ParseTagDefs tests ---

func TestParseTagDefs_NilInput(t *testing.T) {
	result := ParseTagDefs(nil, &testExprFactory{})
	assert.Nil(t, result)
}

func TestParseTagDefs_InvalidTopLevel(t *testing.T) {
	result := ParseTagDefs("not a map", &testExprFactory{})
	assert.Nil(t, result)
}

func TestParseTagDefs_MissingMapping(t *testing.T) {
	raw := map[string]interface{}{
		"other": "value",
	}
	result := ParseTagDefs(raw, &testExprFactory{})
	assert.Nil(t, result)
}

func TestParseTagDefs_MappingNotMap(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": "not a map",
	}
	result := ParseTagDefs(raw, &testExprFactory{})
	assert.Nil(t, result)
}

func TestParseTagDefs_MissingTags(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": map[string]interface{}{
			"other": "value",
		},
	}
	result := ParseTagDefs(raw, &testExprFactory{})
	assert.Nil(t, result)
}

func TestParseTagDefs_TagsNotSlice(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": map[string]interface{}{
			"tags": "not a slice",
		},
	}
	result := ParseTagDefs(raw, &testExprFactory{})
	assert.Nil(t, result)
}

func TestParseTagDefs_EmptyTags(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": map[string]interface{}{
			"tags": []interface{}{},
		},
	}
	result := ParseTagDefs(raw, &testExprFactory{})
	assert.Nil(t, result)
}

func TestParseTagDefs_SkipsInvalidTagEntry(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": map[string]interface{}{
			"tags": []interface{}{
				"not a map",
			},
		},
	}
	result := ParseTagDefs(raw, &testExprFactory{})
	assert.Nil(t, result)
}

func TestParseTagDefs_SkipsEmptyName(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": map[string]interface{}{
			"tags": []interface{}{
				map[string]interface{}{"name": "", "value": "v1"},
			},
		},
	}
	result := ParseTagDefs(raw, &testExprFactory{})
	assert.Nil(t, result)
}

func TestParseTagDefs_StaticTag(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": map[string]interface{}{
			"tags": []interface{}{
				map[string]interface{}{"name": "env", "value": "prod"},
			},
		},
	}
	result := ParseTagDefs(raw, &testExprFactory{})
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "env", result[0].Name)
	assert.Equal(t, "prod", result[0].Value)
	assert.Nil(t, result[0].nameExpr)
	assert.Nil(t, result[0].expr)
}

func TestParseTagDefs_MultipleTags(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": map[string]interface{}{
			"tags": []interface{}{
				map[string]interface{}{"name": "k1", "value": "v1"},
				map[string]interface{}{"name": "k2", "value": "v2"},
				map[string]interface{}{"name": "k3", "value": "v3"},
			},
		},
	}
	result := ParseTagDefs(raw, &testExprFactory{})
	assert.Equal(t, 3, len(result))
	assert.Equal(t, "k1", result[0].Name)
	assert.Equal(t, "k2", result[1].Name)
	assert.Equal(t, "k3", result[2].Name)
}

func TestParseTagDefs_ValueExpression(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": map[string]interface{}{
			"tags": []interface{}{
				map[string]interface{}{"name": "region", "value": "=myVar"},
			},
		},
	}
	result := ParseTagDefs(raw, &testExprFactory{})
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "region", result[0].Name)
	assert.Equal(t, "=myVar", result[0].Value)
	assert.Nil(t, result[0].nameExpr)
	assert.NotNil(t, result[0].expr)
}

func TestParseTagDefs_NameExpression(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": map[string]interface{}{
			"tags": []interface{}{
				map[string]interface{}{"name": "=dynKey", "value": "staticVal"},
			},
		},
	}
	result := ParseTagDefs(raw, &testExprFactory{})
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "=dynKey", result[0].Name)
	assert.Equal(t, "staticVal", result[0].Value)
	assert.NotNil(t, result[0].nameExpr)
	assert.Nil(t, result[0].expr)
}

func TestParseTagDefs_BothExpressions(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": map[string]interface{}{
			"tags": []interface{}{
				map[string]interface{}{"name": "=dynKey", "value": "=dynVal"},
			},
		},
	}
	result := ParseTagDefs(raw, &testExprFactory{})
	assert.Equal(t, 1, len(result))
	assert.NotNil(t, result[0].nameExpr)
	assert.NotNil(t, result[0].expr)
}

func TestParseTagDefs_NilFactory(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": map[string]interface{}{
			"tags": []interface{}{
				map[string]interface{}{"name": "=dynKey", "value": "=dynVal"},
			},
		},
	}
	result := ParseTagDefs(raw, nil)
	assert.Equal(t, 1, len(result))
	assert.Nil(t, result[0].nameExpr)
	assert.Nil(t, result[0].expr)
}

func TestParseTagDefs_FactoryError(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": map[string]interface{}{
			"tags": []interface{}{
				map[string]interface{}{"name": "=badExpr", "value": "=badExpr"},
			},
		},
	}
	result := ParseTagDefs(raw, &failExprFactory{})
	assert.Equal(t, 1, len(result))
	assert.Nil(t, result[0].nameExpr)
	assert.Nil(t, result[0].expr)
}

func TestParseTagDefs_MixedValidAndInvalid(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": map[string]interface{}{
			"tags": []interface{}{
				map[string]interface{}{"name": "valid", "value": "v1"},
				"invalid entry",
				map[string]interface{}{"name": "", "value": "v2"},
				map[string]interface{}{"name": "also_valid", "value": "v3"},
			},
		},
	}
	result := ParseTagDefs(raw, &testExprFactory{})
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "valid", result[0].Name)
	assert.Equal(t, "also_valid", result[1].Name)
}

func TestParseTagDefs_TagWithNoValue(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": map[string]interface{}{
			"tags": []interface{}{
				map[string]interface{}{"name": "myTag"},
			},
		},
	}
	result := ParseTagDefs(raw, &testExprFactory{})
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "myTag", result[0].Name)
	assert.Equal(t, "", result[0].Value)
}

// --- ResolveTagDefs tests ---

func TestResolveTagDefs_NilDefs(t *testing.T) {
	result := ResolveTagDefs(nil, nil)
	assert.Nil(t, result)
}

func TestResolveTagDefs_EmptyDefs(t *testing.T) {
	result := ResolveTagDefs([]*TagDef{}, nil)
	assert.Nil(t, result)
}

func TestResolveTagDefs_StaticTags(t *testing.T) {
	defs := []*TagDef{
		{Name: "env", Value: "prod"},
		{Name: "region", Value: "us-east"},
	}
	result := ResolveTagDefs(defs, nil)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "prod", result["env"])
	assert.Equal(t, "us-east", result["region"])
}

func TestResolveTagDefs_ValueExpression(t *testing.T) {
	scope := data.NewSimpleScope(map[string]interface{}{
		"myVar": "resolved-value",
	}, nil)

	defs := []*TagDef{
		{Name: "key", Value: "=myVar", expr: &testExpr{str: "myVar"}},
	}
	result := ResolveTagDefs(defs, scope)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "resolved-value", result["key"])
}

func TestResolveTagDefs_NameExpression(t *testing.T) {
	scope := data.NewSimpleScope(map[string]interface{}{
		"dynKey": "resolved-key",
	}, nil)

	defs := []*TagDef{
		{Name: "=dynKey", Value: "staticVal", nameExpr: &testExpr{str: "dynKey"}},
	}
	result := ResolveTagDefs(defs, scope)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "staticVal", result["resolved-key"])
}

func TestResolveTagDefs_BothExpressions(t *testing.T) {
	scope := data.NewSimpleScope(map[string]interface{}{
		"dynKey": "resolved-key",
		"dynVal": "resolved-value",
	}, nil)

	defs := []*TagDef{
		{
			Name:     "=dynKey",
			Value:    "=dynVal",
			nameExpr: &testExpr{str: "dynKey"},
			expr:     &testExpr{str: "dynVal"},
		},
	}
	result := ResolveTagDefs(defs, scope)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "resolved-value", result["resolved-key"])
}

func TestResolveTagDefs_ExprEvalError_FallsBackToStaticValue(t *testing.T) {
	scope := data.NewSimpleScope(map[string]interface{}{}, nil)

	defs := []*TagDef{
		{Name: "key", Value: "fallback", expr: &testExpr{str: "missing"}},
	}
	result := ResolveTagDefs(defs, scope)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "fallback", result["key"])
}

func TestResolveTagDefs_NameExprEvalError_UsesOriginalName(t *testing.T) {
	scope := data.NewSimpleScope(map[string]interface{}{}, nil)

	defs := []*TagDef{
		{Name: "originalKey", Value: "val", nameExpr: &testExpr{str: "missing"}},
	}
	result := ResolveTagDefs(defs, scope)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "val", result["originalKey"])
}

func TestResolveTagDefs_NilScope_SkipsExpressions(t *testing.T) {
	defs := []*TagDef{
		{
			Name:     "=dynKey",
			Value:    "=dynVal",
			nameExpr: &testExpr{str: "dynKey"},
			expr:     &testExpr{str: "dynVal"},
		},
	}
	result := ResolveTagDefs(defs, nil)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "=dynVal", result["=dynKey"])
}

func TestResolveTagDefs_NameExprResolvesToEmpty_UsesOriginalName(t *testing.T) {
	scope := data.NewSimpleScope(map[string]interface{}{
		"dynKey": "",
	}, nil)

	defs := []*TagDef{
		{Name: "originalKey", Value: "val", nameExpr: &testExpr{str: "dynKey"}},
	}
	result := ResolveTagDefs(defs, scope)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "val", result["originalKey"])
}

func TestResolveTagDefs_NameExprResolvesToNonString_UsesOriginalName(t *testing.T) {
	scope := data.NewSimpleScope(map[string]interface{}{
		"dynKey": 12345,
	}, nil)

	defs := []*TagDef{
		{Name: "originalKey", Value: "val", nameExpr: &testExpr{str: "dynKey"}},
	}
	result := ResolveTagDefs(defs, scope)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "val", result["originalKey"])
}

// --- Integration: ParseTagDefs + ResolveTagDefs ---

func TestParseAndResolve_EndToEnd(t *testing.T) {
	raw := map[string]interface{}{
		"mapping": map[string]interface{}{
			"tags": []interface{}{
				map[string]interface{}{"name": "static", "value": "hello"},
				map[string]interface{}{"name": "dynamic", "value": "=myVar"},
				map[string]interface{}{"name": "=keyVar", "value": "=valVar"},
			},
		},
	}

	ef := &testExprFactory{}
	defs := ParseTagDefs(raw, ef)
	assert.Equal(t, 3, len(defs))

	scope := data.NewSimpleScope(map[string]interface{}{
		"myVar":  "world",
		"keyVar": "resolved-key",
		"valVar": "resolved-val",
	}, nil)

	result := ResolveTagDefs(defs, scope)
	assert.Equal(t, 3, len(result))
	assert.Equal(t, "hello", result["static"])
	assert.Equal(t, "world", result["dynamic"])
	assert.Equal(t, "resolved-val", result["resolved-key"])
}

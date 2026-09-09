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
	assert.NotNil(t, result)
	assert.True(t, result.IsEmpty())
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
	assert.NotNil(t, result)
	assert.True(t, result.IsEmpty())
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
	assert.NotNil(t, result)
	assert.True(t, result.IsEmpty())
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
	assert.False(t, result.IsEmpty())
	assert.Equal(t, 1, len(result.staticDefs))
	assert.Equal(t, "env", result.staticDefs[0].Name)
	assert.Equal(t, "prod", result.staticDefs[0].Value)
	assert.Nil(t, result.staticDefs[0].nameExpr)
	assert.Nil(t, result.staticDefs[0].expr)
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
	assert.Equal(t, 3, len(result.staticDefs))
	assert.Equal(t, "k1", result.staticDefs[0].Name)
	assert.Equal(t, "k2", result.staticDefs[1].Name)
	assert.Equal(t, "k3", result.staticDefs[2].Name)
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
	assert.Equal(t, 1, len(result.staticDefs))
	assert.Equal(t, "region", result.staticDefs[0].Name)
	assert.Equal(t, "=myVar", result.staticDefs[0].Value)
	assert.Nil(t, result.staticDefs[0].nameExpr)
	assert.NotNil(t, result.staticDefs[0].expr)
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
	assert.Equal(t, 1, len(result.staticDefs))
	assert.Equal(t, "=dynKey", result.staticDefs[0].Name)
	assert.Equal(t, "staticVal", result.staticDefs[0].Value)
	assert.NotNil(t, result.staticDefs[0].nameExpr)
	assert.Nil(t, result.staticDefs[0].expr)
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
	assert.Equal(t, 1, len(result.staticDefs))
	assert.NotNil(t, result.staticDefs[0].nameExpr)
	assert.NotNil(t, result.staticDefs[0].expr)
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
	assert.Equal(t, 1, len(result.staticDefs))
	assert.Nil(t, result.staticDefs[0].nameExpr)
	assert.Nil(t, result.staticDefs[0].expr)
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
	assert.Equal(t, 1, len(result.staticDefs))
	assert.Nil(t, result.staticDefs[0].nameExpr)
	assert.Nil(t, result.staticDefs[0].expr)
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
	assert.Equal(t, 2, len(result.staticDefs))
	assert.Equal(t, "valid", result.staticDefs[0].Name)
	assert.Equal(t, "also_valid", result.staticDefs[1].Name)
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
	assert.Equal(t, 1, len(result.staticDefs))
	assert.Equal(t, "myTag", result.staticDefs[0].Name)
	assert.Equal(t, "", result.staticDefs[0].Value)
}

// --- TagDefs.IsEmpty tests ---

func TestTagDefs_IsEmpty_Nil(t *testing.T) {
	var td *TagDefs
	assert.True(t, td.IsEmpty())
}

func TestTagDefs_IsEmpty_EmptyStatic(t *testing.T) {
	td := &TagDefs{}
	assert.True(t, td.IsEmpty())
}

func TestTagDefs_IsEmpty_WithStatic(t *testing.T) {
	td := &TagDefs{staticDefs: []*TagDef{{Name: "k", Value: "v"}}}
	assert.False(t, td.IsEmpty())
}

func TestTagDefs_IsEmpty_WithDynamic(t *testing.T) {
	td := &TagDefs{dynamicExpr: &testExpr{str: "test"}}
	assert.False(t, td.IsEmpty())
}

// --- ResolveTagDefs tests ---

func TestResolveTagDefs_NilDefs(t *testing.T) {
	result := ResolveTagDefs(nil, nil)
	assert.Nil(t, result)
}

func TestResolveTagDefs_EmptyDefs(t *testing.T) {
	result := ResolveTagDefs(&TagDefs{}, nil)
	assert.Nil(t, result)
}

func TestResolveTagDefs_StaticTags(t *testing.T) {
	td := &TagDefs{staticDefs: []*TagDef{
		{Name: "env", Value: "prod"},
		{Name: "region", Value: "us-east"},
	}}
	result := ResolveTagDefs(td, nil)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "prod", result["env"])
	assert.Equal(t, "us-east", result["region"])
}

func TestResolveTagDefs_ValueExpression(t *testing.T) {
	scope := data.NewSimpleScope(map[string]interface{}{
		"myVar": "resolved-value",
	}, nil)

	td := &TagDefs{staticDefs: []*TagDef{
		{Name: "key", Value: "=myVar", expr: &testExpr{str: "myVar"}},
	}}
	result := ResolveTagDefs(td, scope)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "resolved-value", result["key"])
}

func TestResolveTagDefs_NameExpression(t *testing.T) {
	scope := data.NewSimpleScope(map[string]interface{}{
		"dynKey": "resolved-key",
	}, nil)

	td := &TagDefs{staticDefs: []*TagDef{
		{Name: "=dynKey", Value: "staticVal", nameExpr: &testExpr{str: "dynKey"}},
	}}
	result := ResolveTagDefs(td, scope)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "staticVal", result["resolved-key"])
}

func TestResolveTagDefs_BothExpressions(t *testing.T) {
	scope := data.NewSimpleScope(map[string]interface{}{
		"dynKey": "resolved-key",
		"dynVal": "resolved-value",
	}, nil)

	td := &TagDefs{staticDefs: []*TagDef{
		{
			Name:     "=dynKey",
			Value:    "=dynVal",
			nameExpr: &testExpr{str: "dynKey"},
			expr:     &testExpr{str: "dynVal"},
		},
	}}
	result := ResolveTagDefs(td, scope)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "resolved-value", result["resolved-key"])
}

func TestResolveTagDefs_ExprEvalError_FallsBackToStaticValue(t *testing.T) {
	scope := data.NewSimpleScope(map[string]interface{}{}, nil)

	td := &TagDefs{staticDefs: []*TagDef{
		{Name: "key", Value: "fallback", expr: &testExpr{str: "missing"}},
	}}
	result := ResolveTagDefs(td, scope)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "fallback", result["key"])
}

func TestResolveTagDefs_NameExprEvalError_UsesOriginalName(t *testing.T) {
	scope := data.NewSimpleScope(map[string]interface{}{}, nil)

	td := &TagDefs{staticDefs: []*TagDef{
		{Name: "originalKey", Value: "val", nameExpr: &testExpr{str: "missing"}},
	}}
	result := ResolveTagDefs(td, scope)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "val", result["originalKey"])
}

func TestResolveTagDefs_NilScope_SkipsExpressions(t *testing.T) {
	td := &TagDefs{staticDefs: []*TagDef{
		{
			Name:     "=dynKey",
			Value:    "=dynVal",
			nameExpr: &testExpr{str: "dynKey"},
			expr:     &testExpr{str: "dynVal"},
		},
	}}
	result := ResolveTagDefs(td, nil)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "=dynVal", result["=dynKey"])
}

func TestResolveTagDefs_NameExprResolvesToEmpty_UsesOriginalName(t *testing.T) {
	scope := data.NewSimpleScope(map[string]interface{}{
		"dynKey": "",
	}, nil)

	td := &TagDefs{staticDefs: []*TagDef{
		{Name: "originalKey", Value: "val", nameExpr: &testExpr{str: "dynKey"}},
	}}
	result := ResolveTagDefs(td, scope)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "val", result["originalKey"])
}

func TestResolveTagDefs_NameExprResolvesToNonString_UsesOriginalName(t *testing.T) {
	scope := data.NewSimpleScope(map[string]interface{}{
		"dynKey": 12345,
	}, nil)

	td := &TagDefs{staticDefs: []*TagDef{
		{Name: "originalKey", Value: "val", nameExpr: &testExpr{str: "dynKey"}},
	}}
	result := ResolveTagDefs(td, scope)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "val", result["originalKey"])
}

// --- Dynamic tags (resolve) tests ---

func TestResolveDynamicTags_ArrayResult(t *testing.T) {
	dynamicResult := []interface{}{
		map[string]interface{}{"name": "tag1", "value": "val1"},
		map[string]interface{}{"name": "tag2", "value": "val2"},
	}
	td := &TagDefs{dynamicExpr: &staticResultExpr{result: dynamicResult}}

	result := ResolveTagDefs(td, nil)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "val1", result["tag1"])
	assert.Equal(t, "val2", result["tag2"])
}

func TestResolveDynamicTags_MapResult(t *testing.T) {
	dynamicResult := map[string]interface{}{
		"tag1": "val1",
		"tag2": "val2",
	}
	td := &TagDefs{dynamicExpr: &staticResultExpr{result: dynamicResult}}

	result := ResolveTagDefs(td, nil)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "val1", result["tag1"])
	assert.Equal(t, "val2", result["tag2"])
}

func TestResolveDynamicTags_NilResult(t *testing.T) {
	td := &TagDefs{dynamicExpr: &staticResultExpr{result: nil}}
	result := ResolveTagDefs(td, nil)
	assert.Nil(t, result)
}

func TestResolveDynamicTags_EvalError(t *testing.T) {
	td := &TagDefs{dynamicExpr: &errorExpr{}}
	result := ResolveTagDefs(td, nil)
	assert.Nil(t, result)
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
	td := ParseTagDefs(raw, ef)
	assert.False(t, td.IsEmpty())
	assert.Equal(t, 3, len(td.staticDefs))

	scope := data.NewSimpleScope(map[string]interface{}{
		"myVar":  "world",
		"keyVar": "resolved-key",
		"valVar": "resolved-val",
	}, nil)

	result := ResolveTagDefs(td, scope)
	assert.Equal(t, 3, len(result))
	assert.Equal(t, "hello", result["static"])
	assert.Equal(t, "world", result["dynamic"])
	assert.Equal(t, "resolved-val", result["resolved-key"])
}

// --- Helper test types ---

type staticResultExpr struct {
	result interface{}
}

func (e *staticResultExpr) Eval(scope data.Scope) (interface{}, error) {
	return e.result, nil
}

type errorExpr struct{}

func (e *errorExpr) Eval(scope data.Scope) (interface{}, error) {
	return nil, fmt.Errorf("eval error")
}

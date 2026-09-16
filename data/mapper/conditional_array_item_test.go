package mapper

import (
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/stretchr/testify/assert"
)

// applyArrayMapping runs a single-field object mapping and returns the field value as an array.
// The `{"mapping": ...}` wrapper is the exact shape the FE emits for object/array fields, so this
// exercises the same path used by the Mapper activity, any activity input, and trigger mappings.
func applyArrayMapping(t *testing.T, field string, mapping interface{}, attrs map[string]interface{}) []interface{} {
	t.Helper()
	mappings := map[string]interface{}{
		field: map[string]interface{}{"mapping": mapping},
	}
	factory := NewFactory(resolver)
	m, err := factory.NewMapper(mappings)
	assert.Nil(t, err)
	assert.NotNil(t, m)

	scope := data.NewSimpleScope(attrs, nil)
	results, err := m.Apply(scope)
	assert.Nil(t, err)

	arr, ok := results[field].([]interface{})
	assert.True(t, ok, "expected field %q to be an array, got %T", field, results[field])
	return arr
}

func ids(t *testing.T, arr []interface{}) []string {
	t.Helper()
	out := make([]string, 0, len(arr))
	for _, item := range arr {
		obj, ok := item.(map[string]interface{})
		assert.True(t, ok, "array item is not an object: %T", item)
		id, _ := obj["id"].(string)
		out = append(out, id)
	}
	return out
}

// A conditional array item whose condition matches is included in the array.
func TestConditionalArrayItem_Match(t *testing.T) {
	mapping := []interface{}{
		map[string]interface{}{"id": "always1"},
		map[string]interface{}{
			"@conditional": []interface{}{
				map[string]interface{}{`"1"=="1"`: map[string]interface{}{"id": "if"}},
				map[string]interface{}{`"2"=="2"`: map[string]interface{}{"id": "elseif"}},
				map[string]interface{}{"@otherwise": map[string]interface{}{"id": "else"}},
			},
		},
		map[string]interface{}{"id": "always2"},
	}
	arr := applyArrayMapping(t, "Employees", mapping, map[string]interface{}{})
	assert.Equal(t, []string{"always1", "if", "always2"}, ids(t, arr))
}

// A conditional array item that matches nothing and has NO @otherwise is omitted entirely
// (not added as a null element). This is the core of FLOGO-19268.
func TestConditionalArrayItem_NoMatchNoOtherwise_Omitted(t *testing.T) {
	mapping := []interface{}{
		map[string]interface{}{"id": "always1"},
		map[string]interface{}{
			"@conditional": []interface{}{
				map[string]interface{}{`"1"=="2"`: map[string]interface{}{"id": "if"}},
				map[string]interface{}{`"3"=="4"`: map[string]interface{}{"id": "elseif"}},
			},
		},
		map[string]interface{}{"id": "always2"},
	}
	arr := applyArrayMapping(t, "Employees", mapping, map[string]interface{}{})
	assert.Equal(t, 2, len(arr), "unmatched conditional item must be omitted, not null")
	assert.Equal(t, []string{"always1", "always2"}, ids(t, arr))
	for _, item := range arr {
		assert.NotNil(t, item)
	}
}

// A conditional array item that matches nothing but has an @otherwise falls through to it.
func TestConditionalArrayItem_NoMatchWithOtherwise(t *testing.T) {
	mapping := []interface{}{
		map[string]interface{}{"id": "always1"},
		map[string]interface{}{
			"@conditional": []interface{}{
				map[string]interface{}{`"1"=="2"`: map[string]interface{}{"id": "if"}},
				map[string]interface{}{"@otherwise": map[string]interface{}{"id": "else"}},
			},
		},
	}
	arr := applyArrayMapping(t, "Employees", mapping, map[string]interface{}{})
	assert.Equal(t, []string{"always1", "else"}, ids(t, arr))
}

// Multiple conditional items mixed with always-on items; order is preserved and only
// the matching/otherwise items are emitted.
func TestConditionalArrayItem_MixedItems(t *testing.T) {
	mapping := []interface{}{
		map[string]interface{}{"id": "a"},
		map[string]interface{}{
			"@conditional": []interface{}{
				map[string]interface{}{`"1"=="1"`: map[string]interface{}{"id": "b"}}, // matches
			},
		},
		map[string]interface{}{
			"@conditional": []interface{}{
				map[string]interface{}{`"1"=="2"`: map[string]interface{}{"id": "skip"}}, // no match, no otherwise
			},
		},
		map[string]interface{}{"id": "c"},
	}
	arr := applyArrayMapping(t, "Employees", mapping, map[string]interface{}{})
	assert.Equal(t, []string{"a", "b", "c"}, ids(t, arr))
}

// Dynamic condition params (references resolved from scope) work the same as static ones.
func TestConditionalArrayItem_DynamicCondition(t *testing.T) {
	mapping := []interface{}{
		map[string]interface{}{"id": "keep"},
		map[string]interface{}{
			"@conditional": []interface{}{
				map[string]interface{}{`$.person.name == "abc"`: map[string]interface{}{"id": "matched-abc"}},
				map[string]interface{}{`$.person.name == "bcd"`: map[string]interface{}{"id": "matched-bcd"}},
			},
		},
	}

	// name == "abc" -> the first conditional item is included
	arr := applyArrayMapping(t, "Employees", mapping,
		map[string]interface{}{"person": map[string]interface{}{"name": "abc"}})
	assert.Equal(t, []string{"keep", "matched-abc"}, ids(t, arr))

	// name == "zzz" -> no branch matches, no otherwise -> conditional item omitted
	arr = applyArrayMapping(t, "Employees", mapping,
		map[string]interface{}{"person": map[string]interface{}{"name": "zzz"}})
	assert.Equal(t, []string{"keep"}, ids(t, arr))
}

// The @merge + @foreach + literal-array shape (Mapper1 in the sample app): a foreach-produced
// array is merged with a literal array that mixes always-on items and conditional items.
func TestConditionalArrayItem_InMergeWithForeach(t *testing.T) {
	mapping := map[string]interface{}{
		"@merge": []interface{}{
			map[string]interface{}{
				"@foreach($.src, loopVar)": map[string]interface{}{"=": "$loop"},
			},
			[]interface{}{
				map[string]interface{}{"id": "literal"},
				map[string]interface{}{
					"@conditional": []interface{}{
						map[string]interface{}{`"1"=="1"`: map[string]interface{}{"id": "cond-yes"}},
					},
				},
				map[string]interface{}{
					"@conditional": []interface{}{
						map[string]interface{}{`"1"=="2"`: map[string]interface{}{"id": "cond-no"}},
					},
				},
			},
		},
	}
	attrs := map[string]interface{}{
		"src": []interface{}{
			map[string]interface{}{"id": "s1"},
			map[string]interface{}{"id": "s2"},
		},
	}
	arr := applyArrayMapping(t, "Employees", mapping, attrs)
	assert.Equal(t, []string{"s1", "s2", "literal", "cond-yes"}, ids(t, arr))
}

// Ordering: branches are if / else-if.../ else in array order. The FIRST true branch wins and
// short-circuits, even when a later branch would also be true.
func TestConditionalArrayItem_ElseIfOrdering(t *testing.T) {
	mapping := []interface{}{
		map[string]interface{}{"id": "always"},
		map[string]interface{}{
			"@conditional": []interface{}{
				map[string]interface{}{`"1"=="2"`: map[string]interface{}{"id": "if"}},       // false
				map[string]interface{}{`"3"=="3"`: map[string]interface{}{"id": "elseif-1"}}, // true -> wins
				map[string]interface{}{`"4"=="4"`: map[string]interface{}{"id": "elseif-2"}}, // true but not reached
				map[string]interface{}{"@otherwise": map[string]interface{}{"id": "else"}},
			},
		},
	}
	arr := applyArrayMapping(t, "Employees", mapping, map[string]interface{}{})
	assert.Equal(t, []string{"always", "elseif-1"}, ids(t, arr))
}

// Adding more branches between the first (if) and the last (@otherwise) makes them extra
// else-if branches, matched in order.
func TestConditionalArrayItem_ExtraMiddleElseIfBranches(t *testing.T) {
	mapping := []interface{}{
		map[string]interface{}{
			"@conditional": []interface{}{
				map[string]interface{}{`$.x == "a"`: map[string]interface{}{"id": "if-a"}},
				map[string]interface{}{`$.x == "b"`: map[string]interface{}{"id": "elseif-b"}},
				map[string]interface{}{`$.x == "c"`: map[string]interface{}{"id": "elseif-c"}},
				map[string]interface{}{"@otherwise": map[string]interface{}{"id": "else"}},
			},
		},
	}
	assert.Equal(t, []string{"if-a"}, ids(t, applyArrayMapping(t, "Employees", mapping, map[string]interface{}{"x": "a"})))
	assert.Equal(t, []string{"elseif-b"}, ids(t, applyArrayMapping(t, "Employees", mapping, map[string]interface{}{"x": "b"})))
	assert.Equal(t, []string{"elseif-c"}, ids(t, applyArrayMapping(t, "Employees", mapping, map[string]interface{}{"x": "c"})))
	assert.Equal(t, []string{"else"}, ids(t, applyArrayMapping(t, "Employees", mapping, map[string]interface{}{"x": "z"})))
}

// Sanity: a plain literal array with no conditional items is unaffected (always emitted, no skips).
func TestLiteralArray_NoConditional_Unaffected(t *testing.T) {
	mapping := []interface{}{
		map[string]interface{}{"id": "1"},
		map[string]interface{}{"id": "2"},
		map[string]interface{}{"id": "3"},
	}
	arr := applyArrayMapping(t, "Employees", mapping, map[string]interface{}{})
	assert.Equal(t, []string{"1", "2", "3"}, ids(t, arr))
}

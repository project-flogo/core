package mapper

import (
	"encoding/json"
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConditionalPrimitive(t *testing.T) {

	testcases := []struct {
		Mapping  interface{}
		Data     []map[string]interface{}
		Expected []string
	}{
		{
			Mapping: map[string]interface{}{
				"@conditional": []map[string]interface{}{
					{
						"$.person.name == \"abc\"": "this is abc",
					},
					{
						"$.person.name == \"bcd\"": "this is bcd",
					},
					{
						"@otherwise": "this is ddd",
					},
				},
			},
			Data: []map[string]interface{}{
				{
					"id":   "abc",
					"name": "abc",
					"address": map[string]interface{}{
						"city":  "sugarLand",
						"state": "tx",
					},
				},
				{
					"id":   "bcd",
					"name": "bcd",
					"address": map[string]interface{}{
						"city":  "sugarLand",
						"state": "tx",
					},
				},
				{
					"id":   "ddd",
					"name": "dddd",
					"address": map[string]interface{}{
						"city":  "sugarLand",
						"state": "tx",
					},
				},
			},
			Expected: []string{"this is abc", "this is bcd", "this is ddd"},
		},
	}

	for _, tt := range testcases {
		assert.True(t, IsConditionalMapping(tt.Mapping))
		mappings := map[string]interface{}{"output": tt.Mapping}

		v, _ := json.Marshal(mappings)
		fmt.Println("====", string(v))
		factory := NewFactory(resolver)
		mapper, err := factory.NewMapper(mappings)
		assert.Nil(t, err)

		for i, input := range tt.Data {
			attrs := map[string]interface{}{"person": input}
			scope := data.NewSimpleScope(attrs, nil)
			results, err := mapper.Apply(scope)
			assert.Nil(t, err)
			assert.Equal(t, tt.Expected[i], results["output"])
		}

	}

}

func TestConditionalSwitchPrimitive(t *testing.T) {

	testcases := []struct {
		Mapping  interface{}
		Data     []map[string]interface{}
		Expected []string
	}{
		{
			Mapping: map[string]interface{}{
				"@conditional($.person.name)": []map[string]interface{}{
					{
						"== \"abc\"": "this is abc",
					},
					{
						"== \"bcd\"": "this is bcd",
					},
					{
						"@otherwise": "this is ddd",
					},
				},
			},
			Data: []map[string]interface{}{
				{
					"id":   "abc",
					"name": "abc",
					"address": map[string]interface{}{
						"city":  "sugarLand",
						"state": "tx",
					},
				},
				{
					"id":   "bcd",
					"name": "bcd",
					"address": map[string]interface{}{
						"city":  "sugarLand",
						"state": "tx",
					},
				},
				{
					"id":   "ddd",
					"name": "dddd",
					"address": map[string]interface{}{
						"city":  "sugarLand",
						"state": "tx",
					},
				},
			},
			Expected: []string{"this is abc", "this is bcd", "this is ddd"},
		},
	}

	for _, tt := range testcases {
		mappings := map[string]interface{}{"output": tt.Mapping}
		factory := NewFactory(resolver)
		mapper, err := factory.NewMapper(mappings)
		assert.Nil(t, err)

		for i, input := range tt.Data {
			attrs := map[string]interface{}{"person": input}
			scope := data.NewSimpleScope(attrs, nil)
			results, err := mapper.Apply(scope)
			assert.Nil(t, err)
			assert.Equal(t, tt.Expected[i], results["output"])
		}

	}

}

func TestConditionalObjectMapper(t *testing.T) {

	testcases := []struct {
		Mapping  interface{}
		Data     []map[string]interface{}
		Expected []string
	}{
		{
			Mapping: map[string]interface{}{
				"@conditional": []map[string]interface{}{
					{
						"$.person.name == \"abc\"": map[string]interface{}{
							"id":      "=$.person.id",
							"name":    "=$.person.id",
							"address": "=$.person.address",
						},
					},
					{
						"$.person.name == \"bcd\"": map[string]interface{}{
							"id":      "=$.person.id",
							"name":    "=$.person.id",
							"address": "=$.person.address",
						},
					},
					{
						"@otherwise": map[string]interface{}{
							"id":      "=$.person.id",
							"name":    "=$.person.id",
							"address": "=$.person.address",
						},
					},
				},
			},
			Data: []map[string]interface{}{
				{
					"id":   "abc",
					"name": "abc",
					"address": map[string]interface{}{
						"city":  "abcsugarLand",
						"state": "abctx",
					},
				},
				{
					"id":   "bcd",
					"name": "bcd",
					"address": map[string]interface{}{
						"city":  "bcdsugarLand",
						"state": "bcdtx",
					},
				},
				{
					"id":   "ddd",
					"name": "dddd",
					"address": map[string]interface{}{
						"city":  "dddsugarLand",
						"state": "dddtx",
					},
				},
			},
			Expected: []string{"abc", "bcd", "ddd"},
		},
	}

	for _, tt := range testcases {
		assert.True(t, IsConditionalMapping(tt.Mapping))
		mappings := map[string]interface{}{"output": tt.Mapping}
		factory := NewFactory(resolver)
		mapper, err := factory.NewMapper(mappings)
		assert.Nil(t, err)

		for i, input := range tt.Data {
			attrs := map[string]interface{}{"person": input}
			scope := data.NewSimpleScope(attrs, nil)
			results, err := mapper.Apply(scope)
			assert.Nil(t, err)
			assert.Equal(t, tt.Expected[i], results["output"].(map[string]interface{})["name"])
		}

	}

}
func TestConditionalShorterObjectMapper(t *testing.T) {

	testcases := []struct {
		Mapping  interface{}
		Data     []map[string]interface{}
		Expected []string
	}{
		{
			Mapping: map[string]interface{}{
				"@conditional($.person.name)": []map[string]interface{}{
					{
						" == \"abc\"": map[string]interface{}{
							"id":      "=$.person.id",
							"name":    "=$.person.id",
							"address": "=$.person.address",
						},
					},
					{
						" == \"bcd\"": map[string]interface{}{
							"id":      "=$.person.id",
							"name":    "=$.person.id",
							"address": "=$.person.address",
						},
					},
					{
						"@otherwise": map[string]interface{}{
							"id":      "=$.person.id",
							"name":    "=$.person.id",
							"address": "=$.person.address",
						},
					},
				},
			},
			Data: []map[string]interface{}{
				{
					"id":   "abc",
					"name": "abc",
					"address": map[string]interface{}{
						"city":  "abcsugarLand",
						"state": "abctx",
					},
				},
				{
					"id":   "bcd",
					"name": "bcd",
					"address": map[string]interface{}{
						"city":  "bcdsugarLand",
						"state": "bcdtx",
					},
				},
				{
					"id":   "ddd",
					"name": "dddd",
					"address": map[string]interface{}{
						"city":  "dddsugarLand",
						"state": "dddtx",
					},
				},
			},
			Expected: []string{"abc", "bcd", "ddd"},
		},
	}

	for _, tt := range testcases {
		assert.True(t, IsConditionalMapping(tt.Mapping))
		mappings := map[string]interface{}{"output": tt.Mapping}
		factory := NewFactory(resolver)
		mapper, err := factory.NewMapper(mappings)
		assert.Nil(t, err)

		for i, input := range tt.Data {
			attrs := map[string]interface{}{"person": input}
			scope := data.NewSimpleScope(attrs, nil)
			results, err := mapper.Apply(scope)
			assert.Nil(t, err)
			assert.Equal(t, tt.Expected[i], results["output"].(map[string]interface{})["name"])
		}

	}

}

func TestConditionalArrayMapper(t *testing.T) {

	testcases := []struct {
		Mapping  interface{}
		Data     []map[string]interface{}
		Expected []string
	}{
		{
			Mapping: map[string]interface{}{
				"@conditional": []map[string]interface{}{
					{
						"$.person.name == \"abc\"": map[string]interface{}{
							"@foreach($.person.address, \"address\")": map[string]interface{}{
								"city":  "=$loop.city",
								"state": "=$loop.state",
							},
						},
					},
					{
						"$.person.name == \"bcd\"": map[string]interface{}{
							"@foreach($.person.address, \"address\")": map[string]interface{}{
								"city":  "=$loop.city",
								"state": "=$loop.state",
							},
						},
					},
					{
						"@otherwise": map[string]interface{}{
							"@foreach($.person.address, \"address\")": map[string]interface{}{
								"city":  "=$loop.city",
								"state": "=$loop.state",
							},
						},
					},
				},
			},
			Data: []map[string]interface{}{
				{
					"name": "abc",
					"address": []map[string]interface{}{
						{
							"city":  "abcsugarLand1",
							"state": "abctx1",
						},
						{
							"city":  "abcsugarLand2",
							"state": "abctx2",
						},
					},
				},
				{
					"name": "bcd",
					"address": []map[string]interface{}{
						{
							"city":  "bcdsugarLand1",
							"state": "bcdtx1",
						},
						{
							"city":  "bcdsugarLand2",
							"state": "bcdtx2",
						},
					},
				}, {
					"name": "ddd",
					"address": []map[string]interface{}{
						{
							"city":  "dddsugarLand1",
							"state": "dddtx1",
						},
						{
							"city":  "dddsugarLand2",
							"state": "dddtx2",
						},
					},
				},
			},
			Expected: []string{"abcsugarLand1", "bcdsugarLand1", "dddsugarLand1"},
		},
	}

	for _, tt := range testcases {
		v, _ := json.Marshal(tt.Mapping)
		fmt.Println("====", string(v))
		assert.True(t, IsConditionalMapping(tt.Mapping))
		mappings := map[string]interface{}{"output": tt.Mapping}
		factory := NewFactory(resolver)
		mapper, err := factory.NewMapper(mappings)
		assert.Nil(t, err)

		for i, input := range tt.Data {
			attrs := map[string]interface{}{"person": input}
			scope := data.NewSimpleScope(attrs, nil)
			results, err := mapper.Apply(scope)
			assert.Nil(t, err)
			assert.Equal(t, tt.Expected[i], results["output"].([]interface{})[0].(map[string]interface{})["city"])
		}

	}

}

func TestNestedConditionalObjectMapper(t *testing.T) {

	testcases := []struct {
		Mapping  interface{}
		Data     []map[string]interface{}
		Expected []string
	}{
		{
			Mapping: map[string]interface{}{
				"@conditional": []map[string]interface{}{
					{
						"$.person.name == \"abc\"": map[string]interface{}{
							"id": "=$.person.id",
							"name": map[string]interface{}{
								"@conditional": []map[string]interface{}{
									{
										"$.person.address.city == \"abcsugarLand\"": "abc",
									},
									{
										"$.person.address.city == \"abcsugarLand2\"": "abc2",
									},
									{
										"@otherwise": "abc3",
									},
								},
							},
							"address": "=$.person.address",
						},
					},
					{
						"$.person.name == \"bcd\"": map[string]interface{}{
							"id":      "=$.person.id",
							"name":    "=$.person.name",
							"address": "=$.person.address",
						},
					},
					{
						"@otherwise": map[string]interface{}{
							"id":      "=$.person.id",
							"name":    "=$.person.name",
							"address": "=$.person.address",
						},
					},
				},
			},
			Data: []map[string]interface{}{
				{
					"id":   "abc",
					"name": "abc",
					"address": map[string]interface{}{
						"city":  "abcsugarLand",
						"state": "abctx",
					},
				},
				{
					"id":   "bcd",
					"name": "abc",
					"address": map[string]interface{}{
						"city":  "abcsugarLand1",
						"state": "bcdtx",
					},
				},
				{
					"id":   "ddd",
					"name": "dddd",
					"address": map[string]interface{}{
						"city":  "dddsugarLand",
						"state": "dddtx",
					},
				},
			},
			Expected: []string{"abc", "abc3", "dddd"},
		},
	}

	for _, tt := range testcases {
		assert.True(t, IsConditionalMapping(tt.Mapping))
		mappings := map[string]interface{}{"output": tt.Mapping}
		factory := NewFactory(resolver)
		mapper, err := factory.NewMapper(mappings)
		assert.Nil(t, err)

		for i, input := range tt.Data {
			attrs := map[string]interface{}{"person": input}
			scope := data.NewSimpleScope(attrs, nil)
			results, err := mapper.Apply(scope)
			assert.Nil(t, err)
			assert.Equal(t, tt.Expected[i], results["output"].(map[string]interface{})["name"])
		}

	}

}

func TestNestedConditionalArrayMapper(t *testing.T) {

	testcases := []struct {
		Mapping  interface{}
		Data     []map[string]interface{}
		Expected []string
	}{
		{
			Mapping: map[string]interface{}{
				"@conditional": []map[string]interface{}{
					{
						"$.person.name == \"abc\"": map[string]interface{}{
							"@foreach($.person.address, \"address\")": map[string]interface{}{
								"city": "=$loop.city",
								"state": map[string]interface{}{
									"@conditional": []map[string]interface{}{
										{
											"$loop.state == \"abctx1\"": "tx1",
										},
										{
											"$loop.state == \"abctx2\"": "tx2",
										},
										{
											"@otherwise": "tx3",
										},
									},
								},
							},
						},
					},
					{
						"$.person.name == \"bcd\"": map[string]interface{}{
							"@foreach($.person.address, \"address\")": map[string]interface{}{
								"city":  "=$loop.city",
								"state": "=$loop.state",
							},
						},
					},
					{
						"@otherwise": map[string]interface{}{
							"@foreach($.person.address, \"address\")": map[string]interface{}{
								"city":  "=$loop.city",
								"state": "=$loop.state",
							},
						},
					},
				},
				"@if($.person.name == \"abc\")": map[string]interface{}{
					"@foreach($.person.address, \"address\")": map[string]interface{}{
						"city": "=$loop.city",
						"state": map[string]interface{}{
							"@if($loop.state == \"abctx1\")":     "tx1",
							"@elseif($loop.state == \"abctx2\")": "tx2",
							"@else":                              "tx3",
						},
					},
				},
				"@elseif($.person.name == \"bcd\")": map[string]interface{}{
					"@foreach($.person.address, \"address\")": map[string]interface{}{
						"city":  "=$loop.city",
						"state": "=$loop.state",
					},
				},
				"@else": map[string]interface{}{
					"@foreach($.person.address, \"address\")": map[string]interface{}{
						"city":  "=$loop.city",
						"state": "=$loop.state",
					},
				},
			},
			Data: []map[string]interface{}{
				{
					"name": "abc",
					"address": []map[string]interface{}{
						{
							"city":  "abcsugarLand1",
							"state": "abctx1",
						},
						{
							"city":  "abcsugarLand2",
							"state": "abctx2",
						},
					},
				},
				{
					"name": "abc",
					"address": []map[string]interface{}{
						{
							"city":  "bcdsugarLand1",
							"state": "dddd",
						},
						{
							"city":  "bcdsugarLand2",
							"state": "abctx2",
						},
					},
				}, {
					"name": "ddd",
					"address": []map[string]interface{}{
						{
							"city":  "dddsugarLand1",
							"state": "dddtx1",
						},
						{
							"city":  "dddsugarLand2",
							"state": "dddtx2",
						},
					},
				},
			},
			Expected: []string{"tx1", "tx3", "dddtx1"},
		},
	}

	for _, tt := range testcases {
		assert.True(t, IsConditionalMapping(tt.Mapping))
		mappings := map[string]interface{}{"output": tt.Mapping}
		factory := NewFactory(resolver)
		mapper, err := factory.NewMapper(mappings)
		assert.Nil(t, err)

		for i, input := range tt.Data {
			attrs := map[string]interface{}{"person": input}
			scope := data.NewSimpleScope(attrs, nil)
			results, err := mapper.Apply(scope)
			assert.Nil(t, err)
			assert.Equal(t, tt.Expected[i], results["output"].([]interface{})[0].(map[string]interface{})["state"])
		}
	}
}

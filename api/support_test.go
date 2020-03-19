package api

import (
	"testing"

	sampleActivity "github.com/project-flogo/core/examples/activity"
	"github.com/stretchr/testify/assert"
)

func TestToMappings(t *testing.T) {

	mappings := []string{"in1=b", "in2= $.blah", "in3 = $.blah2"}

	//todo add additional tests when support for more mapping type is added
	defs, err := toMappings(mappings)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(defs))

	v, exists := defs["in1"]
	assert.True(t, exists)
	assert.Equal(t, "=b", v)

	v, exists = defs["in2"]
	assert.True(t, exists)
	assert.Equal(t, "= $.blah", v)

	v, exists = defs["in3"]
	assert.True(t, exists)
	assert.Equal(t, "= $.blah2", v)

}

func TestActivity(t *testing.T) {
	sampleAct, err := NewActivity(&sampleActivity.Activity{}, map[string]interface{}{"aSetting": "aSetting"})

	assert.Nil(t, err)
	assert.NotNil(t, sampleAct)

	result, err := EvalActivity(sampleAct, map[string]interface{}{"anInput": "a"})
	assert.Nil(t, err)
	assert.Equal(t, "a", result["anOutput"])
}

func TestTriggerConfig(t *testing.T) {
	var handlers []*Handler
	var actions []*Action

	actions = append(actions, &Action{ref: "sampleAction", condition: "if", inputMappings: []string{"=in"}, settings: map[string]interface{}{"aSettings": "aSet"}})
	handlers = append(handlers, &Handler{actions: actions, name: "aSampleAction"})
	trigger := &Trigger{ref: "sampleTrigger", settings: map[string]interface{}{"aSetting": "aSet"}, handlers: handlers}

	cfg := toTriggerConfig("sample", trigger)

	assert.NotNil(t, cfg)
}

func TestActionConfig(t *testing.T) {
	action := &Action{ref: "sampleAction", condition: "if", inputMappings: []string{"=in"}, settings: map[string]interface{}{"aSettings": "aSet"}}

	actCfg := toActionConfig(action)

	assert.NotNil(t, actCfg)
}

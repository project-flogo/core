package sample

import (
	"context"
	"testing"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/support"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	ref := support.GetRef(&Action{})
	f := action.GetFactory(ref)

	assert.NotNil(t, f)
}

func TestFactory(t *testing.T) {
	f := &ActionFactory{}

	config := &action.Config{Settings: map[string]interface{}{"aSetting": "test_setting"}}
	act, err := f.New(config)

	assert.Nil(t, err)
	actInst, ok := act.(*Action)
	assert.True(t, ok)

	assert.NotNil(t, actInst.settings)
	assert.Equal(t, "test_setting", actInst.settings.ASetting)
}

func TestEval(t *testing.T) {

	act := &Action{}
	inputMap := map[string]interface{}{"anInput": "test"}
	results, err := act.Run(context.Background(), inputMap)
	assert.Nil(t, err)

	output := &Output{}
	err = output.FromMap(results)
	assert.Nil(t, err)
	assert.Equal(t, "test", output.AnOutput)
}

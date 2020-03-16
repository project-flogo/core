package api

import (
	"context"
	"testing"

	"github.com/project-flogo/core/data"
	sampleAction "github.com/project-flogo/core/examples/action"
	sampleTrigger "github.com/project-flogo/core/examples/trigger"
	"github.com/project-flogo/core/support/log"
	"github.com/stretchr/testify/assert"
)

type Settings struct {
	ASetting int `md:"aSetting"`
}
type HandlerSettings struct {
	ASetting string `md:"aSetting"`
}

func TestNewApp(t *testing.T) {
	app := NewApp()
	assert.NotNil(t, app)

}
func TestNewTrigger(t *testing.T) {
	app := NewApp()

	trg := app.NewTrigger(&sampleTrigger.Trigger{}, map[string]interface{}{"aSetting": 1})
	assert.NotNil(t, trg)

	trg = app.NewTrigger(&sampleTrigger.Trigger{}, &Settings{})
	assert.NotNil(t, trg)

	handler, err := trg.NewHandler(map[string]interface{}{"aSetting": "aSetting"})
	assert.Nil(t, err)
	assert.NotNil(t, handler)

	handler, err = trg.NewHandler(&HandlerSettings{})
	assert.Nil(t, err)
	assert.NotNil(t, handler)

}

type IndependentAction struct {
	*sampleAction.Action
}

func (ind *IndependentAction) Ref() string {
	return "github.com/project-flogo/core/examples/action"
}
func TestIndependentAction(t *testing.T) {

	app := NewApp()
	eng, err := NewEngine(app)
	assert.NotNil(t, eng)
	var act *IndependentAction

	newAct, err := app.NewIndependentAction(act, map[string]interface{}{"aSetting": "a"})
	assert.Nil(t, err)
	assert.NotNil(t, newAct)

	result, err := RunAction(context.Background(), newAct, map[string]interface{}{"anInput": "a"})
	assert.Nil(t, err)
	assert.Equal(t, "a", result["anOutput"])

	assert.Nil(t, newAct.IOMetadata())
	assert.NotNil(t, newAct.Metadata())

}

func TestNewAction(t *testing.T) {
	app := NewApp()

	trg := app.NewTrigger(&sampleTrigger.Trigger{}, &Settings{})
	assert.NotNil(t, trg)

	handler, err := trg.NewHandler(map[string]interface{}{"aSetting": "aSetting"})
	assert.Nil(t, err)
	assert.NotNil(t, handler)

	newAct, err := handler.NewAction(&sampleAction.Action{}, map[string]interface{}{"aSetting": "a"})
	assert.Nil(t, err)
	assert.NotNil(t, newAct)

	err = app.AddAction("sampleAction", &sampleAction.Action{}, map[string]interface{}{"aSetting": "aSetting"})
	assert.Nil(t, err)
}

func LogMessage(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {
	log.RootLogger().Infof("#v", inputs)
	return nil, nil
}

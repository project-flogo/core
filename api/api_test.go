package api

import (
	"context"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/support/log"
	"testing"
)

func TestNewTrigger(t *testing.T) {

	//todo add dummy test trigger
	//app := NewApp()
	//
	//trg := app.NewTrigger(&rest.RestTrigger{}, map[string]interface{}{"port": 8080})
	//
	//h1 := trg.NewHandler(map[string]interface{}{"method": "GET", "path": "/blah"})
	//a := h1.NewAction(&flow.FlowAction{}, map[string]interface{}{"flowURI": "res://flow:get_git_hub_issues"})
	//a.SetInputMappings("in1='blah'", "in2=1")
	//a.SetOutputMappings("out1='blah'", "out2=$.flowOut")
	//
	////app.AddResource("flow:myflow", flowJson)
	//e, err := NewEngine(app)
	//
	//assert.Nil(t, err)
	//assert.NotNil(t, e)
	//
	//err = e.Start()
	//assert.Nil(t, err)
	//
	//err = e.Stop()
	//assert.Nil(t, err)
}

func TestTrigger_NewFuncHandler(t *testing.T) {

	//todo add dummy test trigger

	//app := NewApp()

	//trg := app.NewTrigger(&rest.RestTrigger{}, map[string]interface{}{"port": 8080})
	//trg.NewFuncHandler(map[string]interface{}{"method": "GET", "path": "/blah"}, LogMessage)
	//
	//e, err := NewEngine(app)
	//
	//assert.Nil(t, err)
	//assert.NotNil(t, e)
	//
	//err = e.Start()
	//assert.Nil(t, err)
	//
	//err = e.Stop()
	//assert.Nil(t, err)
}

func LogMessage(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {
	log.RootLogger().Infof("#v", inputs)
	return nil, nil
}

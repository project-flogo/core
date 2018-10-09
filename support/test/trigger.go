package test

import (
	"context"
	"fmt"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/engine/runner"
	"github.com/project-flogo/core/trigger"
)

func InitTrigger(factory trigger.Factory, tConfig *trigger.Config, actions map[string]action.Action) (trigger.Trigger, error) {

	r := runner.NewDirect()

	if factory == nil {
		return nil, fmt.Errorf("Trigger Factory not provided")
	}

	trg := factory.New(tConfig)

	if trg == nil {
		return nil, fmt.Errorf("cannot create Trigger for id '%s'", tConfig.Id)
	}

	tConfig.FixUp(trg.Metadata())

	initCtx := &initContext{handlers: make([]*trigger.Handler, 0, len(tConfig.Handlers))}

	trgInit, ok := trg.(trigger.Initializable)

	if !ok {
		return nil, fmt.Errorf("Trigger does not implement trigger.Initializable interface")
	}

	//create handlers for that trigger and init
	for _, hConfig := range tConfig.Handlers {

		if hConfig.Action.Id == "" {
			return nil, fmt.Errorf("Action not specified for handler")
		}

		act, exists := actions[hConfig.Action.Id]
		if !exists {
			return nil, fmt.Errorf("specified Action '%s' does not exists", hConfig.Action.Id)
		}

		handler := trigger.NewHandler(hConfig, act, trg.Metadata().Output, trg.Metadata().Reply, r)
		initCtx.handlers = append(initCtx.handlers, handler)
	}

	err := trgInit.Initialize(initCtx)
	if err != nil {
		return nil, err
	}

	return trg, nil
}

//////////////////////////
// Simple Init Context

type initContext struct {
	handlers []*trigger.Handler
}

func (ctx *initContext) GetHandlers() []*trigger.Handler {
	return ctx.handlers
}

//////////////////////////
// Dummy Test Action

func NewDummyAction(f func()) action.Action {
	return &testAction{f: f}
}

type testAction struct {
	f func()
}

// Metadata get the Action's metadata
func (a *testAction) Metadata() *action.Metadata {
	return nil
}

// IOMetadata get the Action's IO metadata
func (a *testAction) IOMetadata() *data.IOMetadata {
	return nil
}

// Run implementation of action.SyncAction.Run
func (a *testAction) Run(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {
	a.f()
	return nil, nil
}

//
//type testRunner struct {
//}
//
////DEPRECATED
//func (tr *testRunner) Run(context context.Context, act action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
//	return 0, nil, nil
//}
//
////DEPRECATED
//func (tr *testRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {
//	return 0, nil, nil
//}
//
//func (tr *testRunner) Execute(ctx context.Context, act action.Action, inputs map[string]*data.Attribute) (results map[string]*data.Attribute, err error) {
//
//
//	return nil, nil
//}

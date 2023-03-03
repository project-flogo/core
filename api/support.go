package api

import (
	"context"
	"reflect"
	"strconv"
	"strings"

	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/trace"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/app"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/trigger"
)

// toAppConfig converts an App to the core app configuration model
func toAppConfig(a *App) *app.Config {

	appCfg := &app.Config{}
	appCfg.Name = "app"
	appCfg.Version = "1.0.0"
	appCfg.Resources = a.resources

	var triggerConfigs []*trigger.Config
	for id, trg := range a.Triggers() {

		triggerConfigs = append(triggerConfigs, toTriggerConfig(strconv.Itoa(id+1), trg))
	}

	appCfg.Triggers = triggerConfigs

	var properties []*data.Attribute
	for name, tv := range a.Properties() {
		attr := data.NewAttribute(name, tv.Type(), tv.Value())
		properties = append(properties, attr)
	}

	appCfg.Properties = properties

	for key, value := range a.actions {
		act := &action.Config{
			Ref:      value.ref,
			Id:       key,
			Settings: value.settings,
		}
		appCfg.Actions = append(appCfg.Actions, act)
	}

	return appCfg
}

// toTriggerConfig converts Trigger to the core Trigger configuration model
func toTriggerConfig(id string, trg *Trigger) *trigger.Config {

	triggerConfig := &trigger.Config{Id: id, Ref: trg.ref, Settings: trg.Settings()}

	var handlerConfigs []*trigger.HandlerConfig
	for _, handler := range trg.Handlers() {
		h := &trigger.HandlerConfig{Settings: handler.Settings()}
		actions := handler.Actions()
		h.Actions = make([]*trigger.ActionConfig, len(actions))
		for i, anAction := range actions {
			h.Actions[i] = toActionConfig(anAction)
		}
		handlerConfigs = append(handlerConfigs, h)
	}

	triggerConfig.Handlers = handlerConfigs
	return triggerConfig
}

// toActionConfig converts Action to the core Action configuration model
func toActionConfig(act *Action) *trigger.ActionConfig {
	actionCfg := &trigger.ActionConfig{
		Config: &action.Config{},
	}

	if act.act != nil {
		actionCfg.Act = act.act
		return actionCfg
	}

	actionCfg.Id = act.id
	actionCfg.Ref = act.ref
	actionCfg.Settings = act.settings
	actionCfg.If = act.condition

	if len(act.inputMappings) > 0 {
		actionCfg.Input, _ = toMappings(act.inputMappings)
	}
	if len(act.outputMappings) > 0 {
		actionCfg.Output, _ = toMappings(act.outputMappings)
	}

	return actionCfg
}

func toMappings(strMappings []string) (map[string]interface{}, error) {

	mappings := make(map[string]interface{}, len(strMappings))
	for _, strMapping := range strMappings {

		idx := strings.Index(strMapping, "=")
		lhs := strings.TrimSpace(strMapping[:idx])
		rhs := strings.TrimSpace(strMapping[idx:])

		mappings[lhs] = rhs
	}
	return mappings, nil
}

// ProxyAction

type proxyAction struct {
	handlerFunc HandlerFunc
	metadata    *action.Metadata
}

func NewProxyAction(f HandlerFunc) action.Action {
	return &proxyAction{
		handlerFunc: f,
		metadata:    &action.Metadata{},
	}
}

// Metadata get the Action's metadata
func (a *proxyAction) Metadata() *action.Metadata {
	return a.metadata
}

// IOMetadata get the Action's IO metadata
func (a *proxyAction) IOMetadata() *metadata.IOMetadata {
	return nil
}

// Run implementation of action.SyncAction.Run
func (a *proxyAction) Run(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	return a.handlerFunc(ctx, input)
}

// NewActivity creates an instance of the specified activity
func NewActivity(act activity.Activity, settings ...interface{}) (activity.Activity, error) {

	ref := activity.GetRef(act)

	if f := activity.GetFactory(ref); f == nil {

		return activity.Get(ref), nil
	} else {

		var settingsMap map[string]interface{}

		if len(settings) == 0 {
			settingsMap = make(map[string]interface{})
		} else {
			inSettings := settings[0]

			if im, ok := inSettings.(map[string]interface{}); ok {
				settingsMap = im
			} else {
				settingsMap = metadata.StructToMap(inSettings)
			}
		}

		f := activity.GetFactory(ref)
		ctx := &initCtx{settings: settingsMap}
		return f(ctx)
	}

}

type initCtx struct {
	settings map[string]interface{}
}

func (ctx *initCtx) Settings() map[string]interface{} {
	return ctx.settings
}

func (ctx *initCtx) MapperFactory() mapper.Factory {
	return nil
}

func (ctx *initCtx) Logger() log.Logger {
	return log.RootLogger()
}

func (ctx *initCtx) Name() string {
	return ""
}

func (ctx *initCtx) HostName() string {
	return ""
}

var activityLogger = log.ChildLogger(log.RootLogger(), "activity")

// EvalActivity evaluates the specified activity using the provided inputs
func EvalActivity(act activity.Activity, input interface{}) (map[string]interface{}, error) {

	var inputMap map[string]interface{}

	if im, ok := input.(map[string]interface{}); ok {
		inputMap = im
	} else if tm, ok := input.(data.StructValue); ok {
		inputMap = tm.ToMap()
	}

	logger := activityLogger

	if act.Metadata() == nil {

		//try loading activity with metadata

		var ref string
		if hr, ok := act.(support.HasRef); ok {
			ref = hr.Ref()
		} else {
			value := reflect.ValueOf(act)
			value = value.Elem()
			ref = value.Type().PkgPath()
		}
		act = activity.Get(ref)
		l := activity.GetLogger(ref)
		if l != nil {
			logger = l
		}
	}

	if act.Metadata() == nil {
		//return error
	}

	ac := &activityContext{input: make(map[string]interface{}), output: make(map[string]interface{}), logger: logger}

	for key, value := range inputMap {
		ac.input[key] = value
	}

	_, evalErr := act.Eval(ac)

	if evalErr != nil {
		return nil, evalErr
	}

	return ac.output, nil
}

/////////////////////////////////////////
// activity.Context Implementation

type activityContext struct {
	input  map[string]interface{}
	output map[string]interface{}
	logger log.Logger
}

func (aCtx *activityContext) ActivityHost() activity.Host {
	return aCtx
}

func (aCtx *activityContext) Name() string {
	return ""
}

// GetInput implements activity.Context.GetInput
func (aCtx *activityContext) GetInput(name string) interface{} {

	val, found := aCtx.input[name]
	if found {
		return val
	}

	return nil
}

// SetOutput implements activity.Context.SetOutput
func (aCtx *activityContext) SetOutput(name string, value interface{}) error {
	aCtx.output[name] = value
	return nil
}

func (aCtx *activityContext) GetSharedTempData() map[string]interface{} {
	return nil
}

func (aCtx *activityContext) Logger() log.Logger {
	return aCtx.logger
}

func (aCtx *activityContext) GetInputObject(input data.StructValue) error {
	err := input.FromMap(aCtx.input)
	return err
}

func (aCtx *activityContext) SetOutputObject(output data.StructValue) error {
	aCtx.output = output.ToMap()
	return nil
}

func (aCtx *activityContext) GetTracingContext() trace.TracingContext {
	return nil
}

/////////////////////////////////////////
// activity.Host Implementation

func (aCtx *activityContext) ID() string {
	//ignore
	return ""
}

func (aCtx *activityContext) IOMetadata() *metadata.IOMetadata {
	return nil
}

func (aCtx *activityContext) Reply(replyData map[string]interface{}, err error) {
	// ignore
}

func (aCtx *activityContext) Return(returnData map[string]interface{}, err error) {
	// ignore
}

func (aCtx *activityContext) Scope() data.Scope {
	return nil
}

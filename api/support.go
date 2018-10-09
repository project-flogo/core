package api

import (
	"context"
	"encoding/json"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/data/resolve"
	"reflect"
	"strconv"
	"strings"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/app"
	"github.com/project-flogo/core/data"
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

	return appCfg
}

// toTriggerConfig converts Trigger to the core Trigger configuration model
func toTriggerConfig(id string, trg *Trigger) *trigger.Config {

	triggerConfig := &trigger.Config{Id: id, Ref: trg.ref, Settings: trg.Settings()}

	var handlerConfigs []*trigger.HandlerConfig
	for _, handler := range trg.Handlers() {
		h := &trigger.HandlerConfig{Settings: handler.Settings()}
		h.Action = toActionConfig(handler.Action())
		handlerConfigs = append(handlerConfigs, h)
	}

	triggerConfig.Handlers = handlerConfigs
	return triggerConfig
}

// toActionConfig converts Action to the core Action configuration model
func toActionConfig(act *Action) *trigger.ActionConfig {
	actionCfg := &trigger.ActionConfig{}

	if act.act != nil {
		actionCfg.Act = act.act
		return actionCfg
	}

	actionCfg.Ref = act.ref

	//todo handle error
	jsonData, _ := json.Marshal(act.Settings())
	actionCfg.Data = jsonData

	if len(act.inputMappings) > 0 {
		actionCfg.Input,_ = toMappings(act.inputMappings)
	}
	if len(act.outputMappings) > 0 {
		actionCfg.Output,_ = toMappings(act.outputMappings)
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

	if len(settings) == 0 {
		return activity.Get(ref), nil
	} else {

		inSettings := settings[0]

		var settingsMap map[string]interface{}
		if im, ok := inSettings.(map[string]interface{}); ok {
			settingsMap = im
		} else {
			settingsMap = metadata.StructToMap(settings)
		}

		f := activity.GetFactory(ref)
		ctx := &initCtx{settings:settingsMap}
		return f(ctx)
	}

	return nil, nil
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

// EvalActivity evaluates the specified activity using the provided inputs
func EvalActivity(act activity.Activity, input interface{}) (map[string]interface{}, error) {

	var inputMap map[string]interface{}

	if im, ok := input.(map[string]interface{}); ok {
		inputMap = im
	} else if tm, ok := input.(metadata.ToMap); ok {
		inputMap = tm.ToMap()
	}

	if act.Metadata() == nil {
		//try loading activity with metadata
		value := reflect.ValueOf(act)
		value = value.Elem()
		ref := value.Type().PkgPath()

		act = activity.Get(ref)
	}

	if act.Metadata() == nil {
		//return error
	}

	ac := &activityContext{input: make(map[string]interface{}), output: make(map[string]interface{})}

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
func (aCtx *activityContext) SetOutput(name string, value interface{}) {
	aCtx.output[name] = value
}

func (aCtx *activityContext) GetSharedTempData() map[string]interface{} {
	return nil
}

func (aCtx *activityContext) GetInputObject(object interface{}, converter activity.InputConverter) error {
	return nil
}

func (aCtx *activityContext) SetOutputObject(object interface{}, converter activity.OutputConverter) error {
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

func (aCtx *activityContext) GetResolver() resolve.CompositeResolver {
	return resolve.GetBasicResolver()
}

func (aCtx *activityContext) WorkingData() data.Scope {
	return nil
}

func (aCtx *activityContext) GetDetails() data.StringsMap {
	return nil
}

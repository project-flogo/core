package test

import (
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/trace"
)

var logger = log.ChildLogger(log.RootLogger(), "test")

// NewActivityContext creates a new TestActivityContext
func NewActivityContext(md *activity.Metadata) *TestActivityContext {

	input := map[string]data.TypedValue{"Input1": data.NewTypedValue(data.TypeString, "")}
	output := map[string]data.TypedValue{"Output1": data.NewTypedValue(data.TypeString, "")}

	ac := &TestActivityHost{
		HostId:     "1",
		HostRef:    "github.com/project-flogo/flow",
		IoMetadata: &metadata.IOMetadata{Input: input, Output: output},
		HostData:   data.NewSimpleScope(nil, nil),
	}

	return NewActivityContextWithAction(md, ac)
}

// NewActivityContextWithAction creates a new TestActivityContext
func NewActivityContextWithAction(md *activity.Metadata, activityHost *TestActivityHost) *TestActivityContext {

	tc := &TestActivityContext{
		metadata:     md,
		activityHost: activityHost,
		TaskNameVal:  "Test TaskOld",
		inputs:       make(map[string]interface{}, len(md.Input)),
		outputs:      make(map[string]interface{}, len(md.Output)),
		settings:     make(map[string]interface{}, len(md.Settings)),
	}

	for name, tv := range md.Input {
		tc.inputs[name] = tv.Value()
	}
	for name, tv := range md.Output {
		tc.outputs[name] = tv.Value()
	}

	return tc
}

func NewActivityInitContext(settings interface{}, f mapper.Factory) activity.InitContext {

	var settingVals map[string]interface{}

	if sm, ok := settings.(map[string]interface{}); ok {
		settingVals = sm
	} else {
		settingVals = metadata.StructToMap(settings)
	}

	if f == nil {
		f = mapper.NewFactory(resolve.GetBasicResolver())
	}

	return &TestActivityInitContext{settings: settingVals, factory: f}
}

type TestActivityInitContext struct {
	settings map[string]interface{}
	factory  mapper.Factory
}

func (ic *TestActivityInitContext) Settings() map[string]interface{} {
	return ic.settings
}

func (ic *TestActivityInitContext) MapperFactory() mapper.Factory {
	return ic.factory
}

func (ic *TestActivityInitContext) Logger() log.Logger {
	return logger
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// TestActivityHost

type TestActivityHost struct {
	HostId  string
	HostRef string

	IoMetadata *metadata.IOMetadata
	HostData   data.Scope
	ReplyData  map[string]interface{}
	ReplyErr   error
}

func (ac *TestActivityHost) IOMetadata() *metadata.IOMetadata {
	return ac.IoMetadata
}

func (ac *TestActivityHost) Reply(replyData map[string]interface{}, err error) {
	ac.ReplyData = replyData
	ac.ReplyErr = err
}

func (ac *TestActivityHost) Return(returnData map[string]interface{}, err error) {
	ac.ReplyData = returnData
	ac.ReplyErr = err
}

func (ac *TestActivityHost) Name() string {
	return ""
}

func (ac *TestActivityHost) ID() string {
	return ac.HostId
}

func (ac *TestActivityHost) Scope() data.Scope {
	return ac.HostData
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// TestActivityContext

// TestActivityContext is a dummy ActivityContext to assist in testing
type TestActivityContext struct {
	TaskNameVal  string
	activityHost activity.Host

	metadata *activity.Metadata
	settings map[string]interface{}
	inputs   map[string]interface{}
	outputs  map[string]interface{}

	shared map[string]interface{}
}

func (c *TestActivityContext) SetInputObject(input data.StructValue) error {
	c.inputs = input.ToMap()
	return nil
}

func (c *TestActivityContext) GetOutputObject(output data.StructValue) error {
	err := output.FromMap(c.outputs)
	return err
}

func (c *TestActivityContext) GetInputObject(input data.StructValue) error {
	err := input.FromMap(c.inputs)
	return err
}

func (c *TestActivityContext) SetOutputObject(output data.StructValue) error {
	c.outputs = output.ToMap()
	return nil
}

func (c *TestActivityContext) ActivityHost() activity.Host {
	return c.activityHost
}

func (c *TestActivityContext) Name() string {
	return c.TaskNameVal
}

// GetSetting implements activity.Context.GetSetting
func (c *TestActivityContext) GetSetting(setting string) (value interface{}, exists bool) {

	attr, found := c.settings[setting]

	if found {
		return attr, true
	}

	return nil, false
}

func (c *TestActivityContext) SetInput(name string, val interface{}) {
	c.inputs[name] = val
}

// GetInput implements activity.Context.GetInput
func (c *TestActivityContext) GetInput(name string) interface{} {

	attr, found := c.inputs[name]

	if found {
		return attr
	}

	return nil
}

// SetOutput implements activity.Context.SetOutput
func (c *TestActivityContext) SetOutput(name string, value interface{}) error {
	c.outputs[name] = value
	return nil
}

// GetOutput implements activity.Context.GetOutput
func (c *TestActivityContext) GetOutput(name string) interface{} {

	attr, found := c.outputs[name]

	if found {
		return attr
	}

	return nil
}

func (c *TestActivityContext) GetSharedTempData() map[string]interface{} {

	if c.shared == nil {
		c.shared = make(map[string]interface{})
	}
	return c.shared
}

func (c *TestActivityContext) Logger() log.Logger {
	return logger
}

func (c *TestActivityContext) GetTracingContext() trace.TracingContext {
	return nil
}

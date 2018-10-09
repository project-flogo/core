package test

import (
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/data/resolve"
)

//todo needs to move to lib
// NewTestActivityContext creates a new TestActivityContext
func NewTestActivityContext(md *activity.Metadata) *TestActivityContext {

	input := map[string]data.TypedValue{"Input1": data.NewTypedValue(data.TypeString, "")}
	output := map[string]data.TypedValue{"Output1": data.NewTypedValue(data.TypeString, "")}

	ac := &TestActivityHost{
		HostId:     "1",
		HostRef:    "github.com/TIBCOSoftware/flogo-contrib/action/flow",
		IoMetadata: &metadata.IOMetadata{Input: input, Output: output},
		HostData:   data.NewSimpleScope(nil, nil),
	}

	return NewTestActivityContextWithAction(md, ac)
}

// NewTestActivityContextWithAction creates a new TestActivityContext
func NewTestActivityContextWithAction(md *activity.Metadata, activityHost *TestActivityHost) *TestActivityContext {

	//fd := &TestFlowDetails{
	//	FlowIDVal:   "1",
	//	FlowNameVal: "Test Flow",
	//}

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



/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// TestFlowDetails

// TestFlowDetails simple FlowDetails for use in testing
type TestFlowDetails struct {
	FlowIDVal   string
	FlowNameVal string
}

// ID implements activity.FlowDetails.ID
func (fd *TestFlowDetails) ID() string {
	return fd.FlowIDVal
}

// Name implements activity.FlowDetails.Name
func (fd *TestFlowDetails) Name() string {
	return fd.FlowNameVal
}

// ReplyHandler implements activity.FlowDetails.ReplyHandler
func (fd *TestFlowDetails) ReplyHandler() activity.ReplyHandler {
	return nil
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

func (ac *TestActivityHost) GetResolver() resolve.CompositeResolver {
	return resolve.GetBasicResolver()
}

func (ac *TestActivityHost) Name() string {
	return ""
}

func (ac *TestActivityHost) ID() string {
	return ac.HostId
}

func (ac *TestActivityHost) WorkingData() data.Scope {
	return ac.HostData
}

func (ac *TestActivityHost) GetDetails() data.StringsMap {
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// TestActivityContext

// TestActivityContext is a dummy ActivityContext to assist in testing
type TestActivityContext struct {
	//details      activity.FlowDetails
	TaskNameVal string
	//Attrs        map[string]*data.Attribute
	activityHost activity.Host

	metadata *activity.Metadata
	settings map[string]interface{}
	inputs   map[string]interface{}
	outputs  map[string]interface{}

	shared map[string]interface{}
}

func (c *TestActivityContext) SetInputObject(input data.ToMap) error {
	c.inputs = input.ToMap()
	return nil
}

func (c *TestActivityContext) GetOutputObject(output data.FromMap) error {
	err := output.FromMap(c.outputs)
	return err
}

func (c *TestActivityContext) GetInputObject(input data.FromMap) error {
	err := input.FromMap(c.inputs)
	return err
}

func (c *TestActivityContext) SetOutputObject(output data.ToMap) error {
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
func (c *TestActivityContext) SetOutput(name string, value interface{}) {

	c.outputs[name] = value
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

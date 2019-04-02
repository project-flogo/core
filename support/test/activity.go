package test

import (
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.LegacyRegister("testlog", NewLogActivity())
	_ = activity.LegacyRegister("testcounter", NewCounterActivity())
}

type TestLogActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewLogActivity() activity.Activity {

	md := &activity.Metadata{IOMetadata: &metadata.IOMetadata{Input: map[string]data.TypedValue{"message": data.NewTypedValue(data.TypeString, "")},
		Output: map[string]data.TypedValue{"message": data.NewTypedValue(data.TypeString, "")}}}
	return &TestLogActivity{metadata: md}
}

// Metadata returns the activity's metadata
func (a *TestLogActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *TestLogActivity) Eval(ctx activity.Context) (done bool, err error) {

	ctx.Logger().Debugf("eval test-log activity")

	message, _ := ctx.GetInput("message").(string)

	ctx.Logger().Infof("message: %s", message)

	err = ctx.SetOutput("message", message)
	if err != nil {
		return false, err
	}

	return true, nil
}

type TestCounterActivity struct {
	metadata *activity.Metadata
	counters map[string]int
}

// NewActivity creates a new AppActivity
func NewCounterActivity() activity.Activity {

	md := &activity.Metadata{IOMetadata: &metadata.IOMetadata{Input: map[string]data.TypedValue{"counterName": data.NewTypedValue(data.TypeString, "")},
		Output: map[string]data.TypedValue{"value": data.NewTypedValue(data.TypeInt, "")}}}

	return &TestCounterActivity{metadata: md, counters: make(map[string]int)}
}

// Metadata returns the activity's metadata
func (a *TestCounterActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *TestCounterActivity) Eval(ctx activity.Context) (done bool, err error) {

	ctx.Logger().Debugf("eval test-counter activity")

	counterName, _ := ctx.GetInput("counterName").(string)

	ctx.Logger().Debugf("counterName: %s", counterName)

	count := 1

	if counter, exists := a.counters[counterName]; exists {
		count = counter + 1
	}

	a.counters[counterName] = count

	ctx.Logger().Debugf("value: %s", count)

	err = ctx.SetOutput("value", count)
	if err != nil {
		return false, err
	}

	return true, nil
}

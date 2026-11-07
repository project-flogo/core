package trigger

import (
	"context"
	"testing"

	"github.com/project-flogo/core/support/log"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/data/property"
	"github.com/project-flogo/core/data/resolve"
	"github.com/stretchr/testify/assert"
)

var defResolver = resolve.NewCompositeResolver(map[string]resolve.Resolver{
	".": &resolve.ScopeResolver{},
})

type MockAction struct {
}

func (t *MockAction) Metadata() *action.Metadata {
	return nil
}

func (t *MockAction) IOMetadata() *metadata.IOMetadata {
	return nil
}

func (t *MockAction) Run(context context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
	return inputs, nil
}

func TestNewHandler(t *testing.T) {
	actionCfg := &ActionConfig{Input: map[string]interface{}{"anInput": "input"}, Output: map[string]interface{}{"anOutput": "output"}}
	var actionCfgArr []*ActionConfig
	actionCfgArr = append(actionCfgArr, actionCfg)

	hCfg := &HandlerConfig{Name: "sampleConfig", Settings: map[string]interface{}{"aSetting": "aSetting"}, Actions: actionCfgArr}

	mf := mapper.NewFactory(defResolver)
	expf := expression.NewFactory(defResolver)

	//Action not specified
	handler, err := NewHandler(hCfg, nil, mf, expf, nil, log.RootLogger())
	assert.NotNil(t, err, "Actions not specified.")

	//HandlerSettings configured
	handler, err = NewHandler(hCfg, []action.Action{&MockAction{}}, mf, expf, nil, log.RootLogger())
	assert.NotNil(t, handler.Settings())
}

func TestHandlerContext(t *testing.T) {

	actionCfg := &ActionConfig{Input: map[string]interface{}{"anInput": "input"}, Output: map[string]interface{}{"anOutput": "output"}}
	var actionCfgArr []*ActionConfig
	actionCfgArr = append(actionCfgArr, actionCfg)

	hCfg := &HandlerConfig{Name: "sampleConfig", Settings: map[string]interface{}{"aSetting": "aSetting"}, Actions: actionCfgArr}
	hCfg.Parent = &Config{Id: "sampleTrig"}

	assert.NotNil(t, NewHandlerContext(context.Background(), hCfg))

}

type captureRunner struct {
	inputs map[string]interface{}
}

func (r *captureRunner) RunAction(ctx context.Context, act action.Action, inputs map[string]interface{}) (map[string]interface{}, error) {
	r.inputs = inputs
	return nil, nil
}

func TestHandler_PropertySnapshot(t *testing.T) {
	t.Setenv(property.EnvAppPropertySnapshotEnabled, "true")
	previous := property.DefaultManager()
	t.Cleanup(func() { property.SetDefaultManager(previous) })
	property.SetDefaultManager(property.NewManager(map[string]property.Property{
		"p1": {Value: "a", Overridable: true},
		"p2": {Value: 2},
	}))

	mf := mapper.NewFactory(defResolver)
	expf := expression.NewFactory(defResolver)
	for name, seqKey := range map[string]string{"handler": "", "sequence-key handler": "k1"} {
		t.Run(name, func(t *testing.T) {
			hCfg := &HandlerConfig{Name: "h", Actions: []*ActionConfig{{SequenceKey: seqKey}}, Parent: &Config{Id: "t"}}
			runner := &captureRunner{}
			handler, err := NewHandler(hCfg, []action.Action{&MockAction{}}, mf, expf, runner, log.RootLogger())
			assert.NoError(t, err)
			_, isSeqKeyHandler := handler.(*seqKeyHandlerImpl)
			assert.Equal(t, seqKey != "", isSeqKeyHandler)

			_, err = handler.Handle(context.Background(), map[string]interface{}{})
			assert.NoError(t, err)

			snapshot, ok := runner.inputs["_PROPERTIES"].(map[string]interface{})
			assert.True(t, ok)
			assert.Equal(t, map[string]interface{}{"p1": "a", "p2": 2}, snapshot)
			snapshot["p1"] = "changed"
			v, _ := property.DefaultManager().GetProperty("p1")
			assert.Equal(t, "a", v, "the action's snapshot must be its own copy")
		})
	}
}

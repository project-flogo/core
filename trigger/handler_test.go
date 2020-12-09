package trigger

import (
	"context"
	"github.com/project-flogo/core/support/log"
	"testing"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/engine/runner"
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
	handler, err := NewHandler(hCfg, nil, mf, expf, runner.NewDirect(), log.RootLogger())
	assert.NotNil(t, err, "Actions not specified.")

	//Parent not defined in the Handler Config
	handler, err = NewHandler(hCfg, []action.Action{&MockAction{}}, mf, expf, runner.NewDirect(), log.RootLogger())
	_, err = handler.Handle(context.Background(), map[string]interface{}{"anInput": "input"})
	assert.NotNil(t, err, "Parent not defined.")

	//Parent defined.
	hCfg.parent = &Config{Id: "sampleTrig"}
	handler, err = NewHandler(hCfg, []action.Action{&MockAction{}}, mf, expf, runner.NewDirect(), log.RootLogger())
	assert.Nil(t, err)
	assert.NotNil(t, handler)

	assert.NotNil(t, handler.Settings())
	out, err := handler.Handle(context.Background(), map[string]interface{}{"anInput": "input"})

	assert.Equal(t, "output", out["anOutput"])
}

func TestHandlerContext(t *testing.T) {

	actionCfg := &ActionConfig{Input: map[string]interface{}{"anInput": "input"}, Output: map[string]interface{}{"anOutput": "output"}}
	var actionCfgArr []*ActionConfig
	actionCfgArr = append(actionCfgArr, actionCfg)

	hCfg := &HandlerConfig{Name: "sampleConfig", Settings: map[string]interface{}{"aSetting": "aSetting"}, Actions: actionCfgArr}
	hCfg.parent = &Config{Id: "sampleTrig"}

	assert.NotNil(t, NewHandlerContext(context.Background(), hCfg))

}

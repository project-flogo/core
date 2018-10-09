package runner

import (
	"context"
	"errors"
	"testing"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAsyncAction struct {
	mock.Mock
}

func (m *MockAsyncAction) IOMetadata() *data.IOMetadata {
	return nil
}

func (m *MockAsyncAction) Config() *action.Config {
	return nil
}

func (m *MockAsyncAction) Metadata() *action.Metadata {
	return nil
}

func (m *MockAsyncAction) Run(context context.Context, inputs map[string]*data.Attribute, handler action.ResultHandler) error {
	args := m.Called(context, inputs, handler)
	if handler != nil {
		dataAttr, _ := data.NewAttribute("data", data.TypeString, "mock")
		codeAttr, _ := data.NewAttribute("code", data.TypeInteger, 200)
		resultData := map[string]*data.Attribute{
			"data": dataAttr,
			"code": codeAttr,
		}
		handler.HandleResult(resultData, nil)
		handler.Done()
	}
	return args.Error(0)
}

//Test that Result returns the expected values
func TestResultOk(t *testing.T) {

	//mockData,_ :=data.CoerceToObject("{\"data\":\"mock data \"}")
	dataAttr, _ := data.NewAttribute("data", data.TypeString, "mock data")
	codeAttr, _ := data.NewAttribute("code", data.TypeInteger, 1)
	resultData := map[string]*data.Attribute{
		"data": dataAttr,
		"code": codeAttr,
	}

	rh := &SyncResultHandler{resultData: resultData, err: errors.New("New Error")}
	data, err := rh.Result()
	assert.Equal(t, 1, data["code"].Value())
	assert.Equal(t, "mock data", data["data"].Value())
	assert.NotNil(t, err)
}

//Test Direct Start method
func TestDirectStartOk(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	err := runner.Start()
	assert.Nil(t, err)
}

//Test Stop method
func TestDirectStopOk(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	err := runner.Stop()
	assert.Nil(t, err)
}

//Test Run method with a nil action
func TestDirectRunNilAction(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	_, err := runner.Execute(nil, nil, nil)
	assert.NotNil(t, err)
}

//Test Run method with error running action
func TestDirectRunErr(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	// Mock Action
	mockAction := new(MockAsyncAction)
	mockAction.On("Run", nil, mock.AnythingOfType("map[string]*data.Attribute"), mock.AnythingOfType("*runner.SyncResultHandler")).Return(errors.New("Action Error"))
	_, err := runner.Execute(nil, mockAction, nil)
	assert.NotNil(t, err)
}

//Test Run method ok
func TestDirectRunOk(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	// Mock Action
	mockAction := new(MockAsyncAction)

	mockAction.On("Run", nil, mock.AnythingOfType("map[string]*data.Attribute"), mock.AnythingOfType("*runner.SyncResultHandler")).Return(nil)
	results, err := runner.Execute(nil, mockAction, nil)
	assert.Nil(t, err)
	assert.NotNil(t, results)
	code, ok := results["code"]
	assert.True(t, ok)
	data, ok := results["data"]
	assert.True(t, ok)
	assert.Equal(t, 200, code.Value())
	assert.Equal(t, "mock", data.Value())
}

//Test Run method with a nil action
func TestDirectRunNilActionOld(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	_, err := runner.Execute(nil, nil, nil)
	assert.NotNil(t, err)
}

//Test Run method with error running action
func TestDirectRunErrOld(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	// Mock Action
	mockAction := new(MockAsyncAction)
	mockAction.On("Run", nil, mock.AnythingOfType("map[string]*data.Attribute"), mock.AnythingOfType("*runner.SyncResultHandler")).Return(errors.New("Action Error"))
	_, err := runner.Execute(nil, mockAction, nil)
	assert.NotNil(t, err)
}

//Test Run method ok
//func TestDirectRunOkOld(t *testing.T) {
//	runner := NewDirect()
//	assert.NotNil(t, runner)
//	// Mock Action
//	mockAction := new(MockAsyncAction)
//
//	mockAction.On("Run", nil, mock.AnythingOfType("map[string]*data.Attribute"), mock.AnythingOfType("*runner.SyncResultHandler")).Return(nil)
//	code, data, err := runner.Execute(nil, mockAction,  nil)
//	assert.Nil(t, err)
//	assert.Equal(t, 200, code)
//	assert.Equal(t, "mock", data)
//}

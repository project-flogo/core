package runner

import (
	"context"
	"errors"
	"testing"

	"github.com/project-flogo/core/data/metadata"

	"github.com/project-flogo/core/action"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAsyncAction struct {
	mock.Mock
}

func (m *MockAsyncAction) IOMetadata() *metadata.IOMetadata {
	return nil
}

func (m *MockAsyncAction) Config() *action.Config {
	return nil
}

func (m *MockAsyncAction) Metadata() *action.Metadata {
	return nil
}

func (m *MockAsyncAction) Run(context context.Context, inputs map[string]interface{}, handler action.ResultHandler) error {
	args := m.Called(context, inputs, handler)

	if handler != nil {
		resultData := make(map[string]interface{})
		resultData["data"] = "mock"
		resultData["code"] = 200
		handler.HandleResult(resultData, nil)
		handler.Done()
	}

	return args.Error(0)
}

//Test that Result returns the expected values
func TestResultOk(t *testing.T) {

	//mockData,_ :=data.CoerceToObject("{\"data\":\"mock data \"}")

	resultData := make(map[string]interface{})
	resultData["data"] = "mock data"
	resultData["code"] = 1

	rh := &SyncResultHandler{resultData: resultData, err: errors.New("new error")}
	data, err := rh.Result()
	assert.Equal(t, 1, data["code"])
	assert.Equal(t, "mock data", data["data"])
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
	_, err := runner.RunAction(nil, nil, nil)
	assert.NotNil(t, err)
}

//Test Run method with error running action
func TestDirectRunErr(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	// Mock Action
	mockAction := new(MockAsyncAction)

	mockAction.On("Run", nil, mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("*runner.SyncResultHandler")).Return(errors.New("action error"))
	_, err := runner.RunAction(nil, mockAction, nil)
	assert.NotNil(t, err)
}

//Test Run method ok
func TestDirectRunOk(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	// Mock Action
	mockAction := new(MockAsyncAction)

	mockAction.On("Run", context.Background(), mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("*runner.SyncResultHandler")).Return(nil)
	results, err := runner.RunAction(context.Background(), mockAction, nil)
	assert.Nil(t, err)
	assert.NotNil(t, results)
	code, ok := results["code"]
	assert.True(t, ok)
	data, ok := results["data"]
	assert.True(t, ok)
	assert.Equal(t, 200, code)
	assert.Equal(t, "mock", data)
}

//Test Run method with a nil action
func TestDirectRunNilActionOld(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	_, err := runner.RunAction(context.Background(), nil, nil)
	assert.NotNil(t, err)
}

//Test Run method with error running action
func TestDirectRunErrOld(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	// Mock Action
	mockAction := new(MockAsyncAction)
	mockAction.On("Run", context.Background(), mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("*runner.SyncResultHandler")).Return(errors.New("action error"))
	_, err := runner.RunAction(context.Background(), mockAction, nil)
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

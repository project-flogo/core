package runner

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

//TestWorkerInvalidRequestType worker returns error for invalid request type
func TestWorkerInvalidRequestType(t *testing.T) {
	worker := createDefaultWorker()
	worker.Start()

	rc := make(chan *ActionResult)
	actionData := &ActionData{arc: rc}

	// Create some work
	invalidWorkRequest := ActionWorkRequest{ReqType: -1, actionData: actionData}

	// Send some work
	worker.Work <- invalidWorkRequest

	// Check work result
	result := <-actionData.arc

	assert.NotNil(t, result.err)
	assert.Equal(t, "unsupported work request type: '-1'", result.err.Error())
}

//TestWorkerErrorInAction returns an error when the action returns error
func TestWorkerErrorInAction(t *testing.T) {
	worker := createDefaultWorker()
	worker.Start()

	rc := make(chan *ActionResult)

	action := new(MockFullAction)
	action.On("Run", nil, mock.AnythingOfType("map[string]*data.Attribute"), mock.AnythingOfType("*runner.AsyncResultHandler")).Return(errors.New("Error in action"))

	actionData := &ActionData{arc: rc, action: action}

	// Create some work
	errorWorkRequest := ActionWorkRequest{ReqType: RtRun, actionData: actionData}

	// Send some work
	worker.Work <- errorWorkRequest

	// Check work result
	result := <-actionData.arc

	assert.NotNil(t, result.err)
	assert.Equal(t, "Error in action", result.err.Error())
}

//TestWorkerStartOk
func TestWorkerStartOk(t *testing.T) {
	worker := createDefaultWorker()
	worker.Start()

	rc := make(chan *ActionResult)

	action := new(MockResultAction)
	action.On("Run", nil, mock.AnythingOfType("map[string]*data.Attribute"), mock.AnythingOfType("*runner.AsyncResultHandler")).Return(nil)

	actionData := &ActionData{arc: rc, action: action}

	// Create some work
	okWorkRequest := ActionWorkRequest{ReqType: RtRun, actionData: actionData}

	// Send some work
	worker.Work <- okWorkRequest

	// Check work result
	result := <-actionData.arc

	assert.Nil(t, result.err)
	assert.NotNil(t, result)
	assert.Equal(t, 200, result.results["code"].Value())
	assert.Equal(t, "mock", result.results["data"].Value())
}

func createDefaultWorker() ActionWorker {
	runner := NewDirect()
	queue := make(chan chan ActionWorkRequest, 2)
	return NewWorker(1, runner, queue)
}

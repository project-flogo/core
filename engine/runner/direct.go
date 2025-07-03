package runner

import (
	"context"
	"errors"
	"fmt"
	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/engine/runner/debugger"
	coreSupport "github.com/project-flogo/core/engine/support"
	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/trigger"
)

// DirectRunner runs an action synchronously
type DirectRunner struct {
	debugMode bool
	mockFile  string
	index     int
	appName   string
	appVer    string
}

var idGenerator *support.Generator

// NewDirectRunner create a new DirectRunner
func NewDirect() *DirectRunner {
	return &DirectRunner{}
}

// NewDirectRunner create a new DirectRunner
func NewDirectWithDebug(debugMode bool, mockFile string, appName string, appVer string) *DirectRunner {
	return &DirectRunner{
		debugMode: debugMode,
		mockFile:  mockFile,
		appName:   appName,
		appVer:    appVer,
	}
}

// Start will start the engine, by starting all of its workers
func (runner *DirectRunner) Start() error {
	//op-op
	return nil
}

// Stop will stop the engine, by stopping all of its workers
func (runner *DirectRunner) Stop() error {
	// check if all actions done till waiting time
	trackDirectRunnerActions.gracefulStop()
	return nil
}

var trackDirectRunnerActions = NewRunnerTracker()

// Execute implements action.Runner.Execute
func (runner *DirectRunner) RunAction(ctx context.Context, act action.Action, inputs map[string]interface{}) (results map[string]interface{}, err error) {

	if idGenerator == nil {
		idGenerator, _ = support.NewGenerator()
	}
	if act == nil {
		return nil, errors.New("action not specified")
	}

	config := inputs["_handler_config"]
	handlerConfig, _ := config.(*trigger.HandlerConfig)

	delete(inputs, "_handler_config")
	var tasks []*coreSupport.TaskInterceptor
	var coverage *coreSupport.Coverage
	var ro *coreSupport.DebugOptions

	if runner.debugMode {
		tasks = []*coreSupport.TaskInterceptor{}
		coverage = &coreSupport.Coverage{
			ActivityCoverage:   make([]*coreSupport.ActivityCoverage, 0),
			TransitionCoverage: make([]*coreSupport.TransitionCoverage, 0),
			SubFlowCoverage:    make([]*coreSupport.SubFlowCoverage, 0),
		}
		interceptor := &coreSupport.Interceptor{TaskInterceptors: tasks, Coverage: coverage, CollectIO: true}

		execOptions := &coreSupport.DebugExecOptions{Interceptor: interceptor}
		ro = &coreSupport.DebugOptions{ExecOptions: execOptions, InstanceId: idGenerator.NextAsString()}
		inputs["_run_options"] = ro
	}

	trackDirectRunnerActions.AddRunner()
	defer trackDirectRunnerActions.RemoveRunner()
	if syncAct, ok := act.(action.SyncAction); ok {
		return syncAct.Run(ctx, inputs)
	} else if asyncAct, ok := act.(action.AsyncAction); ok {
		handler := &SyncResultHandler{done: make(chan bool, 1)}

		err = asyncAct.Run(ctx, inputs, handler)

		if err != nil {
			return nil, err
		}

		<-handler.done

		if runner.debugMode {

			outputs := handler.resultData
			debugger.GenerateReport(handlerConfig, tasks, coverage, ro.InstanceId, inputs, outputs, runner.appName, runner.appVer)
		}

		runner.index++
		return handler.Result()
	} else {
		return nil, fmt.Errorf("unsupported action: %v", act)
	}
}

// SyncResultHandler simple result handler to use in synchronous case
type SyncResultHandler struct {
	done       chan bool
	resultData map[string]interface{}
	err        error
	set        bool
}

// HandleResult implements action.ResultHandler.HandleResult
func (rh *SyncResultHandler) HandleResult(resultData map[string]interface{}, err error) {

	if !rh.set {
		rh.set = true
		rh.resultData = resultData
		rh.err = err
	}
}

// Done implements action.ResultHandler.Done
func (rh *SyncResultHandler) Done() {
	rh.done <- true
}

// Result returns the latest Result set on the handler
func (rh *SyncResultHandler) Result() (resultData map[string]interface{}, err error) {
	return rh.resultData, rh.err
}

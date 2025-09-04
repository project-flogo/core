package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/engine/runner/debugger"
	coreSupport "github.com/project-flogo/core/engine/support"
	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
	"os"
	"path"
	"path/filepath"
)

// DirectRunner runs an action synchronously
type DirectRunner struct {
	debugMode   bool
	mockFile    string
	genMockFile bool
	outputPath  string
	index       int
	mockData    *coreSupport.MockReport
	appPath     string
}

var idGenerator *support.Generator

// NewDirectRunner create a new DirectRunner
func NewDirect() *DirectRunner {
	return &DirectRunner{}
}

// NewDirectRunner create a new DirectRunner
func NewDirectWithDebug(debugMode bool, mockFile string, outputPath string, genMock bool, appPath string) *DirectRunner {
	return &DirectRunner{
		debugMode:   debugMode,
		mockFile:    mockFile,
		genMockFile: genMock,
		outputPath:  outputPath,
		appPath:     appPath,
	}
}

// Start will start the engine, by starting all of its workers
func (runner *DirectRunner) Start() error {
	if runner.debugMode {
		reportPath := runner.outputPath
		if reportPath == "" {
			reportPath = os.Getenv("FLOW_EXECUTION_FILES")
		}

		if reportPath == "" {
			reportPath = path.Join(os.TempDir(), "flow-executions")
		}
		reportPath = filepath.Join(reportPath, debugger.GetAppName())

		log.RootLogger().Infof("Generate Report for Flow Execution: %s", reportPath)

		os.RemoveAll(reportPath)
	}

	if runner.mockFile != "" {
		content, err := os.ReadFile(runner.mockFile)
		if err != nil {
			return err
		}
		var report map[string]interface{}
		if err := json.Unmarshal(content, &report); err != nil {
			return err
		}

		var mockResult = report["mocks"].(map[string]interface{})

		var flows = mockResult["flows"].(map[string]interface{})
		mockReport := &coreSupport.MockReport{
			make(map[string]*coreSupport.FlowMock),
		}
		for _, val := range flows {
			flow := val.(map[string]interface{})
			flowMock := &coreSupport.FlowMock{}
			flowMock.Name = flow["flowName"].(string)
			activities := flow["activities"].([]interface{})
			actList := make([]*coreSupport.ActivityMock, 0)
			for _, val := range activities {
				actMap := val.(map[string]interface{})
				activity := &coreSupport.ActivityMock{}
				activity.ActivityName = actMap["name"].(string)
				activity.MockType = 1
				activity.Mock = actMap["mock"]
				actList = append(actList, activity)
			}
			flowMock.ActivityReport = actList
			mockReport.Flows[flowMock.Name] = flowMock
		}
		runner.mockData = mockReport
	}

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
			SubFlowMap:         make(map[string]*coreSupport.SubFlowCoverage),
		}

		//if runner.mockData != nil && runner.mockData.Flows != nil {
		//	for _, flow := range runner.mockData.Flows {
		//		for _, activity := range flow.ActivityReport {
		//			interceptor := &coreSupport.TaskInterceptor{}
		//			interceptor.ID = flow.Name + "-" + activity.ActivityName
		//			interceptor.Type = coreSupport.MockActivity
		//			interceptor.Skip = true
		//			interceptor.SkipExecution = true
		//			if activity.Mock != nil {
		//				interceptor.Outputs = activity.Mock.(map[string]interface{})
		//			}
		//			tasks = append(tasks, interceptor)
		//		}
		//	}
		//}

		interceptor := &coreSupport.Interceptor{TaskInterceptors: tasks, Coverage: coverage, CollectIO: true}

		execOptions := &coreSupport.DebugExecOptions{Interceptor: interceptor}
		instanceId := idGenerator.NextAsString()
		ro = &coreSupport.DebugOptions{ExecOptions: execOptions, InstanceId: instanceId}

		inputs["_run_options"] = ro
	}
	log.RootLogger().Infof("Executing flow with instanceId %s", ro.InstanceId)
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
			log.RootLogger().Infof("Flow execution completed for instanceId %s", ro.InstanceId)

			var outputs map[string]interface{}
			if handler.resultData != nil {
				outputs = handler.resultData
			} else if handler.err != nil {
				outputs = convertErrorToMap(handler.err)
			}
			debugger.GenerateReport(handlerConfig, tasks, coverage, ro.InstanceId, inputs, outputs, runner.outputPath, runner.appPath)
		}
		if runner.genMockFile {
			debugger.GenerateMock(coverage, runner.outputPath)
		}

		runner.index++
		return handler.Result()
	} else {
		return nil, fmt.Errorf("unsupported action: %v", act)
	}
}

func convertErrorToMap(handlerErr error) map[string]interface{} {
	if handlerErr == nil {
		return nil
	}

	if activityErr, ok := handlerErr.(*activity.Error); ok {
		return activityErrorToMap(activityErr)
	}

	return map[string]interface{}{
		"error": handlerErr.Error(),
		"type":  fmt.Sprintf("%T", handlerErr),
	}
}

func activityErrorToMap(activityErr *activity.Error) map[string]interface{} {
	return map[string]interface{}{
		"error":        activityErr.Error(),
		"activityName": activityErr.ActivityName(),
		"errorCode":    activityErr.Code(),
		"category":     activityErr.Category(),
		"retriable":    activityErr.Retriable(),
		"data":         activityErr.Data(),
		"type":         fmt.Sprintf("%T", activityErr),
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

package runner

import (
	"context"
	"errors"
	"github.com/project-flogo/core/app"

	"sync"
	"time"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/log"
)

// PooledRunner is a action runner that queues and runs a action in a worker pool
type PooledRunner struct {
	workerQueue chan chan ActionWorkRequest
	workQueue   chan ActionWorkRequest
	numWorkers  int
	workers     []*ActionWorker
	active      bool
	logger      log.Logger

	directRunner *DirectRunner
}

// PooledConfig is the configuration object for a PooledRunner
type PooledConfig struct {
	NumWorkers    int `json:"numWorkers"`
	WorkQueueSize int `json:"workQueueSize"`
}

// NewPooledRunner create a new pooled
func NewPooled(config *PooledConfig) *PooledRunner {

	var pooledRunner PooledRunner
	pooledRunner.directRunner = NewDirect()

	// config via engine config
	pooledRunner.numWorkers = config.NumWorkers
	pooledRunner.workQueue = make(chan ActionWorkRequest, config.WorkQueueSize)

	//todo should this be root logger or engine logger?
	pooledRunner.logger = log.RootLogger()

	return &pooledRunner
}

var trackActions sync.WaitGroup

// Start will start the engine, by starting all of its workers
func (runner *PooledRunner) Start() error {

	logger := runner.logger

	if !runner.active {

		runner.workerQueue = make(chan chan ActionWorkRequest, runner.numWorkers)

		runner.workers = make([]*ActionWorker, runner.numWorkers)

		for i := 0; i < runner.numWorkers; i++ {
			id := i + 1
			logger.Debugf("Starting worker with id '%d'", id)
			worker := NewWorker(id, runner.directRunner, runner.workerQueue)
			runner.workers[i] = &worker
			trackActions.Add(1)
			worker.Start()
		}

		go func() {
			for {
				select {
				case work := <-runner.workQueue:
					logger.Debug("Received work request")

					//todo fix, this creates unbounded go routines waiting to be serviced by worker queue
					go func() {
						worker := <-runner.workerQueue

						logger.Debug("Dispatching work request")
						worker <- work
					}()
				}
			}
		}()

		runner.active = true
	}

	return nil
}

// Stop will stop the engine, by stopping all of its workers
func (runner *PooledRunner) Stop() error {

	if runner.active {

		runner.active = false

		for _, worker := range runner.workers {
			runner.logger.Debug("Stopping worker", worker.ID)
			worker.Stop()
		}
		// check if all actions done till shutdown waiting time
		gracefulStop()
	}

	return nil
}

func gracefulStop() {
	logger := log.RootLogger()
	delayedStopInterval := app.GetDelayedStopInterval()
	if delayedStopInterval != "" {
		duration, err := time.ParseDuration(delayedStopInterval)
		if err != nil {
			logger.Errorf("Invalid interval - %s  specified for delayed stop. It must suffix with time unit e.g. %sms, %ss", delayedStopInterval, delayedStopInterval, delayedStopInterval)
		} else {
			logger.Infof("Delaying application stop by max - %s", delayedStopInterval)
			if isTimeout := waitForActionsCompletion(duration); isTimeout {
				logger.Info("All actions not completed before engine shutdown")
			} else {
				logger.Info("All actions completed before engine shutdown")
			}
		}

	}
}

func waitForActionsCompletion(timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		trackActions.Wait()
	}()
	select {
	case <-c:
		return false // actions completed
	case <-time.After(timeout):
		return true // timed out
	}
}

// Execute implements action.Runner.Execute
func (runner *PooledRunner) RunAction(ctx context.Context, act action.Action, inputs map[string]interface{}) (results map[string]interface{}, err error) {

	logger := runner.logger

	if act == nil {
		return nil, errors.New("action not specified")
	}

	if runner.active {

		actionData := &ActionData{context: ctx, action: act, inputs: inputs, arc: make(chan *ActionResult, 1)}
		work := ActionWorkRequest{ReqType: RtRun, actionData: actionData}

		runner.workQueue <- work

		if logger.DebugEnabled() {
			logger.Debugf("Action '%s' queued", support.GetRef(act))
		}

		reply := <-actionData.arc

		if logger.DebugEnabled() {
			logger.Debugf("Action '%s' returned", support.GetRef(act))
		}

		return reply.results, reply.err
	}

	//Run rejected
	return nil, errors.New("runner not active")
}

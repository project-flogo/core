package app

import (
	"fmt"
	"sync"

	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)

const (
	AlreadyControlled = "app is already event flow controlled"
)

var controller Controller
var logger = log.ChildLogger(log.RootLogger(), "events.controller")

type Controller interface {
	StartControl() error
	ReleaseControl() error
}

type controllerData struct {
	flowControlled bool
	triggers       map[string]trigger.Trigger
	lock           sync.Mutex
}

func GetEventFlowController() Controller {
	return controller
}

// StartControl uses to start control the controller, the evaluator must call start control first then release control
func (c *controllerData) StartControl() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.flowControlled {
		return fmt.Errorf(AlreadyControlled)
	} else {
		// Pause trigger
		c.flowControlled = true
		err := c.pauseTriggers()
		if err != nil {
			errMsg := fmt.Errorf("error pausing triggers: %s", err.Error())
			logger.Error(errMsg)
			return errMsg
		}
		return nil
	}
}

// ReleaseControl uses to release control the controller
func (c *controllerData) ReleaseControl() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.flowControlled {
		err := c.resumeTriggers()
		if err != nil {
			// Release control if error occurred here
			c.flowControlled = false
			errMsg := fmt.Errorf("error resume triggers: %s", err.Error())
			logger.Error(errMsg)
			return errMsg
		}
		c.flowControlled = false
	}
	return nil
}

func (a *App) initEventFlowController() {
	controllerData := &controllerData{lock: sync.Mutex{}}
	controllerData.triggers = make(map[string]trigger.Trigger)
	for _, trgW := range a.triggers {
		controllerData.triggers[trgW.id] = trgW.trg
	}
	controller = controllerData
}

// Start triggers
func (c *controllerData) resumeTriggers() error {
	// Resume  triggers
	logger.Info("Resuming Triggers...")
	for id, trg := range c.triggers {
		var err error
		if flowControlAware, ok := trg.(trigger.EventFlowControlAware); ok {
			err = flowControlAware.Resume()
		}

		if err != nil {
			//return err
			//TODO Starting other triggers. Should we stop the app here?
			logger.Errorf("Trigger [%s] failed to resume due to error - %s.", id, err.Error())
			continue
		}
		logger.Infof("Trigger [%s] is resumed.", id)
	}
	logger.Info("Triggers are resumed")
	return nil
}

// Stop triggers
func (c *controllerData) pauseTriggers() error {
	logger.Info("Pausing Triggers...")
	// Pause Triggers
	for id, trg := range c.triggers {
		var err error
		if flowControlAware, ok := trg.(trigger.EventFlowControlAware); ok {
			err = flowControlAware.Pause()
		}
		if err != nil {
			//return err
			//TODO Stopping other triggers. Should we stop the app here?
			logger.Errorf("Trigger [%s] failed to pause due to error - %s.", id, err.Error())
			continue
		}
		logger.Infof("Trigger [%s] is paused.", id)
	}
	logger.Info("Triggers are paused")
	return nil
}

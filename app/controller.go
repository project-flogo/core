package app

import (
	"fmt"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
	"sync"
)

const (
	AlreadyControlled = "app is already controlled"
)

var controller Controller

type Controller interface {
	StartControl() error
	ReleaseControl() error
}

type controllerData struct {
	flowControlled bool
	triggers       map[string]trigger.FlowControlAware
	lock           sync.Mutex
}

func GetFlowController() Controller {
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
			log.RootLogger().Error(errMsg)
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
			log.RootLogger().Error(errMsg)
			return errMsg
		}
		c.flowControlled = false
	}
	return nil
}

func (app *App) initFlowController() {
	controllerData := &controllerData{lock: sync.Mutex{}}
	controllerData.triggers = make(map[string]trigger.FlowControlAware)
	for id, trgW := range app.triggers {
		if t, ok := trgW.trg.(trigger.FlowControlAware); ok {
			controllerData.triggers[id] = t
		}
	}
	controller = controllerData
}

// Resume triggers
func (c *controllerData) resumeTriggers() error {
	// Resume  triggers
	log.RootLogger().Info("Resuming Triggers...")
	for id, trg := range c.triggers {
		err := trg.Resume()
		if err != nil {
			//return err
			//TODO Letting other triggers resume. Should we stop the app here?
			log.RootLogger().Errorf("Trigger [%s] failed to resume due to error - %s.", id, err.Error())
			continue
		}
		log.RootLogger().Infof("Trigger [%s] is resumed.", id)
	}
	log.RootLogger().Info("Triggers Resumed")
	return nil
}

// Pause triggers
func (c *controllerData) pauseTriggers() error {
	log.RootLogger().Info("Pausing Triggers...")
	// Pause Triggers
	for id, trg := range c.triggers {
		err := trg.Pause()
		if err != nil {
			//return err
			//TODO Letting other triggers pause. Should we stop the app here?
			log.RootLogger().Errorf("Trigger [%s] failed to pause due to error - %s.", id, err.Error())
			continue
		}
		log.RootLogger().Infof("Trigger [%s] is paused.", id)
	}
	log.RootLogger().Info("Triggers Paused")
	return nil
}

package app

import (
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/managed"
	"github.com/project-flogo/core/trigger"
)

// Interface for controlling application inbound events
type FlowController interface {
	// Resume/start triggers when engine comes out of the flow control mode
	StartTriggers() error

	// Pause/stop triggers when engine enters into flow control mode
	StopTriggers() error
}


 var controller FlowController
type controllerType struct {
	app *App
}
func GetFlowController() FlowController {
	return controller
}

func setController(app *App)  {
	controller = &controllerType{app}
}

// start/resume triggers
func (c *controllerType) StartTriggers() error {

	normalTriggers := make(map[string]*triggerWrapper)

	for id, trgW := range c.app.triggers {
		if _, ok := trgW.trg.(LifecycleAware); !ok {
			normalTriggers[id] = trgW
		}
	}

	// Start the triggers
	log.RootLogger().Info("Restarting Triggers...")

	// Start normal triggers
	for id, trgW := range normalTriggers {
		_, err := c.app.startTrigger(id, trgW)
		if err != nil {
			return err
		}
	}
	log.RootLogger().Info("Triggers Restarted")
	return nil
}

// stop/pause triggers
func (c *controllerType) StopTriggers() error {
	log.RootLogger().Info("Stopping Triggers...")

	normalTriggers := make(map[string]*triggerWrapper)
	for id, trgW := range c.app.triggers {
		if _, ok := trgW.trg.(LifecycleAware); !ok {
			normalTriggers[id] = trgW
		}
	}
	// Stop Normal Triggers
	for id, trg := range normalTriggers {
		_ = managed.Stop("Trigger [ "+id+" ]", trg.trg)
		trg.status.Status = managed.StatusStopped
		trigger.PostTriggerEvent(trigger.STOPPED, id)
	}
	log.RootLogger().Info("Triggers Stopped")
	return nil
}

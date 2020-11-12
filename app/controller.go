package app

import (
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)

var controller controllerData

type controllerData struct {
	flowControlled bool
	notify         chan bool
	triggers       map[string]trigger.FlowControlAware
}

func (c controllerData) startController() {
	for {
		select {
		case v := <-c.notify:
			if len(c.triggers) > 0 {
				if v && !c.flowControlled {
					_ = c.pauseTriggers()
				} else if !v && c.flowControlled {
					_ = c.resumeTriggers()
				}
			}
		}
	}
}

func GetFlowController() <-chan bool {
	return controller.notify
}

func (app *App) enableFlowController() {
	controller = controllerData{}
	controller.notify = make(chan bool)
	controller.triggers = make(map[string]trigger.FlowControlAware)
	for id, trgW := range app.triggers {
		if t, ok := trgW.trg.(trigger.FlowControlAware); ok {
			controller.triggers[id] = t
		}
	}
	if len(controller.triggers) > 0 {
		go controller.startController()
	}
}

// Resume triggers
func (c controllerData) resumeTriggers() error {
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
	c.flowControlled = false
	return nil
}

// Pause triggers
func (c controllerData) pauseTriggers() error {
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
	c.flowControlled = true
	return nil
}

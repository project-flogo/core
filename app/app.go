package app

import (
	"fmt"
	"runtime/debug"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/app/resource"
	"github.com/project-flogo/core/data/property"
	"github.com/project-flogo/core/support/logger"
	"github.com/project-flogo/core/support/managed"
	"github.com/project-flogo/core/trigger"
)

type Option func(*App) error

func New(config *Config, runner action.Runner, options ...Option) (*App, error) {

	app := &App{stopOnError: true, name: config.Name, version: config.Version}

	properties := make(map[string]interface{}, len(config.Properties))
	for _, attr := range config.Properties {
		properties[attr.Name()] = attr.Value()
	}

	app.propManager = property.NewManager(properties)

	for _, option := range options {
		option(app)
	}

	resources := make(map[string]*resource.Resource, len(config.Resources))
	app.resManager = resource.NewManager(resources)

	for _, actionFactory := range action.Factories() {
		actionFactory.Initialize(app)
	}

	for _, resConfig := range config.Resources {
		resType, err := resource.GetTypeFromID(resConfig.ID)
		if err != nil {
			return nil, err
		}

		loader := resource.GetLoader(resType)
		res, err := loader.LoadResource(resConfig)
		if err != nil {
			return nil, err
		}

		resources[resConfig.ID] = res
	}

	var err error

	app.actions, err = app.createSharedActions(config.Actions)
	if err != nil {
		return nil, fmt.Errorf("error creating shared action instances - %s", err.Error())
	}

	app.triggers, err = app.createTriggers(config.Triggers, runner)
	if err != nil {
		return nil, fmt.Errorf("error creating trigger instances - %s", err.Error())
	}

	for _, option := range options {
		option(app)
	}

	return app, nil
}

func ContinueOnError(a *App) error {
	a.stopOnError = false
	return nil
}

func ExternalProperties(providerId string, overrides string, processors ...property.PostProcessor) func(*App) error {
	return func(a *App) error {
		return a.propManager.AddExternalProperties(providerId, overrides, processors...)
	}
}

type App struct {
	name        string
	version     string
	propManager *property.Manager
	resManager  *resource.Manager
	actions     map[string]action.Action
	triggers    map[string]*triggerWrapper
	stopOnError bool
	started     bool
}

type triggerWrapper struct {
	ref    string
	trg    trigger.Trigger
	status *managed.StatusInfo
}

func (a *App) GetProperty(name string) (interface{}, bool) {
	return a.propManager.GetProperty(name)
}

func (a *App) GetResource(id string) *resource.Resource {
	return a.resManager.GetResource(id)
}

func (a *App) ResourceManager() *resource.Manager {
	return a.resManager
}

func (a *App) Name() interface{} {
	return a.name
}

func (a *App) Version() interface{} {
	return a.version
}

// TriggerStatuses gets the status information for the triggers
func (a *App) TriggerStatuses() []*managed.StatusInfo {
	statuses := make([]*managed.StatusInfo, 0, len(a.triggers))
	for _, trg := range a.triggers {
		statuses = append(statuses, trg.status)
	}

	return statuses
}

func (a *App) Start() error {

	if a.started {
		return fmt.Errorf("app already started")
	}

	// Start the triggers
	logger.Info("Starting Triggers...")

	var failed []string

	for id, trg := range a.triggers {
		statusInfo := trg.status
		err := managed.Start(fmt.Sprintf("Trigger [ %s ]", id), trg.trg)
		if err != nil {
			if a.stopOnError {
				return fmt.Errorf("error starting Trigger[%s] : %s", id, err)
			}
			logger.Infof("Trigger [%s] failed to start due to error [%s]", id, err.Error())
			statusInfo.Status = managed.StatusFailed
			statusInfo.Error = err
			logger.Debugf("StackTrace: %s", debug.Stack())
			failed = append(failed, id)
		} else {
			statusInfo.Status = managed.StatusStarted
			logger.Infof("Trigger [ %s ]: Started", id)
			version := ""
			logger.Debugf("Trigger [ %s ] has ref [ %s ] and version [ %s ]", id, trg.ref, version)
		}
	}

	if len(failed) > 0 {
		//remove failed trigger, we have no use for them - todo will cause a problem if we can start again
		for _, triggerId := range failed {
			delete(a.triggers, triggerId)
		}
	}

	logger.Info("Triggers Started")

	a.started = true
	return nil
}

func (a *App) Stop() error {

	logger.Info("Stopping Triggers...")

	// Stop Triggers
	for id, trg := range a.triggers {
		managed.Stop("Trigger [ "+id+" ]", trg.trg)
		trg.status.Status = managed.StatusStopped
	}

	logger.Info("Triggers Stopped")

	//a.active = false - this will allow restart
	return nil
}
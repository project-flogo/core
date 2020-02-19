package app

import (
	"fmt"
	"path"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/activity"
	appresolve "github.com/project-flogo/core/app/resolve"
	"github.com/project-flogo/core/app/resource"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/data/property"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/data/schema"
	"github.com/project-flogo/core/engine/channels"
	"github.com/project-flogo/core/engine/event"
	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/managed"
	"github.com/project-flogo/core/support/service"
	"github.com/project-flogo/core/trigger"
)

type Option func(*App) error

type Status string

const (
	FAILED   = "Failed"
	STARTED  = "Started"
	STARTING = "Starting"
	STOPPED  = "Stopped"
)

const AppEventType = "appevent"

type AppEvent interface {
	AppName() string
	AppVersion() string
	AppStatus() Status
}

type appEvent struct {
	status        Status
	name, version string
}

func (ae *appEvent) AppName() string {
	return ae.name
}

func (ae *appEvent) AppVersion() string {
	return ae.version
}

func (ae *appEvent) AppStatus() Status {
	return ae.status
}

var flogoImportPattern = regexp.MustCompile(`^(([^ ]*)[ ]+)?([^@:]*)@?([^:]*)?:?(.*)?$`) // extract import path even if there is an alias and/or a version

func New(config *Config, runner action.Runner, options ...Option) (*App, error) {

	app := &App{stopOnError: true, name: config.Name, version: config.Version}

	properties := make(map[string]interface{}, len(config.Properties))
	for _, attr := range config.Properties {
		properties[attr.Name()] = attr.Value()
	}

	app.propManager = property.NewManager(properties)
	property.SetDefaultManager(app.propManager)

	for _, option := range options {
		err := option(app)
		if err != nil {
			return nil, err
		}
	}

	if app.srvManager == nil {
		app.srvManager = service.NewServiceManager()
	}

	if app.actionSettings == nil {
		app.actionSettings = make(map[string]map[string]interface{})
	}

	channelDescriptors := config.Channels
	if len(channelDescriptors) > 0 {
		for _, descriptor := range channelDescriptors {
			name, buffSize := channels.Decode(descriptor)

			log.RootLogger().Debugf("Creating Engine Channel '%s'", name)

			_, err := channels.New(name, buffSize)
			if err != nil {
				return nil, err
			}
		}
	}

	resolver := resolve.NewCompositeResolver(map[string]resolve.Resolver{
		".":        &resolve.ScopeResolver{},
		"env":      &resolve.EnvResolver{},
		"property": &property.Resolver{},
		"loop":     &resolve.LoopResolver{},
	})

	app.resolver = resolver
	appresolve.SetAppResolver(resolver)

	for _, anImport := range config.Imports {
		matches := flogoImportPattern.FindStringSubmatch(anImport)
		err := registerImport(matches[1] + matches[3] + matches[5]) // alias + module path + relative import path
		if err != nil {
			log.RootLogger().Errorf("cannot register import '%s' : %v", anImport, err)
		}
	}

	function.ResolveAliases()

	// register schemas, assumes appropriate schema factories have been registered
	for id, def := range config.Schemas {
		_, err := schema.Register(id, def)
		if err != nil {
			return nil, err
		}
	}

	schema.ResolveSchemas()

	for id, config := range config.Connections {
		_, err := connection.NewSharedManager(id, config)
		if err != nil {
			return nil, err
		}
	}

	resources := make(map[string]*resource.Resource, len(config.Resources))
	app.resManager = resource.NewManager(resources)

	for ref, actionFactory := range action.Factories() {

		var initCtx action.InitContext

		if s, ok := app.actionSettings[ref]; ok {
			initCtx = &actionInitCtx{
				app:      app,
				settings: s,
			}
		} else {
			initCtx = app
		}

		err := actionFactory.Initialize(initCtx)
		if err != nil {
			return nil, err
		}
	}

	for _, resConfig := range config.Resources {
		resType, err := resource.GetTypeFromID(resConfig.ID)
		if err != nil {
			return nil, err
		}

		loader := resource.GetLoader(resType)
		if loader == nil {
			return nil, fmt.Errorf("resource loader for '%s' not registered", resType)
		}

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

	return app, nil
}

type actionInitCtx struct {
	app      *App
	settings map[string]interface{}
}

func (a actionInitCtx) ResourceManager() *resource.Manager {
	return a.app.resManager
}

func (a actionInitCtx) ServiceManager() *service.Manager {
	return a.app.srvManager
}

func (a actionInitCtx) RuntimeSettings() map[string]interface{} {
	return a.settings
}

func EngineSettings(svcManager *service.Manager, actionSettings map[string]map[string]interface{}) func(*App) error {
	return func(a *App) error {
		a.srvManager = svcManager
		a.actionSettings = actionSettings
		return nil
	}
}

func ContinueOnError(a *App) error {
	a.stopOnError = false
	return nil
}

func FinalizeProperties(processors ...property.PostProcessor) func(*App) error {
	return func(a *App) error {
		return a.propManager.Finalize(processors...)
	}
}

type App struct {
	name           string
	version        string
	propManager    *property.Manager
	resManager     *resource.Manager
	srvManager     *service.Manager
	actions        map[string]action.Action
	triggers       map[string]*triggerWrapper
	stopOnError    bool
	started        bool
	resolver       resolve.CompositeResolver
	actionSettings map[string]map[string]interface{}
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

func (a *App) ServiceManager() *service.Manager {
	return a.srvManager
}

func (a *App) RuntimeSettings() map[string]interface{} {
	return nil
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

	logger := log.RootLogger()

	managers := connection.Managers()

	if len(managers) > 0 {
		// Start the connection managers
		logger.Info("Starting Connection Managers...")

		for id, manager := range managers {
			if m, ok := manager.(managed.Managed); ok {
				err := m.Start()
				if err != nil {
					return fmt.Errorf("unable to start connection manager for '%s': %v", id, err)
				}
			}
		}

		logger.Info("Connection Managers Started")
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
			trigger.PostTriggerEvent(trigger.FAILED, id)
		} else {
			statusInfo.Status = managed.StatusStarted
			//logger.Infof("Trigger [ %s ]: Started", id)
			version := ""
			logger.Debugf("Trigger [ %s ] has ref [ %s ] and version [ %s ]", id, trg.ref, version)
			trigger.PostTriggerEvent(trigger.STARTED, id)
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

	logger := log.RootLogger()

	logger.Info("Stopping Triggers...")

	// Stop Triggers
	for id, trg := range a.triggers {
		_ = managed.Stop("Trigger [ "+id+" ]", trg.trg)
		trg.status.Status = managed.StatusStopped
		trigger.PostTriggerEvent(trigger.STOPPED, id)
	}

	logger.Info("Triggers Stopped")

	delayedStopInterval := GetDelayedStopInterval()
	if delayedStopInterval != "" {
		// Delay stopping of connection manager so that in-flight actions can continue until specified interval
		// No new events will be processed as triggers are stopped.
		duration, err := time.ParseDuration(delayedStopInterval)
		if err != nil {
			logger.Errorf("Invalid interval - %s  specified for delayed stop. It must suffix with time unit e.g. %sms, %ss", delayedStopInterval, delayedStopInterval, delayedStopInterval)
		} else {
			logger.Infof("Delaying application stop by - %s", delayedStopInterval)
			time.Sleep(duration)
		}

	}

	managers := connection.Managers()

	if len(managers) > 0 {
		// Stop the connection managers
		logger.Info("Stopping Connection Managers...")

		for id, manager := range managers {
			if m, ok := manager.(managed.Managed); ok {
				err := m.Stop()
				if err != nil {
					logger.Warnf("Unable to start connection manager for '%s': %v", id, err)
				}
			}
		}

		logger.Info("Connection Managers Stopped")
	}

	logger.Debugf("Cleaning up singleton activities")
	activity.CleanupSingletons()

	logger.Debugf("Cleaning up resources")
	a.resManager.CleanupResources()

	return nil
}

func registerImport(anImport string) error {

	parts := strings.Split(anImport, " ")

	var alias string
	var ref string
	numParts := len(parts)
	if numParts == 1 {
		ref = parts[0]
		alias = path.Base(ref)
	} else if numParts == 2 {
		alias = parts[0]
		ref = parts[1]
	} else {
		return fmt.Errorf("invalid import %s", anImport)
	}

	if alias == "" || ref == "" {
		return fmt.Errorf("invalid import %s", anImport)
	}

	ct := getContribType(ref)
	if ct == "other" {
		log.RootLogger().Debugf("Added Non-Contribution Import: %s", ref)
		return nil
		//return fmt.Errorf("invalid import, contribution '%s' not registered", anImport)
	}

	log.RootLogger().Debugf("Registering type alias '%s' for %s [%s]", alias, ct, ref)

	err := support.RegisterAlias(ct, alias, ref)
	if err != nil {
		return err
	}

	if ct == "function" {
		function.SetPackageAlias(ref, alias)
	}

	return nil
}

func getContribType(ref string) string {

	if activity.Get(ref) != nil {
		return "activity"
	} else if action.GetFactory(ref) != nil {
		return "action"
	} else if trigger.GetFactory(ref) != nil {
		return "trigger"
	} else if function.IsFunctionPackage(ref) {
		return "function"
	} else if connection.GetManagerFactory(ref) != nil {
		return "connection"
	} else {
		return "other"
	}
}

func (a *App) PostAppEvent(appStatus Status) {
	if event.HasListener(AppEventType) {
		ae := &appEvent{name: a.name, version: a.version, status: appStatus}
		event.Post(AppEventType, ae)
	}
}

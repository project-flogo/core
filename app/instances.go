package app

import (
	"fmt"
	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/managed"
	"github.com/project-flogo/core/trigger"
)

func (a *App) createSharedActions(actionConfigs []*action.Config) (map[string]action.Action, error) {

	actions := make(map[string]action.Action)

	for _, config := range actionConfigs {

		if config.Ref == "" {
			var ok bool
			config.Ref, ok =support.GetAliasRef("action", config.Type)
			if !ok {
				return nil, fmt.Errorf("Action type '%s' not registered", config.Type)
			}
		}

		actionFactory := action.GetFactory(config.Ref)
		if actionFactory == nil {
			return nil, fmt.Errorf("Action Factory '%s' not registered", config.Ref)
		}

		act, err := actionFactory.New(config)
		if err != nil {
			return nil, err
		}

		actions[config.Id] = act
	}

	return actions, nil
}

func (a *App) createTriggers(tConfigs []*trigger.Config, runner action.Runner) (map[string]*triggerWrapper, error) {

	triggers := make(map[string]*triggerWrapper)

	mapperFactory := mapper.NewFactory(resolve.GetBasicResolver())
	expressionFactory := expression.NewFactory(resolve.GetBasicResolver())
	for _, tConfig := range tConfigs {

		tConfig.AppConfig = map[string]interface{}{"Name": a.name, "Version": a.version, "Description": a.description}

		_, exists := triggers[tConfig.Id]
		if exists {
			return nil, fmt.Errorf("Trigger with id '%s' already registered, trigger ids have to be unique", tConfig.Id)
		}


		if tConfig.Ref == "" {
			var ok bool
			tConfig.Ref, ok =support.GetAliasRef("trigger", tConfig.Type)
			if !ok {
				return nil, fmt.Errorf("Trigger type '%s' not registered", tConfig.Type)
			}
		}

		triggerFactory := trigger.GetFactory(tConfig.Ref)

		if triggerFactory == nil {
			return nil, fmt.Errorf("Trigger Factory '%s' not registered", tConfig.Ref)
		}

		tConfig.FixUp(triggerFactory.Metadata())

		trg, err := triggerFactory.New(tConfig)

		if err != nil {
			return nil, err
		}

		if trg == nil {
			return nil, fmt.Errorf("cannot create Trigger nil for id '%s'", tConfig.Id)
		}

		logger := trigger.GetLogger(tConfig.Ref)

		if log.CtxLoggingEnabled() {
			logger = log.ChildLoggerWithFields(logger, log.String("triggerId", tConfig.Id))
		}

		log.ChildLogger(logger, tConfig.Id)
		initCtx := &initContext{logger: logger, handlers: make([]trigger.Handler, 0, len(tConfig.Handlers))}

		//create handlers for that trigger and init
		for _, hConfig := range tConfig.Handlers {

			var acts []action.Action
			var err error

			if len(hConfig.Actions) == 0 {
				return nil, fmt.Errorf("trigger '%s' has a handler with no action", tConfig.Id)
			}

			//use action if already associated with Handler
			for _, act := range hConfig.Actions {
				if act.Act != nil {
					acts = append(acts, act.Act)
				} else {
					if id := act.Id; id != "" {
						act, _ := a.actions[id]
						if act == nil {
							return nil, fmt.Errorf("shared Action '%s' does not exists", id)
						}
						acts = append(acts, act)
					} else {
						//create the action

						if act.Ref == "" {
							var ok bool
							act.Ref, ok =support.GetAliasRef("action", act.Type)
							if !ok {
								return nil, fmt.Errorf("Action type '%s' not registered", act.Type)
							}
						}

						actionFactory := action.GetFactory(act.Ref)
						if actionFactory == nil {
							return nil, fmt.Errorf("Action Factory '%s' not registered", act.Ref)
						}

						act, err := actionFactory.New(act.Config)
						if err != nil {
							return nil, err
						}
						acts = append(acts, act)
					}
				}
			}

			handler, err := trigger.NewHandler(hConfig, acts, mapperFactory, expressionFactory, runner)
			if err != nil {
				return nil, err
			}

			initCtx.handlers = append(initCtx.handlers, handler)
		}

		err = trg.Initialize(initCtx)
		if err != nil {
			return nil, err
		}

		triggers[tConfig.Id] = &triggerWrapper{ref: tConfig.Ref, trg: trg, status: &managed.StatusInfo{Name: tConfig.Id}}
		tConfig.AppConfig["Trigger"] = map[string]interface{}{tConfig.Id.(string) : tConfig}
	}
	return triggers, nil
}

type initContext struct {
	handlers []trigger.Handler
	logger   log.Logger
}

func (ctx *initContext) GetHandlers() []trigger.Handler {
	return ctx.handlers
}

func (ctx *initContext) Logger() log.Logger {
	return ctx.logger
}

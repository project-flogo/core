package app

import (
	"fmt"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support/managed"
	"github.com/project-flogo/core/trigger"
)

func (a *App) createSharedActions(actionConfigs []*action.Config) (map[string]action.Action, error) {

	actions := make(map[string]action.Action)

	for _, config := range actionConfigs {

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

	for _, tConfig := range tConfigs {

		_, exists := triggers[tConfig.Id]
		if exists {
			return nil, fmt.Errorf("Trigger with id '%s' already registered, trigger ids have to be unique", tConfig.Id)
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

		initCtx := &initContext{handlers: make([]trigger.Handler, 0, len(tConfig.Handlers))}

		//create handlers for that trigger and init
		for _, hConfig := range tConfig.Handlers {

			var act action.Action
			var err error

			//use action if already associated with Handler
			if hConfig.Action.Act != nil {
				act = hConfig.Action.Act
			} else {

				if hConfig.Action.Id != "" {

					act, exists = a.actions[hConfig.Action.Id]
					if act == nil {
						return nil, fmt.Errorf("shared Action '%s' does not exists", hConfig.Action.Id)
					}

				} else {
					//create the action
					actionFactory := action.GetFactory(hConfig.Action.Ref)
					if actionFactory == nil {
						return nil, fmt.Errorf("Action Factory '%s' not registered", hConfig.Action.Ref)
					}

					act, err = actionFactory.New(hConfig.Action.Config)
					if err != nil {
						return nil, err
					}
				}
			}

			handler, err := trigger.NewHandler(hConfig, act, mapperFactory, runner)
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
	}

	return triggers, nil
}

type initContext struct {
	handlers []trigger.Handler
}

func (ctx *initContext) GetHandlers() []trigger.Handler {
	return ctx.handlers
}

package app

import (
	"fmt"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/managed"
	"github.com/project-flogo/core/trigger"
)

func (a *App) createSharedActions(actionConfigs []*action.Config) (map[string]action.Action, error) {

	actions := make(map[string]action.Action)

	for _, config := range actionConfigs {

		if config.Ref == "" && config.Type != "" {
			log.RootLogger().Warnf("action configuration 'type' deprecated, use 'ref' in the future")
			config.Ref = "#" + config.Type
		}

		if config.Ref == "" {
			return nil, fmt.Errorf("ref not specified for action: %s", config.Id)
		}

		ref := config.Ref

		if config.Ref[0] == '#' {
			var ok bool
			ref, ok = support.GetAliasRef("action", config.Ref)
			if !ok {
				return nil, fmt.Errorf("action '%s' not imported", config.Ref)
			}
		}

		actionFactory := action.GetFactory(ref)
		if actionFactory == nil {
			return nil, fmt.Errorf("action factory '%s' not registered", ref)
		}

		act, err := actionFactory.New(config)
		if err != nil {
			return nil, err
		}

		actions[config.Id] = act

		//if needsCleanup, ok := act.(support.NeedsCleanup); ok {
		//	a.toCleanup = append(a.toCleanup, needsCleanup)
		//}
	}

	return actions, nil
}

func (a *App) createTriggers(tConfigs []*trigger.Config, runner action.Runner) (map[string]*triggerWrapper, error) {

	triggers := make(map[string]*triggerWrapper)

	mapperFactory := mapper.NewFactory(a.resolver)
	expressionFactory := expression.NewFactory(a.resolver)

	for _, tConfig := range tConfigs {

		_, exists := triggers[tConfig.Id]
		if exists {
			return nil, fmt.Errorf("trigger with id '%s' already registered, trigger ids have to be unique", tConfig.Id)
		}

		if tConfig.Ref == "" && tConfig.Type != "" {
			log.RootLogger().Warnf("trigger configuration 'type' deprecated, use 'ref' in the future")
			tConfig.Ref = "#" + tConfig.Type
		}

		ref := tConfig.Ref

		if tConfig.Ref == "" {
			return nil, fmt.Errorf("ref not specified for trigger: %s", tConfig.Id)
		}

		if tConfig.Ref[0] == '#' {
			var ok bool
			ref, ok = support.GetAliasRef("trigger", tConfig.Ref)
			if !ok {
				return nil, fmt.Errorf("trigger '%s' not imported", tConfig.Ref)
			}
		}

		triggerFactory := trigger.GetFactory(ref)

		if triggerFactory == nil {
			return nil, fmt.Errorf("trigger factory '%s' not registered", ref)
		}

		err := tConfig.FixUp(triggerFactory.Metadata(), a.resolver)
		if err != nil {
			return nil, err
		}

		trg, err := triggerFactory.New(tConfig)
		if err != nil {
			return nil, err
		}

		if trg == nil {
			return nil, fmt.Errorf("cannot create trigger nil for id '%s'", tConfig.Id)
		}

		logger := trigger.GetLogger(ref)

		if log.CtxLoggingEnabled() {
			logger = log.ChildLoggerWithFields(logger, log.FieldString("triggerId", tConfig.Id))
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
							return nil, fmt.Errorf("shared action '%s' does not exists", id)
						}
						acts = append(acts, act)
					} else {
						//create the action

						if act.Ref == "" && act.Type != "" {
							log.RootLogger().Warnf("action configuration 'type' deprecated, use 'ref' in the future")
							act.Ref = "#" + act.Type
						}

						if act.Ref == "" {
							return nil, fmt.Errorf("ref not specified for action in trigger '%s", tConfig.Id)
						}

						ref := act.Ref

						if act.Ref[0] == '#' {
							var ok bool
							ref, ok = support.GetAliasRef("action", act.Ref)
							if !ok {
								return nil, fmt.Errorf("action '%s' not imported", act.Ref)
							}
						}

						actionFactory := action.GetFactory(ref)
						if actionFactory == nil {
							return nil, fmt.Errorf("action factory '%s' not registered", ref)
						}

						act, err := actionFactory.New(act.Config)
						if err != nil {
							return nil, err
						}
						//if needsDisposal, ok := act.(support.NeedsCleanup); ok {
						//	a.toDispose = append(a.toDispose, needsDisposal)
						//}

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
		trigger.PostTriggerEvent(trigger.INITIALIZING, tConfig.Id)
		err = trg.Initialize(initCtx)
		if err != nil {
			trigger.PostTriggerEvent(trigger.INIT_FAILED, tConfig.Id)
			return nil, err
		}
		trigger.PostTriggerEvent(trigger.INITIALIZED, tConfig.Id)

		triggers[tConfig.Id] = &triggerWrapper{ref: ref, trg: trg, status: &managed.StatusInfo{Name: tConfig.Id}}
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

package app

import (
	"fmt"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/schema"
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

func (a *App) createTriggers(tConfigs []*trigger.Config, runner action.Runner) ([]*triggerWrapper, error) {

	triggers := make([]*triggerWrapper, len(tConfigs))

	mapperFactory := mapper.NewFactory(a.resolver)
	expressionFactory := expression.NewFactory(a.resolver)

	for i, tConfig := range tConfigs {

		for _, t := range triggers {
			if t != nil && t.id == tConfig.Id {
				return nil, fmt.Errorf("trigger with id '%s' already registered, trigger ids have to be unique", tConfig.Id)
			}
		}

		if tConfig.Ref == "" && tConfig.Type != "" {
			log.RootLogger().Warnf("trigger [%s]'s configuration uses deprecated property 'type', use 'ref' in the future", tConfig.Id)
			tConfig.Ref = "#" + tConfig.Type
		}

		ref := tConfig.Ref

		if tConfig.Ref == "" {
			return nil, fmt.Errorf("trigger [%s]'s ref not specified", tConfig.Id)
		}

		if tConfig.Ref[0] == '#' {
			var ok bool
			ref, ok = support.GetAliasRef("trigger", tConfig.Ref)
			if !ok {
				return nil, fmt.Errorf("unable to start trigger [%s], ref alias '%s' has no corresponding installed trigger", tConfig.Id, tConfig.Ref)
			}
		}

		triggerFactory := trigger.GetFactory(ref)

		if triggerFactory == nil {
			return nil, fmt.Errorf("trigger [%s]'s factory '%s' not registered", tConfig.Id, ref)
		}

		err := tConfig.FixUp(triggerFactory.Metadata(), a.resolver)
		if err != nil {
			return nil, fmt.Errorf("error fixing up trigger [%s]'s metadata:%s", tConfig.Id, err.Error())
		}

		trg, err := triggerFactory.New(tConfig)
		if err != nil {
			return nil, fmt.Errorf("error creating trigger [%s]:%s", tConfig.Id, err.Error())
		}

		if trg == nil {
			return nil, fmt.Errorf("cannot create trigger [%s] with nil", tConfig.Id)
		}

		logger := trigger.GetLogger(ref)

		if log.CtxLoggingEnabled() {
			logger = log.ChildLoggerWithFields(logger, log.FieldString("triggerId", tConfig.Id))
		}

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
							return nil, fmt.Errorf("trigger [%s]'s handler [%s] references nonexistent shared action '%s'", tConfig.Id, hConfig.Name, id)
						}
						acts = append(acts, act)
					} else {
						//create the action

						if act.Ref == "" && act.Type != "" {
							log.RootLogger().Warnf("action configuration 'type' deprecated in trigger [%s]'s handler [%s], use 'ref' in the future", tConfig.Id, hConfig.Name)
							act.Ref = "#" + act.Type
						}

						if act.Ref == "" {
							return nil, fmt.Errorf("ref not specified for action in trigger [%s]'s handler [%s]", tConfig.Id, hConfig.Name)
						}

						ref := act.Ref

						if act.Ref[0] == '#' {
							var ok bool
							ref, ok = support.GetAliasRef("action", act.Ref)
							if !ok {
								return nil, fmt.Errorf("action '%s' referenced in trigger [%s]'s handler [%s], has not been imported", act.Ref, tConfig.Id, hConfig.Name)
							}
						}

						actionFactory := action.GetFactory(ref)
						if actionFactory == nil {
							return nil, fmt.Errorf("action factory '%s' referenced by trigger [%s]'s handler [%s], has not been registered\"", ref, tConfig.Id, hConfig.Name)
						}

						act, err := actionFactory.New(act.Config)
						if err != nil {
							return nil, fmt.Errorf("error creating action [%s] for trigger [%s]'s handler [%s]\"", ref, tConfig.Id, hConfig.Name)
						}
						//if needsDisposal, ok := act.(support.NeedsCleanup); ok {
						//	a.toDispose = append(a.toDispose, needsDisposal)
						//}

						acts = append(acts, act)
					}
				}
			}

			// Resolve schema references
			if hConfig.Schemas != nil {
				if out := hConfig.Schemas.Output; out != nil {
					for name, def := range out {
						ref, ok := def.(string)
						if ok {
							s, err := schema.FindOrCreate(ref)
							if err != nil {
								return nil, fmt.Errorf("unable to find or create output schema [%s] for trigger [%s]'s handler [%s]: %s", ref, tConfig.Id, hConfig.Name, err.Error())
							}
							hConfig.Schemas.Output[name] = s
						}
					}
				}

				if reply := hConfig.Schemas.Reply; reply != nil {
					for name, def := range reply {
						ref, ok := def.(string)
						if ok {
							s, err := schema.FindOrCreate(ref)
							if err != nil {
								return nil, fmt.Errorf("unable to find or create reply schema [%s] in trigger [%s]'s handler [%s]: %s", ref, tConfig.Id, hConfig.Name, err.Error())
							}
							hConfig.Schemas.Reply[name] = s
						}
					}
				}
			}

			handler, err := trigger.NewHandler(hConfig, acts, mapperFactory, expressionFactory, runner, logger)
			if err != nil {
				return nil, fmt.Errorf("error creating handler [%s] in trigger [%s]:%s", hConfig.Name, tConfig.Id, err.Error())
			}

			initCtx.handlers = append(initCtx.handlers, handler)
		}
		trigger.PostTriggerEvent(trigger.INITIALIZING, tConfig.Id)
		err = trg.Initialize(initCtx)
		if err != nil {
			trigger.PostTriggerEvent(trigger.INIT_FAILED, tConfig.Id)
			return nil, fmt.Errorf("error initializing trigger [%s]:%s", tConfig.Id, err.Error())
		}
		trigger.PostTriggerEvent(trigger.INITIALIZED, tConfig.Id)

		triggers[i] = &triggerWrapper{id: tConfig.Id, ref: ref, trg: trg, status: &managed.StatusInfo{Name: tConfig.Id}}
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

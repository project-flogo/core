package trigger

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/property"
	"github.com/project-flogo/core/support/log"
)

var handlerLog = log.ChildLogger(log.RootLogger(), "handler")

type Handler interface {
	Name() string
	Settings() map[string]interface{}
	Schemas() *SchemaConfig
	Handle(ctx context.Context, triggerData interface{}) (map[string]interface{}, error)
}

type actImpl struct {
	act                action.Action
	condition          expression.Expr
	actionInputMapper  mapper.Mapper
	actionOutputMapper mapper.Mapper
}

type handlerImpl struct {
	runner    action.Runner
	config    *HandlerConfig
	acts      []actImpl
	eventData map[string]string
}

func (h *handlerImpl) Name() string {
	return h.config.Name
}

func (h *handlerImpl) Schemas() *SchemaConfig {
	return h.config.Schemas
}

func (h *handlerImpl) Settings() map[string]interface{} {
	return h.config.Settings
}

func (h *handlerImpl) SetDefaultEventData(data map[string]string) {
	h.eventData = data
}

func NewHandler(config *HandlerConfig, acts []action.Action, mf mapper.Factory, ef expression.Factory, runner action.Runner) (Handler, error) {

	if len(acts) == 0 {
		return nil, errors.New("no action specified for handler")
	}

	handler := &handlerImpl{config: config, acts: make([]actImpl, len(acts)), runner: runner}

	var err error

	//todo we could filter inputs/outputs based on the metadata, maybe make this an option
	for i, act := range acts {
		handler.acts[i].act = act

		if config.Actions[i].If != "" {
			condition, err := ef.NewExpr(config.Actions[i].If)
			if err != nil {
				return nil, err
			}
			handler.acts[i].condition = condition
		}

		if len(config.Actions[i].Input) != 0 {
			handler.acts[i].actionInputMapper, err = mf.NewMapper(config.Actions[i].Input)
			if err != nil {
				return nil, err
			}
		}

		if len(config.Actions[i].Output) != 0 {
			handler.acts[i].actionOutputMapper, err = mf.NewMapper(config.Actions[i].Output)
			if err != nil {
				return nil, err
			}
		}
	}

	return handler, nil
}

func (h *handlerImpl) GetSetting(setting string) (interface{}, bool) {

	if h.config == nil {
		return nil, false
	}

	val, exists := h.config.Settings[setting]

	if !exists {
		val, exists = h.config.parent.Settings[setting]
	}

	return val, exists
}

func (h *handlerImpl) Handle(ctx context.Context, triggerData interface{}) (results map[string]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			handlerLog.Warnf("Unhandled Error while handling handler [%s]: %v", h.Name(), r)
			if handlerLog.DebugEnabled() {
				handlerLog.Debugf("StackTrace: %s", debug.Stack())
			}
			err = fmt.Errorf("Unhandled Error while handling handler [%s]: %v", h.Name(), r)
		}
	}()

	eventData := h.eventData

	// check if any event data was attached to the context
	if ctxEventData, _ := ExtractEventDataFromContext(ctx); ctxEventData != nil {
		//use this event data values and add missing default event values
		for key, value := range eventData {
			if _, exists := ctxEventData[key]; !exists {
				ctxEventData[key] = value
			}
		}
		eventData = ctxEventData
	}

	var triggerValues map[string]interface{}
	PostHandlerEvent(STARTED, h.Name(), h.config.parent.Id, eventData)

	if triggerData == nil {
		triggerValues = make(map[string]interface{})
	} else if values, ok := triggerData.(map[string]interface{}); ok {
		triggerValues = values
	} else if value, ok := triggerData.(data.StructValue); ok {
		triggerValues = value.ToMap()
	} else {
		return nil, fmt.Errorf("unsupported trigger data: %v", triggerData)
	}

	var act actImpl
	scope := data.NewSimpleScope(triggerValues, nil)
	for _, v := range h.acts {
		if v.condition == nil {
			act = v
			break
		}
		val, err := v.condition.Eval(scope)
		if err != nil {
			return nil, err
		}
		if val == nil {
			return nil, errors.New("expression has nil result")
		}
		condition, ok := val.(bool)
		if !ok {
			return nil, errors.New("expression has a non-bool result")
		}
		if condition {
			act = v
			break
		}
	}

	if act.act == nil {
		log.RootLogger().Warnf("no action to execute")
		return nil, nil
	}

	var inputMap map[string]interface{}

	if act.actionInputMapper != nil {
		inScope := data.NewSimpleScope(triggerValues, nil)

		inputMap, err = act.actionInputMapper.Apply(inScope)
		if err != nil {
			return nil, err
		}
	} else {
		inputMap = triggerValues
	}

	if ioMd := act.act.IOMetadata(); ioMd != nil {
		for name, tv := range ioMd.Input {
			if val, ok := inputMap[name]; ok {
				inputMap[name], err = coerce.ToType(val, tv.Type())
				if err != nil {
					return nil, err
				}
			}
		}
	}

	newCtx := NewHandlerContext(ctx, h.config)

	if property.IsPropertySnapshotEnabled() {
		if inputMap == nil {
			inputMap = make(map[string]interface{})
		}
		// Take snapshot of current app properties
		propSnapShot := make(map[string]interface{}, len(property.DefaultManager().GetProperties()))
		for k, v := range property.DefaultManager().GetProperties() {
			propSnapShot[k] = v
		}
		inputMap["_PROPERTIES"] = propSnapShot
	}

	results, err = h.runner.RunAction(newCtx, act.act, inputMap)
	if err != nil {
		PostHandlerEvent(FAILED, h.Name(), h.config.parent.Id, eventData)
		return nil, err
	}

	PostHandlerEvent(COMPLETED, h.Name(), h.config.parent.Id, eventData)

	if act.actionOutputMapper != nil {
		outScope := data.NewSimpleScope(results, nil)
		retValue, err := act.actionOutputMapper.Apply(outScope)

		return retValue, err
	} else {
		return results, nil
	}
}

func (h *handlerImpl) String() string {

	triggerId := ""
	if h.config.parent != nil {
		triggerId = h.config.parent.Id
	}
	handlerId := "Handler"
	if h.config.Name != "" {
		handlerId = h.config.Name
	}

	return fmt.Sprintf("Trigger[%s].%s", triggerId, handlerId)
}

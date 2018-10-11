package trigger

import (
	"context"
	"fmt"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/resolve"
)

type Handler interface {
	GetSetting(setting string) (interface{}, bool)
	Settings() map[string]interface{}
	Handle(ctx context.Context, triggerData interface{}) (map[string]interface{}, error)
}

type handlerImpl struct {
	runner action.Runner
	act    action.Action

	outputMd map[string]data.TypedValue
	replyMd  map[string]data.TypedValue

	config *HandlerConfig

	settings map[string]interface{}

	actionInputMapper  mapper.Mapper
	actionOutputMapper mapper.Mapper
}

func (h *handlerImpl) Settings() map[string]interface{} {
	return h.settings
}

func NewHandler(trigger Trigger, config *HandlerConfig, act action.Action, mf mapper.Factory, runner action.Runner) (Handler, error) {

	handler := &handlerImpl{config: config, act: act, outputMd: trigger.Metadata().Output, replyMd: trigger.Metadata().Reply, runner: runner}

	cr := resolve.GetBasicResolver()

	handler.settings = make(map[string]interface{}, len(config.Settings))
	for key, value := range config.Settings {

		if toResolve, ok := value.(string); ok {
			//static resolution
			newValue, err := cr.Resolve(toResolve, nil)
			if err != nil {
				return nil, err
			}
			handler.settings[key] = newValue
		} else {
			handler.settings[key] = value
		}
	}

	var err error

	//todo we could filter inputs/outputs based on the metadata, maybe make this an option
	if len(config.Action.Input) != 0 {
		handler.actionInputMapper, err = mf.NewMapper(config.Action.Input)
		if err != nil {
			return nil, err
		}
	}

	if len(config.Action.Output) != 0 {
		handler.actionOutputMapper, err = mf.NewMapper(config.Action.Output)
		if err != nil {
			return nil, err
		}
	}

	return handler, nil
}

func (h *handlerImpl) GetSetting(setting string) (interface{}, bool) {

	if h.config == nil {
		return nil, false
	}

	val, exists := h.settings[setting]

	if !exists {
		val, exists = h.config.parent.Settings[setting]
	}

	return val, exists
}

func (h *handlerImpl) Handle(ctx context.Context, triggerData interface{}) (map[string]interface{}, error) {

	var err error

	var triggerValues map[string]interface{}

	if values, ok := triggerData.(map[string]interface{}); ok {
		triggerValues = values
	} else if value, ok := triggerData.(data.StructValue); ok {
		triggerValues = value.ToMap()
	} else {
		return nil, fmt.Errorf("Unsupport trigger data: %v", triggerData)
	}

	var inputMap map[string]interface{}

	if h.actionInputMapper != nil {
		inScope := data.NewSimpleScope(triggerValues, nil)

		inputMap, err = h.actionInputMapper.Apply(inScope)
		if err != nil {
			return nil, err
		}
	} else {
		inputMap = triggerValues
	}

	newCtx := NewHandlerContext(ctx, h.config)
	results, err := h.runner.RunAction(newCtx, h.act, inputMap)
	if err != nil {
		return nil, err
	}

	if h.actionOutputMapper != nil {
		outScope := data.NewSimpleScope(results, nil)
		retValue, err := h.actionOutputMapper.Apply(outScope)

		return retValue, err
	} else {
		return results, nil
	}
}

func (h *handlerImpl) String() string {
	return fmt.Sprintf("Handler[action:%s]", h.config.Action.Ref)
}

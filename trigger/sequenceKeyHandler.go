package trigger

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/property"
	"github.com/project-flogo/core/support/log"
)

type SeqKayActionWrapper func()

type seqKeyHandlerImpl struct {
	runner           action.Runner
	logger           log.Logger
	config           *HandlerConfig
	acts             []actImpl
	eventData        map[string]string
	seqKeyChannelMap sync.Map
	seqKeyChannleSize int 
}

func (h *seqKeyHandlerImpl) Name() string {
	return h.config.Name
}

func (h *seqKeyHandlerImpl) Schemas() *SchemaConfig {
	return h.config.Schemas
}

func (h *seqKeyHandlerImpl) Settings() map[string]interface{} {
	return h.config.Settings
}

func (h *seqKeyHandlerImpl) Logger() log.Logger {
	return h.logger
}

func (h *seqKeyHandlerImpl) SetDefaultEventData(data map[string]string) {
	h.eventData = data
}

func (h *seqKeyHandlerImpl) GetSetting(setting string) (interface{}, bool) {

	if h.config == nil {
		return nil, false
	}

	val, exists := h.config.Settings[setting]

	if !exists {
		val, exists = h.config.parent.Settings[setting]
	}

	return val, exists
}

func (h *seqKeyHandlerImpl) Handle(ctx context.Context, triggerData interface{}) (results map[string]interface{}, err error) {
	handlerName := "Handler"
	if h.config != nil && h.config.Name != "" {
		handlerName = h.config.Name
	}
	newCtx := NewHandlerContext(ctx, h.config)

	defer func() {
		h.Logger().Debugf("Handler [%s] for event id [%s] completed in %s", handlerName, GetHandlerEventIdFromContext(newCtx), time.Since(GetHandleStartTimeFromContext(newCtx)).String())
		if r := recover(); r != nil {
			h.Logger().Warnf("Unhandled Error while handling handler [%s]: %v", h.Name(), r)
			if h.Logger().DebugEnabled() {
				h.Logger().Debugf("StackTrace: %s", debug.Stack())
			}
			err = fmt.Errorf("Unhandled Error while handling handler [%s]: %v", h.Name(), r)
		}
	}()

	var triggerValues map[string]interface{}

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

	if act.sequenceKey != nil {
		sequenceKeyObj, err := act.sequenceKey.Apply(scope)
		if err != nil {
			return nil, err
		}
		if sequenceKeyObj == nil {
			h.logger.Warnf("SequenceKey is evaluated to nil. Running action in concurrent mode")
			return h.runAction(newCtx, act, scope, triggerValues, handlerName)
		} else {
			sequenceKeyString, _ := sequenceKeyObj["key"].(string)
			if sequenceKeyString == "" {
				h.logger.Warnf("SequenceKey is evaluated to empty string. Running action in concurrent mode")
				return h.runAction(newCtx, act, scope, triggerValues, handlerName)
			}
			// Run actions in sequencial mode for matching key
			// Check if a channel is already created for the sequence key
			runActionChannel, _ := h.seqKeyChannelMap.Load(sequenceKeyString)
			if runActionChannel == nil {
				// Create a new channel for the sequence key
				runActionChannel = make(chan SeqKayActionWrapper, h.seqKeyChannleSize)
				h.seqKeyChannelMap.Store(sequenceKeyString, runActionChannel)
				// Start a go routine to listen on the channel
				go h.seqKeyActionListener(runActionChannel.(chan SeqKayActionWrapper), sequenceKeyString)
			}

			resultChann := make(chan ExecResult)
			runActionWrapper := func() {
				h.runSeqKeyBasedAction(newCtx, act, scope, triggerValues, handlerName, resultChann)
			}
			// Send the action to the channel
			runActionChannel.(chan SeqKayActionWrapper) <- runActionWrapper

			// Wait for the reply
			result := <-resultChann
			return result.results, result.err
		}
	} else {
		// Run action in concurrent mode
		return h.runAction(newCtx, act, scope, triggerValues, handlerName)
	}
}

func (h *seqKeyHandlerImpl) seqKeyActionListener(seqActionChannel chan SeqKayActionWrapper, seqKey string) {
	for seqKayBasedAction := range seqActionChannel {
		h.logger.Infof("Running action[%s] for sequence key [%s]", h.Name(), seqKey)
		seqKayBasedAction()
		h.logger.Infof("Action[%s] for sequence key [%s] completed", h.Name(), seqKey)
	}
}

type ExecResult struct {
	results map[string]interface{}
	err     error
}

func (h *seqKeyHandlerImpl) runAction(ctx context.Context, act actImpl, scope data.Scope, triggerValues map[string]interface{}, handlerName string) (results map[string]interface{}, err error) {

	newCtx := NewHandlerContext(ctx, h.config)
	h.Logger().Infof("Executing handler [%s] for event Id [%s]", handlerName, GetHandlerEventIdFromContext(newCtx))
	eventData := h.eventData

	// check if any event data was attached to the context
	if ctxEventData, _ := ExtractEventDataFromContext(newCtx); ctxEventData != nil {
		//use this event data values and add missing default event values
		for key, value := range eventData {
			if _, exists := ctxEventData[key]; !exists {
				ctxEventData[key] = value
			}
		}
		eventData = ctxEventData
	}

	PostHandlerEvent(STARTED, h.Name(), h.config.parent.Id, eventData)
	var inputMap map[string]interface{}

	if act.actionInputMapper != nil {
		inputMap, err = act.actionInputMapper.Apply(scope)
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

	results, err = h.runner.RunAction(ctx, act.act, inputMap)
	if err != nil {
		PostHandlerEvent(FAILED, h.Name(), h.config.parent.Id, eventData)
		return nil, err
	}

	PostHandlerEvent(COMPLETED, h.Name(), h.config.parent.Id, eventData)

	if act.actionOutputMapper != nil {
		outScope := data.NewSimpleScope(results, nil)
		results, err = act.actionOutputMapper.Apply(outScope)
	}

	return results, err
}

func (h *seqKeyHandlerImpl) runSeqKeyBasedAction(ctx context.Context, act actImpl, scope data.Scope, triggerValues map[string]interface{}, handlerName string, resultChan chan ExecResult) {
	results, err := h.runAction(ctx, act, scope, triggerValues, handlerName)
	resultChan <- ExecResult{results: results, err: err}
}

func (h *seqKeyHandlerImpl) String() string {
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

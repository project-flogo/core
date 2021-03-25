package trigger

import (
	"context"
	"github.com/project-flogo/core/support"
	"time"
)

var handlerKey = "_ctx_handler_info"

var idGenerator *support.Generator

func init() {
	var err error
	idGenerator, err = support.NewGenerator()
	if err != nil {
		panic("initialization uuid generator error:" + err.Error())
	}
}

type HandlerInfo struct {
	Name      string
	EventId   string
	StartTime time.Time
}

// NewHandlerContext add the handler info to a new child context
func NewHandlerContext(parentCtx context.Context, config *HandlerConfig) context.Context {
	return newContext(parentCtx, config)
}

// NewContextWithEventId new context by adding event id to the context
func NewContextWithEventId(parentCtx context.Context, eventId string) context.Context {
	return context.WithValue(parentCtx, handlerKey, &HandlerInfo{EventId: eventId, StartTime: time.Now()})
}

func newContext(parentCtx context.Context, config *HandlerConfig) context.Context {
	var handlerInfo *HandlerInfo
	//Take default name as handler in case no handler name.
	name := "handler"
	if config != nil && config.Name != "" {
		name = config.Name
	} else if config != nil && config.parent != nil && config.parent.Id != "" {
		// Take trigger name if no handler name.
		name = config.parent.Id
	}

	value := parentCtx.Value(handlerKey)
	if value != nil {
		info, ok := value.(*HandlerInfo)
		if ok {
			handlerInfo = info
			//Update trigger info
			if len(info.EventId) > 0 {
				handlerInfo.Name = name
			} else {
				handlerInfo.EventId = idGenerator.NextAsString()
				handlerInfo.Name = name
			}
			if handlerInfo.StartTime.IsZero() {
				handlerInfo.StartTime = time.Now()
			}
			return context.WithValue(parentCtx, handlerKey, handlerInfo)
		}
	}
	return context.WithValue(parentCtx, handlerKey, &HandlerInfo{Name: name, EventId: idGenerator.NextAsString(), StartTime: time.Now()})
}

// HandlerFromContext returns the handler info stored in the context, if any.
func HandlerFromContext(ctx context.Context) (*HandlerInfo, bool) {
	u, ok := ctx.Value(handlerKey).(*HandlerInfo)
	return u, ok
}

func GetHandlerEventIdFromContext(ctx context.Context) string {
	u, ok := ctx.Value(handlerKey).(*HandlerInfo)
	if ok {
		return u.EventId
	}
	return ""
}

func GetHandleStartTimeFromContext(ctx context.Context) time.Time {
	u, ok := ctx.Value(handlerKey).(*HandlerInfo)
	if ok {
		return u.StartTime
	}
	return time.Time{}
}

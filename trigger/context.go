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
	return NewHandlerContextWithEventId(parentCtx, config, idGenerator.NextAsString())
}

// NewHandlerContextWithEventId add the handler info to a new child context with event id and starting time
func NewHandlerContextWithEventId(parentCtx context.Context, config *HandlerConfig, eventId string) context.Context {
	value := parentCtx.Value(handlerKey)
	if value != nil {
		info, ok := value.(*HandlerInfo)
		if ok {
			if len(info.EventId) > 0 {
				return parentCtx
			} else {
				info.EventId = eventId
				info.StartTime = time.Now()
			}
		}
	}
	return context.WithValue(parentCtx, handlerKey, &HandlerInfo{Name: config.Name, EventId: eventId, StartTime: time.Now()})
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

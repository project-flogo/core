package trigger

import "github.com/project-flogo/core/engine/event"

type Status string

const (
	INITIALIZING     = "Initializing"
	INITIALIZED      = "Initialized"
	INIT_FAILED      = "InitFailed"
	STARTED          = "Started"
	STOPPED          = "Stopped"
	FAILED           = "Failed"
	COMPLETED        = "Completed"
	TriggerEventType = "triggerevent"
)

// Trigger Event
type TriggerEvent interface {
	// Name of trigger
	Name() string
	// Status of trigger. Valid status - INITIALIZING, INITIALIZED, STARTED, STOPPED, FAILED
	Status() Status
}

type triggerEvent struct {
	name   string
	status Status
}

func (te triggerEvent) Name() string {
	return te.name
}

func (te triggerEvent) Status() Status {
	return te.status
}

type HandlerEvent interface {
	// Name of trigger this handler belongs to
	TriggerName() string
	// Name of the handler
	HandlerName() string
	// Status of handler. Valid status - INITIALIZED, STARTED, COMPLETED, FAILED
	Status() Status
	// Handler specific tags set by the underlying implementation e.g. method and path by REST trigger handler or
	// topic name by Kafka trigger handler. This is useful when peek view of trigger(and handlers) is desired.
	Tags() map[string]string
}

type handlerEvent struct {
	triggerName string
	name        string
	status      Status
	tags        map[string]string
}

func (he handlerEvent) TriggerName() string {
	return he.triggerName
}

func (he handlerEvent) HandlerName() string {
	return he.name
}

func (he handlerEvent) Status() Status {
	return he.status
}

func (he handlerEvent) Tags() map[string]string {
	return he.tags
}

func (s Status) String() string {
	return string(s)
}

func PostTriggerEvent(tStatus Status, name string) {
	if event.HasListener(TriggerEventType) {
		te := &triggerEvent{name: name, status: tStatus}
		event.Post(TriggerEventType, te)
	}
}

// Publish handler event
func PostHandlerEvent(hStatus Status, hName, tName string, hTags map[string]string) {
	if event.HasListener(TriggerEventType) {
		te := &handlerEvent{name: hName, triggerName: tName, status: hStatus, tags: hTags}
		event.Post(TriggerEventType, te)
	}
}

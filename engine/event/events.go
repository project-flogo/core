package event

type Listener interface {
	// Called when matching event occurs
	HandleEvent(*Context) error
}

// Context is a wrapper over specific event
type Context struct {
	// Type of event
	eventType string

	// Event data
	event interface{}
}

// Returns wrapped event data
func (ec *Context) GetEvent() interface{} {
	return ec.event
}

// Returns event type
func (ec *Context) GetEventType() string {
	return ec.eventType
}

// Buffered channel
var eventQueue = make(chan *Context, 100)

// Puts event with given type and data on the channel
func Post(eventType string, event interface{}) {
	if publishEventsEnabled && publisherRunning && HasListener(eventType) {
		evtCtx := &Context{event: event, eventType: eventType}
		// Put event on the queue
		eventQueue <- evtCtx
	}
}

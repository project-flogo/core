package event

import (
	"os"
	"strconv"

	"github.com/project-flogo/core/support/log"
)

const (
	ENV_PUBLISH_AUDIT_EVENTS_KEY = "FLOGO_PUBLISH_AUDIT_EVENTS"
)

var publishEventsEnabled = PublishEnabled()
var publisherRoutineStarted = false
var shutdown = make(chan bool)

// Buffered channel
var eventQueue = make(chan *Context, 100)

//TODO channel to be passed to actions
// Puts event with given type and data on the channel
func Post(eventType string, event interface{}) {
	if publishEventsEnabled && publisherRoutineStarted && HasListener(eventType) {
		evtCtx := &Context{event: event, eventType: eventType}
		// Put event on the queue
		eventQueue <- evtCtx
	}
}

func startPublisherRoutine() {
	if publisherRoutineStarted {
		return
	}

	go publishEvents()
	publisherRoutineStarted = true
}

func stopPublisherRoutine() {
	if !publisherRoutineStarted {
		return
	}

	hasListeners := false

	if len(emitters) > 0 {
		for _, emitter := range emitters {
			if len(emitter.listeners) > 0 {
				hasListeners = true
				break
			}
		}
	}

	if !hasListeners {
		// No more listeners. Stop go routine
		shutdown <- true
	}
}

func publishEvents() {
	defer func() {
		publisherRoutineStarted = false
	}()

	for {
		select {
		case evtCtx := <-eventQueue:
			publishEvent(evtCtx)
		case <-shutdown:
			log.RootLogger().Infof("Shutting down event publisher routine")
			return
		}
	}
}

func publishEvent(evtCtx *Context) {

	emittersMutex.RLock()
	emitter, ok := emitters[evtCtx.eventType]
	emittersMutex.RUnlock()

	if ok {
		emitter.Publish(evtCtx)
	}
}

func PublishEnabled() bool {
	key := os.Getenv(ENV_PUBLISH_AUDIT_EVENTS_KEY)
	if len(key) > 0 {
		publish, _ := strconv.ParseBool(key)
		return publish
	}
	return true
}

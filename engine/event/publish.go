package event

import (
	"github.com/project-flogo/core/support/log"
)

const ()

var publishEventsEnabled = PublishEventEnabled()
var publisherRunning = false
var shutdown = make(chan bool)

func startPublisherRoutine() {
	if publisherRunning || !publishEventsEnabled {
		return
	}

	go publishEvents()
	publisherRunning = true
}

func stopPublisherRoutine() {
	if !publisherRunning {
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

	log.RootLogger().Infof("Starting event publisher")

	defer func() {
		publisherRunning = false
	}()

	for {
		select {
		case evtCtx := <-eventQueue:
			publishEvent(evtCtx)
		case <-shutdown:
			log.RootLogger().Infof("Shutting down event publisher")
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

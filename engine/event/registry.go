package event

import (
	"errors"
	"sync"

	"github.com/project-flogo/core/support/log"
)

var emitters = make(map[string]*Emitter)
var emittersMutex = &sync.RWMutex{}

//todo do we need to dynamically add/remove listeners at runtime? If not, we can remove all locking

// Registers listener for given event types
func RegisterListener(name string, listener Listener, eventTypes []string) error {
	if name == "" {
		return errors.New("event listener name must be specified")
	}

	if listener == nil {
		return errors.New("event listener must not nil")
	}

	if len(eventTypes) == 0 {
		return errors.New("at least one event type must be provided")
	}

	emittersMutex.Lock()

	for _, eType := range eventTypes {
		emitter, ok := emitters[eType]
		if !ok {
			emitter = &Emitter{eventType: eType, mutex: &sync.RWMutex{}, listeners: make(map[string]Listener)}
			emitters[eType] = emitter
		}
		err := emitter.RegisterListener(name, listener)
		if err != nil {
			log.RootLogger().Debugf("Event listener - '%s' successfully registered for event type - '%s'", name, eType)
		}
	}

	emittersMutex.Unlock()

	startPublisherRoutine()
	return nil
}

// Unregister event listener for given event types .
// To unregister from all event types, set eventTypes to nil
func UnRegisterListener(name string, eventTypes []string) {

	if name == "" || len(emitters) == 0 {
		return
	}

	//unregister doesn't remove emitter, so no lock is required

	if len(eventTypes) > 0 {
		for _, eventType := range eventTypes {
			if emitter, ok := emitters[eventType]; ok {
				err := emitter.UnRegisterListener(name)
				if err != nil {
					log.RootLogger().Debugf("Event listener - '%s' successfully unregistered for event type - '%s'", name, eventType)
				}
			}
		}
	} else {
		//unregister from all emitters
		for _, emitter := range emitters {
			err := emitter.UnRegisterListener(name)
			if err != nil {
				log.RootLogger().Debugf("Event listener - '%s' successfully unregistered for event type - '%s'", name, emitter.eventType)
			}
		}
	}

	stopPublisherRoutine()
}

func HasListener(eventType string) bool {

	emittersMutex.RLock()
	emitter, ok := emitters[eventType]
	emittersMutex.RUnlock()

	if ok {
		return emitter.HasListeners()
	}

	return false
}

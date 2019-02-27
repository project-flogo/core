package event

import (
	"errors"
	"runtime/debug"
	"sync"

	"github.com/project-flogo/core/support/log"
)

type Emitter struct {
	mutex     *sync.RWMutex
	eventType string
	listeners map[string]Listener
}

func (te *Emitter) RegisterListener(name string, listener Listener) error {
	if name == "" {
		return errors.New("event listener name must be specified")
	}

	if listener == nil {
		return errors.New("event listener must not nil")
	}

	te.mutex.Lock()
	te.listeners[name] = listener
	te.mutex.Unlock()

	return nil
}

func (te *Emitter) HasListeners() bool {
	te.mutex.RLock()
	hasListeners := len(te.listeners) > 0
	te.mutex.RUnlock()

	return hasListeners
}

func (te *Emitter) UnRegisterListener(name string) error {
	if name == "" {
		return errors.New("event listener name must be specified")
	}

	te.mutex.Lock()
	delete(te.listeners, name)
	te.mutex.Unlock()

	return nil
}

func (te *Emitter) Publish(evtCtx *Context) {

	listenerName := ""

	te.mutex.RLock()

	//todo consider handling panic one level up, this will improve performance

	defer func() {
		te.mutex.RUnlock()
		if r := recover(); r != nil {
			log.RootLogger().Errorf("Registered event listener - '%s' failed to process event due to error - '%v' ", listenerName, r)
			log.RootLogger().Errorf("StackTrace: %s", debug.Stack())
		}
	}()

	for name, listener := range te.listeners {
		listenerName = name

		err := listener.HandleEvent(evtCtx)
		if err != nil {
			log.RootLogger().Errorf("Registered event listener - '%s' failed to process event due to error - '%s' ", name, err.Error())
		} else {
			log.RootLogger().Debugf("Event - '%s' is successfully delivered to event listener - '%s'", te.eventType, name)
		}
	}
}

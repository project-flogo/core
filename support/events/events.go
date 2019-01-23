package events

import (
	"errors"
	"runtime/debug"
	"strings"
	"sync"

	config "github.com/project-flogo/core/engine"
	"github.com/project-flogo/core/support/log"
)

type EventListener interface {
	// Returns name of the listener
	Name() string

	// Called when matching event occurs
	HandleEvent(*EventContext) error
}

var eventListeners = make(map[string][]EventListener)

// Buffered channel
var eventQueue = make(chan *EventContext, 100)
var publisherRoutineStarted = false
var shutdown = make(chan bool)
var publishEventsEnabled = config.PublishAuditEvents()

var lock = &sync.RWMutex{}

// Registers listener for given event types
func RegisterEventListener(evtListener EventListener, eventTypes []string) error {
	logger := log.RootLogger()
	if evtListener == nil {
		return errors.New("Event listener must not nil")
	}

	if len(eventTypes) == 0 {
		return errors.New("Failed register event listener. At-least one event type must be provided.")
	}

	lock.Lock()
	for _, eType := range eventTypes {
		eventListeners[eType] = append(eventListeners[eType], evtListener)
		logger.Debugf("Event listener - '%s' successfully registered for event type - '%s'", evtListener.Name(), eType)
	}
	lock.Unlock()
	startPublisherRoutine()
	return nil
}

// Unregister event listener for given event types .
// To unregister from all event types, set eventTypes to nil
func UnRegisterEventListener(name string, eventTypes []string) {
	logger := log.RootLogger()
	if name == "" {
		return
	}

	lock.Lock()

	var deleteList []string
	var index = -1

	if eventTypes != nil && len(eventTypes) > 0 {
		for _, eType := range eventTypes {
			evtLs, ok := eventListeners[eType]
			if ok {
				for i, el := range evtLs {
					if strings.EqualFold(el.Name(), name) {
						index = i
						break
					}
				}
				if index > -1 {
					if len(evtLs) > 1 {
						// More than one listeners
						copy(evtLs[index:], evtLs[index+1:])
						evtLs[len(evtLs)-1] = nil
						eventListeners[eType] = evtLs[:len(evtLs)-1]
					} else {
						// Single listener in the map. Remove map entry
						deleteList = append(deleteList, eType)
					}
					logger.Debugf("Event listener - '%s' successfully unregistered for event type - '%s'", name, eType)
					index = -1
				}
			}
		}
	} else {
		for eType, elList := range eventListeners {
			for i, el := range elList {
				if strings.EqualFold(el.Name(), name) {
					index = i
					break
				}
			}
			if index > -1 {
				if len(elList) > 1 {
					// More than one listeners
					copy(elList[index:], elList[index+1:])
					elList[len(elList)-1] = nil
					eventListeners[eType] = elList[:len(elList)-1]
				} else {
					// Single listener in the map. Remove map entry
					deleteList = append(deleteList, eType)
				}
				logger.Debugf("Event listener - '%s' successfully unregistered for event type - '%s'", name, eType)
				index = -1
			}
		}
	}

	if len(deleteList) > 0 {
		for _, evtType := range deleteList {
			delete(eventListeners, evtType)
		}
	}

	lock.Unlock()
	stopPublisherRoutine()
}

func startPublisherRoutine() {
	if publisherRoutineStarted == true {
		return
	}

	if len(eventListeners) > 0 {
		// start publisher routine
		go publishEvents()
		publisherRoutineStarted = true
	}
}

func stopPublisherRoutine() {
	if publisherRoutineStarted == false {
		return
	}

	if len(eventListeners) == 0 {
		// No more listeners. Stop go routine
		shutdown <- true
		publisherRoutineStarted = false
	}
}

//  EventContext is a wrapper over specific event
type EventContext struct {
	// Type of event
	eventType string
	// Event data
	event interface{}
}

// Returns wrapped event data
func (ec *EventContext) GetEvent() interface{} {
	return ec.event
}

// Returns event type
func (ec *EventContext) GetType() string {
	return ec.eventType
}

func publishEvents() {
	logger := log.RootLogger()
	defer func() {
		publisherRoutineStarted = false
	}()

	for {
		select {
		case event := <-eventQueue:
			lock.RLock()
			publishEvent(event)
			lock.RUnlock()
		case <-shutdown:
			logger.Infof("Shutting down event publisher routine")
			return
		}
	}
}

func publishEvent(fe *EventContext) {
	logger := log.RootLogger()
	regListeners, ok := eventListeners[fe.eventType]
	if ok {
		for _, ls := range regListeners {
			func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Errorf("Registered event listener - '%s' failed to process event due to error - '%v' ", ls.Name(), r)
						logger.Errorf("StackTrace: %s", debug.Stack())
					}
				}()
				err := ls.HandleEvent(fe)
				if err != nil {
					logger.Errorf("Registered event listener - '%s' failed to process event due to error - '%s' ", ls.Name(), err.Error())
				} else {
					logger.Debugf("Event - '%s' is successfully delivered to event listener - '%s'", fe.eventType, ls.Name())
				}

			}()
		}
		fe = nil
	}
}

func HasListener(eventType string) bool {
	// event publishing is turned off
	if !publishEventsEnabled {
		return false
	}

	lock.RLock()
	ls, ok := eventListeners[eventType]
	lock.RUnlock()
	return ok && len(ls) > 0
}

//TODO channel to be passed to actions
// Puts event with given type and data on the channel
func PostEvent(eType string, event interface{}) {
	if publishEventsEnabled && publisherRoutineStarted && HasListener(eType) {
		evtContext := &EventContext{event: event, eventType: eType}
		// Put event on the queue
		eventQueue <- evtContext
	}
}

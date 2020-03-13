package managed

import (
	"fmt"
	"github.com/project-flogo/core/support/log"
)

// Managed is an interface that is implemented by an object that needs to be
// managed via start/stop
type Managed interface {

	// Start starts the managed object
	Start() error

	// Stop stops the manged object
	Stop() error
}

// start starts a "Managed" object
func start(managed Managed) (err error) {

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(error); !ok {
				err = fmt.Errorf("%v", r)
			} else {
				err = r.(error)
			}
		}
	}()

	return managed.Start()
}

// stop stops a "Managed" object
func stop(managed Managed) (err error) {

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(error); !ok {
				err = fmt.Errorf("%v", r)
			} else {
				err = r.(error)
			}
		}
	}()

	return managed.Stop()
}

// Start starts a Managed object, handles panics and logs details
func Start(name string, managed Managed) error {

	log.RootLogger().Debugf("%s: Starting...", name)
	err := start(managed)

	if err != nil {
		log.RootLogger().Errorf("%s: Error Starting", name)
		return err
	}

	log.RootLogger().Infof("%s: Started", name)
	return nil
}

// Stop stops a Managed object, handles panics and logs details
func Stop(name string, managed Managed) error {

	log.RootLogger().Debugf("%s: Stopping...", name)

	err := stop(managed)

	if err != nil {
		log.RootLogger().Errorf("Error stopping '%s': %s", name, err.Error())
		return err
	}

	log.RootLogger().Infof("%s: Stopped", name)
	return nil
}

// Initializable is an interface that is implemented by an object that needs to be
// initialized
type Initializable interface {

	// Initializes the object
	Init() error
}

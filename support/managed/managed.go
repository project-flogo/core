package managed

import (
	"fmt"

	"github.com/project-flogo/core/support/logger"
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
func start(managed Managed) error {

	defer func() error {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}

			return err
		}

		return nil
	}()

	return managed.Start()
}

// stop stops a "Managed" object
func stop(managed Managed) error {

	defer func() error {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}

			return err
		}

		return nil
	}()

	return managed.Stop()
}

// Start starts a Managed object, handles panics and logs details
func Start(name string, managed Managed) error {

	logger.Debugf("%s: Starting...", name)
	err := start(managed)

	if err != nil {
		logger.Errorf("%s: Error Starting", name)
		return err
	}

	logger.Debugf("%s: Started", name)
	return nil
}

// Stop stops a Managed object, handles panics and logs details
func Stop(name string, managed Managed) error {

	logger.Debugf("%s: Stopping...", name)

	err := stop(managed)

	if err != nil {
		logger.Errorf("Error stopping '%s': %s", name, err.Error())
		return err
	}

	logger.Debugf("%s: Stopped", name)
	return nil
}

// Initializable is an interface that is implemented by an object that needs to be
// initialized
type Initializable interface {

	// Initializes the object
	Init() error
}

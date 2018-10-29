package trigger

import (
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/managed"
)

// Trigger is object that triggers/starts flow instances and
// is managed by an engine
type Trigger interface {
	managed.Managed

	// Initialize is called to initialize the Trigger
	Initialize(ctx InitContext) error
}

// InitContext is the initialization context for the trigger instance
type InitContext interface {

	// Logger the logger for the trigger
	Logger() log.Logger

	// GetHandlers gets the handlers associated with the trigger
	GetHandlers() []Handler
}

// Factory is used to create new instances of a trigger
type Factory interface {

	// Metadata returns the metadata of the trigger
	Metadata() *Metadata

	// New create a new Trigger
	New(config *Config) (Trigger, error)
}

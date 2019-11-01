package action

import (
	"context"

	"github.com/project-flogo/core/app/resource"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/service"
)

// Action is an action to perform as a result of a trigger
type Action interface {

	// Metadata get the Action's metadata
	Metadata() *Metadata

	// IOMetadata get the Action's IO metadata
	IOMetadata() *metadata.IOMetadata
}

// Factory s used to create new instances for an action
type Factory interface {

	// Initialize is called to initialize the action factory
	Initialize(ctx InitContext) error

	// New create a new Action
	New(config *Config) (Action, error)
}

// InitContext is the initialization context for the action factory
type InitContext interface {

	// ResourceManager gets the resource manager for the app
	ResourceManager() *resource.Manager

	// ResourceManager gets the service manager for the engine
	ServiceManager() *service.Manager

	// RuntimeSettings are any runtime setting provided to the engine
	RuntimeSettings() map[string]interface{}
}

// SyncAction is a synchronous action to perform as a result of a trigger
type SyncAction interface {
	Action

	// Run this Action
	Run(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)
}

// AsyncAction is an asynchronous action to perform as a result of a trigger, the action can asynchronously
// return results as it runs.  It returns immediately, but will continue to run.
type AsyncAction interface {
	Action

	// Run this Action
	Run(ctx context.Context, input map[string]interface{}, handler ResultHandler) error
}

// Runner runs actions
type Runner interface {

	// RunAction the specified Action
	RunAction(ctx context.Context, act Action, input map[string]interface{}) (results map[string]interface{}, err error)
}

// ResultHandler used to handle results from the Action
type ResultHandler interface {

	// HandleResult is invoked when there are results available
	HandleResult(results map[string]interface{}, err error)

	// Done indicates that the action has completed
	Done()
}

package activity

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/data/resolve"
)

// Context describes the execution context for an Activity.
// It provides access to attributes, task and Flow information.
type Context interface {
	// ActivityHost gets the "host" under with the activity is executing
	ActivityHost() Host

	//Name the name of the activity that is currently executing
	Name() string

	// GetInput gets the value of the specified input attribute
	GetInput(name string) interface{}

	// SetOutput sets the value of the specified output attribute
	SetOutput(name string, value interface{})

	// GetSharedTempData get shared temporary data for activity, lifespan
	// of the data dependent on the activity host implementation
	GetSharedTempData() map[string]interface{}

	// GetInputObject gets all the activity input as the specified object.
	GetInputObject(input data.FromMap) error

	// SetOutputObject sets the activity output as the specified object.
	SetOutputObject(output data.ToMap) error
}

type Host interface {
	// ID returns the ID of the Activity Host
	ID() string

	// Name the name of the Activity Host
	Name() string

	// IOMetadata get the input/output metadata of the activity host
	IOMetadata() *metadata.IOMetadata

	// Reply is used to reply to the activity Host with the results of the execution
	Reply(replyData map[string]interface{}, err error)

	// Return is used to indicate to the activity Host that it should complete and return the results of the execution
	Return(returnData map[string]interface{}, err error)

	//todo rename, essentially the flow's attrs for now
	WorkingData() data.Scope

	// GetResolver gets the resolver associated with the activity host
	GetResolver() resolve.CompositeResolver

	// GetDetails gets a StringsMap with host specific details/properties, ie. "type":"flow", "id":"2134", etc.
	GetDetails() data.StringsMap
}
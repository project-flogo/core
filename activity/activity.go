package activity

import (
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
)

// LogActivity is an interface for defining a custom LogActivity Execution
type Activity interface {

	// Metadata returns the metadata of the activity
	Metadata() *Metadata

	// Eval is called when an LogActivity is being evaluated.  Returning true indicates
	// that the task is done.
	Eval(ctx Context) (done bool, err error)
}

type Factory func(ctx InitContext) (Activity, error)

type InitContext interface {

	// Settings
	Settings() map[string]interface{}

	// MapperFactory gets the mapper factory associated with the activity host
	MapperFactory() mapper.Factory

	// Logger logger to using during initialization, activity implementations should not
	// keep a reference to this
	Logger() log.Logger
}

// HasDynamicMd is an optional interface that can be implemented by an activity.  If implemented,
// DynamicMd() will be invoked to determine the inputs/outputs of the activity instead of
// relying on the static information from the LogActivity's Metadata
type HasDynamicMd interface {

	// DynamicMd get the input/output metadata
	DynamicMd(ctx Context) (*metadata.IOMetadata, error)
}

type Details struct {
	IsReturn  bool
	IsReply   bool
	DynamicIO bool
}

type HasDetails interface {
	Details() *Details
}

package activity

import (
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/support/log"
)

// Activity is an interface for defining a custom Activity Execution
type Activity interface {

	// Metadata returns the metadata of the activity
	Metadata() *Metadata

	// Eval is called when an Activity is being evaluated.  Returning true indicates
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

type Details struct {
	IsReturn bool
	IsReply  bool
}

type HasDetails interface {
	Details() *Details
}

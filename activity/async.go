package activity

// AsyncActivity is an interface for defining a custom Task that supports asynchronous callback
type AsyncActivity interface {
	Activity

	// PostEval is called when a activity that didn't complete during the Eval
	// needs to be notified.  Returning true indicates that the activity is done.
	PostEval(context Context, userData interface{}) (done bool, err error)
}

package activity

// Error is an activity error
type Error struct {
	activityName string
	errorStr     string
	errorCode    string
	errorData    interface{}
}

func NewError(errorText string, code string, errorData interface{}) *Error {
	return &Error{errorStr: errorText, errorData: errorData, errorCode: code}
}

// Error implements error.Error()
func (e *Error) Error() string {
	return e.errorStr
}

// ActivityName the activity name
func (e *Error) ActivityName() string {
	return e.activityName
}

// Set the activity name
func (e *Error) SetActivityName(name string) {
	e.activityName = name
}

// Data returns any associated error data
func (e *Error) Data() interface{} {
	return e.errorData
}

// Code returns any associated error code
func (e *Error) Code() string {
	return e.errorCode
}

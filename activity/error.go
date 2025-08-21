package activity

type ErrorCategory string

const (
	// ErrorCategory ConfigError is the error category for configuration errors
	ConfigError ErrorCategory = "CONFIG-ERROR"
	// ErrorCategory ActivityError is the error category for activity errors
	ActivityError ErrorCategory = "ACTIVITY-ERROR"
	// ErrorCategory ServerError is the error category for server errors
	ServerError ErrorCategory = "SERVER-ERROR"
	// ErrorCategory ClientError is the error category for client errors
	ClientError ErrorCategory = "CLIENT-ERROR"
	// ErrorCategory TimeoutError is the error category for timeout errors
	TimeoutError ErrorCategory = "TIMEOUT-ERROR"
	// ErrorCategory RetryError is the error category for retry errors
	RetryError ErrorCategory = "RETRY-ERROR"
	// ErrorCategory ConnectionError is the error category for connection errors
	ConnectionError ErrorCategory = "CONNECTION-ERROR"
)

// Error is an activity error
type Error struct {
	activityName  string
	errorStr      string
	errorCode     string
	errorCategory ErrorCategory
	errorData     interface{}
	retriable     bool
}

// ErrorData is the data structure for error details reported by an activity
type ErrorData struct {
	Details string `json:"details,omitempty"`
	Code    string `json:"code,omitempty"`
}

// NewActivityError creates a new activity error with the specified message, category, and details
func NewActivityError(errorMsg string, errorCategory ErrorCategory, errorDetails ErrorData) *Error {
	return &Error{errorStr: errorMsg, errorData: errorDetails, errorCode: errorDetails.Code, errorCategory: errorCategory, retriable: false}
}

// NewRetriableActivityError creates a new retriable activity error with the specified message, category, and details
func NewRetriableActivityError(errorMsg string, errorCategory ErrorCategory, errorDetails ErrorData) *Error {
	return &Error{errorStr: errorMsg, errorData: errorDetails, errorCode: errorDetails.Code, errorCategory: errorCategory, retriable: true}
}

func NewError(errorText string, code string, errorData interface{}) *Error {
	return &Error{errorStr: errorText, errorData: errorData, errorCode: code, retriable: false, errorCategory: ActivityError}
}

func NewRetriableError(errorText string, code string, errorData interface{}) *Error {
	return &Error{errorStr: errorText, errorData: errorData, errorCode: code, retriable: true, errorCategory: ActivityError}
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
	if e.errorData == nil {
		return ErrorData{Code: e.errorCode, Details: e.errorStr}
	}

	if err, ok := e.errorData.(error); ok {
		return ErrorData{Details: err.Error(), Code: e.errorCode}
	}
	return e.errorData
}

// Code returns any associated error code
func (e *Error) Code() string {
	return e.errorCode
}

// Retriable returns wether error is retriable
func (e *Error) Retriable() bool {
	return e.retriable
}

// Category returns any associated error category
func (e *Error) Category() string {
	return string(e.errorCategory)
}

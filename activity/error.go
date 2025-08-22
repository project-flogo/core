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

// ErrorData is the structure for additional error data
type ErrorData struct {
	// Details provides additional information about the error
	// This field can be used to provide context or specific details about the error
	Details string `json:"details,omitempty"`
}

// NewActivityError creates a new activity error with the specified message, category, and details
// This error is not retriable
// errorMsg: the error message
// errorCode: the error code
// errorCategory: the error category (e.g., ConfigError, ActivityError)
// errorData: any additional data associated with the error
// Returns: a pointer to the created Error instance
// Example usage:
//
//	err := NewActivityError("Failed to execute activity", "ACTIVITY-001", ConfigError, map[string]interface{}{"details": "Invalid input"})
func NewActivityError(errorMsg string, errorCode string, errorCategory ErrorCategory, errorData interface{}) *Error {
	return &Error{errorStr: errorMsg, errorData: errorData, errorCode: errorCode, errorCategory: errorCategory, retriable: false}
}

// NewRetriableActivityError creates a new retriable activity error with the specified message, category, and details
// errorMsg: the error message
// errorCode: the error code
// errorCategory: the error category (e.g., ConfigError, ActivityError)
// errorData: any additional data associated with the error
// Returns: a pointer to the created Error instance
// Example usage:
//
//	err := NewRetriableActivityError("Temporary failure, please retry", "ACTIVITY-002", ConnectionError, map[string]interface{}{"details": "Network issue"})
//
// This error indicates that the activity can be retried
// and is suitable for scenarios where transient issues may occur, such as network failures or temporary unavailability of resources.
func NewRetriableActivityError(errorMsg string, errorCode string, errorCategory ErrorCategory, errorData interface{}) *Error {
	return &Error{errorStr: errorMsg, errorData: errorData, errorCode: errorCode, errorCategory: errorCategory, retriable: true}
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
	switch v := e.errorData.(type) {
	case nil:
		return ErrorData{}
	case error:
		return ErrorData{Details: v.Error()}
	default:
		return v
	}
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

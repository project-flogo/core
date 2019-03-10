package schema

type Schema interface {
	Type() string

	Value() string

	Validate(data interface{}) error
}

type Factory interface {
	New(def *Def) (Schema, error)
}

func NewValidationError(msg string, errors []error) *ValidationError {
	return &ValidationError{msg: msg, errors: errors}
}

type ValidationError struct {
	msg    string // description of error
	errors []error
}

func (e *ValidationError) Error() string {
	return e.msg
}

func (e *ValidationError) Errors() []error {
	return e.errors
}

type Def struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

var enabled bool
var validationEnabled bool

func Enable() {
	enabled = true
	validationEnabled = true
}

func Enabled() bool {
	return enabled
}

func DisableValidation() {
	validationEnabled = false
}

func ValidationEnabled() bool {
	return validationEnabled
}

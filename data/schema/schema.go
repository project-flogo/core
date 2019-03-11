package schema

import (
	"fmt"
	"strings"
)

type Schema interface {
	Type() string

	Value() string

	Validate(data interface{}) error
}

type HasSchema interface {
	Schema() Schema
}

type Factory interface {
	New(def *Def) (Schema, error)
}

func NewValidationError(msg string, errors []error)  *ValidationError {
	return &ValidationError{msg:msg, errors:errors}
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

func New(schemaDef *Def) (Schema, error) {

	if !enabled {
		return emptySchema, nil
	}

	if !validationEnabled {
		// validation disabled, so return non-validating schema
		return &schemaSansValidation{def:schemaDef}, nil
	}

	factory := GetFactory(schemaDef.Type)

	if factory == nil {
		return nil, fmt.Errorf("support for schema type '%s' not installed", schemaDef.Type)
	}

	s, err := factory.New(schemaDef)

	if err != nil {
		return nil, err
	}

	return s, nil
}

func FindOrCreate(schemaRep interface{}) (Schema, error) {

	switch t := schemaRep.(type) {
	case HasSchema:
		return t.Schema(), nil
	case Def:
		return New(&t)
	case *Def:
		return New(t)
	case string:
		if strings.HasPrefix(t, "schema://") {
			id := t[9:]
			s := Get(t[9:])
			if s == nil {
				sh := &schemaHolder{id: id}
				toResolve = append(toResolve, sh)
				s = sh
			}

			return s, nil
		}
		return nil, fmt.Errorf("invalid schema reference: %s", t)
	case map[string]string:

		def := &Def{}
		if sType, ok := t["type"]; ok {
			def.Type = sType
		} else {
			return nil, fmt.Errorf("invalid schema definition, type not specified: %s", t)
		}

		if sValue, ok := t["value"]; ok {
			def.Value = sValue
		} else {
			return nil, fmt.Errorf("invalid schema definition, value not specified: %s", t)
		}

		return New(def)
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid schema definition, %v", t)
	}
}

var emptySchema = &emptySchemaImpl{}

type emptySchemaImpl struct {
}

func (*emptySchemaImpl) Type() string {
	return ""
}

func (*emptySchemaImpl) Value() string {
	return ""
}

func (*emptySchemaImpl) Validate(data interface{}) error {
	return nil
}

// schemaSansValidation holds the schema information and ignores validation
type schemaSansValidation struct {
	def *Def
}

func (s *schemaSansValidation) Type() string {
	return s.def.Type
}

func (s *schemaSansValidation) Value() string {
	return s.def.Value
}

func (*schemaSansValidation) Validate(data interface{}) error {
	return nil
}

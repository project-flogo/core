package json

import (
	"errors"

	"github.com/project-flogo/core/data/schema"
	"github.com/xeipuuv/gojsonschema"
)

func init() {
	_ = schema.RegisterFactory("json", &factory{})
}

type factory struct {
}

func (factory) New(def *schema.Def) (schema.Schema, error) {

	schemaLoader := gojsonschema.NewStringLoader(string(def.Value))
	s, err := gojsonschema.NewSchema(schemaLoader)

	if err != nil {
		return nil, err
	}

	return &jsonSchema{def: def, js: s}, nil
}

type jsonSchema struct {
	def *schema.Def
	js  *gojsonschema.Schema
}

func (s *jsonSchema) Type() string {
	return s.def.Type
}

func (s *jsonSchema) Value() string {
	return s.def.Value
}

func (s *jsonSchema) Validate(data interface{}) error {

	loader := gojsonschema.NewGoLoader(data)
	r, err := s.js.Validate(loader)

	if err != nil {
		return err
	}

	if !r.Valid() {

		var errs []error
		for _, err := range r.Errors() {
			errs = append(errs, errors.New(err.String()))
		}

		return schema.NewValidationError("validation failed", errs)
	}

	return nil
}

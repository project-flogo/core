package schema

import (
	"fmt"
	"strings"
)

type HasSchema interface {
	Schema() Schema
}

type HasSchemaIO interface {
	GetInputSchema(name string) Schema

	GetOutputSchema(name string) Schema
}

//DEPRECATED
type ValidationBypass interface {
	BypassValidation() bool
}

func New(schemaDef *Def) (Schema, error) {

	if !enabled {
		return emptySchema, nil
	}

	if !validationEnabled {
		// validation disabled, so return non-validating schema
		return &schemaSansValidation{def: schemaDef}, nil
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
			return nil, fmt.Errorf("invalid schema definition, type not specified: %+v", t)
		}

		if sValue, ok := t["value"]; ok {
			def.Value = sValue
		} else {
			return nil, fmt.Errorf("invalid schema definition, value not specified: %+v", t)
		}
		return New(def)
	case map[string]interface{}:
		def := &Def{}
		if sType, ok := t["type"]; ok {
			def.Type, ok = sType.(string)
			if !ok {
				return nil, fmt.Errorf("invalid schema definition, type is not a string specified: %+v", sType)
			}
		} else {
			return nil, fmt.Errorf("invalid schema definition, type not specified: %+v", t)
		}

		if sValue, ok := t["value"]; ok {
			def.Value, ok = sValue.(string)
			if !ok {
				return nil, fmt.Errorf("invalid schema definition, value is not a string specified: %+v", sValue)
			}
		} else {
			return nil, fmt.Errorf("invalid schema definition, value not specified: %+v", t)
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

package schema

import (
	"fmt"

	"github.com/project-flogo/core/support/log"
)

var (
	schemas   = make(map[string]Schema)
	factories = make(map[string]Factory)
)

func Register(id string, def *Def) (Schema, error) {

	log.RootLogger().Debugf("Registering schema: %s", id)

	if id == "" {
		return nil, fmt.Errorf("id is required to register schema")
	}

	if _, dup := schemas[id]; dup {
		return nil, fmt.Errorf("schema with id '%s' already registered", id)
	}

	if def == nil {
		return nil, fmt.Errorf("cannot register 'nil' schema")
	}

	s, err := New(def)

	if err != nil {
		return nil, err
	}

	schemas[id] = s

	log.RootLogger().Debugf("Registered schema: %s", id)

	return s, nil
}

// Get gets specified schema by id
func Get(id string) Schema {
	return schemas[id]
}

func RegisterFactory(schemaType string, factory Factory) error {

	if schemaType == "" {
		return fmt.Errorf("schemaType is required to register schema factory")
	}

	if _, dup := factories[schemaType]; dup {
		return fmt.Errorf("schema foctory for type '%s' already registered", schemaType)
	}

	if factory == nil {
		return fmt.Errorf("cannot register 'nil' schema factory")
	}

	log.RootLogger().Debugf("Registering schema factory: %s", schemaType)

	factories[schemaType] = factory

	return nil
}

// Get gets specified schema by id
func GetFactory(schemaType string) Factory {
	return factories[schemaType]
}

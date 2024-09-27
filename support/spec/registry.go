package spec

import (
	"fmt"

	"github.com/project-flogo/core/support/log"
)

var (
	specs = make(map[string]Spec)
)

func Register(id string, def *Def) (Spec, error) {

	log.RootLogger().Debugf("Registering spec: %s", id)

	if id == "" {
		return nil, fmt.Errorf("id is required to register spec")
	}

	if _, dup := specs[id]; dup {
		return nil, fmt.Errorf("spec with id '%s' already registered", id)
	}

	if def == nil {
		return nil, fmt.Errorf("cannot register 'nil' spec")
	}

	s, err := New(def)

	if err != nil {
		return nil, err
	}

	specs[id] = s

	log.RootLogger().Debugf("Registered spec: %s", id)

	return s, nil
}

// Get gets specified schema by id
func Get(id string) Spec {
	return specs[id]
}

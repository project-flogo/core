package function

import (
	"fmt"
	"github.com/project-flogo/core/support/log"
)

var (
	functions = make(map[string]Function)
)

func Register(function Function) error {

	if function == nil {
		return fmt.Errorf("cannot register 'nil' function")
	}

	if _, dup := functions[function.Name()]; dup {
		return fmt.Errorf("function '%s' already registered", function.Name())
	}

	log.RootLogger().Debugf("Registering function: %s", function.Name())

	functions[function.Name()] = function

	return nil
}

// Get gets specified function by id
func Get(id string) Function {
	return functions[id]
}

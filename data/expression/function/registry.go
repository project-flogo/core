package function

import (
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/project-flogo/core/support/log"
)

var (
	functions    = make(map[string]Function)
	functionsTmp = make(map[string]Function)
	packages     = make(map[string]string)
)

func Register(function Function) error {

	if function == nil {
		return fmt.Errorf("cannot register 'nil' function")
	}

	goPkg, err := getGoPackage(function)
	if err != nil {
		return err
	}

	alias := path.Base(goPkg)

	if _, dup := functions[alias+"."+function.Name()]; dup {
		return fmt.Errorf("function '%s' already registered", function.Name())
	}

	if _, exists := packages[goPkg]; !exists {
		log.RootLogger().Debugf("Registering function package: %s", goPkg)
	}

	packages[goPkg] = alias

	log.RootLogger().Debugf("Registering function: %s:%s", goPkg, function.Name())

	functionsTmp[goPkg+":"+function.Name()] = function

	return nil
}

func getGoPackage(function Function) (string, error) {
	value := reflect.ValueOf(function)
	// unwrap pointer
	value = value.Elem()
	if value.Kind() != reflect.Struct {
		return "", fmt.Errorf("unable to determine fo package of '%v'", function)
	}

	return value.Type().PkgPath(), nil
}

// Get gets specified function by id
func Get(id string) Function {
	return functions[id]
}

func IsFunctionPackage(pkg string) bool {
	_, ok := packages[pkg]
	return ok
}

func SetPackageAlias(pkg string, alias string) {
	packages[pkg] = alias
}

func ResolveAliases() {

	if functionsTmp == nil {
		return
	}

	for key, f := range functionsTmp {

		parts := strings.Split(key, ":")
		pkg := parts[0]
		name := parts[1]
		alias := packages[pkg]
		id := alias + "." + name
		functions[id] = f
		log.RootLogger().Debugf("Resolved function '%s' to '%s", key, id)
	}

	//remove temp function holder
	functionsTmp = nil
}

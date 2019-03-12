package property

import (
	"errors"
	"fmt"
	"strings"

	"github.com/project-flogo/core/support/log"
)

var (
	RegisteredResolvers = make(map[string]Resolver)
	EnabledResolvers    []Resolver
)

// Resolver used to resolve property value from external configuration like env, file etc
type Resolver interface {
	// Should return value and true if the given key exists in the external configuration otherwise should return nil and false.
	LookupValue(key string) (interface{}, bool)
}

func RegisterPropertyResolver(resolverType string, resolver Resolver) error {

	logger := log.RootLogger()

	if resolverType == "" {
		return fmt.Errorf("'resolverType' must be specified when registering a property resolver")
	}

	if resolver == nil {
		return fmt.Errorf("cannot register 'nil' property resolver")
	}

	if _, dup := RegisteredResolvers[resolverType]; dup {
		return fmt.Errorf("property resolver already registered: %s", resolverType)
	}

	logger.Debugf("Registering property resolver [ %s ]", resolverType)

	RegisteredResolvers[resolverType] = resolver

	return nil
}

func GetPropertyResolver(resolverType string) Resolver {
	return RegisteredResolvers[resolverType]
}

func EnablePropertyResolvers(resolverTypes string) error {

	for _, resolverType := range strings.Split(resolverTypes, ",") {
		resolver := RegisteredResolvers[resolverType]
		if resolver == nil {
			errMag := fmt.Sprintf("Unsupported property resolver type - %s. Resolver not registered.", resolverType)
			return errors.New(errMag)
		}
		EnabledResolvers = append(EnabledResolvers, resolver)
	}

	return nil
}

func ResolveProperty(propertyName string) (interface{}, bool) {

	for _, resolver := range EnabledResolvers {
		// Use resolver
		value, resolved := resolver.LookupValue(propertyName)
		if resolved {
			return value, true
		}
	}

	return nil, false
}

func PropertyResolverProcessor(properties map[string]interface{}) error {

	logger := log.RootLogger()

	for name := range properties {
		newVal, found := ResolveProperty(name)

		if !found {
			logger.Warnf("Property '%s' could not be resolved using property resolver(s). Using default value from flogo.json.", name)
		} else {
			properties[name] = newVal
		}
	}

	return nil
}

package property

import (
	"errors"
	"fmt"
	"strings"

	"github.com/project-flogo/core/support/log"
)

var (
	externalResolverMap = make(map[string]ExternalResolver)
	externalResolvers   []ExternalResolver
)

// Resolver used to resolve property value from external configuration like env, file etc
type ExternalResolver interface {
	// Should return value and true if the given key exists in the external configuration otherwise should return nil and false.
	LookupValue(key string) (interface{}, bool)
}

func RegisterExternalResolver(resolverType string, resolver ExternalResolver) error {

	logger := log.RootLogger()

	if resolverType == "" {
		return fmt.Errorf("'resolverType' must be specified when registering external property resolver")
	}

	if resolver == nil {
		return fmt.Errorf("cannot register 'nil' external property resolver")
	}

	if _, dup := externalResolverMap[resolverType]; dup {
		return fmt.Errorf("external property resolver already registered: %s", resolverType)
	}

	logger.Debugf("Registering external property resolver [ %s ]", resolverType)

	externalResolverMap[resolverType] = resolver

	return nil
}

func GetExternalResolver(resolverType string) ExternalResolver {
	return externalResolverMap[resolverType]
}

func EnableExternalResolvers(resolverTypes string) error {

	for _, resolverType := range strings.Split(resolverTypes, ",") {
		resolver := externalResolverMap[resolverType]
		if resolver == nil {
			errMag := fmt.Sprintf("Unsupported external property resolver type - %s. Resolver not registered.", resolverType)
			return errors.New(errMag)
		}
		externalResolvers = append(externalResolvers, resolver)
	}

	return nil
}

func ResolveExternally(propertyName string) (interface{}, bool) {

	for _, resolver := range externalResolvers {
		// Use resolver
		value, resolved := resolver.LookupValue(propertyName)
		if resolved {
			return value, true
		}
	}

	return nil, false
}

func ExternalPropertyResolverProcessor(properties map[string]interface{}) error {

	logger := log.RootLogger()

	for name := range properties {
		newVal, found := ResolveExternally(name)

		if !found {
			logger.Warnf("Property '%s' could not be resolved using external resolver(s) '%s'. Using default value.", name)
		} else {
			properties[name] = newVal
		}
	}

	return nil
}

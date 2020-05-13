package property

import (
	"errors"
	"fmt"
	"strings"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/log"
)

var (
	RegisteredResolvers = make(map[string]ExternalResolver)
	EnabledResolvers    []ExternalResolver
)

// Resolver used to resolve property value from external configuration like env, file etc
type ExternalResolver interface {
	// Name of the resolver (e.g., consul)
	Name() string
	// Should return value and true if the given key exists in the external configuration otherwise should return nil and false.
	LookupValue(key string) (interface{}, bool)
}

//DEPRECATED
func RegisterPropertyResolver(resolver ExternalResolver) error {
	return RegisterExternalResolver(resolver)
}

func RegisterExternalResolver(resolver ExternalResolver) error {

	logger := log.RootLogger()

	resolverName := resolver.Name()

	if resolverName == "" {
		return fmt.Errorf("an external property resolver must have a non-empty name")
	}

	if resolver == nil {
		return fmt.Errorf("cannot register a 'nil' external property resolver")
	}

	if _, dup := RegisteredResolvers[resolverName]; dup {
		return fmt.Errorf("external property resolver already registered: %s", resolverName)
	}

	logger.Debugf("Registering external property resolver [ %s ]", resolverName)

	RegisteredResolvers[resolverName] = resolver

	return nil
}

func GetExternalPropertyResolver(resolverType string) ExternalResolver {
	return RegisteredResolvers[resolverType]
}

func EnableExternalPropertyResolvers(resolverTypes string) error {

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

func ResolvePropertyExternally(propertyName string) (interface{}, bool) {

	for _, resolver := range EnabledResolvers {
		// Use resolver
		value, resolved := resolver.LookupValue(propertyName)
		if resolved {
			return value, true
		}
	}

	return nil, false
}

func ExternalResolverProcessor(properties map[string]interface{}) error {

	logger := log.RootLogger()

	var enabledResolvers []string
	var resolver ExternalResolver
	for _, resolver = range EnabledResolvers {
		enabledResolvers = append(enabledResolvers, resolver.Name())
	}
	if len(enabledResolvers) == 1 {
		logger.Infof("Properties will be resolved with the '%s' resolver", EnabledResolvers[0].Name())
	} else {
		logger.Infof("Properties will be resolved with these resolvers (in decreasing order of priority): %v", enabledResolvers)
	}

	for name := range properties {
		newVal, found := ResolvePropertyExternally(name)

		if !found {
			logger.Warnf("Property '%s' could not be resolved using property resolver(s). Using default value from flogo.json.", name)
		} else {
			// Get datatype of old value
			dType, _ := data.GetType(properties[name])
			if dType != data.TypeUnknown {
				coercedVal, err := coerce.ToType(newVal, dType)
				if err == nil {
					properties[name] = coercedVal
					continue
				}
			}
			properties[name] = newVal
		}
	}

	return nil
}

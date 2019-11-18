package resolve

import (
	"fmt"
	"strings"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/path"
)

var defaultResolver = NewCompositeResolver(map[string]Resolver{
	".":        &ScopeResolver{},
	"env":      &EnvResolver{},
})

func GetBasicResolver() CompositeResolver {
	return defaultResolver
}

func NewCompositeResolver(resolvers map[string]Resolver, options ...func(r *basicResolver)) CompositeResolver {

	resolver := &basicResolver{resolvers: resolvers}
	resolver.cleanDirective = defaultDereferenceCleaner

	for _, option := range options {
		option(resolver)
	}

	return resolver
}

func NoDereferencing(r *basicResolver) {
	r.cleanDirective = nil
}

func CustomDereferenceCleaner(f DereferenceCleaner) func(*basicResolver) {
	return func(r *basicResolver) {
		r.cleanDirective = f
	}
}

// DereferenceCleaner removes the dereference characters from the resolve directive
type DereferenceCleaner func(string) (string, bool)

func defaultDereferenceCleaner(value string) (string, bool) {
	if strings.HasPrefix(value, "$") {
		return value[1:], true
	}

	return value, false
}

type basicResolver struct {
	resolvers      map[string]Resolver
	cleanDirective DereferenceCleaner
}

// Resolve implements CompositeResolver.Resolve
func (r *basicResolver) Resolve(directive string, scope data.Scope) (value interface{}, err error) {

	if r.cleanDirective != nil {

		sansDereferencer, ok := r.cleanDirective(directive)
		if !ok {
			if scope == nil {
				//todo should we throw an error in this circumstance?
				return directive, nil
			}

			val, _ := scope.GetValue(directive)
			return val, nil
		}

		directive = sansDereferencer
	}

	resolverName, nextIdx := GetResolverInfo(directive)

	resolver, exists := r.resolvers[resolverName]
	if !exists {
		return nil, fmt.Errorf("unable to find a '%s' resolver", resolverName)
	}

	details, err := GetResolveDirectiveDetails(directive[nextIdx:], resolver.GetResolverInfo().UsesItemFormat(), resolver.GetResolverInfo().IsImplicit())
	if err != nil {
		return nil, err
	}

	value, err = resolver.Resolve(scope, details.ItemName, details.ValueName)
	if err != nil {
		return nil, err
	}

	if details.Path != "" {
		value, err = path.GetValue(value, details.Path)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}

// GetResolution implements CompositeResolver.GetResolution
func (r *basicResolver) GetResolution(directive string) (Resolution, error) {
	if r.cleanDirective != nil {

		sansDereferencer, ok := r.cleanDirective(directive)
		if !ok {
			return nil, fmt.Errorf("invalid resolution directive '%s' unable to remove deferencing characters", directive)
		}

		directive = sansDereferencer
	}

	resolverName, nextIdx := GetResolverInfo(directive)

	resolver, exists := r.resolvers[resolverName]
	if !exists {
		return nil, fmt.Errorf("unable to find a '%s' resolver", resolverName)
	}

	details, err := GetResolveDirectiveDetails(directive[nextIdx:], resolver.GetResolverInfo().UsesItemFormat(), resolver.GetResolverInfo().IsImplicit())
	if err != nil {
		return nil, fmt.Errorf("cannot resolve '%s' : %v", directive, err)
	}

	if resolver.GetResolverInfo().IsStatic() {
		// its a static resolver, so we can go ahead and resolve the value now
		val, err := resolver.Resolve(nil, details.ItemName, details.ValueName)
		if err != nil {
			return nil, err
		}

		if details.Path != "" {
			val, err = path.GetValue(val, details.Path)
			if err != nil {
				return nil, err
			}
		}

		return &staticResolution{val}, nil
	}

	return &resolution{resolver: resolver, details: details}, nil
}

//staticResolution is a resolution that has a static value
type staticResolution struct {
	value interface{}
}

// IsStatic implements data.IsStatic
func (*staticResolution) IsStatic() bool {
	return true
}

// GetValue implements data.Resolution.GetValue
func (r *staticResolution) GetValue(scope data.Scope) (interface{}, error) {
	return r.value, nil
}

// resolution implements data.Resolution
type resolution struct {
	resolver Resolver
	details  *ResolveDirectiveDetails
}

// IsStatic implements data.IsStatic
func (*resolution) IsStatic() bool {
	return false
}

// GetValue implements data.Resolution.GetValue
func (r *resolution) GetValue(scope data.Scope) (interface{}, error) {

	value, err := r.resolver.Resolve(scope, r.details.ItemName, r.details.ValueName)
	if err != nil {
		return nil, err
	}

	if r.details.Path != "" {
		value, err = path.GetValue(value, r.details.Path)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}

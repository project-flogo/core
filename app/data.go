package app

import (
	"fmt"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support/log"
)

var appData = data.NewSimpleSyncScope(nil, nil)

var resolverInfo = resolve.NewResolverInfo(false, false)

type AppResolver struct {
}

func (r *AppResolver) GetResolverInfo() *resolve.ResolverInfo {
	return resolverInfo
}

func (r *AppResolver) Resolve(scope data.Scope, itemName, valueName string) (interface{}, error) {

	value, exists := appData.GetValue(valueName)
	if !exists {
		return nil, fmt.Errorf("failed to resolve app attr: '%s', not found in app", valueName)
	}

	return value, nil
}

// GetValue gets an app attribute value
func GetValue(name string) (value interface{}, exists bool) {
	if log.RootLogger().TraceEnabled() {
		log.RootLogger().Tracef("Getting App Value '%s': %v", name)
	}
	return appData.GetValue(name)
}

// SetValue sets an app attribute value
func SetValue(name string, value interface{}) error {
	if log.RootLogger().TraceEnabled() {
		log.RootLogger().Tracef("Set App Value '%s': %v", name, value)
	}
	return appData.SetValue(name, value)
}

// DeleteValue remove an app attribute
func DeleteValue(name string) {
	if log.RootLogger().TraceEnabled() {
		log.RootLogger().Tracef("Deleting App Value '%s'", name)
	}
	if d, ok := appData.(data.NeedsDelete); ok {
		d.Delete(name)
	}
}

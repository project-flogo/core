package activity

import (
	"fmt"
	"github.com/project-flogo/core/support/log"
)

type void struct{}

var (
	hasLegacy     = false
	empty         void
	legacyTracker = make(map[string]void)
)

//DEPRECATED
func HasLegacyActivities() bool {
	return hasLegacy
}

//DEPRECATED
func IsLegacyActivity(ref string) bool {
	_, ok := legacyTracker[ref]
	return ok
}

//DEPRECATED
func LegacyRegister(ref string, activity Activity) error {

	if ref == "" {
		return fmt.Errorf("'ref' must be specified when registering")
	}

	if activity == nil {
		return fmt.Errorf("cannot register 'nil' activity")
	}

	if _, dup := activities[ref]; dup {
		return fmt.Errorf("activity already registered: %s", ref)
	}

	log.RootLogger().Debugf("Registering legacy activity: %s", ref)

	hasLegacy = true
	activities[ref] = activity
	legacyTracker[ref] = empty
	activityLoggers[ref] = log.CreateLoggerFromRef(rootLogger, "activity", ref)
	return nil
}

type LegacyCtx interface {

	// GetOutput gets the value of the specified output attribute
	GetOutput(name string) interface{}
	GetSetting(name string) (value interface{}, exists bool)
}

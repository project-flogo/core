package activity

import (
	"fmt"

	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/logger"
)

var (
	activities        = make(map[string]Activity)
	activityFactories = make(map[string]Factory)
)

func Register(activity Activity, f ...Factory) error {

	if activity == nil {
		return fmt.Errorf("cannot register 'nil' activity")
	}

	ref := GetRef(activity)

	if _, dup := activities[ref]; dup {
		return fmt.Errorf("activity already registered: %s", ref)
	}

	logger.Debugf("Registering activity [ %s ]", ref)

	activities[ref] = activity

	if len(f) > 1 {
		logger.Debugf("Only one factory can be associated with activity [ %s ]", ref)
	}

	if len(f) == 1 {
		activityFactories[ref] = f[0]
	}

	return nil
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

	logger.Debugf("Registering trigger [ %s ]", ref)

	activities[ref] = activity

	return nil
}

func GetRef(activity Activity) string {
	return support.GetRef(activity)
}

// Get gets specified activity by ref
func Get(ref string) Activity {
	return activities[ref]
}

// GetFactory gets activity factory by ref
func GetFactory(ref string) Factory {
	return activityFactories[ref]
}

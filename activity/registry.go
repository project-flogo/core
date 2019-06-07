package activity

import (
	"fmt"
	"path"

	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/log"
)

var (
	activities        = make(map[string]Activity)
	activityFactories = make(map[string]Factory)
	activityLoggers   = make(map[string]log.Logger)
)

var activityLogger = log.ChildLogger(log.RootLogger(), "activity")

func Register(activity Activity, f ...Factory) error {

	if activity == nil {
		return fmt.Errorf("cannot register 'nil' activity")
	}

	ref := GetRef(activity)

	if _, dup := activities[ref]; dup {
		return fmt.Errorf("activity already registered: %s", ref)
	}

	log.RootLogger().Debugf("Registering activity: %s", ref)

	activities[ref] = activity
	name := path.Base(ref) //todo should we use this or the alias?
	activityLoggers[ref] = log.ChildLogger(activityLogger, name)

	if len(f) > 1 {
		log.RootLogger().Warnf("Only one factory can be associated with activity: %s", ref)
	}

	if len(f) == 1 {
		activityFactories[ref] = f[0]
	}

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

// GetLogger gets activity logger by ref
func GetLogger(ref string) log.Logger {
	if ref[0] == '#' {
		ref, _ = support.GetAliasRef("activity", ref[1:])
	}

	logger, ok := activityLoggers[ref]
	if ok {
		return logger
	} else {
		return log.RootLogger()
	}
}

func CleanupSingletons() {
	for ref, activity := range activities {

		if _, ok := activityFactories[ref]; !ok {
			//singleton activities don't have factories
			if needsCleanup, ok := activity.(support.NeedsCleanup); ok {
				err := needsCleanup.Cleanup()
				if err != nil {
					log.RootLogger().Errorf("Error cleaning up activity '%s' : ", ref, err)
				}
			}
		}
	}
}

func IsSingleton(activity Activity) bool {
	ref := support.GetRef(activity)
	_, hasFactory := activityFactories[ref]

	//if it doesn't have a factory, it is a singleton/shared activity
	return !hasFactory
}

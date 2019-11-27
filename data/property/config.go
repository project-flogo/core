package property

import (
	"os"
	"strings"
)

var EnvAppPropertyDynamicUpdate = "FLOGO_APP_PROP_SNAPSHOTS"

func IsPropertySnapshotEnabled() bool {
	dynamicUpdateEnv := os.Getenv(EnvAppPropertyDynamicUpdate)
	if strings.EqualFold(dynamicUpdateEnv, "true") {
		return true
	}

	return false
}
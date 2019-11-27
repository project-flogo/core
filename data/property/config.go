package property

import (
	"os"
	"strings"
)

var EnvAppPropertySnapshotEnabled = "FLOGO_APP_PROP_SNAPSHOTS"

func IsPropertySnapshotEnabled() bool {
	dynamicUpdateEnv := os.Getenv(EnvAppPropertySnapshotEnabled)
	if strings.EqualFold(dynamicUpdateEnv, "true") {
		return true
	}

	return false
}
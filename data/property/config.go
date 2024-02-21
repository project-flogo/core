package property

import (
	"os"
	"strings"
)

var EnvAppPropertySnapshotEnabled = "FLOGO_APP_PROP_SNAPSHOTS"
var EnvAppPropertyReconfigure = "FLOGO_APP_PROP_RECONFIGURE"

func IsPropertySnapshotEnabled() bool {
	appPropertySnapshotEnabled := os.Getenv(EnvAppPropertySnapshotEnabled)
	if strings.EqualFold(appPropertySnapshotEnabled, "true") {
		return true
	}
	return false
}

func IsPropertyReconfigureEnabled() bool {
	appPropertyAutoReconfigureEnabled := os.Getenv(EnvAppPropertyReconfigure)
	return strings.EqualFold(appPropertyAutoReconfigureEnabled, "true")
}

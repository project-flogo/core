package core

import _ "embed"

//go:embed VERSION
var version string

// Version will return the core version
func Version() string {
	return version
}

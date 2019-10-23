package config

import (
	"os"
	"strconv"
)

const (
	EnvMappingIgnoreError        = "FLOGO_MAPPING_IGNORE_ERRORS"
	EnvMappingIgnoreErrorDefault = false

	EnvMappingSkipMissing        = "FLOGO_MAPPING_SKIP_MISSING"
	EnvMappingSkipMissingDefault = false
)

func IsMappingIgnoreError() bool {
	relaxed := os.Getenv(EnvMappingIgnoreError)
	if len(relaxed) <= 0 {
		return EnvMappingIgnoreErrorDefault
	}
	b, _ := strconv.ParseBool(relaxed)
	return b
}

func IsMappingSkipMissing() bool {
	relaxed := os.Getenv(EnvMappingSkipMissing)
	if len(relaxed) <= 0 {
		return EnvMappingSkipMissingDefault
	}
	b, _ := strconv.ParseBool(relaxed)
	return b
}

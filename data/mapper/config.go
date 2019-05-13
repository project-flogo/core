package mapper

import (
	"os"
	"strconv"
)

const (
	EnvMappingRelexed        = "FLOGO_MAPPING_RELAXED"
	EnvMappingRelexedDefault = false
)

func IsMappingRelaxed() bool {
	relaxed := os.Getenv(EnvMappingRelexed)
	if len(relaxed) <= 0 {
		return EnvMappingRelexedDefault
	}
	b, _ := strconv.ParseBool(relaxed)
	return b
}

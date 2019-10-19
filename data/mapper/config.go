package mapper

import (
	"os"
	"strconv"
)

const (
	EnvMappingRelexed        = "FLOGO_MAPPING_RELAXED"
	EnvMappingRelexedDefault = false
	EnvMappingOmitNull        = "FLOGO_MAPPING_OMIT_NULL"
	EnvMappingOmitNullDefault = false
)

func IsMappingRelaxed() bool {
	relaxed := os.Getenv(EnvMappingRelexed)
	if len(relaxed) <= 0 {
		return EnvMappingRelexedDefault
	}
	b, _ := strconv.ParseBool(relaxed)
	return b
}

func OmitNull() bool {
	omitNull := os.Getenv(EnvMappingOmitNull)
	if len(omitNull) <= 0 {
		return EnvMappingOmitNullDefault
	}
	b, _ := strconv.ParseBool(omitNull)
	return b
}

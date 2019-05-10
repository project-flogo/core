package mapper

import (
	"github.com/project-flogo/core/support/log"
	"os"
	"strings"
)

const (
	EnvDataResolving        = "FLOGO_RESOLVING"
	EnvDataResolvingStrict  = "STRICT"
	EnvDataResolvingRelaxed = "RELAXED"
)

func GetDataResolving() string {
	dataResolver := os.Getenv(EnvDataResolving)
	if len(dataResolver) > 0 {
		if dataResolver != EnvDataResolvingRelaxed && dataResolver != EnvDataResolvingStrict {
			log.RootLogger().Warnf("Unknow FLOGO_DATA_RESOLVING [%s], we only support [STRICT] or [RELAXED]", dataResolver)
			return EnvDataResolvingStrict
		}
		return dataResolver
	}
	return EnvDataResolvingStrict
}

func isRelexed() bool {
	return strings.EqualFold(GetDataResolving(), EnvDataResolvingRelaxed)
}

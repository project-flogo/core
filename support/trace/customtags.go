package trace

import (
	"os"
	"strconv"
)

var traceCustomTagsEnabled = envBoolDefault("FLOGO_OTEL_TRACE_USE_CUSTOM_TAGS", true)
var metricsCustomTagsEnabled = envBoolDefault("FLOGO_OTEL_METRICS_USE_CUSTOM_TAGS", true)

func TraceCustomTagsEnabled() bool {
	return traceCustomTagsEnabled
}

func MetricsCustomTagsEnabled() bool {
	return metricsCustomTagsEnabled
}

func envBoolDefault(key string, defaultVal bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return b
}

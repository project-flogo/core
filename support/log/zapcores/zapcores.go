package zapcores

import (
	"go.uber.org/zap/zapcore"
)

// zapCores holds log cores
var zapCores map[string]zapcore.Core

// zapTraceCores holds trace log cores
var zapTraceCores map[string]zapcore.Core

func init() {
	zapCores = make(map[string]zapcore.Core)
	zapTraceCores = make(map[string]zapcore.Core)
}

// RegisterLogCore adds core to zapCores
func RegisterLogCore(name string, core zapcore.Core) {
	zapCores[name] = core
}

// RegisterTraceLogCore adds trace core to zapTraceCores
func RegisterTraceLogCore(name string, core zapcore.Core) {
	zapTraceCores[name] = core
}

// RegisteredCores returns complete log core map
func RegisteredCores() map[string]zapcore.Core {
	return zapCores
}

// RegisteredTraceCores returns complete trace log core map
func RegisteredTraceCores() map[string]zapcore.Core {
	return zapTraceCores
}

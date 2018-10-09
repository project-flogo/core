package logger

import (
	"fmt"
)

func Debug(args ...interface{}) {
	GetDefaultLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	GetDefaultLogger().Debugf(format, args...)
}

func Info(args ...interface{}) {
	GetDefaultLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	GetDefaultLogger().Infof(format, args...)
}

func Warn(args ...interface{}) {
	GetDefaultLogger().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	GetDefaultLogger().Warnf(format, args...)
}

func Error(args ...interface{}) {
	GetDefaultLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	GetDefaultLogger().Errorf(format, args...)
}

func SetLogLevel(level Level) {
	GetDefaultLogger().SetLogLevel(level)
}

func GetLogLevel() Level {
	return GetDefaultLogger().GetLogLevel()
}

func DebugEnabled() bool {
	return GetDefaultLogger().DebugEnabled()
}

var defaultLoggerName = "flogo"
var defaultLogLevel = "INFO"

func SetDefaultLogger(name string) {
	defaultLoggerName = name
}

func GetDefaultLogger() Logger {
	defLogger := GetLogger(defaultLoggerName)
	if defLogger == nil {
		errorMsg := fmt.Sprintf("error getting default logger '%s'", defaultLoggerName)
		panic(errorMsg)
	}
	return defLogger
}

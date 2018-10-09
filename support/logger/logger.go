package logger

import (
	"fmt"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	SetLogLevel(Level)
	GetLogLevel() Level
	DebugEnabled() bool
}

type LoggerFactory interface {
	GetLogger(name string) Logger
}

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

var levelNames = initLevelNames()

func initLevelNames() map[string]Level {
	newLevelNames := make(map[string]Level, 4)
	newLevelNames["DEBUG"] = DebugLevel
	newLevelNames["INFO"] = InfoLevel
	newLevelNames["WARN"] = WarnLevel
	newLevelNames["ERROR"] = ErrorLevel
	return newLevelNames
}

var logFactory LoggerFactory

func RegisterLoggerFactory(factory LoggerFactory) {
	logFactory = factory
}

// GetLogger returns the Logger using the logFactory registered.
// Returns nil if no factory is found
func GetLogger(name string) Logger {
	if logFactory == nil {
		return nil
	}
	return logFactory.GetLogger(name)
}

func GetLevelForName(name string) (Level, error) {
	levelForName, ok := levelNames[name]
	if !ok {
		return 0, fmt.Errorf("unsupported Log Level '%s'", name)
	}
	return levelForName, nil
}

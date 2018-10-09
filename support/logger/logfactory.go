package logger

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	EnvKeyLogDateFormat  = "FLOGO_LOG_DTFORMAT"
	DefaultLogDateFormat = "2006-01-02 15:04:05.000"
	EnvKeyLogLevel       = "FLOGO_LOG_LEVEL"
	DefaultLogLevel      = "INFO"
)

var loggerMap = make(map[string]Logger)
var mutex = &sync.RWMutex{}
var logLevel = InfoLevel
var timeFormat = DefaultLogDateFormat

type DefaultLoggerFactory struct {
}

func init() {

	RegisterLoggerFactory(&DefaultLoggerFactory{})

	logLevelName := getLogLevel()
	// Get log level for name
	getLogLevel, err := GetLevelForName(logLevelName)
	if err != nil {
		println("Unsupported Log Level - [" + logLevelName + "]. Set to Log Level - [" + defaultLogLevel + "]")
	} else {
		logLevel = getLogLevel
	}

	timeFormat = getLogDateTimeFormat()
}

//GetLogLevel returns the log level
func getLogLevel() string {
	logLevelEnv := os.Getenv(EnvKeyLogLevel)
	if len(logLevelEnv) > 0 {
		return logLevelEnv
	}
	return DefaultLogLevel
}

func getLogDateTimeFormat() string {
	logLevelEnv := os.Getenv(EnvKeyLogDateFormat)
	if len(logLevelEnv) > 0 {
		return logLevelEnv
	}
	return DefaultLogDateFormat
}

type DefaultLogger struct {
	loggerName string
	loggerImpl *logrus.Logger
}

type LogFormatter struct {
	loggerName string
}

func SetLogDateTimeFormat(format string) {
	timeFormat = format
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	logEntry := fmt.Sprintf("%s %-6s [%s] - %s\n", entry.Time.Format(timeFormat), getLevel(entry.Level), f.loggerName, entry.Message)
	return []byte(logEntry), nil
}

func getLevel(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "DEBUG"
	case logrus.InfoLevel:
		return "INFO"
	case logrus.ErrorLevel:
		return "ERROR"
	case logrus.WarnLevel:
		return "WARN"
	case logrus.PanicLevel:
		return "PANIC"
	case logrus.FatalLevel:
		return "FATAL"
	}

	return "UNKNOWN"
}

// Debug logs message at Debug level.
func (l *DefaultLogger) Debug(args ...interface{}) {
	l.loggerImpl.Debug(args...)
}

// DebugEnabled checks if Debug level is enabled.
func (l *DefaultLogger) DebugEnabled() bool {
	return l.loggerImpl.Level >= logrus.DebugLevel
}

// Info logs message at Info level.
func (l *DefaultLogger) Info(args ...interface{}) {
	l.loggerImpl.Info(args...)
}

// InfoEnabled checks if Info level is enabled.
func (l *DefaultLogger) InfoEnabled() bool {
	return l.loggerImpl.Level >= logrus.InfoLevel
}

// Warn logs message at Warning level.
func (l *DefaultLogger) Warn(args ...interface{}) {
	l.loggerImpl.Warn(args...)
}

// WarnEnabled checks if Warning level is enabled.
func (l *DefaultLogger) WarnEnabled() bool {
	return l.loggerImpl.Level >= logrus.WarnLevel
}

// Error logs message at Error level.
func (l *DefaultLogger) Error(args ...interface{}) {
	l.loggerImpl.Error(args...)
}

// ErrorEnabled checks if Error level is enabled.
func (l *DefaultLogger) ErrorEnabled() bool {
	return l.loggerImpl.Level >= logrus.ErrorLevel
}

// Debug logs message at Debug level.
func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	l.loggerImpl.Debugf(format, args...)
}

// Info logs message at Info level.
func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	l.loggerImpl.Infof(format, args...)
}

// Warn logs message at Warning level.
func (l *DefaultLogger) Warnf(format string, args ...interface{}) {
	l.loggerImpl.Warnf(format, args...)
}

// Error logs message at Error level.
func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	l.loggerImpl.Errorf(format, args...)
}

//SetLog Level
func (l *DefaultLogger) SetLogLevel(logLevel Level) {
	switch logLevel {
	case DebugLevel:
		l.loggerImpl.Level = logrus.DebugLevel
	case InfoLevel:
		l.loggerImpl.Level = logrus.InfoLevel
	case ErrorLevel:
		l.loggerImpl.Level = logrus.ErrorLevel
	case WarnLevel:
		l.loggerImpl.Level = logrus.WarnLevel
	default:
		l.loggerImpl.Level = logrus.ErrorLevel
	}
}

func (l *DefaultLogger) GetLogLevel() Level {
	levelStr := getLevel(l.loggerImpl.Level)
	level, _ := GetLevelForName(levelStr)
	return level
}

func (logfactory *DefaultLoggerFactory) GetLogger(name string) Logger {
	mutex.RLock()
	l := loggerMap[name]
	mutex.RUnlock()
	if l == nil {
		logImpl := logrus.New()
		logImpl.Formatter = &LogFormatter{
			loggerName: name,
		}
		l = &DefaultLogger{
			loggerName: name,
			loggerImpl: logImpl,
		}

		l.SetLogLevel(logLevel)

		mutex.Lock()
		loggerMap[name] = l
		mutex.Unlock()
	}
	return l
}

package log

import (
	"os"
	"path/filepath"
	"strings"
)

type Level int
type Format int

const (
	EnvKeyLogCtx         = "FLOGO_LOG_CTX"
	EnvKeyLogDateFormat  = "FLOGO_LOG_DTFORMAT"
	DefaultLogDateFormat = "2006-01-02 15:04:05.000"
	EnvKeyLogLevel       = "FLOGO_LOG_LEVEL"
	DefaultLogLevel      = InfoLevel
	EnvKeyLogFormat      = "FLOGO_LOG_FORMAT"
	DefaultLogFormat     = FormatConsole

	EnvKeyLogSeparator  = "FLOGO_LOG_SEPARATOR"
	DefaultLogSeparator = "\t"

	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel

	FormatConsole Format = iota
	FormatJson
)

type Logger interface {
	DebugEnabled() bool
	TraceEnabled() bool

	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})

	Tracef(template string, args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})

	Structured() StructuredLogger
}

type StructuredLogger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

type Field = interface{}

var (
	rootLogger Logger
	ctxLogging bool
)

func init() {
	configureLogging()
}

func CtxLoggingEnabled() bool {
	return ctxLogging
}

func RootLogger() Logger {
	return rootLogger
}

func SetLogLevel(logger Logger, level Level) {
	setZapLogLevel(logger, level)
}

func ChildLogger(logger Logger, name string) Logger {

	childLogger, err := newZapChildLogger(logger, name)
	if err != nil {
		rootLogger.Warnf("unable to create child logger named: %s - %s", name, err.Error())
		childLogger = logger
	}

	return childLogger
}

func ChildLoggerWithFields(logger Logger, fields ...Field) Logger {
	childLogger, err := newZapChildLoggerWithFields(logger, fields...)
	if err != nil {
		rootLogger.Warnf("unable to create child logger with fields: %s", err.Error())
		childLogger = logger
	}

	return childLogger
}

func CreateLoggerFromRef(logger Logger, contributionType, ref string) Logger {
	ref = strings.TrimSpace(ref)
	if strings.HasSuffix(ref, "/") {
		ref = ref[:len(ref)-1]
	}
	dirs := strings.Split(ref, "/")
	if len(dirs) >= 3 {
		name := dirs[len(dirs)-1]
		acType := dirs[len(dirs)-2]
		if acType == "activity" || acType == "trigger" || acType == "connector" {
			categoryName := dirs[len(dirs)-3]
			return ChildLogger(logger, strings.ToLower(categoryName+"."+acType+"."+name))
		} else {
			return ChildLogger(logger, strings.ToLower(acType+"."+contributionType+"."+name))
		}
	} else {
		return ChildLogger(logger, strings.ToLower(contributionType+"."+filepath.Base(ref)))
	}
}

func Sync() {
	zapSync(rootLogger)
}

var traceEnabled = false

func configureLogging() {
	envLogCtx := os.Getenv(EnvKeyLogCtx)
	if strings.ToLower(envLogCtx) == "true" {
		ctxLogging = true
	}

	rootLogLevel := DefaultLogLevel

	envLogLevel := strings.ToUpper(os.Getenv(EnvKeyLogLevel))
	switch envLogLevel {
	case "TRACE":
		rootLogLevel = DebugLevel
		traceEnabled = true
	case "DEBUG":
		rootLogLevel = DebugLevel
	case "INFO":
		rootLogLevel = InfoLevel
	case "WARN":
		rootLogLevel = WarnLevel
	case "ERROR":
		rootLogLevel = ErrorLevel
	default:
		rootLogLevel = DefaultLogLevel
	}

	logFormat := DefaultLogFormat
	envLogFormat := strings.ToUpper(os.Getenv(EnvKeyLogFormat))
	if envLogFormat == "JSON" {
		logFormat = FormatJson
	}

	rootLogger = newZapRootLogger("flogo", logFormat, rootLogLevel)
	SetLogLevel(rootLogger, rootLogLevel)
}

func ToLogLevel(lvlStr string) Level {

	lvl := DefaultLogLevel

	switch lvlStr {
	case "TRACE":
		lvl = DebugLevel
	case "DEBUG":
		lvl = DebugLevel
	case "INFO":
		lvl = InfoLevel
	case "WARN":
		lvl = WarnLevel
	case "ERROR":
		lvl = ErrorLevel
	default:
		lvl = DefaultLogLevel
	}

	return lvl
}

func getLogSeparator() string {
	v, ok := os.LookupEnv(EnvKeyLogSeparator)
	if ok && len(v) > 0 {
		return v
	}
	return DefaultLogSeparator
}

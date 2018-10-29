package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

type Level int
type Format int

const (
	EnvKeyLogCtx         = "FLOGO_LOG_CTX"
	EnvKeyLogDateFormat  = "FLOGO_LOG_DTFORMAT"
	DefaultLogDateFormat = "2006-01-02 15:04:05.000"
	EnvKeyLogLevel       = "FLOGO_LOG_LEVEL"
	DefaultLogLevel      = LevelInfo
	EnvKeyLogFormat      = "FLOGO_LOG_FORMAT"
	DefaultLogFormat     = FormatConsole

	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError

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
	Trace(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

var	(
	rootLogger Logger

	ctxLogging bool
	rootLogLevel zapcore.Level
	logFormat = DefaultLogFormat
)


func init() {

	configureLogging()

	zl, lvl, _ := newZapLogger()
	lvl.SetLevel(rootLogLevel)

	rootLogger = &loggerImpl{loggerLevel: lvl, zapLogger: zl.Named("flogo").Sugar()}
}

func CtxLoggingEnabled() bool {
	return ctxLogging
}

func RootLogger() Logger {
	return rootLogger
}

func SetLogLevel(level Level, logger Logger) {

	impl := logger.(*loggerImpl)

	switch level {
	case LevelDebug:
		impl.loggerLevel.SetLevel(zapcore.DebugLevel)
	case LevelInfo:
		impl.loggerLevel.SetLevel(zapcore.InfoLevel)
	case LevelWarn:
		impl.loggerLevel.SetLevel(zapcore.WarnLevel)
	case LevelError:
		impl.loggerLevel.SetLevel(zapcore.ErrorLevel)
	}
}

func ChildLogger(l Logger, name string) Logger {

	impl := l.(*loggerImpl)

	zapLogger := impl.zapLogger
	newZl := zapLogger.Named(name)

	return &loggerImpl{loggerLevel: impl.loggerLevel, zapLogger: newZl}
}

func ChildLoggerWith(l Logger, fields ...Field) Logger {

	impl := l.(*loggerImpl)

	zapLogger := impl.zapLogger
	newZl := zapLogger.With(fields...)

	return &loggerImpl{loggerLevel: impl.loggerLevel, zapLogger: newZl}
}

func Sync() {
	impl := rootLogger.(*loggerImpl)
	impl.zapLogger.Sync()
}

func newZapLogger() (*zap.Logger, *zap.AtomicLevel, error) {
	cfg := zap.NewProductionConfig()
	cfg.DisableCaller = true

	eCfg := cfg.EncoderConfig
	eCfg.TimeKey = "timestamp"
	eCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	if logFormat == FormatConsole {
		eCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.Encoding = "console"
		eCfg.EncodeName = nameEncoder
	}

	cfg.EncoderConfig = eCfg

	lvl := cfg.Level
	zl, err := cfg.Build(zap.AddCallerSkip(1))

	return zl, &lvl, err
}

func configureLogging()  {
	envLogCtx := os.Getenv(EnvKeyLogCtx)
	if strings.ToLower(envLogCtx) == "true" {
		ctxLogging = true
	}

	envLogLevel := strings.ToUpper(os.Getenv(EnvKeyLogLevel))
	switch envLogLevel {
	case "DEBUG":
		rootLogLevel = zapcore.DebugLevel
	case "INFO":
		rootLogLevel = zapcore.InfoLevel
	case "WARN":
		rootLogLevel = zapcore.WarnLevel
	case "ERROR":
		rootLogLevel = zapcore.ErrorLevel
	default:
		rootLogLevel = zapcore.InfoLevel
	}

	envLogFormat := strings.ToUpper(os.Getenv(EnvKeyLogFormat))
	if envLogFormat == "JSON" {
		logFormat = FormatJson
	}
}


func nameEncoder(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + loggerName + "] -")
}
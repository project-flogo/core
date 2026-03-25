package log

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var traceLogger *zap.SugaredLogger

type zapLoggerImpl struct {
	loggerLevel *zap.AtomicLevel
	mainLogger  *zap.SugaredLogger
	// FLOGO-17735: add traceID and spanID attributes to log message
	tracePrefix  string
	traceContext map[string]string
}

func (l *zapLoggerImpl) DebugEnabled() bool {
	return l.loggerLevel.Enabled(zapcore.DebugLevel)
}

func (l *zapLoggerImpl) TraceEnabled() bool {
	return traceEnabled && l.loggerLevel.Enabled(zapcore.DebugLevel)
}

func (l *zapLoggerImpl) Trace(args ...any) {
	if traceEnabled {
		traceLogger.Debug(args...)
	}
}

func (l *zapLoggerImpl) Debug(args ...any) {
	if l.tracePrefix != "" {
		s := make([]any, 1, 1+len(args))
		s[0] = l.tracePrefix
		l.mainLogger.Debug(append(s, args...)...)
	} else {
		l.mainLogger.Debug(args...)
	}
}

func (l *zapLoggerImpl) Info(args ...any) {
	if l.tracePrefix != "" {
		s := make([]any, 1, 1+len(args))
		s[0] = l.tracePrefix
		l.mainLogger.Info(append(s, args...)...)
	} else {
		l.mainLogger.Info(args...)
	}
}

func (l *zapLoggerImpl) Warn(args ...any) {
	if l.tracePrefix != "" {
		s := make([]any, 1, 1+len(args))
		s[0] = l.tracePrefix
		l.mainLogger.Warn(append(s, args...)...)
	} else {
		l.mainLogger.Warn(args...)
	}
}

func (l *zapLoggerImpl) Error(args ...any) {
	if l.tracePrefix != "" {
		s := make([]any, 1, 1+len(args))
		s[0] = l.tracePrefix
		l.mainLogger.Error(append(s, args...)...)
	} else {
		l.mainLogger.Error(args...)
	}
}

func (l *zapLoggerImpl) Tracef(template string, args ...any) {
	if traceEnabled {
		traceLogger.Debugf(template, args...)
	}
}

func (l *zapLoggerImpl) Debugf(template string, args ...any) {
	if l.tracePrefix != "" {
		l.mainLogger.Debugf(l.tracePrefix+template, args...)
	} else {
		l.mainLogger.Debugf(template, args...)
	}
}

func (l *zapLoggerImpl) Infof(template string, args ...any) {
	if l.tracePrefix != "" {
		l.mainLogger.Infof(l.tracePrefix+template, args...)
	} else {
		l.mainLogger.Infof(template, args...)
	}
}

func (l *zapLoggerImpl) Warnf(template string, args ...any) {
	if l.tracePrefix != "" {
		l.mainLogger.Warnf(l.tracePrefix+template, args...)
	} else {
		l.mainLogger.Warnf(template, args...)
	}
}

func (l *zapLoggerImpl) Errorf(template string, args ...any) {
	if l.tracePrefix != "" {
		l.mainLogger.Errorf(l.tracePrefix+template, args...)
	} else {
		l.mainLogger.Errorf(template, args...)
	}
}

func (l *zapLoggerImpl) Structured() StructuredLogger {
	return &zapStructuredLoggerImpl{zl: l.mainLogger.Desugar()}
}

func (l *zapLoggerImpl) GetTracingContext() map[string]string {
	return l.traceContext
}

func (l *zapLoggerImpl) SetTracingContext(traceContext map[string]string) {
	if !traceContextLogging {
		return
	}
	l.traceContext = traceContext
	if traceContext != nil && traceContext[KeyTraceID] != "" {
		l.tracePrefix = fmt.Sprintf("[%s: %s] [%s: %s] ", KeyTraceID, traceContext[KeyTraceID], KeySpanID, traceContext[KeySpanID])
	} else {
		l.tracePrefix = ""
	}
}

type zapStructuredLoggerImpl struct {
	lvl *zap.AtomicLevel
	zl  *zap.Logger
}

func (l *zapStructuredLoggerImpl) Debug(msg string, fields ...Field) {

	fs := make([]zap.Field, len(fields))
	for i, f := range fields {
		fs[i] = f.(zap.Field)
	}

	l.zl.Debug(msg, fs...)
}

func (l *zapStructuredLoggerImpl) Info(msg string, fields ...Field) {
	fs := make([]zap.Field, len(fields))
	for i, f := range fields {
		fs[i] = f.(zap.Field)
	}

	l.zl.Info(msg, fs...)
}

func (l *zapStructuredLoggerImpl) Warn(msg string, fields ...Field) {
	fs := make([]zap.Field, len(fields))
	for i, f := range fields {
		fs[i] = f.(zap.Field)
	}

	l.zl.Warn(msg, fs...)
}

func (l *zapStructuredLoggerImpl) Error(msg string, fields ...Field) {
	fs := make([]zap.Field, len(fields))
	for i, f := range fields {
		fs[i] = f.(zap.Field)
	}

	l.zl.Error(msg, fs...)
}

func setZapLogLevel(logger Logger, level Level) {
	impl, ok := logger.(*zapLoggerImpl)

	if ok {
		zapLevel := toZapLogLevel(level)
		impl.loggerLevel.SetLevel(zapLevel)
	}
}

func toZapLogLevel(level Level) zapcore.Level {
	switch level {
	case DebugLevel, TraceLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	}

	return zapcore.InfoLevel
}

func newZapRootLogger(name string, format Format, level Level) Logger {

	zl, lvl, _ := newZapLogger(format, level)

	var rootLogger Logger
	if name == "" {
		rootLogger = &zapLoggerImpl{loggerLevel: lvl, mainLogger: zl.Sugar()}
	} else {
		rootLogger = &zapLoggerImpl{loggerLevel: lvl, mainLogger: zl.Named(name).Sugar()}
	}

	if traceEnabled {
		tl, _, _ := newZapTraceLogger(format)
		traceLogger = tl.Sugar()
	}

	return rootLogger
}

func newZapLogger(logFormat Format, level Level) (*zap.Logger, *zap.AtomicLevel, error) {
	cfg := zap.NewProductionConfig()
	cfg.DisableCaller = true

	// change the default output paths for the logger if the env variable is set and has the supported values (stderr and stdout).
	// Otherwise, the logger will use the default value of stderr.
	logstream, ok := os.LookupEnv(EnvLogConsoleStream)
	if ok {
		if strings.ToLower(logstream) == "stdout" || strings.ToLower(logstream) == "stderr" {
			cfg.OutputPaths = []string{strings.ToLower(logstream)}
		}
	}

	eCfg := cfg.EncoderConfig
	eCfg.TimeKey = "timestamp"
	eCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	if logFormat == FormatConsole {
		eCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.Encoding = "console"
		eCfg.EncodeName = nameEncoder
	}

	eCfg.ConsoleSeparator = getLogSeparator()
	//Don't print stacktrace for log level lower than debug
	if level > DebugLevel {
		eCfg.StacktraceKey = ""
	}

	cfg.EncoderConfig = eCfg

	lvl := cfg.Level
	zl, err := cfg.Build(zap.AddCallerSkip(1))

	return zl, &lvl, err
}

func newZapTraceLogger(logFormat Format) (*zap.Logger, *zap.AtomicLevel, error) {
	cfg := zap.NewProductionConfig()
	eCfg := cfg.EncoderConfig
	eCfg.TimeKey = "timestamp"
	eCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	if logFormat == FormatConsole {
		eCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.Encoding = "console"
		eCfg.EncodeName = nameEncoder
		eCfg.EncodeLevel = traceLevelEncoder
	}

	eCfg.ConsoleSeparator = getLogSeparator()
	cfg.EncoderConfig = eCfg

	lvl := cfg.Level
	lvl.SetLevel(zapcore.DebugLevel)
	zl, err := cfg.Build(zap.AddCallerSkip(1))

	return zl, &lvl, err
}

func traceLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[TRACE]")
}

func nameEncoder(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + loggerName + "] -")
}

func newZapChildLogger(logger Logger, name string) (Logger, error) {

	impl, ok := logger.(*zapLoggerImpl)

	if ok {
		zapLogger := impl.mainLogger
		newZl := zapLogger.Named(name)

		return &zapLoggerImpl{loggerLevel: impl.loggerLevel, mainLogger: newZl}, nil
	} else {
		return nil, fmt.Errorf("unable to create child logger")
	}
}

func newZapChildLoggerWithFields(logger Logger, fields ...Field) (Logger, error) {

	impl, ok := logger.(*zapLoggerImpl)

	if ok {
		zapLogger := impl.mainLogger
		newZl := zapLogger.With(fields...)

		return &zapLoggerImpl{loggerLevel: impl.loggerLevel, mainLogger: newZl}, nil
	} else {
		return nil, fmt.Errorf("unable to create child logger")
	}
}

func zapSync(logger Logger) {
	impl, ok := logger.(*zapLoggerImpl)

	if ok {
		_ = impl.mainLogger.Sync()
	}
}

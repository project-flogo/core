package log

import (
	"fmt"

	"github.com/nareshkumarthota/core/support/log/zapconfig"
	"github.com/nareshkumarthota/core/support/log/zapcores"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var traceLogger *zap.SugaredLogger

type zapLoggerImpl struct {
	loggerLevel *zap.AtomicLevel
	mainLogger  *zap.SugaredLogger
}

func (l *zapLoggerImpl) DebugEnabled() bool {
	return l.loggerLevel.Enabled(zapcore.DebugLevel)
}

func (l *zapLoggerImpl) TraceEnabled() bool {
	return traceEnabled && l.loggerLevel.Enabled(zapcore.DebugLevel)
}

func (l *zapLoggerImpl) Trace(args ...interface{}) {
	if traceEnabled {
		traceLogger.Debug(args...)
	}
}

func (l *zapLoggerImpl) Debug(args ...interface{}) {
	l.mainLogger.Debug(args...)
}

func (l *zapLoggerImpl) Info(args ...interface{}) {
	l.mainLogger.Info(args...)
}

func (l *zapLoggerImpl) Warn(args ...interface{}) {
	l.mainLogger.Warn(args...)
}

func (l *zapLoggerImpl) Error(args ...interface{}) {
	l.mainLogger.Error(args...)
}

func (l *zapLoggerImpl) Tracef(template string, args ...interface{}) {
	if traceEnabled {
		traceLogger.Debugf(template, args...)
	}
}

func (l *zapLoggerImpl) Debugf(template string, args ...interface{}) {
	l.mainLogger.Debugf(template, args...)
}

func (l *zapLoggerImpl) Infof(template string, args ...interface{}) {
	l.mainLogger.Infof(template, args...)
}

func (l *zapLoggerImpl) Warnf(template string, args ...interface{}) {
	l.mainLogger.Warnf(template, args...)
}

func (l *zapLoggerImpl) Errorf(template string, args ...interface{}) {
	l.mainLogger.Errorf(template, args...)
}

func (l *zapLoggerImpl) Structured() StructuredLogger {
	return &zapStructuredLoggerImpl{zl: l.mainLogger.Desugar()}
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

func newZapRootLogger(name string) Logger {

	zl, lvl, _ := newZapLogger()

	// appending all available cores together
	if len(zapcores.RegisteredCores()) != 0 {
		for _, value := range zapcores.RegisteredCores() {
			zl = zl.WithOptions(
				zap.WrapCore(
					func(c zapcore.Core) zapcore.Core {
						return zapcore.NewTee(value, zl.Core())
					}))
		}
	}

	var rootLogger Logger
	if name == "" {
		rootLogger = &zapLoggerImpl{loggerLevel: lvl, mainLogger: zl.Sugar()}
	} else {
		rootLogger = &zapLoggerImpl{loggerLevel: lvl, mainLogger: zl.Named(name).Sugar()}
	}

	if traceEnabled {
		tl, _, _ := newZapTraceLogger()

		// appending all available cores together for tracing logs
		if len(zapcores.RegisteredTraceCores()) != 0 {
			for _, value := range zapcores.RegisteredTraceCores() {
				tl = tl.WithOptions(
					zap.WrapCore(
						func(c zapcore.Core) zapcore.Core {
							return zapcore.NewTee(value, tl.Core())
						}))
			}
		}

		traceLogger = tl.Sugar()
	}

	return rootLogger
}

func newZapLogger() (*zap.Logger, *zap.AtomicLevel, error) {

	zl, err := zapconfig.DefaultCfg().LogCfg().Build(zap.AddCallerSkip(1))

	return zl, zapconfig.DefaultCfg().LogLvl(), err
}

func newZapTraceLogger() (*zap.Logger, *zap.AtomicLevel, error) {

	zl, err := zapconfig.DefaultCfg().TraceLogCfg().Build(zap.AddCallerSkip(1))

	return zl, zapconfig.DefaultCfg().TraceLogLvl(), err
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
		impl.mainLogger.Sync()
	}
}

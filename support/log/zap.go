package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerImpl struct {
	loggerLevel *zap.AtomicLevel
	zapLogger   *zap.SugaredLogger
	traceLogger *zap.SugaredLogger
}

func (l *loggerImpl) DebugEnabled() bool {
	return l.loggerLevel.Enabled(zapcore.DebugLevel)
}

func (l *loggerImpl) TraceEnabled() bool {
	return false
}

func (l *loggerImpl) Trace(args ...interface{}) {
	panic("implement me")
}

func (l *loggerImpl) Debug(args ...interface{}) {
	l.zapLogger.Debug(args...)
}

func (l *loggerImpl) Info(args ...interface{}) {
	l.zapLogger.Info(args...)
}

func (l *loggerImpl) Warn(args ...interface{}) {
	l.zapLogger.Warn(args...)
}

func (l *loggerImpl) Error(args ...interface{}) {
	l.zapLogger.Error(args...)
}

func (l *loggerImpl) Tracef(template string, args ...interface{}) {
	panic("implement me")
}

func (l *loggerImpl) Debugf(template string, args ...interface{}) {
	l.zapLogger.Debugf(template, args...)
}

func (l *loggerImpl) Infof(template string, args ...interface{}) {
	l.zapLogger.Infof(template, args...)
}

func (l *loggerImpl) Warnf(template string, args ...interface{}) {
	l.zapLogger.Warnf(template, args...)
}

func (l *loggerImpl) Errorf(template string, args ...interface{}) {
	l.zapLogger.Errorf(template, args...)
}

func (l *loggerImpl) Structured() StructuredLogger {
	return &structuredLoggerImpl{zl:l.zapLogger.Desugar()}
}

type structuredLoggerImpl struct {
	lvl *zap.AtomicLevel
	zl  *zap.Logger
}

func (l *structuredLoggerImpl) Trace(msg string, fields ...Field) {
	panic("implement me")
}

func (l *structuredLoggerImpl) Debug(msg string, fields ...Field) {

	fs := make([]zap.Field, len(fields))
	for i, f := range fields {
		fs[i] = f.(zap.Field)
	}

	l.zl.Debug(msg, fs...)
}

func (l *structuredLoggerImpl) Info(msg string, fields ...Field) {
	fs := make([]zap.Field, len(fields))
	for i, f := range fields {
		fs[i] = f.(zap.Field)
	}

	l.zl.Info(msg, fs...)
}

func (l *structuredLoggerImpl) Warn(msg string, fields ...Field) {
	fs := make([]zap.Field, len(fields))
	for i, f := range fields {
		fs[i] = f.(zap.Field)
	}

	l.zl.Warn(msg, fs...)
}

func (l *structuredLoggerImpl) Error(msg string, fields ...Field) {
	fs := make([]zap.Field, len(fields))
	for i, f := range fields {
		fs[i] = f.(zap.Field)
	}

	l.zl.Error(msg, fs...)
}

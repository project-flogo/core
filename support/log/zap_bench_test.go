package log

import (
	"io"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newDiscardLogger(b testing.TB, level zapcore.Level, traced bool) Logger {
	enableTraceContextLogging(b)
	lvl := zap.NewAtomicLevelAt(level)
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(io.Discard), lvl)
	l := &zapLoggerImpl{loggerLevel: &lvl, mainLogger: zap.New(core).Sugar()}
	if traced {
		l.SetTracingContext(testTraceCtx)
	}
	return l
}

var tracedCases = []struct {
	name   string
	traced bool
}{{"NoPrefix", false}, {"Prefix", true}}

func BenchmarkDebugfSuppressed(b *testing.B) {
	for _, c := range tracedCases {
		b.Run(c.name, func(b *testing.B) {
			l := newDiscardLogger(b, zapcore.InfoLevel, c.traced)
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				l.Debugf("Step: %d", i)
			}
		})
	}
}

func BenchmarkDebugSuppressed(b *testing.B) {
	for _, c := range tracedCases {
		b.Run(c.name, func(b *testing.B) {
			l := newDiscardLogger(b, zapcore.InfoLevel, c.traced)
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				l.Debug("Scheduling task ", i)
			}
		})
	}
}

func BenchmarkDebugfGuarded(b *testing.B) {
	for _, c := range tracedCases {
		b.Run(c.name, func(b *testing.B) {
			l := newDiscardLogger(b, zapcore.InfoLevel, c.traced)
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if l.DebugEnabled() {
					l.Debugf("Step: %d", i)
				}
			}
		})
	}
}

func BenchmarkInfofEnabled(b *testing.B) {
	for _, c := range tracedCases {
		b.Run(c.name, func(b *testing.B) {
			l := newDiscardLogger(b, zapcore.InfoLevel, c.traced)
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				l.Infof("Flow Instance [%d] completed", i)
			}
		})
	}
}

func BenchmarkNewLogger(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewLogger("flogo.general.activity.log")
	}
}

func BenchmarkChildLoggerWithFields(b *testing.B) {
	parent := newDiscardLogger(b, zapcore.InfoLevel, false)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ChildLoggerWithFields(parent, FieldString("trace.id", testTraceCtx[KeyTraceID]), FieldString("span.id", testTraceCtx[KeySpanID]))
	}
}

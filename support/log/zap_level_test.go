package log

import (
	"fmt"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var testTraceCtx = map[string]string{KeyTraceID: "4bf92f3577b34da6a3ce929d0e0e4736", KeySpanID: "00f067aa0ba902b7"} // gitleaks:allow

func newObservedLogger(level zapcore.Level) (*zapLoggerImpl, *observer.ObservedLogs) {
	lvl := zap.NewAtomicLevelAt(level)
	core, logs := observer.New(lvl)
	return &zapLoggerImpl{loggerLevel: &lvl, mainLogger: zap.New(core).Sugar()}, logs
}

func enableTraceContextLogging(t testing.TB) {
	prev := traceContextLogging
	traceContextLogging = true
	t.Cleanup(func() { traceContextLogging = prev })
}

// A call below the enabled level must cost nothing, even when the logger carries a trace prefix.
func TestSuppressedLevelsDoNotBuildTracePrefix(t *testing.T) {
	enableTraceContextLogging(t)

	l, logs := newObservedLogger(zapcore.ErrorLevel)
	l.SetTracingContext(testTraceCtx)
	if l.tracePrefix() == "" {
		t.Fatal("test setup: expected a trace prefix")
	}

	calls := []struct {
		name string
		fn   func()
	}{
		{"Debugf", func() { l.Debugf("Step") }},
		{"Infof", func() { l.Infof("Step") }},
		{"Warnf", func() { l.Warnf("Step") }},
		{"Debug", func() { l.Debug() }},
		{"Info", func() { l.Info() }},
		{"Warn", func() { l.Warn() }},
	}
	for _, c := range calls {
		if allocs := testing.AllocsPerRun(200, c.fn); allocs != 0 {
			t.Errorf("%s below the enabled level: %.0f allocations per call, want 0", c.name, allocs)
		}
	}

	l.loggerLevel.SetLevel(zapcore.DPanicLevel)
	for _, c := range []struct {
		name string
		fn   func()
	}{
		{"Errorf", func() { l.Errorf("Step") }},
		{"Error", func() { l.Error() }},
	} {
		if allocs := testing.AllocsPerRun(200, c.fn); allocs != 0 {
			t.Errorf("%s below the enabled level: %.0f allocations per call, want 0", c.name, allocs)
		}
	}

	if n := logs.Len(); n != 0 {
		t.Errorf("suppressed calls wrote %d entries, want 0", n)
	}
}

// Pins the output of every leveled method, with and without a trace prefix, when the level is enabled.
func TestEnabledLevelsWritePrefixOnce(t *testing.T) {
	enableTraceContextLogging(t)

	for _, traced := range []bool{false, true} {
		l, logs := newObservedLogger(zapcore.DebugLevel)
		prefix := ""
		if traced {
			l.SetTracingContext(testTraceCtx)
			prefix = l.tracePrefix()
		}

		l.Debugf("a %d", 1)
		l.Infof("b %s", "x")
		l.Warnf("c")
		l.Errorf("d %v", true)
		l.Debug("e", 2)
		l.Info("f")
		l.Warn("g", "h")
		l.Error(3, 4)

		want := []struct {
			level zapcore.Level
			msg   string
		}{
			{zapcore.DebugLevel, prefix + "a 1"},
			{zapcore.InfoLevel, prefix + "b x"},
			{zapcore.WarnLevel, prefix + "c"},
			{zapcore.ErrorLevel, prefix + "d true"},
			{zapcore.DebugLevel, fmt.Sprint(prefix, "e", 2)},
			{zapcore.InfoLevel, fmt.Sprint(prefix, "f")},
			{zapcore.WarnLevel, fmt.Sprint(prefix, "g", "h")},
			{zapcore.ErrorLevel, fmt.Sprint(prefix, 3, 4)},
		}

		got := logs.AllUntimed()
		if len(got) != len(want) {
			t.Fatalf("traced=%t: %d entries, want %d", traced, len(got), len(want))
		}
		for i, w := range want {
			if got[i].Level != w.level || got[i].Message != w.msg {
				t.Errorf("traced=%t entry %d: got %s %q, want %s %q", traced, i, got[i].Level, got[i].Message, w.level, w.msg)
			}
			if traced && strings.Count(got[i].Message, KeyTraceID) != 1 {
				t.Errorf("traced=%t entry %d: %q must carry the trace id exactly once", traced, i, got[i].Message)
			}
		}
	}
}

// A method that tests another level than its own either drops enabled lines or builds the prefix for suppressed ones.
func TestLeveledMethodsCheckTheirOwnLevel(t *testing.T) {
	enableTraceContextLogging(t)

	methods := []struct {
		name   string
		level  zapcore.Level
		suffix string
		call   func(l *zapLoggerImpl)
	}{
		{"Debug", zapcore.DebugLevel, "", func(l *zapLoggerImpl) { l.Debug() }},
		{"Info", zapcore.InfoLevel, "", func(l *zapLoggerImpl) { l.Info() }},
		{"Warn", zapcore.WarnLevel, "", func(l *zapLoggerImpl) { l.Warn() }},
		{"Error", zapcore.ErrorLevel, "", func(l *zapLoggerImpl) { l.Error() }},
		{"Debugf", zapcore.DebugLevel, "m", func(l *zapLoggerImpl) { l.Debugf("m") }},
		{"Infof", zapcore.InfoLevel, "m", func(l *zapLoggerImpl) { l.Infof("m") }},
		{"Warnf", zapcore.WarnLevel, "m", func(l *zapLoggerImpl) { l.Warnf("m") }},
		{"Errorf", zapcore.ErrorLevel, "m", func(l *zapLoggerImpl) { l.Errorf("m") }},
	}
	for _, level := range []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel, zapcore.DPanicLevel} {
		l, logs := newObservedLogger(level)
		l.SetTracingContext(testTraceCtx)
		for _, m := range methods {
			if m.level >= level {
				m.call(l)
				want := l.tracePrefix() + m.suffix
				if got := logs.TakeAll(); len(got) != 1 || got[0].Level != m.level || got[0].Message != want {
					t.Errorf("%s at %s: entries %+v, want one %s entry %q", m.name, level, got, m.level, want)
				}
				continue
			}
			if allocs := testing.AllocsPerRun(20, func() { m.call(l) }); allocs != 0 {
				t.Errorf("%s at %s: %.0f allocations per call, want 0", m.name, level, allocs)
			}
			if n := logs.Len(); n != 0 {
				t.Errorf("%s at %s: %d entries written, want 0", m.name, level, n)
			}
		}
	}
}

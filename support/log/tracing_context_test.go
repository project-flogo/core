package log

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestWithTracingContextLeavesParentUnchanged(t *testing.T) {
	enableTraceContextLogging(t)
	parent, logs := newObservedLogger(zapcore.DebugLevel)

	child := WithTracingContext(parent, testTraceCtx)
	if child == Logger(parent) {
		t.Fatal("WithTracingContext must return a new logger")
	}
	if parent.GetTracingContext() != nil || parent.tracePrefix() != "" {
		t.Fatalf("parent was modified: context %v, prefix %q", parent.GetTracingContext(), parent.tracePrefix())
	}
	if got := child.GetTracingContext(); got[KeyTraceID] != testTraceCtx[KeyTraceID] || got[KeySpanID] != testTraceCtx[KeySpanID] {
		t.Fatalf("child context = %v, want %v", got, testTraceCtx)
	}

	child.Infof("from child")
	parent.Infof("from parent")
	entries := logs.AllUntimed()
	if len(entries) != 2 {
		t.Fatalf("%d entries, want 2", len(entries))
	}
	wantPrefix := fmt.Sprintf("[%s: %s] [%s: %s] ", KeyTraceID, testTraceCtx[KeyTraceID], KeySpanID, testTraceCtx[KeySpanID])
	if entries[0].Message != wantPrefix+"from child" {
		t.Errorf("child line = %q, want %q", entries[0].Message, wantPrefix+"from child")
	}
	if entries[1].Message != "from parent" {
		t.Errorf("parent line = %q, want no prefix", entries[1].Message)
	}
}

func TestWithTracingContextReplacesInheritedContext(t *testing.T) {
	enableTraceContextLogging(t)
	parent, logs := newObservedLogger(zapcore.DebugLevel)

	first := WithTracingContext(parent, map[string]string{KeyTraceID: "1111", KeySpanID: "aaaa"})
	second := WithTracingContext(first, map[string]string{KeyTraceID: "2222", KeySpanID: "bbbb"})
	second.Infof("x")

	msg := logs.AllUntimed()[0].Message
	if strings.Count(msg, KeyTraceID) != 1 || !strings.Contains(msg, "2222") || strings.Contains(msg, "1111") {
		t.Errorf("line = %q, want only the second context, once", msg)
	}
	if first.GetTracingContext()[KeyTraceID] != "1111" {
		t.Errorf("deriving from first changed it: %v", first.GetTracingContext())
	}
}

func TestWithTracingContextNilClearsContext(t *testing.T) {
	enableTraceContextLogging(t)
	parent, logs := newObservedLogger(zapcore.DebugLevel)
	traced := WithTracingContext(parent, testTraceCtx)

	cleared := WithTracingContext(traced, nil)
	if cleared.GetTracingContext() != nil {
		t.Errorf("context = %v, want nil", cleared.GetTracingContext())
	}
	cleared.Infof("plain")
	if msg := logs.AllUntimed()[0].Message; msg != "plain" {
		t.Errorf("line = %q, want %q", msg, "plain")
	}
	if traced.GetTracingContext() == nil {
		t.Error("deriving a cleared logger modified its parent")
	}
}

func TestWithTracingContextNilWithoutContextReturnsSameLogger(t *testing.T) {
	enableTraceContextLogging(t)
	parent, _ := newObservedLogger(zapcore.DebugLevel)

	if got := WithTracingContext(parent, nil); got != Logger(parent) {
		t.Error("a nil context on a logger without one must return the logger itself")
	}
	if allocs := testing.AllocsPerRun(200, func() { _ = WithTracingContext(parent, nil) }); allocs != 0 {
		t.Errorf("%.0f allocations per call, want 0", allocs)
	}
}

func TestWithTracingContextDisabledReturnsSameLogger(t *testing.T) {
	prev := traceContextLogging
	traceContextLogging = false
	defer func() { traceContextLogging = prev }()

	parent, _ := newObservedLogger(zapcore.DebugLevel)
	if got := WithTracingContext(parent, testTraceCtx); got != Logger(parent) {
		t.Error("with trace context logging disabled the logger itself must be returned")
	}
	if parent.GetTracingContext() != nil {
		t.Error("parent was modified")
	}
}

func TestWithTracingContextSharesLevelAndName(t *testing.T) {
	enableTraceContextLogging(t)
	parent, logs := newObservedLogger(zapcore.DebugLevel)
	parent.mainLogger = parent.mainLogger.Named("flow")

	child := WithTracingContext(parent, testTraceCtx)
	SetLogLevel(parent, WarnLevel)
	if child.DebugEnabled() {
		t.Error("the derived logger must follow the parent's level")
	}
	child.Infof("dropped")
	child.Warnf("kept")

	entries := logs.AllUntimed()
	if len(entries) != 1 || entries[0].LoggerName != "flow" {
		t.Fatalf("entries = %+v, want one WARN entry from logger %q", entries, "flow")
	}
}

func TestWithTracingContextOnChildLoggers(t *testing.T) {
	enableTraceContextLogging(t)
	parent, logs := newObservedLogger(zapcore.DebugLevel)

	withFields := ChildLoggerWithFields(parent, FieldString("flow.id", "f1"))
	traced := WithTracingContext(withFields, testTraceCtx)
	traced.Infof("x")

	e := logs.AllUntimed()[0]
	if !strings.Contains(e.Message, testTraceCtx[KeyTraceID]) {
		t.Errorf("line = %q, want the trace prefix", e.Message)
	}
	if len(e.Context) != 1 || e.Context[0].Key != "flow.id" {
		t.Errorf("fields = %+v, want the parent's single field", e.Context)
	}
	if ChildLogger(traced, "sub").GetTracingContext() != nil {
		t.Error("named children start without a tracing context, as before")
	}
}

type recordingLogger struct {
	Logger
	mu  sync.Mutex
	set []map[string]string
}

func (r *recordingLogger) SetTracingContext(tc map[string]string) {
	r.mu.Lock()
	r.set = append(r.set, tc)
	r.mu.Unlock()
}

func TestWithTracingContextOtherLoggerImplementations(t *testing.T) {
	enableTraceContextLogging(t)
	r := &recordingLogger{}

	if got := WithTracingContext(r, testTraceCtx); got != Logger(r) {
		t.Error("a logger that cannot derive a copy must be returned as is")
	}
	if len(r.set) != 1 || r.set[0][KeyTraceID] != testTraceCtx[KeyTraceID] {
		t.Errorf("SetTracingContext calls = %v, want one with the context", r.set)
	}
	if WithTracingContext(nil, testTraceCtx) != nil {
		t.Error("a nil logger must stay nil")
	}
}

func TestWithTracingContextConcurrentDerivations(t *testing.T) {
	enableTraceContextLogging(t)
	parent, logs := newObservedLogger(zapcore.InfoLevel)

	const goroutines, iterations = 16, 200
	var wg sync.WaitGroup
	errs := make(chan string, goroutines)
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			own := fmt.Sprintf("%032x", g+1)
			for i := 0; i < iterations; i++ {
				l := WithTracingContext(parent, map[string]string{KeyTraceID: own, KeySpanID: fmt.Sprintf("%016x", i)})
				l.Infof("owner=%s", own)
				if got := l.GetTracingContext()[KeyTraceID]; got != own {
					errs <- fmt.Sprintf("goroutine %d read trace id %s", g, got)
					return
				}
			}
		}(g)
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		t.Error(e)
	}
	if parent.GetTracingContext() != nil {
		t.Errorf("parent context = %v, want nil", parent.GetTracingContext())
	}

	entries := logs.AllUntimed()
	if len(entries) != goroutines*iterations {
		t.Fatalf("%d entries, want %d", len(entries), goroutines*iterations)
	}
	for _, e := range entries {
		owner := e.Message[strings.Index(e.Message, "owner=")+len("owner="):]
		if !strings.HasPrefix(e.Message, fmt.Sprintf("[%s: %s]", KeyTraceID, owner)) {
			t.Fatalf("line %q carries another goroutine's trace id", e.Message)
		}
	}
}

func tracePrefixFor(traceID, spanID string) string {
	return fmt.Sprintf("[%s: %s] [%s: %s] ", KeyTraceID, traceID, KeySpanID, spanID)
}

func TestWithTracingContextReplacesWithoutMergingOrNesting(t *testing.T) {
	enableTraceContextLogging(t)
	parent, logs := newObservedLogger(zapcore.DebugLevel)
	parent.SetTracingContext(map[string]string{KeyTraceID: "root", KeySpanID: "root", "baggage": "root-only"}) // gitleaks:allow

	var l Logger = parent
	for i := 1; i <= 50; i++ {
		tc := map[string]string{KeyTraceID: fmt.Sprintf("%032x", i), KeySpanID: fmt.Sprintf("%016x", i)}
		l = WithTracingContext(l, tc)
		if got := l.GetTracingContext(); len(tc) != 2 || !reflect.DeepEqual(got, tc) {
			t.Fatalf("derivation %d: context %v, caller's map %v; want exactly the caller's two keys", i, got, tc)
		}
		l.Infof("x")
		want := tracePrefixFor(tc[KeyTraceID], tc[KeySpanID]) + "x"
		if e := logs.TakeAll(); len(e) != 1 || e[0].Message != want || len(e[0].Context) != 0 {
			t.Fatalf("derivation %d: entries %+v, want one %q without fields", i, e, want)
		}
	}
	if got := parent.GetTracingContext(); got[KeyTraceID] != "root" || len(got) != 3 {
		t.Errorf("parent context = %v, want the root context", got)
	}
}

// A derived logger must not alias its parent, even when it carries the parent's trace id or the parent's own map.
func TestDerivedLoggerIsIndependentOfItsParent(t *testing.T) {
	enableTraceContextLogging(t)
	parent, logs := newObservedLogger(zapcore.DebugLevel)
	parentCtx := map[string]string{KeyTraceID: "t1", KeySpanID: "s1"}
	parent.SetTracingContext(parentCtx)

	sameTrace := WithTracingContext(parent, map[string]string{KeyTraceID: "t1", KeySpanID: "s2"})
	sameMap := WithTracingContext(parent, parentCtx)
	parent.SetTracingContext(map[string]string{KeyTraceID: "t9", KeySpanID: "s9"})
	sameTrace.SetTracingContext(map[string]string{KeyTraceID: "t3", KeySpanID: "s3"})

	sameTrace.Infof("a")
	sameMap.Infof("b")
	parent.Infof("c")
	want := []string{tracePrefixFor("t3", "s3") + "a", tracePrefixFor("t1", "s1") + "b", tracePrefixFor("t9", "s9") + "c"}
	entries := logs.AllUntimed()
	if len(entries) != len(want) {
		t.Fatalf("%d entries, want %d", len(entries), len(want))
	}
	for i, w := range want {
		if entries[i].Message != w {
			t.Errorf("line %d = %q, want %q", i, entries[i].Message, w)
		}
	}
}

type derivingLogger struct {
	recordingLogger
	derived Logger
	with    []map[string]string
}

func (d *derivingLogger) WithTracingContext(tc map[string]string) Logger {
	d.with = append(d.with, tc)
	return d.derived
}

func TestWithTracingContextUsesTheLoggersOwnDerivation(t *testing.T) {
	derived := &recordingLogger{}
	d := &derivingLogger{derived: derived}

	if got := WithTracingContext(d, testTraceCtx); got != Logger(derived) {
		t.Error("the logger's own derivation must be returned")
	}
	if len(d.with) != 1 || d.with[0][KeyTraceID] != testTraceCtx[KeyTraceID] || len(d.set) != 0 {
		t.Errorf("derivations %v, SetTracingContext calls %v; want one derivation and no SetTracingContext", d.with, d.set)
	}
}

func TestWithTracingContextFallbackPassesNilContext(t *testing.T) {
	r := &recordingLogger{}

	if got := WithTracingContext(r, nil); got != Logger(r) {
		t.Error("a logger that cannot derive a copy must be returned as is")
	}
	if len(r.set) != 1 || r.set[0] != nil {
		t.Errorf("SetTracingContext calls = %v, want one with nil", r.set)
	}
}

// Callers not yet moved to derivation still call SetTracingContext on shared loggers; derivations must stay isolated.
func TestWithTracingContextConcurrentWithSetTracingContext(t *testing.T) {
	enableTraceContextLogging(t)
	shared, logs := newObservedLogger(zapcore.InfoLevel)
	writerCtx := map[string]string{KeyTraceID: strings.Repeat("f", 32), KeySpanID: strings.Repeat("f", 16)}

	const derivers, iterations = 8, 200
	done := make(chan struct{})
	var warm, background sync.WaitGroup
	warm.Add(2)
	background.Add(2)
	go func() {
		defer background.Done()
		defer shared.SetTracingContext(nil)
		for n := 0; ; n++ {
			shared.SetTracingContext(writerCtx)
			if n%64 == 0 {
				shared.Infof("shared")
			}
			shared.SetTracingContext(nil)
			if n == 0 {
				warm.Done()
			}
			select {
			case <-done:
				return
			default:
			}
		}
	}()
	go func() {
		defer background.Done()
		for n := 0; ; n++ {
			SetLogLevel(shared, DebugLevel)
			SetLogLevel(shared, InfoLevel)
			if n == 0 {
				warm.Done()
			}
			select {
			case <-done:
				return
			default:
			}
		}
	}()
	// Derivations start only once both background goroutines are running, so they overlap.
	warm.Wait()

	var wg sync.WaitGroup
	for g := 0; g < derivers; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			traceID := fmt.Sprintf("%032x", g+1)
			if g%2 == 0 {
				traceID = writerCtx[KeyTraceID]
			}
			for i := 0; i < iterations; i++ {
				tc := map[string]string{KeyTraceID: traceID, KeySpanID: fmt.Sprintf("%08x%08x", g, i)}
				l := WithTracingContext(shared, tc)
				if l == Logger(shared) {
					t.Errorf("goroutine %d: derivation returned the shared logger", g)
					return
				}
				l.Infof("want=%s/%s", tc[KeyTraceID], tc[KeySpanID])
				if got := l.GetTracingContext(); got[KeyTraceID] != tc[KeyTraceID] || got[KeySpanID] != tc[KeySpanID] {
					t.Errorf("goroutine %d: context %v, want %v", g, got, tc)
					return
				}
			}
		}(g)
	}
	wg.Wait()
	close(done)
	background.Wait()

	derivedLines, sharedLines := 0, 0
	for _, e := range logs.AllUntimed() {
		_, ids, derived := strings.Cut(e.Message, "want=")
		if !derived {
			if want := tracePrefixFor(writerCtx[KeyTraceID], writerCtx[KeySpanID]) + "shared"; e.Message != want {
				t.Fatalf("shared logger line %q, want %q", e.Message, want)
			}
			sharedLines++
			continue
		}
		traceID, spanID, _ := strings.Cut(ids, "/")
		if want := tracePrefixFor(traceID, spanID) + "want=" + ids; e.Message != want {
			t.Fatalf("derived logger line %q, want %q", e.Message, want)
		}
		derivedLines++
	}
	if derivedLines != derivers*iterations || sharedLines == 0 {
		t.Errorf("%d derived lines (want %d), %d shared lines (want at least 1)", derivedLines, derivers*iterations, sharedLines)
	}
	if got := shared.GetTracingContext(); got != nil {
		t.Errorf("shared context = %v, want nil", got)
	}
}

func BenchmarkWithTracingContext(b *testing.B) {
	enableTraceContextLogging(b)
	parent, _ := newObservedLogger(zapcore.InfoLevel)
	var l Logger = parent
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = WithTracingContext(l, testTraceCtx)
	}
}

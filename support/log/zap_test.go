package log

import (
	"fmt"
	"sync"
	"testing"
)

// TestSetTracingContext pins the observable behaviour of the trace prefix so the
// atomic publication introduced for FLOGO-19401 cannot silently change it.
func TestSetTracingContext(t *testing.T) {
	prev := traceContextLogging
	traceContextLogging = true
	defer func() { traceContextLogging = prev }()

	logger := NewLogger("flogo.test.tracingcontext").(*zapLoggerImpl)

	if got := logger.tracePrefix(); got != "" {
		t.Fatalf("new logger: tracePrefix = %q, want empty", got)
	}
	if got := logger.GetTracingContext(); got != nil {
		t.Fatalf("new logger: GetTracingContext = %v, want nil", got)
	}

	tc := map[string]string{KeyTraceID: "4bf92f3577b34da6a3ce929d0e0e4736", KeySpanID: "00f067aa0ba902b7"}
	logger.SetTracingContext(tc)

	want := fmt.Sprintf("[%s: %s] [%s: %s] ", KeyTraceID, tc[KeyTraceID], KeySpanID, tc[KeySpanID])
	if got := logger.tracePrefix(); got != want {
		t.Errorf("after set: tracePrefix = %q, want %q", got, want)
	}
	if got := logger.GetTracingContext(); got[KeyTraceID] != tc[KeyTraceID] {
		t.Errorf("after set: GetTracingContext[%s] = %q, want %q", KeyTraceID, got[KeyTraceID], tc[KeyTraceID])
	}

	logger.SetTracingContext(nil)
	if got := logger.tracePrefix(); got != "" {
		t.Errorf("after clear: tracePrefix = %q, want empty", got)
	}
	if got := logger.GetTracingContext(); got != nil {
		t.Errorf("after clear: GetTracingContext = %v, want nil", got)
	}

	// A context carrying no trace id must not produce a prefix.
	logger.SetTracingContext(map[string]string{KeySpanID: "00f067aa0ba902b7"})
	if got := logger.tracePrefix(); got != "" {
		t.Errorf("no trace id: tracePrefix = %q, want empty", got)
	}
}

// TestSetTracingContextDisabled covers FLOGO_LOG_TRACE_CTX_ENABLED=false, which
// must leave the logger untouched.
func TestSetTracingContextDisabled(t *testing.T) {
	prev := traceContextLogging
	traceContextLogging = false
	defer func() { traceContextLogging = prev }()

	logger := NewLogger("flogo.test.tracingcontext.disabled").(*zapLoggerImpl)
	logger.SetTracingContext(map[string]string{KeyTraceID: "abc", KeySpanID: "def"})

	if got := logger.tracePrefix(); got != "" {
		t.Errorf("tracing context logging disabled: tracePrefix = %q, want empty", got)
	}
	if got := logger.GetTracingContext(); got != nil {
		t.Errorf("tracing context logging disabled: GetTracingContext = %v, want nil", got)
	}
}

// TestSetTracingContextConcurrentWithLogging is the regression test for FLOGO-19401.
//
// A single logger is shared by every concurrently executing flow instance and
// activity - flow/action.go does `instLogger := logger` and only replaces it with a
// per-instance child when FLOGO_LOG_CTX=true, which is off by default. So
// SetTracingContext (flow/action.go:355,392,397 and flow/instance/taskinst.go:349)
// runs concurrently with log calls on that same logger.
//
// Before the fix the trace prefix was a plain string field. A Go string is a
// two-word {data, len} value and assigning one is not atomic, so a racing reader
// could observe {data: nil, len: N}: the `!= ""` guard passed because len was
// non-zero, and the prefix concatenation then copied N bytes from address 0 and
// segfaulted, taking the whole app down.
//
// Run with -race to catch the underlying data race deterministically. Without
// -race this still exercises the crash, but only probabilistically.
func TestSetTracingContextConcurrentWithLogging(t *testing.T) {
	prev := traceContextLogging
	traceContextLogging = true
	defer func() { traceContextLogging = prev }()

	logger := NewLogger("flogo.test.tracingcontext.race")
	// Keep the test output quiet. The prefix is an argument to Debugf/Infof, so it is
	// still built on every call regardless of the level - that is exactly why a
	// production app running at ERROR was still crashing inside Debugf.
	SetLogLevel(logger, ErrorLevel)

	traceCtx := map[string]string{KeyTraceID: "4bf92f3577b34da6a3ce929d0e0e4736", KeySpanID: "00f067aa0ba902b7"}

	const (
		writers    = 2
		readers    = 8
		iterations = 5000
	)

	var wg sync.WaitGroup

	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				logger.SetTracingContext(traceCtx)
				logger.SetTracingContext(nil)
			}
		}()
	}

	for r := 0; r < readers; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				logger.Debugf("Step: %d", i)
				logger.Infof("Flow Instance [%d] completed", i)
				logger.Debug("Scheduling task ", i)
				logger.Info("Flow Instance ", i)
			}
		}()
	}

	wg.Wait()
}

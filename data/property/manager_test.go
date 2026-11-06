package property

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	raceTestProps    = 32
	raceTestReaders  = 4
	raceTestDuration = 300 * time.Millisecond
	raceTestMaxIter  = 5000
)

// newRaceTestManager returns a manager with raceTestProps properties, all set to "v0". Even-numbered
// properties are overridable, odd-numbered ones are read-only.
func newRaceTestManager() *Manager {
	props := make(map[string]Property, raceTestProps)
	for i := 0; i < raceTestProps; i++ {
		props["p"+strconv.Itoa(i)] = Property{Value: "v0", Overridable: i%2 == 0}
	}
	return NewManager(props)
}

// setAll returns a post-processor that sets every property to the same value, so that a reader can tell
// whether it saw a half-applied refresh.
func setAll(value string) PostProcessor {
	return func(properties map[string]interface{}) error {
		for k := range properties {
			properties[k] = value
		}
		return nil
	}
}

// checkSnapshot reports an error unless snapshot holds every property and all of them carry the value of
// one and the same refresh.
func checkSnapshot(snapshot map[string]interface{}) error {
	if len(snapshot) != raceTestProps {
		return fmt.Errorf("GetProperties returned %d properties, want %d", len(snapshot), raceTestProps)
	}
	first := snapshot["p0"]
	for k, v := range snapshot {
		if v != first {
			return fmt.Errorf("GetProperties returned a half-applied refresh: %s=%v but p0=%v", k, v, first)
		}
	}
	return nil
}

// useDefaultManager makes m the default manager until the test ends.
func useDefaultManager(t *testing.T, m *Manager) {
	previous := DefaultManager()
	SetDefaultManager(m)
	t.Cleanup(func() { SetDefaultManager(previous) })
}

// TestManager_FinalizeWhileReading reproduces FLOGO-19663: a config refresh re-runs Finalize while flows read
// properties. Without the fix, the race detector reports a data race, and without -race the Go runtime aborts
// the process with "fatal error: concurrent map read and map write".
func TestManager_FinalizeWhileReading(t *testing.T) {
	m := newRaceTestManager()

	stop := make(chan struct{})
	errs := make(chan error, raceTestReaders)
	var wg sync.WaitGroup
	for r := 0; r < raceTestReaders; r++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			name := "p" + strconv.Itoa(idx)
			wantOverridable := idx%2 == 0
			for {
				select {
				case <-stop:
					return
				default:
				}
				if v, ok := m.GetProperty(name); !ok || v == nil {
					errs <- fmt.Errorf("GetProperty(%q) = %v, %v during a refresh", name, v, ok)
					return
				}
				if got := m.IsOverridable(name); got != wantOverridable {
					errs <- fmt.Errorf("IsOverridable(%q) = %v during a refresh, want %v", name, got, wantOverridable)
					return
				}
				if err := checkSnapshot(m.GetProperties()); err != nil {
					errs <- err
					return
				}
			}
		}(r)
	}

	deadline := time.Now().Add(raceTestDuration)
	last := 0
	for i := 1; i <= raceTestMaxIter && time.Now().Before(deadline); i++ {
		if err := m.Finalize(setAll("v" + strconv.Itoa(i))); err != nil {
			t.Errorf("Finalize returned %v", err)
		}
		last = i
	}
	close(stop)
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}

	want := "v" + strconv.Itoa(last)
	for i := 0; i < raceTestProps; i++ {
		name := "p" + strconv.Itoa(i)
		if v, ok := m.GetProperty(name); !ok || v != want {
			t.Errorf("after %d refreshes GetProperty(%q) = %v, %v; want %q, true", last, name, v, ok, want)
		}
	}
}

// overlapCounter counts calls that run while another one is still running; it needs no -race.
type overlapCounter struct {
	running, overlaps int32
}

func (c *overlapCounter) enter() {
	if atomic.AddInt32(&c.running, 1) > 1 {
		atomic.AddInt32(&c.overlaps, 1)
	}
	runtime.Gosched()
}

func (c *overlapCounter) leave() { atomic.AddInt32(&c.running, -1) }

// TestManager_ConcurrentFinalize covers two config refreshes that overlap: the refresh endpoint serves each
// request on its own goroutine. Without the fix, the two Finalize calls write the same map at once
// ("fatal error: concurrent map writes").
func TestManager_ConcurrentFinalize(t *testing.T) {
	m := newRaceTestManager()

	var gen int64
	var processors overlapCounter
	refresh := func(properties map[string]interface{}) error {
		processors.enter()
		defer processors.leave()
		v := "v" + strconv.FormatInt(atomic.AddInt64(&gen, 1), 10)
		for k := range properties {
			properties[k] = v
		}
		return nil
	}

	var wg sync.WaitGroup
	for w := 0; w < 2; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			deadline := time.Now().Add(raceTestDuration)
			for i := 0; i < raceTestMaxIter && time.Now().Before(deadline); i++ {
				_ = m.Finalize(refresh)
			}
		}()
	}
	wg.Wait()

	if err := checkSnapshot(m.GetProperties()); err != nil {
		t.Error(err)
	}
	assert.Zero(t, atomic.LoadInt32(&processors.overlaps), "processors of two Finalize calls ran at the same time")
}

// TestManager_FinalizeNoLostUpdate checks that Finalize calls run one at a time, each on the values the previous one
// published: concurrent calls that each increment a counter must not lose an increment. Works without -race.
func TestManager_FinalizeNoLostUpdate(t *testing.T) {
	m := NewManager(map[string]Property{"n": {Value: 0, Overridable: true}})
	const workers, calls = 4, 200

	var processors overlapCounter
	increment := func(properties map[string]interface{}) error {
		processors.enter()
		defer processors.leave()
		n := properties["n"].(int)
		runtime.Gosched()
		properties["n"] = n + 1
		return nil
	}

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < calls; i++ {
				_ = m.Finalize(increment)
			}
		}()
	}
	wg.Wait()

	v, _ := m.GetProperty("n")
	assert.Equal(t, workers*calls, v, "concurrent Finalize calls lost updates")
	assert.Zero(t, atomic.LoadInt32(&processors.overlaps), "processors of two Finalize calls ran at the same time")
}

// TestManager_FinalizeReleasesLockOnPanic checks that a processor that panics, an external resolver for example,
// publishes nothing and does not leave Finalize's lock held: the next Finalize must still run.
func TestManager_FinalizeReleasesLockOnPanic(t *testing.T) {
	m := NewManager(map[string]Property{"p1": {Value: "a", Overridable: true}})
	assert.Panics(t, func() {
		_ = m.Finalize(func(properties map[string]interface{}) error {
			properties["p1"] = "half-done"
			panic("processor failed")
		})
	})
	v, _ := m.GetProperty("p1")
	assert.Equal(t, "a", v, "a Finalize that panicked must not publish anything")

	done := make(chan error, 1)
	go func() {
		done <- m.Finalize(func(properties map[string]interface{}) error {
			properties["p1"] = "b"
			return nil
		})
	}()
	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("a panicking processor left Finalize's lock held")
	}
	v, _ = m.GetProperty("p1")
	assert.Equal(t, "b", v)
}

// TestManager_FinalizeProcessorReadsManager checks that a processor may read the Manager it runs in, as
// ExternalResolverProcessor does through DefaultManager().IsOverridable, and that it sees the values from before
// the Finalize.
func TestManager_FinalizeProcessorReadsManager(t *testing.T) {
	m := NewManager(map[string]Property{"p1": {Value: "a", Overridable: true}})
	useDefaultManager(t, m)

	var seen []interface{}
	processor := func(properties map[string]interface{}) error {
		v, _ := m.GetProperty("p1")
		seen = append(seen, v, m.GetProperties()["p1"], m.IsOverridable("p1"), DefaultManager().IsOverridable("p1"))
		properties["p1"] = "b"
		return nil
	}

	done := make(chan error, 1)
	go func() { done <- m.Finalize(processor, processor) }()
	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Finalize deadlocked when a processor read the Manager")
	}

	assert.Equal(t, []interface{}{"a", "a", true, true, "a", "a", true, true}, seen)
	v, _ := m.GetProperty("p1")
	assert.Equal(t, "b", v)
}

func TestManager_Finalize(t *testing.T) {
	set := func(key string, value interface{}) PostProcessor {
		return func(properties map[string]interface{}) error {
			properties[key] = value
			return nil
		}
	}
	initial := func() map[string]Property {
		return map[string]Property{"p1": {Value: "a"}, "p2": {Value: "b", Overridable: true}}
	}

	tests := []struct {
		name       string
		processors []PostProcessor
		want       map[string]interface{}
	}{
		{"updates existing", []PostProcessor{set("p1", "x")}, map[string]interface{}{"p1": "x", "p2": "b"}},
		{"ignores added keys", []PostProcessor{set("new", "n")}, map[string]interface{}{"p1": "a", "p2": "b"}},
		{"removed key keeps value", []PostProcessor{func(properties map[string]interface{}) error {
			delete(properties, "p1")
			return nil
		}}, map[string]interface{}{"p1": "a", "p2": "b"}},
		{"explicit nil", []PostProcessor{set("p1", nil)}, map[string]interface{}{"p1": nil, "p2": "b"}},
		{"runs processors in order", []PostProcessor{set("p1", "b1"), func(properties map[string]interface{}) error {
			if properties["p1"] == "b1" {
				properties["p1"] = "c"
			}
			return nil
		}}, map[string]interface{}{"p1": "c", "p2": "b"}},
		{"ignores processor errors", []PostProcessor{func(properties map[string]interface{}) error {
			properties["p1"] = "x"
			return errors.New("processor failed")
		}, set("p2", "y")}, map[string]interface{}{"p1": "x", "p2": "y"}},
		{"no processors", nil, map[string]interface{}{"p1": "a", "p2": "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(initial())
			published := m.snapshot()

			assert.NoError(t, m.Finalize(tt.processors...))

			assert.Equal(t, tt.want, m.GetProperties())
			for k, want := range tt.want {
				v, ok := m.GetProperty(k)
				assert.True(t, ok, k)
				assert.Equal(t, want, v, k)
			}
			_, ok := m.GetProperty("new")
			assert.False(t, ok, "a key added by a processor must not appear")
			assert.False(t, m.IsOverridable("p1"))
			assert.True(t, m.IsOverridable("p2"))
			assert.True(t, m.IsOverridable("unknown"))
			assert.Equal(t, initial(), published, "Finalize modified the snapshot it replaced")
		})
	}
}

// fakeResolver resolves every property in values; it is read only by ExternalResolverProcessor during Finalize.
type fakeResolver struct {
	values map[string]interface{}
}

func (r *fakeResolver) Name() string { return "FLOGO-19663-test" }

func (r *fakeResolver) LookupValue(key string) (interface{}, bool) {
	v, ok := r.values[key]
	return v, ok
}

func TestManager_FinalizeWithExternalResolverProcessor(t *testing.T) {
	previous := EnabledResolvers
	EnabledResolvers = []ExternalResolver{&fakeResolver{values: map[string]interface{}{"ro": "z", "rw": "z", "num": "42"}}}
	t.Cleanup(func() { EnabledResolvers = previous })

	m := NewManager(map[string]Property{
		"ro":      {Value: "a"},
		"rw":      {Value: "b", Overridable: true},
		"num":     {Value: 1, Overridable: true},
		"missing": {Value: "c", Overridable: true},
	})
	useDefaultManager(t, m)

	assert.NoError(t, m.Finalize(ExternalResolverProcessor))
	// ExternalResolverProcessor coerces "42" to the type of the old value
	assert.Equal(t, map[string]interface{}{"ro": "a", "rw": "z", "num": int64(42), "missing": "c"}, m.GetProperties())
}

func TestManager_GetPropertiesReturnsCopy(t *testing.T) {
	m := NewManager(map[string]Property{"p1": {Value: "a", Overridable: true}})

	got := m.GetProperties()
	got["p1"] = "x"
	got["new"] = "n"

	assert.Equal(t, map[string]interface{}{"p1": "a"}, m.GetProperties())
	v, _ := m.GetProperty("p1")
	assert.Equal(t, "a", v)
}

func TestNewManager_CopiesInput(t *testing.T) {
	in := map[string]Property{"p1": {Value: "a", Overridable: true}}
	m := NewManager(in)

	in["p1"] = Property{Value: "x"}
	in["p2"] = Property{Value: "y"}

	assert.Equal(t, map[string]interface{}{"p1": "a"}, m.GetProperties())
	assert.True(t, m.IsOverridable("p1"))
}

// TestManager_Empty covers NewManager(nil) and a zero Manager, which stays usable as before.
func TestManager_Empty(t *testing.T) {
	for name, m := range map[string]*Manager{"NewManager(nil)": NewManager(nil), "zero value": {}} {
		t.Run(name, func(t *testing.T) {
			v, ok := m.GetProperty("p1")
			assert.Nil(t, v)
			assert.False(t, ok)
			assert.True(t, m.IsOverridable("p1"))
			assert.NotNil(t, m.GetProperties())
			assert.Empty(t, m.GetProperties())

			var received map[string]interface{}
			assert.NoError(t, m.Finalize(func(properties map[string]interface{}) error {
				received = properties
				properties["p1"] = "x"
				return nil
			}))
			assert.NotNil(t, received)
			assert.Empty(t, m.GetProperties())
		})
	}
}

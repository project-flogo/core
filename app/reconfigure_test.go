package app

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/property"
	"github.com/project-flogo/core/engine/secret"
	"github.com/project-flogo/core/support/service"
	"github.com/stretchr/testify/assert"
)

const (
	refreshDuration = 300 * time.Millisecond
	refreshMaxIter  = 500
)

// newReconfigurableApp creates and starts an app with refresh enabled. When the test ends it stops the app and
// restores the default property manager that New replaced.
func newReconfigurableApp(t *testing.T, props []*data.Attribute, options ...Option) *App {
	t.Setenv(property.EnvAppPropertyReconfigure, "true")
	previous := property.DefaultManager()
	t.Cleanup(func() { property.SetDefaultManager(previous) })

	a, err := New(&Config{Name: "FLOGO-19663", Version: "1.0.0", Properties: props}, nil, options...)
	if err != nil {
		t.Fatalf("New returned %v", err)
	}
	if err = a.Start(); err != nil {
		t.Fatalf("Start returned %v", err)
	}
	t.Cleanup(func() { _ = a.Stop() })
	return a
}

func stringProps(prefix string, n int, value string) []*data.Attribute {
	props := make([]*data.Attribute, 0, n)
	for i := 0; i < n; i++ {
		props = append(props, data.NewAttribute(prefix+strconv.Itoa(i), data.TypeString, value))
	}
	return props
}

// readWhileRefreshing runs readers against a until refresh has run for refreshDuration (or refreshMaxIter times),
// then reports every reader error. check validates one GetProperties result.
func readWhileRefreshing(t *testing.T, a *App, names []string, check func(map[string]interface{}) error, refresh func() error) {
	resolver := &property.Resolver{}
	stop := make(chan struct{})
	errs := make(chan error, len(names))
	var wg sync.WaitGroup
	for _, name := range names {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				if v, ok := a.GetProperty(name); !ok || v == nil {
					errs <- fmt.Errorf("App.GetProperty(%q) = %v, %v during a refresh", name, v, ok)
					return
				}
				if v, err := resolver.Resolve(nil, name, ""); err != nil || v == nil {
					errs <- fmt.Errorf("$property[%s] resolved to %v, %v during a refresh", name, v, err)
					return
				}
				property.DefaultManager().IsOverridable(name)
				if err := check(property.DefaultManager().GetProperties()); err != nil {
					errs <- err
					return
				}
			}
		}(name)
	}

	deadline := time.Now().Add(refreshDuration)
	for i := 0; i < refreshMaxIter && time.Now().Before(deadline); i++ {
		if err := refresh(); err != nil {
			t.Errorf("Reconfigure returned %v", err)
			break
		}
	}
	close(stop)
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
}

// allEqual checks that every value with the prefix came from one and the same refresh.
func allEqual(prefix string, want int) func(map[string]interface{}) error {
	return func(snapshot map[string]interface{}) error {
		var first interface{}
		count := 0
		for k, v := range snapshot {
			if !strings.HasPrefix(k, prefix) {
				continue
			}
			if count == 0 {
				first = v
			} else if v != first {
				return fmt.Errorf("GetProperties returned a half-applied refresh: %s=%v but another one is %v", k, v, first)
			}
			count++
		}
		if count != want {
			return fmt.Errorf("GetProperties returned %d properties named %s*, want %d", count, prefix, want)
		}
		return nil
	}
}

// TestApp_ReconfigureWhileReadingProperties reproduces FLOGO-19663 through App.Reconfigure, the method that a host's
// config-refresh endpoint calls. Flows keep resolving $property[...] while the refresh re-runs FinalizeProperties.
// Without the fix, the race detector reports a data race, and without -race the Go runtime aborts the process with
// "fatal error: concurrent map read and map write".
func TestApp_ReconfigureWhileReadingProperties(t *testing.T) {
	// refresh stands in for an external resolver: each run gives every property a new value. Only the test
	// goroutine runs it (through New and Reconfigure), so gen needs no synchronisation.
	gen := 0
	refresh := func(properties map[string]interface{}) error {
		gen++
		v := "v" + strconv.Itoa(gen)
		for k := range properties {
			properties[k] = v
		}
		return nil
	}
	a := newReconfigurableApp(t, stringProps("p", 16, "v0"), ContinueOnError, FinalizeProperties(refresh))

	readWhileRefreshing(t, a, []string{"p0", "p1", "p2", "p3"}, allEqual("p", 16), a.Reconfigure)

	want := "v" + strconv.Itoa(gen)
	for i := 0; i < 16; i++ {
		name := "p" + strconv.Itoa(i)
		if v, ok := a.GetProperty(name); !ok || v != want {
			t.Errorf("after the last refresh App.GetProperty(%q) = %v, %v; want %q, true", name, v, ok, want)
		}
	}
}

// genResolver gives every property it knows the value "v<gen>"; the test raises gen before each refresh.
type genResolver struct {
	gen int64
}

func (r *genResolver) Name() string { return "FLOGO-19663-test" }

func (r *genResolver) LookupValue(key string) (interface{}, bool) {
	if !strings.HasPrefix(key, "p") && !strings.HasPrefix(key, "r") {
		return nil, false
	}
	return "v" + strconv.FormatInt(atomic.LoadInt64(&r.gen), 10), true
}

// TestApp_ReconfigureWithEngineProcessors runs the processor chain that engine.New installs (external resolvers,
// then secrets) under concurrent reads: overridable properties change together, read-only ones and a decoded
// secret never change.
func TestApp_ReconfigureWithEngineProcessors(t *testing.T) {
	resolver := &genResolver{}
	previous := property.EnabledResolvers
	property.EnabledResolvers = []property.ExternalResolver{resolver}
	t.Cleanup(func() { property.EnabledResolvers = previous })

	encoded, err := secret.GetSecretValueHandler().EncodeValue("plain")
	if err != nil {
		t.Fatalf("EncodeValue returned %v", err)
	}
	props := stringProps("p", 8, "v0")
	for i := 0; i < 4; i++ {
		props = append(props, data.NewAttributeWithOverridable("r"+strconv.Itoa(i), data.TypeString, false, "read-only"))
	}
	props = append(props, data.NewAttributeWithOverridable("s0", data.TypeString, false, "SECRET:"+encoded))

	a := newReconfigurableApp(t, props, FinalizeProperties(property.ExternalResolverProcessor, secret.PropertyProcessor))

	check := func(snapshot map[string]interface{}) error {
		if err := allEqual("p", 8)(snapshot); err != nil {
			return err
		}
		for i := 0; i < 4; i++ {
			if v := snapshot["r"+strconv.Itoa(i)]; v != "read-only" {
				return fmt.Errorf("read-only property r%d = %v after a refresh", i, v)
			}
		}
		if v := snapshot["s0"]; v != "plain" {
			return fmt.Errorf("secret property s0 = %v, want the decoded value", v)
		}
		return nil
	}
	readWhileRefreshing(t, a, []string{"p0", "p1", "r0", "s0"}, check, func() error {
		atomic.AddInt64(&resolver.gen, 1)
		return a.Reconfigure()
	})

	want := "v" + strconv.FormatInt(atomic.LoadInt64(&resolver.gen), 10)
	v, _ := a.GetProperty("p0")
	assert.Equal(t, want, v)
	v, _ = a.GetProperty("r0")
	assert.Equal(t, "read-only", v)
}

// TestApp_ConcurrentReconfigure covers two config refreshes that overlap, for example a client that retries while the
// first refresh still runs: the endpoint serves each request on its own goroutine. The app gets the options that
// engine.New passes (ContinueOnError, FinalizeProperties, EngineSettings), and Reconfigure re-runs all of them.
func TestApp_ConcurrentReconfigure(t *testing.T) {
	var gen int64
	refresh := func(properties map[string]interface{}) error {
		v := "v" + strconv.FormatInt(atomic.AddInt64(&gen, 1), 10)
		for k := range properties {
			properties[k] = v
		}
		return nil
	}
	// overlapProbe is one more option, so every Reconfigure runs it: it counts calls that overlap, without -race.
	var running, overlaps int32
	overlapProbe := func(*App) error {
		if atomic.AddInt32(&running, 1) > 1 {
			atomic.AddInt32(&overlaps, 1)
		}
		runtime.Gosched()
		atomic.AddInt32(&running, -1)
		return nil
	}
	a := newReconfigurableApp(t, stringProps("p", 16, "v0"), ContinueOnError, FinalizeProperties(refresh),
		EngineSettings(service.NewServiceManager(), nil), overlapProbe)

	errs := make(chan error, 2)
	var wg sync.WaitGroup
	for w := 0; w < 2; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			deadline := time.Now().Add(refreshDuration)
			for i := 0; i < refreshMaxIter && time.Now().Before(deadline); i++ {
				if err := a.Reconfigure(); err != nil {
					errs <- fmt.Errorf("Reconfigure returned %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}

	assert.NoError(t, allEqual("p", 16)(property.DefaultManager().GetProperties()))
	assert.Zero(t, atomic.LoadInt32(&overlaps), "two Reconfigure calls ran at the same time")
}

// TestApp_ReconfigureReleasesLockOnFailure checks that a refresh that fails, or panics, does not block later ones.
func TestApp_ReconfigureReleasesLockOnFailure(t *testing.T) {
	tests := []struct {
		name    string
		fail    func() error
		wantErr string
	}{
		{"error", func() error { return errors.New("option failed") }, "option failed"},
		{"panic", func() error { panic("option failed") }, "panic: option failed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			failSecondCall := func(*App) error {
				calls++
				if calls == 2 {
					return tt.fail()
				}
				return nil
			}
			a := newReconfigurableApp(t, stringProps("p", 1, "v0"), failSecondCall)

			reconfigure := func() (err error) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("panic: %v", r)
					}
				}()
				return a.Reconfigure()
			}
			done := make(chan [2]error, 1)
			go func() {
				first := reconfigure()
				done <- [2]error{first, reconfigure()}
			}()
			select {
			case errs := <-done:
				assert.EqualError(t, errs[0], tt.wantErr)
				assert.NoError(t, errs[1])
			case <-time.After(10 * time.Second):
				t.Fatal("Reconfigure still held its lock after a failed refresh")
			}
		})
	}
}

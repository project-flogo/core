package property

import (
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

// The benchmarks measure the read path before and after FLOGO-19663. They use only the exported API, so the same
// file runs against the code from before the fix.

var (
	benchValue interface{}
	benchOK    bool
	benchMap   map[string]interface{}
)

// newBenchManager returns a manager with 32 overridable properties; NewManager logs nothing for it.
func newBenchManager() *Manager {
	props := make(map[string]Property, 32)
	for i := 0; i < 32; i++ {
		props["p"+strconv.Itoa(i)] = Property{Value: "v0", Overridable: true}
	}
	return NewManager(props)
}

func BenchmarkManager_GetProperty(b *testing.B) {
	m := newBenchManager()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchValue, benchOK = m.GetProperty("p7")
	}
}

// readParallel calls read b.N times in total, split over GOMAXPROCS goroutines. It does not use b.RunParallel:
// RunParallel's per-goroutine counters (testing.PB, written on every iteration) can share a cache line, which makes
// its results bimodal on some machines.
func readParallel(b *testing.B, read func() bool) {
	procs := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	b.ResetTimer()
	for p := 0; p < procs; p++ {
		n := b.N / procs
		if p < b.N%procs {
			n++
		}
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			hits := 0
			for i := 0; i < n; i++ {
				if read() {
					hits++
				}
			}
			if hits != n {
				b.Errorf("%d of %d reads missed", n-hits, n)
			}
		}(n)
	}
	wg.Wait()
}

func BenchmarkManager_GetPropertyParallel(b *testing.B) {
	m := newBenchManager()
	b.ReportAllocs()
	readParallel(b, func() bool {
		_, ok := m.GetProperty("p7")
		return ok
	})
}

// BenchmarkManager_GetPropertyParallelDuringRefresh reads while a refresh runs every millisecond. Before the fix
// this crashes the process, so it has no base value.
func BenchmarkManager_GetPropertyParallelDuringRefresh(b *testing.B) {
	m := newBenchManager()
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(time.Millisecond)
		defer ticker.Stop()
		for i := 1; ; i++ {
			select {
			case <-stop:
				return
			case <-ticker.C:
				v := "v" + strconv.Itoa(i)
				_ = m.Finalize(func(properties map[string]interface{}) error {
					for k := range properties {
						properties[k] = v
					}
					return nil
				})
			}
		}
	}()
	b.ReportAllocs()
	readParallel(b, func() bool {
		_, ok := m.GetProperty("p7")
		return ok
	})
	b.StopTimer()
	close(stop)
	<-done
}

func BenchmarkManager_IsOverridable(b *testing.B) {
	m := newBenchManager()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchOK = m.IsOverridable("p7")
	}
}

func BenchmarkManager_GetProperties(b *testing.B) {
	m := newBenchManager()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchMap = m.GetProperties()
	}
}

func BenchmarkManager_Finalize(b *testing.B) {
	m := newBenchManager()
	noop := func(map[string]interface{}) error { return nil }
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Finalize(noop)
	}
}

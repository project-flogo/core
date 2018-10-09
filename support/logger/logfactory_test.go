package logger

import (
	"math"
	"sync"
	"testing"

	"fmt"
	"github.com/stretchr/testify/assert"
)

// TestConcurrentGetLoggerOk tests that the GetLogger function is concurrent
func TestConcurrentGetLoggerOk(t *testing.T) {
	w := sync.WaitGroup{}
	var recovered interface{}
	//Create factory
	f := &DefaultLoggerFactory{}

	for r := 0; r < 100000; r++ {
		w.Add(1)
		go func(i int) {
			defer w.Done()
			defer func() {
				if r := recover(); r != nil {
					recovered = r
				}
			}()
			loggerName := fmt.Sprintf("logger%f", math.Mod(float64(i), 5))
			f.GetLogger(loggerName)
		}(r)

	}
	w.Wait()
	assert.NotNil(t, f, "Recovered not nil, some problem getting logger")
}

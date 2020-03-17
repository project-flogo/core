package managed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestMangedObject struct {
	started bool
	stopped bool
}

func (t *TestMangedObject) Start() error {
	t.started = true
	return nil
}

func (t *TestMangedObject) Stop() error {
	t.stopped = true
	return nil
}

type TestMangedPanicObject struct {
}

func (t *TestMangedPanicObject) Start() error {
	panic("test panic")
}

func (t *TestMangedPanicObject) Stop() error {
	panic("test panic")
}

func TestStart(t *testing.T) {
	managed1 := &TestMangedObject{}
	err := Start("test", managed1)
	assert.Nil(t, err)
	assert.True(t, managed1.started)
}

func TestStartPanic(t *testing.T) {
	managed1 := &TestMangedPanicObject{}
	err := Start("test", managed1)
	assert.NotNil(t, err)
}

func TestStop(t *testing.T) {
	managed1 := &TestMangedObject{}
	err := Stop("test", managed1)
	assert.Nil(t, err)
	assert.True(t, managed1.stopped)
}

func TestStopPanic(t *testing.T) {
	managed1 := &TestMangedPanicObject{}
	err := Stop("test", managed1)
	assert.NotNil(t, err)
}

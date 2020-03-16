package support

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSyncQueue(t *testing.T) {

	q := NewSyncQueue()
	assert.NotNil(t, q)
	assert.NotNil(t, q.List)
	assert.NotNil(t, q.lock)
}

func TestSyncQueue_IsEmpty(t *testing.T) {

	q := NewSyncQueue()
	assert.NotNil(t, q)

	assert.True(t, q.IsEmpty())
	q.Push(&TestStruct{})
	assert.False(t, q.IsEmpty())
}

func TestSyncQueue_Pop(t *testing.T) {
	q := NewSyncQueue()
	assert.NotNil(t, q)

	testVal := &TestStruct{}
	q.List.PushFront(testVal)

	v, hasVal := q.Pop()
	assert.True(t, hasVal)
	assert.Equal(t, testVal, v)

	v, hasVal = q.Pop()
	assert.False(t, hasVal)
	assert.Nil(t, v)
}

func TestSyncQueue_Push(t *testing.T) {
	q := NewSyncQueue()
	assert.NotNil(t, q)

	testVal := &TestStruct{}
	q.Push(testVal)

	assert.Equal(t, 1, q.List.Len())
	assert.Equal(t, testVal, q.List.Front().Value)
}

func TestSyncQueue_Size(t *testing.T) {

	q := NewSyncQueue()
	assert.NotNil(t, q)

	assert.Equal(t, 0, q.Size())

	testVal := &TestStruct{}
	q.List.PushFront(testVal)

	assert.Equal(t, 1, q.Size())
}

package support

import (
	"container/list"
	"sync"
)

// SyncQueue is a List backed queue
type SyncQueue struct {
	List *list.List
	lock sync.Mutex
}

//NewSyncQueue creates a new SyncQueue
func NewSyncQueue() *SyncQueue {
	return &SyncQueue{List: list.New(), lock: sync.Mutex{}}
}

// Push push item on to queue
func (sq *SyncQueue) Push(item interface{}) {
	sq.lock.Lock()
	defer sq.lock.Unlock()

	sq.List.PushFront(item)
}

// Pop pop item off of queue
func (sq *SyncQueue) Pop() (interface{}, bool) {
	sq.lock.Lock()
	defer sq.lock.Unlock()

	if sq.List.Len() == 0 {
		return nil, false
	}

	item := sq.List.Front()
	sq.List.Remove(item)

	return item.Value, true
}

// Size get the size of the queue
func (sq *SyncQueue) Size() int {
	sq.lock.Lock()
	defer sq.lock.Unlock()

	return sq.List.Len()
}

// IsEmpty indicates if the queue is empty
func (sq *SyncQueue) IsEmpty() bool {
	sq.lock.Lock()
	defer sq.lock.Unlock()

	return sq.List.Len() == 0
}

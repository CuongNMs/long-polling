package model

import "sync"

type CappedQueue[T any] struct {
	items    []T
	lock     *sync.RWMutex
	capacity int
}

func NewCappedQueue[T any](capacity int) *CappedQueue[T] {
	return &CappedQueue[T]{
		items:    make([]T, 0, capacity),
		lock:     new(sync.RWMutex),
		capacity: capacity,
	}
}

func (q *CappedQueue[T]) Append(item T) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if l := len(q.items); l == 0 {
		q.items = append(q.items, item)
	} else {
		to := q.capacity - 1
		if l < q.capacity {
			to = l
		}
		q.items = append([]T{item}, q.items[:to]...)
	}
}

func (q *CappedQueue[T]) Copy() []T {
	q.lock.Lock()
	defer q.lock.Unlock()
	copied := make([]T, len(q.items))
	for i, item := range q.items {
		copied[i] = item
	}
	return copied
}

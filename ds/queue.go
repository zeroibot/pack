package ds

import "fmt"

type Queue[T any] struct {
	items List[T]
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return new(Queue[T]{items: make(List[T], 0)})
}

// NewQueueFrom creates a new queue from list of items
func NewQueueFrom[T any](items []T) *Queue[T] {
	return new(Queue[T]{items: items})
}

// String returns the string representation of the queue
func (q *Queue[T]) String() string {
	return fmt.Sprintf("%v", q.items)
}

// Items returns the List of queue items
func (q *Queue[T]) Items() List[T] {
	return q.items
}

// Len returns the number of items in the queue
func (q *Queue[T]) Len() int {
	return q.items.Len()
}

// IsEmpty checks if queue is empty
func (q *Queue[T]) IsEmpty() bool {
	return q.items.IsEmpty()
}

// NotEmpty if queue is not empty
func (q *Queue[T]) NotEmpty() bool {
	return q.items.NotEmpty()
}

// Clear removes all queue items
func (q *Queue[T]) Clear() {
	q.items.Clear()
}

// Copy creates a new Queue with copied items
func (q *Queue[T]) Copy() *Queue[T] {
	return NewQueueFrom[T](q.items.Copy())
}

// Enqueue

// Front

// MustFront

// Dequeue

// MustDequeue

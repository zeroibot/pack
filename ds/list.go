package ds

type List[T any] []T

// NewList creates an empty List with given capacity
func NewList[T any](capacity int) List[T] {
	return make(List[T], 0, capacity)
}

// Len returns the List length
func (l List[T]) Len() int {
	return len(l)
}

// Cap returns the List capacity
func (l List[T]) Cap() int {
	return cap(l)
}

// IsEmpty checks if List is empty
func (l List[T]) IsEmpty() bool {
	return len(l) == 0
}

// NotEmpty checks if List is not empty
func (l List[T]) NotEmpty() bool {
	return len(l) > 0
}

// Clear removes all List items
func (l List[T]) Clear() {
	clear(l)
}

// Copy creates a new List with copied items
func (l List[T]) Copy() List[T] {
	return append(List[T]{}, l...)
}

// ToAny

// Remove

// Has

// HasFunc

// GetOrDefault

// Last

// MustLast

// Shuffle

// GetRandom

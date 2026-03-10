package ds

import (
	"fmt"
	"strings"
)

type Set[T comparable] struct {
	items Map[T, struct{}]
}

// NewSet creates a new empty Set
func NewSet[T comparable]() Set[T] {
	return Set[T]{items: make(Map[T, struct{}])}
}

// NewSetFrom creates a new Set from given items
func NewSetFrom[T comparable](items []T) Set[T] {
	set := NewSet[T]()
	set.Add(items...)
	return set
}

// NewSetFunc creates a new Set from given items, using the keyFn
func NewSetFunc[T any, K comparable](items []T, keyFn func(T) K) Set[K] {
	set := NewSet[K]()
	for _, item := range items {
		set.Add(keyFn(item))
	}
	return set
}

// String returns the string representation of Set
func (s Set[T]) String() string {
	out := fmt.Sprintf("%v", s.items.Keys())
	out = strings.Trim(out, "[]")
	return "{" + out + "}"
}

// Len returns the Set size
func (s Set[T]) Len() int {
	return s.items.Len()
}

// IsEmpty checks if the Set is empty
func (s Set[T]) IsEmpty() bool {
	return s.items.IsEmpty()
}

// NotEmpty checks if the Set is not empty
func (s Set[T]) NotEmpty() bool {
	return s.items.NotEmpty()
}

// Clear removes all Set items
func (s Set[T]) Clear() {
	s.items.Clear()
}

// Copy creates a new Set with copied items
func (s Set[T]) Copy() Set[T] {
	return Set[T]{items: s.items.Copy()}
}

// Items returns the Set items, in arbitrary order
func (s Set[T]) Items() List[T] {
	return s.items.Keys()
}

// Add adds items to the Set
func (s Set[T]) Add(items ...T) {
	for _, item := range items {
		s.items[item] = struct{}{}
	}
}

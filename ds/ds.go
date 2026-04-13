// Package ds contains data structures
package ds

import "fmt"

type Option[T any] struct {
	value T
	isNil bool
}

// NewOption creates a new Option
func NewOption[T any](value *T) Option[T] {
	if value == nil {
		var zero T
		return Option[T]{zero, true}
	}
	return Option[T]{*value, false}
}

// Nil creates a new nil Option
func Nil[T any]() Option[T] {
	var zero T
	return Option[T]{zero, true}
}

// IsNil checks if Option is nil
func (o Option[T]) IsNil() bool {
	return o.isNil
}

// NotNil checks if Option is not nil
func (o Option[T]) NotNil() bool {
	return !o.isNil
}

// Get returns the value and a Boolean indicating value is not nil
func (o Option[T]) Get() (T, bool) {
	return o.value, !o.isNil
}

// Value returns the stored Option value, could be the zero value if Option is nil
func (o Option[T]) Value() T {
	return o.value
}

// String returns the string representation of option
func (o Option[T]) String() string {
	if o.isNil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", o.value)
}

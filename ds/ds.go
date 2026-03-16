// Package ds contains data structures
package ds

import "fmt"

type Option[T any] struct {
	value T
	isNil bool
}

type Result[T any] struct {
	value T
	err   error
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

// NewResult creates a new Result
func NewResult[T any](value T, err error) Result[T] {
	return Result[T]{value, err}
}

// Error creates a new Result with error
func Error[T any](err error) Result[T] {
	var zero T
	return Result[T]{zero, err}
}

// IsError checks if Result has an error
func (r Result[T]) IsError() bool {
	return r.err != nil
}

// NotError checks if Result has no error
func (r Result[T]) NotError() bool {
	return r.err == nil
}

// Error gets the Result error
func (r Result[T]) Error() error {
	return r.err
}

// Get returns the value and a Boolean indicating Result has no error
func (r Result[T]) Get() (T, bool) {
	return r.value, r.err == nil
}

// Value returns the stored Result value, could be the zero value if Result is error
func (r Result[T]) Value() T {
	return r.value
}

// String returns the string representation of Result
func (r Result[T]) String() string {
	if r.err != nil {
		return fmt.Sprintf("error: %s", r.err.Error())
	}
	return fmt.Sprintf("%v", r.value)
}

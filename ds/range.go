package ds

import (
	"fmt"

	"github.com/roidaradal/pack/number"
)

// Range represents a range of Integers from [start, end)
type Range[T number.Integer] struct {
	start      T
	end        T
	isReversed bool
}

// NewRange creates a new Range from [start, end)
func NewRange[T number.Integer](start, end T) Range[T] {
	isReversed := start > end
	return Range[T]{start, end, isReversed}
}

// NewInclusiveRange creates a new Range from [first, last]
func NewInclusiveRange[T number.Integer](first, last T) Range[T] {
	if first > last {
		// Reversed Range
		return Range[T]{first, last - 1, true}
	}
	return Range[T]{first, last + 1, false}
}

// Limits returns the start, end limits of the Range
func (r Range[T]) Limits() (start, end T) {
	return r.start, r.end
}

// IsReversed checks if the Range is reversed
func (r Range[T]) IsReversed() bool {
	return r.isReversed
}

// Len returns the size of the Range
func (r Range[T]) Len() int {
	if r.isReversed {
		return int(r.start - r.end)
	}
	return int(r.end - r.start)
}

// Copy creates a new Range copy
func (r Range[T]) Copy() Range[T] {
	return NewRange[T](r.start, r.end)
}

// Has checks if Integer is included in the Range
func (r Range[T]) Has(item T) bool {
	if r.isReversed {
		return r.start >= item && item > r.end
	}
	return r.start <= item && item < r.end
}

// ToSlice expands the range into a slice of Integers
func (r Range[T]) ToSlice() []T {
	items := make([]T, 0, r.Len())
	if r.isReversed {
		for x := r.start; x > r.end; x-- {
			items = append(items, x)
		}
	} else {
		for x := r.start; x < r.end; x++ {
			items = append(items, x)
		}
	}
	return items
}

// Sum computes the sum of the Range
func (r Range[T]) Sum() T {
	var total T = 0
	if r.isReversed {
		for x := r.start; x > r.end; x-- {
			total += x
		}
	} else {
		for x := r.start; x < r.end; x++ {
			total += x
		}
	}
	return total
}

// Product computes the product of the Range
func (r Range[T]) Product() T {
	var product T = 1
	if r.isReversed {
		for x := r.start; x > r.end; x-- {
			product *= x
		}
	} else {
		for x := r.start; x < r.end; x++ {
			product *= x
		}
	}
	return product
}

// String returns the string representation of the Range
func (r Range[T]) String() string {
	return fmt.Sprintf("[%d, %d)", r.start, r.end)
}

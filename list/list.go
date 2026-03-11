package list

import (
	"slices"

	"github.com/roidaradal/pack/number"
)

// NewEmpty creates an empty list with given capacity
func NewEmpty[T any](capacity int) []T {
	return make([]T, 0, capacity)
}

// Range creates a new list, with numbers from [start, end)
func Range[T number.Integer](start, end T) []T {
	items := NewEmpty[T](int(end - start))
	for i := start; i < end; i++ {
		items = append(items, i)
	}
	return items
}

// InclusiveRange creates a new list, with numbers from [first, last]
func InclusiveRange[T number.Integer](first, last T) []T {
	return Range(first, last+1)
}

// RepeatedItem creates a new list, with <value> repeated <count> times
func RepeatedItem[T any](value T, count int) []T {
	return slices.Repeat([]T{value}, count)
}

// Len returns the list length
func Len[T any](items []T) int {
	return len(items)
}

// Cap returns the list capacity
func Cap[T any](items []T) int {
	return cap(items)
}

// IsEmpty checks if list is empty
func IsEmpty[T any](items []T) bool {
	return len(items) == 0
}

// NotEmpty checks if list is not empty
func NotEmpty[T any](items []T) bool {
	return len(items) > 0
}

// Clear removes all list items
func Clear[T any](items []T) {
	clear(items)
}

// Copy creates a new list with copied items
func Copy[T any](items []T) []T {
	return append([]T{}, items...)
}

// ToAny

// Remove

// Has

// HasFunc

// HasNo

// HasNoFunc

// GetOrDefault

// Last

// MustLast

// Shuffle

// GetRandom

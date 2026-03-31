// Package list contains list functions
package list

import (
	"math/rand/v2"
	"slices"

	"github.com/zeroibot/pack/number"
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

// LastIndex returns the list last index
func LastIndex[T any](items []T) int {
	return len(items) - 1
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

// Copy creates a new list with copied items
func Copy[T any](items []T) []T {
	return append([]T{}, items...)
}

// ToAny creates a list of <any> items from the list
func ToAny[T any](items []T) []any {
	items2 := make([]any, len(items))
	for i, item := range items {
		items2[i] = item
	}
	return items2
}

// IndexFunc returns the index of item (or -1 if not in list), using the item function
func IndexFunc[T any](items []T, itemFn func(T) bool) int {
	return slices.IndexFunc(items, itemFn)
}

// AllIndexFunc returns all indexes of item in the list, using the item function
func AllIndexFunc[T any](items []T, itemFn func(T) bool) []int {
	indexes := make([]int, 0, len(items))
	for i, item := range items {
		if itemFn(item) {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

// RemoveFunc removes the first item from list that passes the item function
func RemoveFunc[T any](items []T, itemFn func(T) bool) ([]T, bool) {
	index := IndexFunc(items, itemFn)
	if index < 0 {
		return items, false
	}
	result := slices.Delete(items, index, index+1)
	return result, true
}

// RemoveAllFunc removes all items from list that passes the item function
func RemoveAllFunc[T any](items []T, itemFn func(T) bool) []T {
	return slices.DeleteFunc(items, itemFn)
}

// GetFuncOrDefault returns the first item that passes the item function, or returns the default value
func GetFuncOrDefault[T any](items []T, itemFn func(T) bool, defaultValue T) T {
	index := IndexFunc(items, itemFn)
	if index < 0 {
		return defaultValue
	}
	return items[index]
}

// Last returns the nth item from the back of the list (starts at 1), and a flag which indicates if it is valid
func Last[T any](items []T, rank int) (T, bool) {
	numItems := len(items)
	if rank > numItems || rank <= 0 {
		var zero T
		return zero, false
	}
	return items[numItems-rank], true
}

// MustLast returns the nth item from the back of the list (starts at 1).
// Panics if rank is not 1 <= rank <= N, where N = length of List.
func MustLast[T any](items []T, rank int) T {
	item, ok := Last(items, rank)
	if !ok {
		panic("invalid rank")
	}
	return item
}

// GetRandom gets a random item from List, and a flag which indicates if it is valid
func GetRandom[T any](items []T) (T, bool) {
	numItems := len(items)
	if numItems == 0 {
		var zero T
		return zero, false
	}
	return items[rand.IntN(numItems)], true
}

// MustGetRandom gets a random item from List, and panics if list is empty
func MustGetRandom[T any](items []T) T {
	item, ok := GetRandom(items)
	if !ok {
		panic("empty list")
	}
	return item
}

// Shuffle shuffles the List in place
func Shuffle[T any](items []T) {
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
}

// Any checks if list has an item that passes the ok function
func Any[T any](items []T, ok func(T) bool) bool {
	return slices.ContainsFunc(items, ok)
}

// NotAny checks if list has no item that passes the ok function
func NotAny[T any](items []T, ok func(T) bool) bool {
	return !slices.ContainsFunc(items, ok)
}

// All checks if all list items pass the ok function
func All[T any](items []T, ok func(T) bool) bool {
	if len(items) == 0 {
		return false
	}
	for _, item := range items {
		if !ok(item) {
			return false
		}
	}
	return true
}

// AnyIndexed checks if any list item passes the ok function: (index, item)
func AnyIndexed[T any](items []T, ok func(int, T) bool) bool {
	for i, item := range items {
		if ok(i, item) {
			return true
		}
	}
	return false
}

// NotAnyIndexed checks if no list item passes the ok function: (index, item)
func NotAnyIndexed[T any](items []T, ok func(int, T) bool) bool {
	return !AnyIndexed(items, ok)
}

// AllIndexed checks if all list item passes the ok function: (index, item)
func AllIndexed[T any](items []T, ok func(int, T) bool) bool {
	if len(items) == 0 {
		return false
	}
	for i, item := range items {
		if !ok(i, item) {
			return false
		}
	}
	return true
}

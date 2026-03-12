package ds

import (
	"math/rand/v2"
	"slices"

	"github.com/roidaradal/pack/number"
)

type List[T any] []T

type NumList[T number.Type] List[T]

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

// ToAny creates a List of <any> items from the List
func (l List[T]) ToAny() List[any] {
	items := make([]any, len(l))
	for i, item := range l {
		items[i] = item
	}
	return items
}

// IndexFunc returns the index of item (or -1 if not in List), using the item function
func (l List[T]) IndexFunc(itemFn func(T) bool) int {
	return slices.IndexFunc(l, itemFn)
}

// AllIndexFunc returns all indexes of item in the List, using the item function
func (l List[T]) AllIndexFunc(itemFn func(T) bool) List[int] {
	indexes := make(List[int], 0, len(l))
	for i, item := range l {
		if itemFn(item) {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

// RemoveFunc removes the first item from List that passes the item function
func (l List[T]) RemoveFunc(itemFn func(T) bool) (List[T], bool) {
	index := l.IndexFunc(itemFn)
	if index < 0 {
		return l, false
	}
	result := slices.Delete(l, index, index+1)
	return result, true
}

// RemoveAllFunc removes all items from List that passes the item function
func (l List[T]) RemoveAllFunc(itemFn func(T) bool) List[T] {
	return slices.DeleteFunc(l, itemFn)
}

// Get returns an Option with the List item if valid index, otherwise nil
func (l List[T]) Get(index int) Option[T] {
	if index < 0 || index >= len(l) {
		return Nil[T]()
	}
	return NewOption(&l[index])
}

// GetFuncOrDefault returns the first item that passes the item function, or returns the default value
func (l List[T]) GetFuncOrDefault(itemFn func(T) bool, defaultValue T) T {
	index := l.IndexFunc(itemFn)
	if index < 0 {
		return defaultValue
	}
	return l[index]
}

// Last returns the nth item from the back of the List (starts at 1), and a flag which indicates if it is valid
func (l List[T]) Last(rank int) Option[T] {
	numItems := len(l)
	if rank > numItems || rank <= 0 {
		return Nil[T]()
	}
	return NewOption(&l[numItems-rank])
}

// MustLast returns the nth item from the back of the List (starts at 1).
// Panics if rank is not 1 <= rank <= N, where N = length of List.
func (l List[T]) MustLast(rank int) T {
	item := l.Last(rank)
	if item.IsNil() {
		panic("invalid rank")
	}
	return item.Value()
}

// GetRandom gets a random item from List, and a flag which indicates if it is valid
func (l List[T]) GetRandom() Option[T] {
	numItems := len(l)
	if numItems == 0 {
		return Nil[T]()
	}
	return NewOption(&l[rand.IntN(numItems)])
}

// MustGetRandom gets a random item from List, and panics if list is empty
func (l List[T]) MustGetRandom() T {
	item := l.GetRandom()
	if item.IsNil() {
		panic("empty list")
	}
	return item.Value()
}

// Shuffle shuffles the List in place
func (l List[T]) Shuffle() {
	rand.Shuffle(len(l), func(i, j int) {
		l[i], l[j] = l[j], l[i]
	})
}

// Any checks if List has an item that passes the ok function
func (l List[T]) Any(ok func(T) bool) bool {
	return slices.ContainsFunc(l, ok)
}

// NotAny checks if List has no item that passes the ok function
func (l List[T]) NotAny(ok func(T) bool) bool {
	return !slices.ContainsFunc(l, ok)
}

// All checks if all List items pass the ok function
func (l List[T]) All(ok func(T) bool) bool {
	if len(l) == 0 {
		return false
	}
	for _, item := range l {
		if !ok(item) {
			return false
		}
	}
	return true
}

// AnyIndexed checks if any List item passes the ok function: (index, item)
func (l List[T]) AnyIndexed(ok func(int, T) bool) bool {
	for i, item := range l {
		if ok(i, item) {
			return true
		}
	}
	return false
}

// NotAnyIndexed checks if no List item passes the ok function: (index, item)
func (l List[T]) NotAnyIndexed(ok func(int, T) bool) bool {
	return !l.AnyIndexed(ok)
}

// AllIndexed checks if all List item passes the ok function: (index, item)
func (l List[T]) AllIndexed(ok func(int, T) bool) bool {
	if len(l) == 0 {
		return false
	}
	for i, item := range l {
		if !ok(i, item) {
			return false
		}
	}
	return true
}

// CountFunc counts the number of item that passes the ok function
func (l List[T]) CountFunc(ok func(T) bool) int {
	count := 0
	for _, item := range l {
		if ok(item) {
			count += 1
		}
	}
	return count
}

// MapList maps the indexes to List items.
// Can have zero values for invalid indexes
func (l List[T]) MapList(indexes []int) List[T] {
	results := make(List[T], len(indexes))
	numItems := len(l)
	for i, index := range indexes {
		if 0 <= index && index < numItems {
			results[i] = l[index]
		}
	}
	return results
}

// Filter filters the List by only keeping items that pass the keep function
func (l List[T]) Filter(keep func(T) bool) List[T] {
	results := make(List[T], 0, len(l))
	for _, item := range l {
		if keep(item) {
			results = append(results, item)
		}
	}
	return results
}

// FilterIndexed filters the List by only keeping items that pass the keep function: (index, item)
func (l List[T]) FilterIndexed(keep func(int, T) bool) List[T] {
	results := make(List[T], 0, len(l))
	for i, item := range l {
		if keep(i, item) {
			results = append(results, item)
		}
	}
	return results
}

// Reduce applies the reducer to each item to get the final result.
// The reducer function has the signature (result, item) => result
func (l List[T]) Reduce(reducer func(T, T) T, initial T) T {
	current := initial
	for _, item := range l {
		current = reducer(current, item)
	}
	return current
}

// Apply applies the task function to each item
func (l List[T]) Apply(task func(T) T) List[T] {
	results := make(List[T], len(l))
	for i, item := range l {
		results[i] = task(item)
	}
	return results
}

// ToList type coerces the NumList to List
func (l NumList[T]) ToList() List[T] {
	return List[T](l)
}

// Sum computes the sum of number items
func (l NumList[T]) Sum() T {
	var total T = 0
	for _, x := range l {
		total += x
	}
	return total
}

// Product computes the product of number items
func (l NumList[T]) Product() T {
	var product T = 1
	for _, x := range l {
		product *= x
	}
	return product
}

package ds

import (
	"math/rand/v2"
	"slices"
)

// List extends the slice collection
type List[T any] []T

// NewEmptyList creates an empty List, with given capacity
func NewEmptyList[T any](capacity int) List[T] {
	return make(List[T], 0, capacity)
}

// NewRangeList creates a new List, containing numbers from [start, end)
func NewRangeList[T Integer](start, end T) List[T] {
	items := NewEmptyList[T](int(end - start))
	for i := start; i < end; i++ {
		items = append(items, i)
	}
	return items
}

// NewInclusiveRangeList creates a new List, containing numbers from [first, last]
func NewInclusiveRangeList[T Integer](first, last T) List[T] {
	return NewRangeList[T](first, last+1)
}

// NewRepeatedItemList creates a new List, with <value> repeated <count> times
func NewRepeatedItemList[T any](value T, count int) List[T] {
	return slices.Repeat(List[T]{value}, count)
}

// Len returns the List length
func (l List[T]) Len() Int {
	return Int(len(l))
}

// Cap returns the List capacity
func (l List[T]) Cap() Int {
	return Int(cap(l))
}

// IsEmpty checks if the List is empty
func (l List[T]) IsEmpty() Boolean {
	return len(l) == 0
}

// NotEmpty checks if the List is not empty
func (l List[T]) NotEmpty() Boolean {
	return len(l) > 0
}

// Copy creates a new List with copied items
func (l List[T]) Copy() List[T] {
	items := append(List[T]{}, l...)
	return items
}

// ToAnyList creates a List of <any> items from List
func (l List[T]) ToAnyList() List[any] {
	items := make(List[any], len(l))
	for i, item := range l {
		items[i] = item
	}
	return items
}

// Last returns the nth item from the back of the list (starts at 1)
func (l List[T]) Last(rank int) (T, Boolean) {
	numItems := len(l)
	if rank > numItems || rank <= 0 {
		var item T
		return item, false
	}
	return l[numItems-rank], true
}

// Shuffle shuffles the List in place
func (l List[T]) Shuffle() {
	rand.Shuffle(len(l), func(i, j int) {
		l[i], l[j] = l[j], l[i]
	})
}

// GetRandom gets a random item from List
func (l List[T]) GetRandom() T {
	return l[rand.IntN(len(l))]
}

// All checks if all List items pass the ok function
func (l List[T]) All(ok func(T) Boolean) Boolean {
	for _, item := range l {
		if !ok(item) {
			return false
		}
	}
	return true
}

// AllWithIndex checks if all List items pass the ok function: (index, item)
func (l List[T]) AllWithIndex(ok func(int, T) Boolean) Boolean {
	for i, item := range l {
		if !ok(i, item) {
			return false
		}
	}
	return true
}

// Any checks if any List item passes the ok function
func (l List[T]) Any(ok func(T) Boolean) Boolean {
	for _, item := range l {
		if ok(item) {
			return true
		}
	}
	return false
}

// AnyWithIndex checks if any List item passes the ok function: (index, item)
func (l List[T]) AnyWithIndex(ok func(int, T) Boolean) Boolean {
	for i, item := range l {
		if ok(i, item) {
			return true
		}
	}
	return false
}

// Filter filters the list by only keeping items that pass the keep function
func (l List[T]) Filter(keep func(T) Boolean) List[T] {
	results := NewEmptyList[T](len(l))
	for _, item := range l {
		if keep(item) {
			results = append(results, item)
		}
	}
	return results
}

// FilterWithIndex filters the list by only keeping items that pass the keep function: (index, item)
func (l List[T]) FilterWithIndex(keep func(int, T) Boolean) List[T] {
	results := NewEmptyList[T](len(l))
	for i, item := range l {
		if keep(i, item) {
			results = append(results, item)
		}
	}
	return results
}

// Reduce applies the reduce function to each item to get the final result
// Reducer function has the signature (result, item) => result
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

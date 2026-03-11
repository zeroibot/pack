package dict

import "fmt"

type (
	Object                      = map[string]any
	Ints                        = map[string]int
	Uints                       = map[string]uint
	Floats                      = map[string]float64
	Bools                       = map[string]bool
	Strings                     = map[string]string
	StringLists                 = map[string][]string
	Lookup[K comparable, V any] = map[K]V
)

type (
	Counter[T comparable] = map[T]int
	IntCounter            = map[int]int
	UintCounter           = map[uint]int
	StringCounter         = map[string]int
)

type (
	Flags[T comparable] = map[T]bool
	IntFlags            = map[int]bool
	UintFlags           = map[uint]bool
	StringFlags         = map[string]bool
)

// Entry represents a map Key-Value pair
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Tuple returns the unpacked Key, Value of Entry
func (e Entry[K, V]) Tuple() (K, V) {
	return e.Key, e.Value
}

// String returns the string representation of Entry
func (e Entry[K, V]) String() string {
	return fmt.Sprintf("<%v: %v>", e.Key, e.Value)
}

// NewCounterFor creates a new counter, with each item initialized to count = 0
func NewCounterFor[T comparable](items []T) Counter[T] {
	counter := make(Counter[T], len(items))
	for _, item := range items {
		counter[item] = 0
	}
	return counter
}

// NewCounterFunc creates a new counter, using the keys produced from keyFn
func NewCounterFunc[T any, K comparable](items []T, keyFn func(T) K) Counter[K] {
	counter := make(Counter[K], len(items))
	for _, item := range items {
		counter[keyFn(item)] = 0
	}
	return counter
}

// UpdateCounter updates the counter with incoming items
func UpdateCounter[T comparable](counter Counter[T], items []T) {
	for _, item := range items {
		counter[item] += 1
	}
}

// UpdateCounterFunc updates the counter with incoming items, using the keyFn
func UpdateCounterFunc[T any, K comparable](counter Counter[K], items []T, keyFn func(T) K) {
	for _, item := range items {
		counter[keyFn(item)] += 1
	}
}

// CounterUpdate updates the old counter with counts from new counter
func CounterUpdate[T comparable](oldCounter, newCounter Counter[T]) {
	for key, count := range newCounter {
		oldCounter[key] += count
	}
}

// MergeCounters merges the counts from given counters into one Counter
func MergeCounters[T comparable](counters ...Counter[T]) Counter[T] {
	total := make(Counter[T])
	for _, counter := range counters {
		for key, count := range counter {
			total[key] += count
		}
	}
	return total
}

// NewFlagsFor creates a new Flags map, with each item initialized to given flag
func NewFlagsFor[T comparable](items []T, flag bool) Flags[T] {
	flags := make(Flags[T], len(items))
	for _, item := range items {
		flags[item] = flag
	}
	return flags
}

// NewFlagsFunc creates a new Flags map, using the keyFn, with each key initialized to given flag
func NewFlagsFunc[T any, K comparable](items []T, flag bool, keyFn func(T) K) Flags[K] {
	flags := make(Flags[K], len(items))
	for _, item := range items {
		flags[keyFn(item)] = flag
	}
	return flags
}

// LookupFunc creates a converter function that returns associated value and exists flag for given key
func LookupFunc[K comparable, V any](items map[K]V) func(K) (V, bool) {
	return func(key K) (V, bool) {
		value, ok := items[key]
		return value, ok
	}
}

// MustLookupFunc creates a converter function that returns associated value for given key, and panics if key is not found
func MustLookupFunc[K comparable, V any](items map[K]V) func(K) V {
	return func(key K) V {
		value, ok := items[key]
		if !ok {
			panic("key not found")
		}
		return value
	}
}

// Get retrieves the value from Object, then type coerces into type T
func Get[T any](obj Object, key string) (T, bool) {
	var item T
	value, ok := obj[key]
	if !ok {
		return item, false
	}
	item, ok = value.(T)
	return item, ok
}

// GetRef retrieves the ref value from Object, then type coerces into *T
func GetRef[T any](obj Object, key string) *T {
	itemRef, ok := Get[*T](obj, key)
	if !ok {
		return nil
	}
	return itemRef
}

// GetList retrieves the list value from Object, then type coerces into []T
func GetList[T any](obj Object, key string) []T {
	items, ok := Get[[]T](obj, key)
	if !ok {
		return nil
	}
	return items
}

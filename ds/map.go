package ds

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"
)

type Map[K comparable, V any] map[K]V

// ZipMap creates a new Map by zipping the keys and values
func ZipMap[K comparable, V any](keys List[K], values List[V]) Map[K, V] {
	m := make(Map[K, V], len(keys))
	numValues := len(values)
	for i, key := range keys {
		if i >= numValues {
			break // stop if no more values
		}
		m[key] = values[i]
	}
	return m
}

// Unzip returns the list of Map keys and values, where order of keys is same as corresponding values
func (m Map[K, V]) Unzip() (List[K], List[V]) {
	numItems := len(m)
	keys := make(List[K], numItems)
	values := make(List[V], numItems)
	i := 0
	for k, v := range m {
		keys[i] = k
		values[i] = v
		i++
	}
	return keys, values
}

// String returns the string representation of Map, where keys are sorted
func (m Map[K, V]) String() string {
	out := make([]string, 0, len(m))
	for k, v := range m {
		out = append(out, fmt.Sprintf("%v: %v", k, v))
	}
	slices.Sort(out)
	return "{" + strings.Join(out, ", ") + "}"
}

// Len returns the Map size
func (m Map[K, V]) Len() int {
	return len(m)
}

// IsEmpty checks if Map is empty
func (m Map[K, V]) IsEmpty() bool {
	return len(m) == 0
}

// NotEmpty checks if Map is not empty
func (m Map[K, V]) NotEmpty() bool {
	return len(m) > 0
}

// Clear removes all Map entries
func (m Map[K, V]) Clear() {
	clear(m)
}

// Copy creates a new Map with copied entries
func (m Map[K, V]) Copy() Map[K, V] {
	m2 := make(Map[K, V], len(m))
	maps.Copy(m2, m)
	return m2
}

// Update adds the entries of new Map to current Map.
// If there are overlapping keys, new Map entries overwrite the old Map entries.
func (m Map[K, V]) Update(newMap Map[K, V]) {
	maps.Copy(m, newMap)
}

// Delete removes key from the Map
func (m Map[K, V]) Delete(key K) {
	delete(m, key)
}

// KeysIter returns an iterator for the Map keys
func (m Map[K, V]) KeysIter() iter.Seq[K] {
	return maps.Keys(m)
}

// ValuesIter returns an iterator for the Map values
func (m Map[K, V]) ValuesIter() iter.Seq[V] {
	return maps.Values(m)
}

// Keys returns the Map keys, in arbitrary order
func (m Map[K, V]) Keys() List[K] {
	return slices.Collect(m.KeysIter())
}

// SortedKeysFunc returns the Map keys in sorted order, using sortFn
func (m Map[K, V]) SortedKeysFunc(sortFn func(K, K) int) List[K] {
	keys := m.Keys()
	slices.SortFunc(keys, sortFn)
	return keys
}

// Values returns the Map values, in arbitrary order
func (m Map[K, V]) Values() List[V] {
	return slices.Collect(m.ValuesIter())
}

// SortedValuesFunc returns the Map values in sorted order, using SortFn
func (m Map[K, V]) SortedValuesFunc(sortFn func(V, V) int) List[V] {
	values := m.Values()
	slices.SortFunc(values, sortFn)
	return values
}

// Entries returns the Map entries, in arbitrary order
func (m Map[K, V]) Entries() List[Tuple2[K, V]] {
	entries := make([]Tuple2[K, V], 0, len(m))
	for k, v := range m {
		entries = append(entries, Tuple2[K, V]{k, v})
	}
	return entries
}

// SortedEntriesFunc returns the Map entries in sorted order, using SortFn
func (m Map[K, V]) SortedEntriesFunc(sortFn func(Tuple2[K, V], Tuple2[K, V]) int) List[Tuple2[K, V]] {
	entries := m.Entries()
	slices.SortFunc(entries, sortFn)
	return entries
}

// HasKey checks if Map has given key
func (m Map[K, V]) HasKey(key K) bool {
	_, ok := m[key]
	return ok
}

// HasKeyFunc checks if any Map key passes the test function
func (m Map[K, V]) HasKeyFunc(test func(K) bool) bool {
	for k := range m {
		if test(k) {
			return true
		}
	}
	return false
}

// NoKey checks if Map does not have given key
func (m Map[K, V]) NoKey(key K) bool {
	return !m.HasKey(key)
}

// NoKeyFunc checks that no Map key passes the test function
func (m Map[K, V]) NoKeyFunc(test func(K) bool) bool {
	return !m.HasKeyFunc(test)
}

// HasValueFunc checks if any Map value passes the test function
func (m Map[K, V]) HasValueFunc(test func(V) bool) bool {
	for _, v := range m {
		if test(v) {
			return true
		}
	}
	return false
}

// NoValueFunc checks that no Map value passes the test function
func (m Map[K, V]) NoValueFunc(test func(V) bool) bool {
	return !m.HasValueFunc(test)
}

// SetDefault assigns default value to key, if key is not in Map
func (m Map[K, V]) SetDefault(key K, defaultValue V) {
	if _, ok := m[key]; !ok {
		m[key] = defaultValue
	}
}

// GetOrDefault gets the value associated with key if it exists, otherwise returns the default value
func (m Map[K, V]) GetOrDefault(key K, defaultValue V) V {
	if value, ok := m[key]; ok {
		return value
	}
	return defaultValue
}

// Filter filters the Map, only keeping entries that pass the keep function
func (m Map[K, V]) Filter(keep func(K, V) bool) Map[K, V] {
	result := make(Map[K, V], len(m))
	for k, v := range m {
		if keep(k, v) {
			result[k] = v
		}
	}
	return result
}

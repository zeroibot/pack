// Package dict contains Map functions
package dict

import (
	"iter"
	"maps"
	"slices"
)

// Len returns the map size
func Len[K comparable, V any](items map[K]V) int {
	return len(items)
}

// IsEmpty checks if map is empty
func IsEmpty[K comparable, V any](items map[K]V) bool {
	return len(items) == 0
}

// NotEmpty checks if map is not empty
func NotEmpty[K comparable, V any](items map[K]V) bool {
	return len(items) > 0
}

// Clear removes all map entries
func Clear[K comparable, V any](items map[K]V) {
	clear(items)
}

// Copy creates a new map with copied entries
func Copy[K comparable, V any](items map[K]V) map[K]V {
	items2 := make(map[K]V, len(items))
	maps.Copy(items2, items)
	return items2
}

// Update adds the entries of new Map to current Map.
// If there are overlapping keys, new Map entries overwrite the old Map entries.
func Update[K comparable, V any](oldMap, newMap map[K]V) {
	maps.Copy(oldMap, newMap)
}

// KeysIter returns an iterator for map keys
func KeysIter[K comparable, V any](items map[K]V) iter.Seq[K] {
	return maps.Keys(items)
}

// ValuesIter returns an iterator for map values
func ValuesIter[K comparable, V any](items map[K]V) iter.Seq[V] {
	return maps.Values(items)
}

// Keys returns map keys, in arbitrary order
func Keys[K comparable, V any](items map[K]V) []K {
	return slices.Collect(KeysIter(items))
}

// Values returns map values, in arbitrary order
func Values[K comparable, V any](items map[K]V) []V {
	return slices.Collect(ValuesIter(items))
}

// Entries returns map entries, in arbitrary order
func Entries[K comparable, V any](items map[K]V) []Entry[K, V] {
	entries := make([]Entry[K, V], 0, len(items))
	for k, v := range items {
		entries = append(entries, Entry[K, V]{k, v})
	}
	return entries
}

// HasKey checks if map has given key
func HasKey[K comparable, V any](items map[K]V, key K) bool {
	_, ok := items[key]
	return ok
}

// HasKeyFunc checks if any map key passes the test function
func HasKeyFunc[K comparable, V any](items map[K]V, test func(K) bool) bool {
	for k := range items {
		if test(k) {
			return true
		}
	}
	return false
}

// NoKey checks if map does not have given key
func NoKey[K comparable, V any](items map[K]V, key K) bool {
	return !HasKey(items, key)
}

// NoKeyFunc checks that no map key passes the test function
func NoKeyFunc[K comparable, V any](items map[K]V, test func(K) bool) bool {
	return !HasKeyFunc(items, test)
}

// HasValue checks if map has given value
func HasValue[K, V comparable](items map[K]V, value V) bool {
	for _, v := range items {
		if v == value {
			return true
		}
	}
	return false
}

// HasValueFunc checks if any map value passes the test function
func HasValueFunc[K comparable, V any](items map[K]V, test func(V) bool) bool {
	for _, v := range items {
		if test(v) {
			return true
		}
	}
	return false
}

// NoValue checks if map does not have given value
func NoValue[K, V comparable](items map[K]V, value V) bool {
	return !HasValue(items, value)
}

// NoValueFunc checks that no Map value passes the test function
func NoValueFunc[K comparable, V any](items map[K]V, test func(V) bool) bool {
	return !HasValueFunc(items, test)
}

// SetDefault assigns default value to key, if key is not in Map
func SetDefault[K comparable, V any](items map[K]V, key K, defaultValue V) {
	if _, ok := items[key]; !ok {
		items[key] = defaultValue
	}
}

// GetOrDefault gets the value associated with key if it exists, otherwise returns the default value
func GetOrDefault[K comparable, V any](items map[K]V, key K, defaultValue V) V {
	if value, ok := items[key]; ok {
		return value
	}
	return defaultValue
}

// Filter filters the map, only keeping entries that pass the keep function
func Filter[K comparable, V any](items map[K]V, keep func(K, V) bool) map[K]V {
	result := make(map[K]V, len(items))
	for k, v := range items {
		if keep(k, v) {
			result[k] = v
		}
	}
	return result
}

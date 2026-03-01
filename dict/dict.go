// Package dict contains Map functions
package dict

import (
	"cmp"
	"maps"
	"slices"

	"github.com/roidaradal/pack/ds"
)

// HasKey checks if Map has given key
func HasKey[M ~map[K]V, K comparable, V any](items M, key K) bool {
	_, hasKey := items[key]
	return hasKey
}

// HasValue checks if Map has given value
func HasValue[M ~map[K]V, K, V comparable](items M, value V) bool {
	for _, v := range items {
		if v == value {
			return true
		}
	}
	return false
}

// NoKey checks if Map does not have the given key
func NoKey[M ~map[K]V, K comparable, V any](items M, key K) bool {
	return !HasKey(items, key)
}

// NoValue checks if Map does not have the given value
func NoValue[M ~map[K]V, K, V comparable](items M, value V) bool {
	return !HasValue(items, value)
}

// SortedKeys returns the Map keys in sorted order
func SortedKeys[M ~map[K]V, K cmp.Ordered, V any](items M) ds.List[K] {
	var keys ds.List[K] = slices.Collect(maps.Keys(items))
	slices.Sort(keys)
	return keys
}

// SortedValues returns the Map values in sorted order
func SortedValues[M ~map[K]V, K comparable, V cmp.Ordered](items M) ds.List[V] {
	var values ds.List[V] = slices.Collect(maps.Values(items))
	slices.Sort(values)
	return values
}

// SortedEntries returns the Map entries in sorted key order
func SortedEntries[M ~map[K]V, K cmp.Ordered, V any](items M) ds.List[ds.Entry[K, V]] {
	keys := SortedKeys(items)
	entries := make(ds.List[ds.Entry[K, V]], len(keys))
	for i, k := range keys {
		entries[i] = ds.Entry[K, V]{Key: k, Value: items[k]}
	}
	return entries
}

// SortValueLists sorts the list of values in place, for each key in the Map
func SortValueLists[M ~map[K][]V, K comparable, V cmp.Ordered](items M) {
	for k, values := range items {
		slices.Sort(values)
		items[k] = values
	}
}

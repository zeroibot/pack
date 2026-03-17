package dict

import (
	"cmp"
	"slices"
)

// SortedKeys returns the map keys in sorted order
func SortedKeys[K cmp.Ordered, V any](items map[K]V) []K {
	return slices.Sorted(KeysIter(items))
}

// SortedKeysFunc returns the map keys in sorted order, using the sortFn
func SortedKeysFunc[K comparable, V any](items map[K]V, sortFn func(K, K) int) []K {
	return slices.SortedFunc(KeysIter(items), sortFn)
}

// SortedValues returns the map values in sorted order
func SortedValues[K comparable, V cmp.Ordered](items map[K]V) []V {
	return slices.Sorted(ValuesIter(items))
}

// SortedValuesFunc returns the map values in sorted order, using the sortFn
func SortedValuesFunc[K comparable, V any](items map[K]V, sortFn func(V, V) int) []V {
	return slices.SortedFunc(ValuesIter(items), sortFn)
}

// SortedEntries returns the map entries in sorted key order
func SortedEntries[K cmp.Ordered, V any](items map[K]V) []Entry[K, V] {
	keys := SortedKeys(items)
	entries := make([]Entry[K, V], len(keys))
	for i, key := range keys {
		entries[i] = Entry[K, V]{key, items[key]}
	}
	return entries
}

// SortedEntriesFunc returns the map entries in sorted order, using the sortFn
func SortedEntriesFunc[K comparable, V any](items map[K]V, sortFn func(Entry[K, V], Entry[K, V]) int) []Entry[K, V] {
	entries := Entries(items)
	slices.SortFunc(entries, sortFn)
	return entries
}

// SortValueLists sorts the list of values in place, for each map key
func SortValueLists[K comparable, V cmp.Ordered](items map[K][]V) {
	for k, values := range items {
		slices.Sort(values)
		items[k] = values
	}
}

// SortValueListsFunc sorts the list of values in place, using the sortFn, for each map key
func SortValueListsFunc[K comparable, V any](items map[K][]V, sortFn func(V, V) int) {
	for k, values := range items {
		slices.SortFunc(values, sortFn)
		items[k] = values
	}
}

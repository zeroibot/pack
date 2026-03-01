// Package list contains List functions
package list

import "github.com/roidaradal/pack/ds"

// Map applies the convert function for each List item
func Map[T, V any](items ds.List[T], convert func(T) V) ds.List[V] {
	results := make(ds.List[V], len(items))
	for i, item := range items {
		results[i] = convert(item)
	}
	return results
}

// MapWithIndex applies the convert function for each List item, with index
func MapWithIndex[T, V any](items ds.List[T], convert func(int, T) V) ds.List[V] {
	results := make(ds.List[V], len(items))
	for i, item := range items {
		results[i] = convert(i, item)
	}
	return results
}

// MapList maps the list indexes to given List.
// Can have zero values for invalid indexes
func MapList[T any](indexes ds.List[ds.Int], items ds.List[T]) ds.List[T] {
	results := make(ds.List[T], len(items))
	numItems := ds.Int(len(items))
	for i, idx := range indexes {
		if 0 <= idx && idx < numItems {
			results[i] = items[idx]
		}
	}
	return results
}

// MapLookup maps the keys to the given lookup Map.
// Can have zero values for invalid indexes
func MapLookup[K comparable, V any](keys ds.List[K], lookup ds.Map[K, V]) ds.List[V] {
	results := make(ds.List[V], len(keys))
	for i, key := range keys {
		results[i] = lookup[key]
	}
	return results
}

// MapIf combines Map and Filter: apply the convert function, and filter items based on the result flag
func MapIf[T, V any](items ds.List[T], convert func(T) (V, bool)) ds.List[V] {
	results := ds.NewEmptyList[V](len(items))
	for _, item := range items {
		if item2, ok := convert(item); ok {
			results = append(results, item2)
		}
	}
	return results
}

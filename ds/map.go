package ds

import (
	"maps"
	"slices"
)

// Map extends the map collection
type Map[K comparable, V any] map[K]V

// Entry represents a Key-Value pair
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Len returns the Map size
func (m Map[K, V]) Len() Int {
	return Int(len(m))
}

// IsEmpty checks if the Map is empty
func (m Map[K, V]) IsEmpty() Boolean {
	return len(m) == 0
}

// NotEmpty checks if the Map is not empty
func (m Map[K, V]) NotEmpty() Boolean {
	return len(m) > 0
}

// Copy creates a new Map with copied entries
func (m Map[K, V]) Copy() Map[K, V] {
	items := make(Map[K, V], len(m))
	maps.Copy(items, m)
	return items
}

// Keys returns the Map keys, in arbitrary order
func (m Map[K, V]) Keys() List[K] {
	return slices.Collect(maps.Keys(m))
}

// Values returns the Map values, in arbitrary order
func (m Map[K, V]) Values() List[V] {
	return slices.Collect(maps.Values(m))
}

// Entries returns the Map entries, in arbitrary order
func (m Map[K, V]) Entries() List[Entry[K, V]] {
	entries := NewEmptyList[Entry[K, V]](len(m))
	for k, v := range m {
		entries = append(entries, Entry[K, V]{k, v})
	}
	return entries
}

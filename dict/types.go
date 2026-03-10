package dict

import "fmt"

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

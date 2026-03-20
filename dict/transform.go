package dict

import (
	"cmp"
	"encoding/json"
)

// Zip creates a new map by zipping the keys and values
func Zip[K comparable, V any](keys []K, values []V) map[K]V {
	m := make(map[K]V, len(keys))
	numValues := len(values)
	for i, key := range keys {
		if i >= numValues {
			break // stop if no more values
		}
		m[key] = values[i]
	}
	return m
}

// Unzip returns the list of Map keys and values, where the order of keys is the same as corresponding values
func Unzip[K comparable, V any](items map[K]V) ([]K, []V) {
	numItems := len(items)
	keys := make([]K, numItems)
	values := make([]V, numItems)
	i := 0
	for k, v := range items {
		keys[i] = k
		values[i] = v
		i++
	}
	return keys, values
}

// SortedUnzip returns the list of Map keys and values, where the keys are sorted and the value order corresponds to the key order
func SortedUnzip[K cmp.Ordered, V any](items map[K]V) ([]K, []V) {
	keys := SortedKeys(items)
	values := make([]V, len(keys))
	for i, key := range keys {
		values[i] = items[key]
	}
	return keys, values
}

// Swap swaps the map keys and values, converting map[K]V to map[V]K.
// This can lose data if values are not unique
func Swap[K, V comparable](items map[K]V) map[V]K {
	inverse := make(map[V]K, len(items))
	for k, v := range items {
		inverse[v] = k
	}
	return inverse
}

// SwapList converts the map[K][]V to map[V]K.
// This can lose data if values are not unique
func SwapList[K, V comparable](items map[K][]V) map[V]K {
	inverse := make(map[V]K)
	for k, values := range items {
		for _, v := range values {
			inverse[v] = k
		}
	}
	return inverse
}

// FromStruct creates a map[string]V from given struct pointer.
// Struct field values must all be of the same type.
func FromStruct[V, T any](structRef *T) (map[string]V, error) {
	output := make(map[string]V)
	if structRef == nil {
		return output, nil
	}
	bytes, err := json.Marshal(structRef)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// ToStruct creates a struct reference from given Object
func ToStruct[T any](obj Object) (*T, error) {
	var output T
	if obj == nil {
		return &output, nil
	}
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &output)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

// ToObject creates an Object from given struct pointer
func ToObject[T any](structRef *T) (Object, error) {
	return FromStruct[any](structRef)
}

// Pruned creates an Object from given struct pointer, but only keeps given fieldNames
func Pruned[T any](structRef *T, fieldNames ...string) (Object, error) {
	fullObj, err := ToObject(structRef)
	if err != nil {
		return nil, err
	}
	obj := make(Object, len(fieldNames))
	for _, fieldName := range fieldNames {
		if value, ok := fullObj[fieldName]; ok {
			obj[fieldName] = value
		}
	}
	return obj, nil
}

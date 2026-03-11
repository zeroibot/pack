package dict

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

// Unzip returns the list of Map keys and values, where order of keys is same as corresponding values
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

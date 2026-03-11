package dict

// GroupByValue groups the map by values
func GroupByValue[K, V comparable](items map[K]V) map[V][]K {
	groups := make(map[V][]K)
	for k, v := range items {
		groups[v] = append(groups[v], k)
	}
	return groups
}

// GroupByFunc groups the map by values, using key and value transformers
func GroupByFunc[K, B comparable, V, A any](items map[K]V, keyFn func(K) A, valueFn func(V) B) map[B][]A {
	groups := make(map[B][]A)
	for k, v := range items {
		a, b := keyFn(k), valueFn(v)
		groups[b] = append(groups[b], a)
	}
	return groups
}

// GroupByValueList groups the map[K][]V by values, produces map[V][]K
func GroupByValueList[K, V comparable](items map[K][]V) map[V][]K {
	groups := make(map[V][]K)
	for k := range items {
		for _, v := range items[k] {
			groups[v] = append(groups[v], k)
		}
	}
	return groups
}

// GroupByFuncList groups the map[K][]V by values, produces map[B][]A using key and value transformers
func GroupByFuncList[K, B comparable, V, A any](items map[K][]V, keyFn func(K) A, valueFn func(V) B) map[B][]A {
	groups := make(map[B][]A)
	for k := range items {
		a := keyFn(k)
		for _, v := range items[k] {
			b := valueFn(v)
			groups[b] = append(groups[b], a)
		}
	}
	return groups
}

// TallyValues creates a tally of how many times each value appears in the map
func TallyValues[K, V comparable](items map[K]V, values []V) Counter[V] {
	counter := NewCounterFor(values)
	for _, value := range items {
		if NoKey(counter, value) {
			continue // skip value if not in counter
		}
		counter[value] += 1
	}
	return counter
}

// TallyFunc creates a tally of how many times a transformed value appears in the map, using value transformer
func TallyFunc[K, T comparable, V any](items map[K]V, valueFn func(V) T) Counter[T] {
	counter := make(Counter[T])
	for _, value := range items {
		countKey := valueFn(value)
		counter[countKey] += 1
	}
	return counter
}

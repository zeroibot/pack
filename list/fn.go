package list

import "github.com/zeroibot/pack/number"

// CountFunc counts the number of items that passes the ok function
func CountFunc[T any](items []T, ok func(T) bool) int {
	count := 0
	for _, item := range items {
		if ok(item) {
			count += 1
		}
	}
	return count
}

// Map applies the convert function for each list item
func Map[T, V any](items []T, convert func(T) V) []V {
	results := make([]V, len(items))
	for i, item := range items {
		results[i] = convert(item)
	}
	return results
}

// MapIndexed applies the convert function for each indexed item
func MapIndexed[T, V any](items []T, convert func(int, T) V) []V {
	results := make([]V, len(items))
	for i, item := range items {
		results[i] = convert(i, item)
	}
	return results
}

// MapIf combines Map and Filter: apply the convert function, and filter items based on result flag
func MapIf[T, V any](items []T, convert func(T) (V, bool)) []V {
	results := NewEmpty[V](len(items))
	for _, item := range items {
		if item2, ok := convert(item); ok {
			results = append(results, item2)
		}
	}
	return results
}

// MapIndexedIf combines Map and Filter: apply the convert function (index, item), and filter items based on result flag
func MapIndexedIf[T, V any](items []T, convert func(int, T) (V, bool)) []V {
	results := NewEmpty[V](len(items))
	for i, item := range items {
		if item2, ok := convert(i, item); ok {
			results = append(results, item2)
		}
	}
	return results
}

// MapList maps the indexes to list items.
// Can have zero values for invalid indexes
func MapList[T any](indexes []int, items []T) []T {
	results := make([]T, len(indexes))
	numItems := len(items)
	for i, index := range indexes {
		if 0 <= index && index < numItems {
			results[i] = items[index]
		}
	}
	return results
}

// MapLookup maps the keys to given lookup map.
// Can have zero values for invalid keys
func MapLookup[K comparable, V any](keys []K, lookup map[K]V) []V {
	results := make([]V, len(keys))
	for i, key := range keys {
		results[i] = lookup[key]
	}
	return results
}

// Filter filters the list by only keeping items that pass the keep function
func Filter[T any](items []T, keep func(T) bool) []T {
	results := NewEmpty[T](len(items))
	for _, item := range items {
		if keep(item) {
			results = append(results, item)
		}
	}
	return results
}

// FilterIndexed filters the list by only keeping items that pass the keep function: (index, item)
func FilterIndexed[T any](items []T, keep func(int, T) bool) []T {
	results := NewEmpty[T](len(items))
	for i, item := range items {
		if keep(i, item) {
			results = append(results, item)
		}
	}
	return results
}

// Reduce applies the reducer to each item to get the final result.
// The reducer function has the signature (result, item) => result
func Reduce[T any](items []T, initial T, reducer func(T, T) T) T {
	current := initial
	for _, item := range items {
		current = reducer(current, item)
	}
	return current
}

// Apply applies the task function to each item
func Apply[T any](items []T, task func(T) T) []T {
	results := make([]T, len(items))
	for i, item := range items {
		results[i] = task(item)
	}
	return results
}

// Sum computes the sum of number items
func Sum[T number.Type](numbers []T) T {
	var sum T = 0
	for _, x := range numbers {
		sum += x
	}
	return sum
}

// SumOf computes the sum of mapped number items
func SumOf[T any, V number.Type](items []T, convert func(T) V) V {
	var sum V = 0
	for _, item := range items {
		sum += convert(item)
	}
	return sum
}

// SumIndex combines SumOf and MapList
func SumIndex[T number.Type](indexes []int, items []T) T {
	var sum T = 0
	numItems := len(items)
	for _, idx := range indexes {
		if 0 <= idx && idx < numItems {
			sum += items[idx]
		}
	}
	return sum
}

// SumKey combines SumOf and MapLookup
func SumKey[K comparable, V number.Type](keys []K, items map[K]V) V {
	var sum V = 0
	for _, key := range keys {
		sum += items[key]
	}
	return sum
}

// Product computes the product of number items
func Product[T number.Type](numbers []T) T {
	var product T = 1
	for _, x := range numbers {
		product *= x
	}
	return product
}

// ProductOf comptues the product of mapped number items
func ProductOf[T any, V number.Type](items []T, convert func(T) V) V {
	var product V = 1
	for _, item := range items {
		product *= convert(item)
	}
	return product
}

// ProductIndex combines ProductOf and MapList
func ProductIndex[T number.Type](indexes []int, items []T) T {
	var product T = 1
	numItems := len(items)
	for _, idx := range indexes {
		if 0 <= idx && idx < numItems {
			product *= items[idx]
		}
	}
	return product
}

// ProductKey combines ProductOf and MapLookup
func ProductKey[K comparable, V number.Type](keys []K, items map[K]V) V {
	var product V = 1
	for _, key := range keys {
		product *= items[key]
	}
	return product
}

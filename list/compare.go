package list

import "slices"

// IndexLookup creates a lookup of { item : index } from list.
// This loses data if list items are not unique
func IndexLookup[T comparable](items []T) map[T]int {
	lookup := make(map[T]int)
	for i, item := range items {
		lookup[item] = i
	}
	return lookup
}

// IndexOf gets the index of given item, returns -1 if not in list
func IndexOf[T comparable](items []T, item T) int {
	return slices.Index(items, item)
}

// AllIndexOf gets all indexes of given item
func AllIndexOf[T comparable](items []T, item T) []int {
	indexes := make([]int, 0, len(items))
	for i, x := range items {
		if x == item {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

// Remove removes the first item from list that matches given item
func Remove[T comparable](items []T, item T) ([]T, bool) {
	index := IndexOf(items, item)
	if index < 0 {
		return items, false
	}
	result := slices.Delete(items, index, index+1)
	return result, true
}

// RemoveAll removes all items from list that matches given item
func RemoveAll[T comparable](items []T, item T) []T {
	return slices.DeleteFunc(items, func(x T) bool {
		return x == item
	})
}

// Has checks if any list items are equal to given value
func Has[T comparable](items []T, item T) bool {
	return slices.Contains(items, item)
}

// HasNo checks if no list items are equal to given value
func HasNo[T comparable](items []T, item T) bool {
	return !slices.Contains(items, item)
}

// GetOrDefault returns the first item that matches given item, or returns default value
func GetOrDefault[T comparable](items []T, item T, defaultValue T) T {
	index := IndexOf(items, item)
	if index < 0 {
		return defaultValue
	}
	return items[index]
}

// AllEqual checks if all list items are equal to given value
func AllEqual[T comparable](items []T, item T) bool {
	if len(items) == 0 {
		return false
	}
	for _, x := range items {
		if x != item {
			return false
		}
	}
	return true
}

// AllTrue checks if all list items are true
func AllTrue(items []bool) bool {
	return AllEqual(items, true)
}

// AllFalse checks if all list items are false
func AllFalse(items []bool) bool {
	return AllEqual(items, false)
}

// AnyTrue checks if any list item is true
func AnyTrue(items []bool) bool {
	return Has(items, true)
}

// AnyFalse checks if any list item is false
func AnyFalse(items []bool) bool {
	return Has(items, false)
}

// AllSame checks if all list items are same
func AllSame[T comparable](items []T) bool {
	return len(Tally(items)) == 1
}

// AllSameFunc checks if all list key values are same
func AllSameFunc[T any, K comparable](items []T, keyFn func(T) K) bool {
	return len(TallyFunc(items, keyFn)) == 1
}

// AllUnique checks if all list items are unique
func AllUnique[T comparable](items []T) bool {
	if len(items) == 0 {
		return false
	}
	return len(Tally(items)) == len(items)
}

// AllUniqueFunc checks if all list key values are unique
func AllUniqueFunc[T any, K comparable](items []T, keyFn func(T) K) bool {
	if len(items) == 0 {
		return false
	}
	return len(TallyFunc(items, keyFn)) == len(items)
}

// CountUnique counts the number of unique items in list
func CountUnique[T comparable](items []T) int {
	return len(Tally(items))
}

// CountUniqueFunc counts the number of unique key values in the list, using the key function
func CountUniqueFunc[T any, K comparable](items []T, keyFn func(T) K) int {
	return len(TallyFunc(items, keyFn))
}

// Tally computes the number of occurrence of each item in list
func Tally[T comparable](items []T) map[T]int {
	return TallyFunc(items, func(x T) T { return x })
}

// TallyFunc computes the number of occurrence of each key value in the list, using the key function
func TallyFunc[T any, K comparable](items []T, keyFn func(T) K) map[K]int {
	count := make(map[K]int)
	for _, item := range items {
		count[keyFn(item)] += 1
	}
	return count
}

// Deduplicate removes duplicates from the list, preserving item order
func Deduplicate[T comparable](items []T) []T {
	return DeduplicateFunc(items, func(x T) T { return x })
}

// DeduplicateFunc removes duplicate key values from the list, preserving item order, using the key function
func DeduplicateFunc[T any, K comparable](items []T, keyFn func(T) K) []T {
	unique := make([]T, 0, len(items))
	done := make(map[K]bool)
	for _, item := range items {
		key := keyFn(item)
		if done[key] {
			continue
		}
		unique = append(unique, item)
		done[key] = true
	}
	return unique
}

// Count counts the number of occurrence of item in the list
func Count[T comparable](items []T, item T) int {
	return CountFunc(items, func(x T) bool {
		return x == item
	})
}

// GroupByFunc groups the list items using the key function
func GroupByFunc[T any, K comparable](items []T, keyFn func(T) K) map[K][]T {
	group := make(map[K][]T)
	for _, item := range items {
		key := keyFn(item)
		group[key] = append(group[key], item)
	}
	return group
}

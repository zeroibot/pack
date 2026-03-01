package list

import (
	"slices"

	"github.com/roidaradal/pack/ds"
)

// AllEqual checks if all List items are equal to given value
func AllEqual[T comparable](items ds.List[T], value T) ds.Boolean {
	if len(items) == 0 {
		return false
	}
	for _, item := range items {
		if item != value {
			return false
		}
	}
	return true
}

// AllNotEqual checks if all List items are not equal to given value
func AllNotEqual[T comparable](items ds.List[T], value T) ds.Boolean {
	return ds.Boolean(!slices.Contains(items, value))
}

// AnyEqual checks if any List items are equal to given value
func AnyEqual[T comparable](items ds.List[T], value T) ds.Boolean {
	return ds.Boolean(slices.Contains(items, value))
}

// AllTrue checks if all List items are true
func AllTrue(items ds.List[ds.Boolean]) ds.Boolean {
	return AllEqual(items, true)
}

// AllFalse checks if all List items are false
func AllFalse(items ds.List[ds.Boolean]) ds.Boolean {
	return AllEqual(items, false)
}

// AnyTrue checks if any List item is true
func AnyTrue(items ds.List[ds.Boolean]) ds.Boolean {
	return AnyEqual(items, true)
}

// AnyFalse checks if any List item is false
func AnyFalse(items ds.List[ds.Boolean]) ds.Boolean {
	return AnyEqual(items, false)
}

// AllSame checks if all List items are the same
func AllSame[T comparable](items ds.List[T]) ds.Boolean {
	return len(Tally(items)) == 1
}

// AllUnique checks if all List items are unique
func AllUnique[T comparable](items ds.List[T]) ds.Boolean {
	return len(Tally(items)) == len(items)
}

// CountUnique counts the unique items in the List
func CountUnique[T comparable](items ds.List[T]) ds.Int {
	return ds.Int(len(Tally(items)))
}

// Deduplicate removes duplicates from the List, preserving the order of items
func Deduplicate[T comparable](items ds.List[T]) ds.List[T] {
	unique := ds.NewEmptyList[T](len(items))
	done := make(map[T]bool)
	for _, item := range items {
		if done[item] {
			continue
		}
		unique = append(unique, item)
		done[item] = true
	}
	return unique
}

// Tally computes the number of occurrence of each item in the List
func Tally[T comparable](items ds.List[T]) ds.Map[T, ds.Int] {
	count := make(ds.Map[T, ds.Int])
	for _, item := range items {
		count[item] += 1
	}
	return count
}

// Count counts the number of occurrence of given value in the List
func Count[T comparable](items ds.List[T], value T) ds.Int {
	var count ds.Int = 0
	for _, item := range items {
		if item == value {
			count += 1
		}
	}
	return count
}

// IndexLookup creates a lookup Map of { item : index } from the List.
// This loses data if items are not unique
func IndexLookup[T comparable](items ds.List[T]) ds.Map[T, int] {
	lookup := make(ds.Map[T, int])
	for i, item := range items {
		lookup[item] = i
	}
	return lookup
}

// GroupBy groups List items using the key function
func GroupBy[T any, K comparable](items ds.List[T], keyFn func(T) K) ds.Map[K, ds.List[T]] {
	group := make(ds.Map[K, ds.List[T]])
	for _, item := range items {
		key := keyFn(item)
		group[key] = append(group[key], item)
	}
	return group
}

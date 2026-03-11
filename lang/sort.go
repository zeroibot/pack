package lang

import "cmp"

// SortAscending creates a comparison function for ascending order
func SortAscending[T cmp.Ordered]() func(T, T) int {
	return func(x1, x2 T) int {
		return cmp.Compare(x1, x2)
	}
}

// SortDescending creates a comparison function for descending order
func SortDescending[T cmp.Ordered]() func(T, T) int {
	return func(x1, x2 T) int {
		return cmp.Compare(x2, x1)
	}
}

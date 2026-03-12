package list

import (
	"cmp"
)

// AllGreater checks if all list items are greater than given value
func AllGreater[T cmp.Ordered](items []T, value T) bool {
	return All(items, func(x T) bool { return x > value })
}

// AllGreaterEqual checks if all list items are greater or equal to given value
func AllGreaterEqual[T cmp.Ordered](items []T, value T) bool {
	return All(items, func(x T) bool { return x >= value })
}

// AllLesser checks if all list items are lesser than given value
func AllLesser[T cmp.Ordered](items []T, value T) bool {
	return All(items, func(x T) bool { return x < value })
}

// AllLesserEqual checks if all list items are lesser or equal to given value
func AllLesserEqual[T cmp.Ordered](items []T, value T) bool {
	return All(items, func(x T) bool { return x <= value })
}

// ArgMin finds the index of minimum item on the list.
// Returns -1 if list is empty
func ArgMin[T cmp.Ordered](items []T) int {
	if len(items) == 0 {
		return -1
	}
	idxMin, currMin := 0, items[0]
	for i := 1; i < len(items); i++ {
		if items[i] < currMin {
			idxMin, currMin = i, items[i]
		}
	}
	return idxMin
}

// ArgMax finds the index of maximum item on the list.
// Returns -1 if list is empty
func ArgMax[T cmp.Ordered](items []T) int {
	if len(items) == 0 {
		return -1
	}
	idxMax, currMax := 0, items[0]
	for i := 1; i < len(items); i++ {
		if items[i] > currMax {
			idxMax, currMax = i, items[i]
		}
	}
	return idxMax
}

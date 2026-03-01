package list

import (
	"cmp"
	"slices"

	"github.com/roidaradal/pack/ds"
)

// AllGreater checks if all List items are greater than given value
func AllGreater[T cmp.Ordered](items ds.List[T], value T) ds.Boolean {
	return items.All(func(x T) ds.Boolean {
		return x > value
	})
}

// AllGreaterEqual checks if all List items are greater or equal to given value
func AllGreaterEqual[T cmp.Ordered](items ds.List[T], value T) ds.Boolean {
	return items.All(func(x T) ds.Boolean {
		return x >= value
	})
}

// AllLess checks if all List items are lesser than given value
func AllLess[T cmp.Ordered](items ds.List[T], value T) ds.Boolean {
	return items.All(func(x T) ds.Boolean {
		return x < value
	})
}

// AllLessEqual checks if all List items are lesser or equal to given value
func AllLessEqual[T cmp.Ordered](items ds.List[T], value T) ds.Boolean {
	return items.All(func(x T) ds.Boolean {
		return x <= value
	})
}

// Min finds the minimum item of the List
func Min[T cmp.Ordered](items ds.List[T]) T {
	return slices.Min(items)
}

// Max finds the maximum item of the List
func Max[T cmp.Ordered](items ds.List[T]) T {
	return slices.Max(items)
}

// ArgMin finds the index of the minimum item of the List
func ArgMin[T cmp.Ordered](items ds.List[T]) int {
	index, currMin := 0, items[0]
	for i := 1; i < len(items); i++ {
		if items[i] < currMin {
			index, currMin = i, items[i]
		}
	}
	return index
}

// ArgMax finds the index of the maximum item of the List
func ArgMax[T cmp.Ordered](items ds.List[T]) int {
	index, currMax := 0, items[0]
	for i := 1; i < len(items); i++ {
		if items[i] > currMax {
			index, currMax = i, items[i]
		}
	}
	return index
}

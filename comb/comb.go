// Package comb contains combinatorics functions
package comb

import (
	"iter"

	"github.com/zeroibot/pack/list"
)

// Factorial computes N!
func Factorial(n uint64) uint64 {
	if n == 0 {
		return 1
	}
	return list.Product(list.InclusiveRange(1, n))
}

// Combinations generates the combinations of items taken k at a time in lexicographic order.
func Combinations[T any](items []T, k int) iter.Seq2[uint64, []T] {
	return func(yield func(uint64, []T) bool) {
		n := len(items)
		if k < 0 || k > n {
			return
		}
		if k == 0 {
			yield(0, []T{})
			return
		}

		// indices stores the indices of the current combination.
		// Initial combination: [0, 1, ..., k-1]
		indices := make([]int, k)
		for i := range k {
			indices[i] = i
		}

		var count uint64 = 0
		for {
			// Yield the current combination.
			if !yield(count, list.MapList(indices, items)) {
				return
			}
			count++

			// Find the rightmost index that can be incremented.
			i := k - 1
			for i >= 0 && indices[i] == i+n-k {
				i--
			}

			// If no such index exists, we have generated all combinations.
			if i < 0 {
				return
			}

			// Increment this index and reset all subsequent indices.
			indices[i]++
			for j := i + 1; j < k; j++ {
				indices[j] = indices[j-1] + 1
			}
		}
	}
}

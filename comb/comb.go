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
	n := len(items)
	return RangeCombinations(0, nCr(uint64(n), uint64(k)), items, k)
}

// RangeCombinations generates the combinations of items taken k at a time in lexicographic order for the given range [start, end).
func RangeCombinations[T any](start, end uint64, items []T, k int) iter.Seq2[uint64, []T] {
	return func(yield func(uint64, []T) bool) {
		n := len(items)
		if k < 0 || k > n {
			return
		}

		for i := start; i < end; i++ {
			indices := combinationIndices(uint64(n), uint64(k), i)
			if !yield(i, list.MapList(indices, items)) {
				return
			}
		}
	}
}

// nCr computes the number of combinations of n items taken r at a time.
func nCr(n, r uint64) uint64 {
	if r > n {
		return 0
	}
	if r == 0 || r == n {
		return 1
	}
	if r > n/2 {
		r = n - r
	}
	var res uint64 = 1
	for i := uint64(1); i <= r; i++ {
		res = res * (n - i + 1) / i
	}
	return res
}

// combinationIndices computes the indices of the idx-th combination in lexicographic order.
func combinationIndices(n, k, idx uint64) []int {
	indices := make([]int, k)
	var next int = 0
	for i := uint64(0); i < k; i++ {
		for {
			count := nCr(n-uint64(next)-1, k-i-1)
			if idx < count {
				indices[i] = next
				next++
				break
			}
			idx -= count
			next++
		}
	}
	return indices
}

package comb

import (
	"iter"
	"slices"

	"github.com/zeroibot/pack/list"
)

// Permutations function creates an iterator for all possible permutations in lexicographic order
// This uses the Factoradic System to compute the Permutation based on the index
func Permutations[T any](items []T) iter.Seq2[uint64, []T] {
	return func(yield func(uint64, []T) bool) {
		n := len(items)
		product, factors := computeFactors(n)
		for i := range product {
			indexes := permutationIndexes(factors, i, n)
			tuple := list.MapList(indexes, items)
			if !yield(i, tuple) {
				return
			}
		}
	}
}

// RangePermutations function creates an iterator for permutations for the given range, in lexicographic order.
// This uses the Factoradic System to compute the Permutation based on the index
func RangePermutations[T any](start, end uint64, items []T) iter.Seq2[uint64, []T] {
	return func(yield func(uint64, []T) bool) {
		n := len(items)
		_, factors := computeFactors(n)
		for i := start; i < end; i++ {
			indexes := permutationIndexes(factors, i, n)
			tuple := list.MapList(indexes, items)
			if !yield(i, tuple) {
				return
			}
		}
	}
}

// computeFactors computes 0! to n! used as factors in computing the permutations
func computeFactors(n int) (uint64, []uint64) {
	factors := []uint64{1} // for 0!
	var product uint64 = 1
	for i := 1; i <= n; i++ {
		product *= uint64(i)
		if i < n {
			factors = append(factors, product)
		}
	}
	slices.Reverse(factors)
	return product, factors
}

// permutationIndexes computes the indexes of the ith permutation
func permutationIndexes(factors []uint64, idx uint64, n int) []int {
	digits := make([]int, len(factors))
	current := idx
	for i, factor := range factors {
		d := current / factor
		current = current % factor
		digits[i] = int(d)
	}
	choices := list.Range(0, n)
	indexes := make([]int, len(digits))
	for i, d := range digits {
		indexes[i] = choices[d]
		slices.Delete(choices, d, d+1)
	}
	return indexes
}

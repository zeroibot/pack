package comb

import (
	"iter"

	"github.com/zeroibot/pack/list"
)

// CartesianProduct creates an iterator for the Cartesian product of the given domains
func CartesianProduct[T any](domains ...[]T) iter.Seq2[uint64, []T] {
	return func(yield func(uint64, []T) bool) {
		sizes := list.Map(domains, func(domain []T) uint64 {
			return uint64(len(domain))
		})
		total := list.Product(sizes)
		for i := range total {
			tuple := productDomainCombo(domains, sizes, i)
			if !yield(i, tuple) {
				return
			}
		}
	}
}

// RangeCartesianProduct creates an iterator for the Cartesian product of the given domains from index [start, end)
func RangeCartesianProduct[T any](start, end uint64, domains ...[]T) iter.Seq2[uint64, []T] {
	return func(yield func(uint64, []T) bool) {
		sizes := list.Map(domains, func(domain []T) uint64 {
			return uint64(len(domain))
		})
		for i := start; i < end; i++ {
			tuple := productDomainCombo(domains, sizes, i)
			if !yield(i, tuple) {
				return
			}
		}
	}
}

// productDomainCombo computes the Cartesian product combination from the given domains for the given index
func productDomainCombo[T any](domains [][]T, sizes []uint64, index uint64) []T {
	numSizes := len(sizes)
	indexes := make([]uint64, numSizes)
	for i := range numSizes {
		denom := list.Product(sizes[i+1:])
		num := sizes[i] * denom
		indexes[i] = (index % num) / denom
	}
	numDomains := len(domains)
	combo := make([]T, numDomains)
	for i := range numDomains {
		combo[i] = domains[i][indexes[i]]
	}
	return combo
}

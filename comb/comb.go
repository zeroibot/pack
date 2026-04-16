// Package comb contains combinatorics functions
package comb

import "github.com/zeroibot/pack/list"

// Factorial computes N!
func Factorial(n uint64) uint64 {
	if n == 0 {
		return 1
	}
	return list.Product(list.InclusiveRange(1, n))
}

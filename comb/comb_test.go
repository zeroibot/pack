package comb

import (
	"testing"

	"github.com/zeroibot/tst"
)

func TestFactorial(t *testing.T) {
	pairs := []tst.P1W1[uint64, uint64]{
		{0, 1}, {1, 1}, {2, 2}, {3, 6}, {4, 24}, {5, 120},
		{6, 720}, {7, 5040}, {8, 40320}, {9, 362880}, {10, 3628800},
	}
	tst.AllP1W1(t, pairs, "Factorial", Factorial, tst.AssertEqual)
}

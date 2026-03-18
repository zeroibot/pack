// Package tst contains unit testing functions and TestCase structs
package tst

import (
	"testing"
)

type P1W1[P, W any] struct {
	P1 P
	W1 W
}

// AllP1W1 tests all P1W1 test cases
func AllP1W1[P, W any](t *testing.T, name string, testCases []P1W1[P, W], testFn func(P) W, assert AssertFn[W]) {
	for _, x := range testCases {
		actual := testFn(x.P1)
		assert(t, name, actual, x.W1)
	}
}

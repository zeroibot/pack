package lang

import (
	"testing"

	"github.com/roidaradal/tst"
)

func TestTernary(t *testing.T) {
	testCases := []tst.P3W1[bool, int, int, int]{
		{true, 1, 0, 1},
		{false, 1, 0, 0},
		{true, 0, 1, 0},
		{false, 0, 1, 1},
	}
	tst.AllP3W1(t, testCases, "Ternary", Ternary, tst.AssertEqual)
}

func TestRef(t *testing.T) {
	a, b := 1, 2
	c, d := "c", "d"
	testCases1 := []tst.P1W1[int, int]{
		{a, a},
		{b, b},
	}
	testCases2 := []tst.P1W1[string, string]{
		{c, c},
		{d, d},
	}
	tst.AllP1W1(t, testCases1, "Ref", func(x int) int { return *Ref(x) }, tst.AssertEqual)
	tst.AllP1W1(t, testCases2, "Ref", func(x string) string { return *Ref(x) }, tst.AssertEqual)
}

func TestDeref(t *testing.T) {
	a, b := 1, 2
	testCases1 := []tst.P1W1[*int, int]{
		{&a, a},
		{&b, b},
		{nil, 0},
	}
	c, d := "c", "d"
	testCases2 := []tst.P1W1[*string, string]{
		{&c, c},
		{&d, d},
		{nil, ""},
	}
	tst.AllP1W1(t, testCases1, "Deref", Deref, tst.AssertEqual)
	tst.AllP1W1(t, testCases2, "Deref", Deref, tst.AssertEqual)
}

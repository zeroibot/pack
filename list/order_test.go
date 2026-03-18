package list

import (
	"testing"

	"github.com/roidaradal/tst"
)

func TestOrderFunctions(t *testing.T) {
	type testCase = tst.P2W1[[]int, int, bool]
	items := InclusiveRange(1, 10)

	testCases := []testCase{{items, 0, true}, {items, 5, false}}
	tst.AllP2W1(t, testCases, "AllGreater", AllGreater, tst.AssertEqual)
	testCases = []testCase{{items, 1, true}, {items, 20, false}}
	tst.AllP2W1(t, testCases, "AllGreaterEqual", AllGreaterEqual, tst.AssertEqual)
	testCases = []testCase{{items, 20, true}, {items, 10, false}}
	tst.AllP2W1(t, testCases, "AllLesser", AllLesser, tst.AssertEqual)
	testCases = []testCase{{items, 10, true}, {items, 7, false}}
	tst.AllP2W1(t, testCases, "AllLesserEqual", AllLesserEqual, tst.AssertEqual)
}

func TestArgMinMax(t *testing.T) {
	var empty []int
	items1 := []int{2, 1, 3, 5, 4, 4, 3}
	items2 := []int{1, 2, 3, 4}
	items3 := []int{4, 3, 2, 1}

	testCases := []tst.P1W1[[]int, int]{
		{empty, -1},
		{items1, 1},
		{items2, 0},
		{items3, 3},
	}
	tst.AllP1W1(t, testCases, "ArgMin", ArgMin, tst.AssertEqual)

	testCases = []tst.P1W1[[]int, int]{
		{empty, -1},
		{items1, 3},
		{items2, 3},
		{items3, 0},
	}
	tst.AllP1W1(t, testCases, "ArgMax", ArgMax, tst.AssertEqual)
}

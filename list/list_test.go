package list

import (
	"fmt"
	"slices"
	"testing"

	"github.com/roidaradal/tst"
)

func TestList(t *testing.T) {
	// Range, InclusiveRange
	pairs := [][2][]int{
		{Range(1, 5), []int{1, 2, 3, 4}},
		{InclusiveRange(1, 5), []int{1, 2, 3, 4, 5}},
	}
	tst.All(t, pairs, "Range", tst.AssertListEqual)

	// RepeatedItem
	testCases1 := []tst.P2W1[int, int, []int]{
		{5, 3, []int{5, 5, 5}},
		{3, 5, []int{3, 3, 3, 3, 3}},
	}
	tst.AllP2W1(t, testCases1, "RepeatedItem", RepeatedItem, tst.AssertListEqual)

	// NewEmpty, Len, Cap, LastIndex
	l1 := NewEmpty[string](5)
	l2 := []string{"a", "b", "c"}
	testCases2 := []tst.P1W1[[]string, int]{{l1, 0}, {l2, 3}}
	tst.AllP1W1(t, testCases2, "Len", Len, tst.AssertEqual)
	testCases2 = []tst.P1W1[[]string, int]{{l1, 5}, {l2, 3}}
	tst.AllP1W1(t, testCases2, "Cap", Cap, tst.AssertEqual)
	testCases2 = []tst.P1W1[[]string, int]{{l1, -1}, {l2, 2}}
	tst.AllP1W1(t, testCases2, "LastIndex", LastIndex, tst.AssertEqual)

	// IsEmpty, NotEmpty
	tst.AssertEqual(t, "IsEmpty", IsEmpty(l1), true)
	tst.AssertEqual(t, "NotEmpty", NotEmpty(l2), true)

	// Copy
	l3 := Copy(l2)
	tst.AssertListEqual(t, "Copy", l3, l2)

	// String
	l2[0] = "x"
	l3[1] = "r"
	tst.AssertEqual(t, "List.String", fmt.Sprintf("%v", l2), "[x b c]")
	tst.AssertEqual(t, "List.String", fmt.Sprintf("%v", l3), "[a r c]")
}

func TestListRandom(t *testing.T) {
	// GetRandom and MustGetRandom
	l1 := NewEmpty[int](3) // empty
	for range 5 {
		item, ok := GetRandom(l1)
		tst.AssertEqualAnd(t, "EmptyList.GetRandom", item, 0, ok, false)
	}
	l := InclusiveRange(1, 100)
	for range 100 {
		value, ok := GetRandom(l)
		tst.AssertTrue(t, "GetRandom", ok && 1 <= value && value <= 100)
		value = MustGetRandom(l)
		tst.AssertTrue(t, "MustGetRandom", 1 <= value && value <= 100)
	}

	// Shuffle
	l2 := []int{1, 2, 3, 4, 5, 6, 7}
	l3 := Copy(l2)
	Shuffle(l3)
	tst.AssertFalse(t, "Shuffle", slices.Equal(l2, l3))

	defer tst.AssertPanic(t, "MustGetRandom")
	MustGetRandom(l1) // should panic (empty list)
}

func TestListMethods(t *testing.T) {
	// ToAny
	items := []int{1, 2, 3}
	anyItems := ToAny(items)
	actualString, wantString := fmt.Sprintf("%v", anyItems), "[1 2 3]"
	tst.AssertEqual(t, "ToAny", actualString, wantString)

	// IndexFunc
	items = []int{1, 2, 3, 4, 1, 2, 4, 2, 5, 3}
	testCases1 := []tst.P2W1[[]int, func(int) bool, int]{
		{items, func(x int) bool { return x == 3 }, 2},
		{items, func(x int) bool { return x == 5 }, 8},
	}
	tst.AllP2W1(t, testCases1, "IndexFunc", IndexFunc, tst.AssertEqual)

	// AllIndexFunc
	testCases2 := []tst.P2W1[[]int, func(int) bool, []int]{
		{items, func(x int) bool { return x == 2 }, []int{1, 5, 7}},
		{items, func(x int) bool { return x%2 == 0 }, []int{1, 3, 5, 6, 7}},
	}
	tst.AllP2W1(t, testCases2, "AllIndexFunc", AllIndexFunc, tst.AssertListEqual)

	// RemoveFunc
	items = []int{1, 2, 1, 2, 3, 2}
	items2 := Copy(items)
	wantList := []int{1, 1, 2, 3, 2}
	items2, ok := RemoveFunc(items2, func(x int) bool { return x == 2 })
	tst.AssertListEqualAnd(t, "RemoveFunc", items2, wantList, ok, true)
	items2, ok = RemoveFunc(items2, func(x int) bool { return x == 4 })
	tst.AssertListEqualAnd(t, "RemoveFunc", items2, wantList, ok, false)

	// RemoveAllFunc
	items2 = Copy(items)
	items2 = RemoveAllFunc(items2, func(x int) bool { return x == 2 })
	tst.AssertListEqual(t, "RemoveAllFunc", items2, []int{1, 1, 3})

	// GetFuncOrDefault
	items = []int{1, 2, 3}
	defaultValue := 69
	actual := GetFuncOrDefault(items, func(x int) bool { return x == 3 }, defaultValue)
	tst.AssertEqual(t, "GetFuncOrDefault", actual, 3)
	actual = GetFuncOrDefault(items, func(x int) bool { return x == 4 }, defaultValue)
	tst.AssertEqual(t, "GetFuncOrDefault", actual, defaultValue)

	// Last
	items = []int{1, 2, 3}
	testCases3 := []tst.P2W2[[]int, int, int, bool]{
		{items, 1, 3, true},
		{items, 3, 1, true},
		{items, 0, 0, false},
		{items, 4, 0, false},
	}
	tst.AllP2W2(t, testCases3, "Last", Last, tst.AssertEqual[int], tst.AssertEqual[bool])

	// MustLast
	testCases4 := []tst.P2W1[[]int, int, int]{
		{items, 1, 3},
		{items, 3, 1},
		{items, 2, 2},
	}
	tst.AllP2W1(t, testCases4, "MustLast", MustLast, tst.AssertEqual)

	defer tst.AssertPanic(t, "MustLast")
	MustLast(items, 4) // should panic
}

func TestListCheck(t *testing.T) {
	var empty []int
	items := []int{1, 2, 3, 4, 5, 6}
	fn1 := func(x int) bool { return x%2 == 0 && x%3 == 0 }
	fn2 := func(x int) bool { return x > 10 }
	fn3 := func(x int) bool { return x <= 10 }

	// Any, NotAny
	type testCase = tst.P2W1[[]int, func(int) bool, bool]
	testCases := []testCase{
		{items, fn1, true}, {items, fn2, false},
	}
	tst.AllP2W1(t, testCases, "Any", Any, tst.AssertEqual)

	testCases = tst.Convert(testCases, func(tc testCase) testCase {
		return testCase{P1: tc.P1, P2: tc.P2, W1: !tc.W1}
	})
	tst.AllP2W1(t, testCases, "NotAny", NotAny, tst.AssertEqual)

	// All
	testCases = []testCase{
		{empty, fn1, false}, {items, fn1, false}, {items, fn3, true},
	}
	tst.AllP2W1(t, testCases, "All", All, tst.AssertEqual)

	// AnyIndexed, NotAnyIndexed
	type testCase2 = tst.P2W1[[]int, func(int, int) bool, bool]
	fn4 := func(i, x int) bool { return i >= 0 && x%2 == 0 && x%3 == 0 }
	fn5 := func(i, x int) bool { return i > 10 && x > 10 }
	fn6 := func(i, x int) bool { return i < 10 && x <= 10 }

	testCases2 := []testCase2{
		{items, fn4, true}, {items, fn5, false},
	}
	tst.AllP2W1(t, testCases2, "AnyIndexed", AnyIndexed, tst.AssertEqual)
	testCases2 = tst.Convert(testCases2, func(tc testCase2) testCase2 {
		return testCase2{P1: tc.P1, P2: tc.P2, W1: !tc.W1}
	})
	tst.AllP2W1(t, testCases2, "NotAnyIndexed", NotAnyIndexed, tst.AssertEqual)

	// AllIndexed
	testCases2 = []testCase2{
		{empty, fn4, false}, {items, fn4, false}, {items, fn6, true},
	}
	tst.AllP2W1(t, testCases2, "AllIndexed", AllIndexed, tst.AssertEqual)
}

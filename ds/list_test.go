package ds

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/roidaradal/tst"
)

func TestList(t *testing.T) {
	l1 := NewList[int](5)
	l2 := List[string]{"a", "b", "c"}
	// Len, Cap, LastIndex
	testCases := []Tuple3[string, int, int]{
		{"Len", l1.Len(), 0},
		{"Cap", l1.Cap(), 5},
		{"LastIndex", l1.LastIndex(), -1},
		{"Len", l2.Len(), 3},
		{"Cap", l2.Cap(), 3},
		{"LastIndex", l2.LastIndex(), 2},
	}
	for _, x := range testCases {
		name, actual, want := x.Unpack()
		tst.AssertEqual(t, name, actual, want)
	}
	// IsEmpty, NotEmpty, Copy
	l3 := l2.Copy()
	tst.AssertTrue(t, "List.IsEmpty", l1.IsEmpty())
	tst.AssertTrue(t, "List.NotEmpty", l2.NotEmpty())
	tst.AssertListEqual(t, "List.Copy", l3, l2)
	// String
	l2[0] = "x"
	l3[1] = "r"
	tst.AssertEqual(t, "List.String", fmt.Sprintf("%v", l2), "[x b c]")
	tst.AssertEqual(t, "List.String", fmt.Sprintf("%v", l3), "[a r c]")
}

func TestListMethods(t *testing.T) {
	// ToAny
	items := List[int]{1, 2, 3}
	anyItems := items.ToAny()
	tst.AssertEqual(t, "List.ToAny", fmt.Sprintf("%v", anyItems), "[1 2 3]")
	// IndexFunc
	items = List[int]{1, 2, 3, 4, 1, 2, 4, 2, 5, 3}
	tst.AssertEqual(t, "List.IndexFunc", items.IndexFunc(func(x int) bool { return x == 3 }), 2)
	tst.AssertEqual(t, "List.Index", items.IndexFunc(func(x int) bool { return x == 5 }), 8)
	// AllIndexFunc
	actualList := items.AllIndexFunc(func(x int) bool { return x == 2 })
	tst.AssertListEqual(t, "List.AllIndexFunc", actualList, List[int]{1, 5, 7})
	actualList = items.AllIndexFunc(func(x int) bool { return x%2 == 0 })
	tst.AssertListEqual(t, "List.AllIndexFunc", actualList, List[int]{1, 3, 5, 6, 7})
	// RemoveFunc
	items = List[int]{1, 2, 1, 2, 3, 2}
	items2 := items.Copy()
	items2, ok := items2.RemoveFunc(func(x int) bool { return x == 2 })
	wantList := List[int]{1, 1, 2, 3, 2}
	tst.AssertListEqualAnd(t, "List.RemoveFunc", items2, wantList, ok, true)
	items2, ok = items2.RemoveFunc(func(x int) bool { return x == 4 })
	tst.AssertListEqualAnd(t, "List.RemoveFunc", items2, wantList, ok, false)
	// RemoveAllFunc
	items2 = items.Copy()
	items2 = items2.RemoveAllFunc(func(x int) bool { return x == 2 })
	tst.AssertListEqual(t, "List.RemoveAllFunc", items2, List[int]{1, 1, 3})
	// Get
	items = List[int]{1, 2, 3}
	option := items.Get(1)
	tst.AssertEqualAnd(t, "List.Get", option.Value(), 2, option.IsNil(), false)
	option = items.Get(3)
	tst.AssertEqualAnd(t, "List.Get", option.Value(), 0, option.NotNil(), false)
	option = items.Get(-1)
	tst.AssertEqualAnd(t, "List.Get", option.Value(), 0, option.IsNil(), true)
	// GetFuncOrDefault
	defaultValue := 69
	actual := items.GetFuncOrDefault(func(x int) bool { return x == 3 }, defaultValue)
	tst.AssertEqual(t, "List.GetFuncOrDefault", actual, 3)
	actual = items.GetFuncOrDefault(func(x int) bool { return x == 4 }, defaultValue)
	tst.AssertEqual(t, "List.GetFuncOrDefault", actual, defaultValue)
	// Last
	option = items.Last(1)
	tst.AssertEqualAnd(t, "List.Last", option.Value(), 3, option.IsNil(), false)
	option = items.Last(3)
	tst.AssertEqualAnd(t, "List.Last", option.Value(), 1, option.IsNil(), false)
	option = items.Last(0)
	tst.AssertEqualAnd(t, "List.Last", option.Value(), 0, option.IsNil(), true)
	option = items.Last(4)
	tst.AssertEqualAnd(t, "List.Last", option.Value(), 0, option.IsNil(), true)
	// MustLast
	tst.AssertEqual(t, "List.MustLast", items.MustLast(1), 3)
	tst.AssertEqual(t, "List.MustLast", items.MustLast(3), 1)

	defer tst.AssertPanic(t, "List.MustLast")
	items.MustLast(4) // should panic
}

func TestListRandom(t *testing.T) {
	l0 := NewList[int](3) // empty
	l1 := NewInclusiveRange(1, 100).ToList()
	// GetRandom and MustGetRandom
	for range 5 {
		tst.AssertTrue(t, "EmptyList.GetRandom.IsNil", l0.GetRandom().IsNil())
	}
	for range 100 {
		item := l1.GetRandom()
		value := item.Value()
		tst.AssertTrue(t, "List.GetRandom", !item.IsNil() && (1 <= value && value <= 100))
		value = l1.MustGetRandom()
		tst.AssertTrue(t, "List.MustGetRandom", 1 <= value && value <= 100)
	}
	// Shuffle
	l2 := List[int]{1, 2, 3, 4, 5, 6, 7}
	l3 := l2.Copy()
	l3.Shuffle()
	tst.AssertFalse(t, "List.Shuffle", slices.Equal(l2, l3))

	defer tst.AssertPanic(t, "MustGetRandom")
	l0.MustGetRandom() // should panic (empty list)
}

func TestListCheck(t *testing.T) {
	type testCase = tst.P2W1[List[int], func(int) bool, bool]
	empty := List[int]{}
	items := List[int]{1, 2, 3, 4, 5, 6}
	fn1 := func(x int) bool { return x%2 == 0 && x%3 == 0 }
	fn2 := func(x int) bool { return x > 10 }
	fn3 := func(x int) bool { return x <= 10 }
	// Any, NotAny
	testCases := []testCase{
		{items, fn1, true}, {items, fn2, false},
	}
	tst.AllP2W1(t, testCases, "List.Any", List[int].Any, tst.AssertEqual)
	testCases = tst.FlipP2W1(testCases)
	tst.AllP2W1(t, testCases, "List.NotAny", List[int].NotAny, tst.AssertEqual)
	// All
	testCases = []testCase{
		{empty, fn1, false}, {items, fn1, false}, {items, fn3, true},
	}
	tst.AllP2W1(t, testCases, "List.All", List[int].All, tst.AssertEqual)
	// AnyIndexed, NotAnyIndexed
	type testCase2 = tst.P2W1[List[int], func(int, int) bool, bool]
	fn4 := func(i, x int) bool { return i >= 0 && x%2 == 0 && x%3 == 0 }
	fn5 := func(i, x int) bool { return i > 10 && x > 10 }
	fn6 := func(i, x int) bool { return i < 10 && x <= 10 }
	testCases2 := []testCase2{
		{items, fn4, true}, {items, fn5, false},
	}
	tst.AllP2W1(t, testCases2, "List.AnyIndexed", List[int].AnyIndexed, tst.AssertEqual)
	testCases2 = tst.FlipP2W1(testCases2)
	tst.AllP2W1(t, testCases2, "List.NotAnyIndexed", List[int].NotAnyIndexed, tst.AssertEqual)
	// AllIndexed
	testCases2 = []testCase2{
		{empty, fn4, false}, {items, fn4, false}, {items, fn6, true},
	}
	tst.AllP2W1(t, testCases2, "List.AllIndexed", List[int].AllIndexed, tst.AssertEqual)
}

func TestListFn(t *testing.T) {
	// MapList
	items := List[string]{" ", "A", "B", "C", "D", "E"}
	indexes := []int{2, 5, 1, 4, 0, 3, 1, 2}
	result := items.MapList(indexes)
	tst.AssertEqual(t, "List.MapList", strings.Join(result, ""), "BEAD CAB")
	// Filter, CountFunc
	numbers := List[int]{1, 2, 3, 4, 5, 6, 7}
	fn1 := func(x int) bool { return x%2 == 0 }
	fn2 := func(x int) bool { return x > 10 }
	fn3 := func(x int) bool { return x <= 10 }

	testCases := []tst.P2W1[List[int], func(int) bool, List[int]]{
		{numbers, fn1, List[int]{2, 4, 6}},
		{numbers, fn2, List[int]{}},
		{numbers, fn3, numbers},
	}
	testCases2 := []tst.P2W1[List[int], func(int) bool, int]{
		{numbers, fn1, 3},
		{numbers, fn2, 0},
		{numbers, fn3, len(numbers)},
	}
	tst.AllP2W1(t, testCases, "List.Filter", List[int].Filter, tst.AssertListEqual)
	tst.AllP2W1(t, testCases2, "List.CountFunc", List[int].CountFunc, tst.AssertEqual)
	// FilterIndexed
	actual := numbers.FilterIndexed(func(i, x int) bool { return x%2 == 0 || i%3 == 0 })
	tst.AssertListEqual(t, "List.FilterIndexed", actual, List[int]{1, 2, 4, 6, 7})
	// Reduce
	actualSum := numbers.Reduce(0, func(result, item int) int {
		return result + item
	})
	tst.AssertEqual(t, "List.Reduce", actualSum, 28)
	// Apply
	actual = numbers.Apply(func(x int) int { return x * 2 })
	tst.AssertListEqual(t, "List.Apply", actual, List[int]{2, 4, 6, 8, 10, 12, 14})
}

func TestNumList(t *testing.T) {
	// ToList
	n := NumList[int]{1, 2, 3, 4, 5, 6}
	l := n.ToList()
	tst.AssertEqual(t, "NumList.ToList.Len", l.Len(), 6)
	// Sum, Product
	tst.AssertEqual(t, "NumList.Sum", n.Sum(), 21)
	tst.AssertEqual(t, "NumList.Product", n.Product(), 720)
}

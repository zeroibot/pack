package list

import (
	"testing"

	"github.com/roidaradal/tst"
)

func TestCompareAllAny(t *testing.T) {
	type testCase1 = tst.P2W1[[]int, int, bool]
	type testCase2 = tst.P1W1[[]bool, bool]

	var ints0 []int
	ints1 := []int{1, 1, 1, 1}
	ints2 := []int{1, 2, 3, 4}
	ints3 := []int{2, 2, 1, 2}
	var bools0 []bool
	bools1 := []bool{true, true, true}
	bools2 := []bool{false, false}
	bools3 := []bool{false, true, false}

	// AllEqual
	testCases1 := []testCase1{
		{ints0, 5, false},
		{ints1, 1, true},
		{ints2, 3, false},
		{ints3, 2, false},
	}
	tst.AllP2W1(t, testCases1, "AllEqual", AllEqual, tst.AssertEqual)

	// AllTrue, All False
	testCases2 := []testCase2{
		{bools0, false},
		{bools1, true},
		{bools3, false},
	}
	tst.AllP1W1(t, testCases2, "AllTrue", AllTrue, tst.AssertEqual)

	testCases2 = []tst.P1W1[[]bool, bool]{
		{bools0, false},
		{bools2, true},
		{bools3, false},
	}
	tst.AllP1W1(t, testCases2, "AllFalse", AllFalse, tst.AssertEqual)

	// Has, HasNo
	testCases1 = []testCase1{
		{ints0, 1, false},
		{ints1, 1, true},
		{ints3, 1, true},
		{ints2, 5, false},
	}
	tst.AllP2W1(t, testCases1, "Has", Has, tst.AssertEqual)

	testCases1 = tst.Convert(testCases1, func(tc testCase1) testCase1 {
		return testCase1{P1: tc.P1, P2: tc.P2, W1: !tc.W1}
	})
	tst.AllP2W1(t, testCases1, "HasNo", HasNo, tst.AssertEqual)

	// AnyTrue, AnyFalse
	testCases2 = []testCase2{
		{bools1, true},
		{bools2, false},
		{bools3, true},
	}
	tst.AllP1W1(t, testCases2, "AnyTrue", AnyTrue, tst.AssertEqual)

	testCases2 = []testCase2{
		{bools1, false},
		{bools2, true},
		{bools3, true},
	}
	tst.AllP1W1(t, testCases2, "AnyFalse", AnyFalse, tst.AssertEqual)
}

func TestIndexFunctions(t *testing.T) {
	// IndexLookup
	items := []string{" ", "A", "B", "C"}
	wantMap := map[string]int{" ": 0, "A": 1, "B": 2, "C": 3}
	tst.AssertMapEqual(t, "IndexLookup", IndexLookup(items), wantMap)

	// IndexOf
	testCases := []tst.P2W1[[]string, string, int]{
		{items, "A", 1},
		{items, "C", 3},
		{items, "X", -1},
	}
	tst.AllP2W1(t, testCases, "IndexOf", IndexOf, tst.AssertEqual)

	// AllIndexOf
	ints := []int{1, 2, 3, 1, 2, 3, 1}
	testCases2 := []tst.P2W1[[]int, int, []int]{
		{ints, 1, []int{0, 3, 6}},
		{ints, 3, []int{2, 5}},
		{ints, 69, []int{}},
	}
	tst.AllP2W1(t, testCases2, "AllIndexOf", AllIndexOf, tst.AssertListEqual)

	// GetOrDefault
	defaultValue := 69
	testCases3 := []tst.P3W1[[]int, int, int, int]{
		{ints, 3, defaultValue, 3},
		{ints, 4, defaultValue, defaultValue},
	}
	tst.AllP3W1(t, testCases3, "GetOrDefault", GetOrDefault, tst.AssertEqual)

	// Remove
	ints2 := Copy(ints)
	wantInts := []int{1, 2, 1, 2, 3, 1}
	ints2, ok := Remove(ints2, 3)
	tst.AssertListEqualAnd(t, "Remove", ints2, wantInts, ok, true)
	ints2, ok = Remove(ints2, 69)
	tst.AssertListEqualAnd(t, "Remove", ints2, wantInts, ok, false)

	// RemoveAll
	ints2 = Copy(ints)
	wantInts = []int{2, 3, 2, 3}
	ints2 = RemoveAll(ints2, 1)
	tst.AssertListEqual(t, "RemoveAll", ints2, wantInts)
	ints2 = RemoveAll(ints2, 5)
	tst.AssertListEqual(t, "RemoveAll", ints2, wantInts)
}

func TestTallyFunctions(t *testing.T) {
	// Count
	chars := []byte{'a', 'b', 'a', 'a', 'c', 'd', 'b', 'c'}
	countCases := []tst.P2W1[[]byte, byte, int]{
		{chars, 'a', 3},
		{chars, 'd', 1},
		{chars, 'x', 0},
	}
	tst.AllP2W1(t, countCases, "Count", Count, tst.AssertEqual)

	// GroupByFunc
	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := map[string][]int{
		"odd":  {1, 3, 5, 7, 9},
		"even": {2, 4, 6, 8},
	}
	oddOrEven := func(x int) string {
		if x%2 == 0 {
			return "even"
		}
		return "odd"
	}
	tst.AssertDeepEqual(t, "GroupByFunc", GroupByFunc(ints, oddOrEven), want)

	// Tally
	wantTally := map[byte]int{'a': 3, 'b': 2, 'c': 2, 'd': 1}
	tst.AssertMapEqual(t, "Tally", Tally(chars), wantTally)

	// TallyFunc
	wantCounts := map[string]int{
		"odd":  5,
		"even": 4,
	}
	tst.AssertMapEqual(t, "TallyFunc", TallyFunc(ints, oddOrEven), wantCounts)
}

func TestUniqueFunctions(t *testing.T) {
	type person struct {
		name  string
		kind  int
		score int
	}
	getName := func(p person) string { return p.name }
	getKind := func(p person) int { return p.kind }
	getScore := func(p person) int { return p.score }

	var ints0 []int
	ints1 := []int{1, 1, 1, 1}
	ints2 := []int{1, 2, 3, 4}
	ints3 := []int{1, 2, 1, 3, 2, 4}
	var persons0 []person
	a, b, c := person{"A", 1, 10}, person{"B", 1, 15}, person{"C", 1, 10}
	persons1 := []person{a, b, c}

	// AllSame
	testCases1 := []tst.P1W1[[]int, bool]{
		{ints0, false}, {ints1, true}, {ints2, false}, {ints3, false},
	}
	tst.AllP1W1(t, testCases1, "AllSame", AllSame, tst.AssertEqual)

	// AllUnique
	testCases1 = []tst.P1W1[[]int, bool]{
		{ints0, false}, {ints1, false}, {ints2, true}, {ints3, false},
	}
	tst.AllP1W1(t, testCases1, "AllUnique", AllUnique, tst.AssertEqual)

	// AllSameFunc
	testCases2 := []tst.P2W1[[]person, func(person) string, bool]{
		{persons0, getName, false},
		{persons1, getName, false},
	}
	testCases3 := []tst.P2W1[[]person, func(person) int, bool]{
		{persons1, getKind, true},
		{persons1, getScore, false},
	}
	tst.AllP2W1(t, testCases2, "AllSameFunc", AllSameFunc, tst.AssertEqual)
	tst.AllP2W1(t, testCases3, "AllSameFunc", AllSameFunc, tst.AssertEqual)

	// AllUniqueFunc
	testCases2 = []tst.P2W1[[]person, func(person) string, bool]{
		{persons0, getName, false},
		{persons1, getName, true},
	}
	testCases3 = []tst.P2W1[[]person, func(person) int, bool]{
		{persons1, getKind, false},
		{persons1, getScore, false},
	}
	tst.AllP2W1(t, testCases2, "AllUniqueFunc", AllUniqueFunc, tst.AssertEqual)
	tst.AllP2W1(t, testCases3, "AllUniqueFunc", AllUniqueFunc, tst.AssertEqual)

	// CountUnique
	testCases4 := []tst.P1W1[[]int, int]{
		{ints0, 0}, {ints1, 1}, {ints2, 4}, {ints3, 4},
	}
	tst.AllP1W1(t, testCases4, "CountUnique", CountUnique, tst.AssertEqual)

	// CountUniqueFunc
	testCases5 := []tst.P2W1[[]person, func(person) string, int]{
		{persons0, getName, 0}, {persons1, getName, 3},
	}
	testCases6 := []tst.P2W1[[]person, func(person) int, int]{
		{persons1, getKind, 1}, {persons1, getScore, 2},
	}
	tst.AllP2W1(t, testCases5, "CountUniqueFunc", CountUniqueFunc, tst.AssertEqual)
	tst.AllP2W1(t, testCases6, "CountUniqueFunc", CountUniqueFunc, tst.AssertEqual)

	// Deduplicate
	testCases7 := []tst.P1W1[[]int, []int]{
		{ints0, ints0},
		{ints1, []int{1}},
		{ints2, ints2},
		{ints3, ints2},
	}
	tst.AllP1W1(t, testCases7, "Deduplicate", Deduplicate, tst.AssertListEqual)

	// DeduplicateFunc
	testCases8 := []tst.P2W1[[]person, func(person) int, []person]{
		{persons1, getKind, []person{a}},
		{persons1, getScore, []person{a, b}},
	}
	tst.AllP2W1(t, testCases8, "DeduplicateFunc", DeduplicateFunc, tst.AssertListEqual)
	tst.AssertListEqual(t, "DeduplicateFunc", DeduplicateFunc(persons1, getName), []person{a, b, c})
}

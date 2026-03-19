package ds

import (
	"slices"
	"testing"

	"github.com/roidaradal/tst"
)

func TestNewSet(t *testing.T) {
	type person struct {
		name  string
		score int
	}
	// NewSet, NewSetFrom, NewSetFunc
	set1 := NewSet[int]()
	items := []int{1, 2, 3, 1, 2, 4, 5, 1, 3}
	set2 := NewSetFrom(items)
	persons := []person{
		{"Jack", 3},
		{"John", 5},
		{"Jill", 3},
		{"Gus", 4},
	}
	set3 := NewSetFunc(persons, func(p person) int { return p.score })
	// Len, IsEmpty, NotEmpty, Items
	type testCase = tst.P1W1[Set[int], bool]
	lenCases := []tst.P1W1[Set[int], int]{{set1, 0}, {set2, 5}, {set3, 3}}
	emptyCases := []testCase{{set1, true}, {set2, false}, {set3, false}}
	notEmptyCases := tst.FlipP1W1(emptyCases)
	itemCases := []tst.P1W1[Set[int], List[int]]{
		{set1, List[int]{}}, {set2, List[int]{1, 2, 3, 4, 5}}, {set3, List[int]{3, 4, 5}},
	}
	tst.AllP1W1(t, lenCases, "Set.Len", Set[int].Len, tst.AssertEqual)
	tst.AllP1W1(t, emptyCases, "Set.IsEmpty", Set[int].IsEmpty, tst.AssertEqual)
	tst.AllP1W1(t, notEmptyCases, "Set.NotEmpty", Set[int].NotEmpty, tst.AssertEqual)
	for _, x := range itemCases {
		set, want := x.P1, x.W1
		actual := set.Items()
		slices.Sort(actual)
		tst.AssertListEqual(t, "Set.Items", actual, want)
	}
	tst.AssertEqual(t, "Set.String", set1.String(), "{}")
}

func TestSetFunctions(t *testing.T) {
	m := NewSet[int]()
	m.Add(1)
	m.Add(1, 2, 3)
	m.Add(2, 3, 4)
	// Has and HasNo
	type testCase = tst.P2W1[Set[int], int, bool]
	testCases := []testCase{
		{m, 1, true}, {m, 0, false}, {m, 3, true},
		{m, 5, false}, {m, 4, true},
	}
	tst.AllP2W1(t, testCases, "Set.Has", Set[int].Has, tst.AssertEqual)
	testCases = tst.FlipP2W1(testCases)
	tst.AllP2W1(t, testCases, "Set.HasNo", Set[int].HasNo, tst.AssertEqual)
	// Add and Delete
	m.Add(5)
	tst.AssertEqual(t, "Set.Add.Has", m.Has(5), true)
	m.Delete(5)
	tst.AssertEqual(t, "Set.Delete.HasNo", m.HasNo(5), true)

	// Copy and Clear
	mc := m.Copy()
	copyItems := mc.Items()
	slices.Sort(copyItems)
	tst.AssertListEqual(t, "Set.Copy.Items", copyItems, []int{1, 2, 3, 4})

	mc.Clear()
	tst.AssertEqual(t, "Set.Clear.Len", mc.Len(), 0)
	tst.AssertEqual(t, "Set.Len", m.Len(), 4)
}

func TestSetMethods(t *testing.T) {
	s1 := NewSetFrom([]int{1, 2, 3, 4})
	s2 := NewSetFrom([]int{3, 4, 5, 6})
	s3 := NewSetFrom([]int{6, 7, 8, 9})
	s4 := NewSetFrom([]int{1, 2, 3, 4, 4, 3, 2, 1})

	testCases := []Tuple3[string, List[int], List[int]]{
		{"Union", s1.Union(s2).Items(), List[int]{1, 2, 3, 4, 5, 6}},
		{"Union", s1.Union(s4).Items(), List[int]{1, 2, 3, 4}},
		{"Intersection", s1.Intersection(s2).Items(), List[int]{3, 4}},
		{"Intersection", s1.Intersection(s3).Items(), List[int]{}},
		{"Difference", s1.Difference(s2).Items(), List[int]{1, 2}},
		{"Difference", s1.Difference(s4).Items(), List[int]{}},
	}
	for _, x := range testCases {
		name, actual, want := x.Unpack()
		slices.Sort(actual)
		tst.AssertListEqual(t, name, actual, want)
	}
	tst.AssertTrue(t, "HasIntersection", s1.HasIntersection(s2))
	tst.AssertTrue(t, "HasNoIntersection", s1.HasNoIntersection(s3))
	tst.AssertTrue(t, "HasDifference", s1.HasDifference(s2))
	tst.AssertTrue(t, "HasNoDifference", s1.HasNoDifference(s4))
}

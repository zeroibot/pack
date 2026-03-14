package ds

import (
	"slices"
	"testing"
)

func TestNewSet(t *testing.T) {
	type testCase[T any] struct {
		name           string
		actual1, want1 int
		actual2, want2 bool
		actual3, want3 bool
		actual4, want4 List[T]
	}
	type person struct {
		name  string
		score int
	}
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
	testCases := []testCase[int]{
		{
			"NewSet",
			set1.Len(), 0,
			set1.IsEmpty(), true,
			set1.NotEmpty(), false,
			set1.Items(), List[int]{},
		},
		{
			"NewSetFrom",
			set2.Len(), 5,
			set2.IsEmpty(), false,
			set2.NotEmpty(), true,
			set2.Items(), List[int]{1, 2, 3, 4, 5},
		},
		{
			"NewSetFunc",
			set3.Len(), 3,
			set3.IsEmpty(), false,
			set3.NotEmpty(), true,
			set3.Items(), List[int]{3, 4, 5},
		},
	}
	for _, x := range testCases {
		if x.actual1 != x.want1 {
			t.Errorf("%s.Len = %d, want %d", x.name, x.actual1, x.want1)
		}
		if x.actual2 != x.want2 {
			t.Errorf("%s.NotEmpty = %v, want %v", x.name, x.actual2, x.want2)
		}
		if x.actual3 != x.want3 {
			t.Errorf("%s.IsEmpty = %v, want %v", x.name, x.actual3, x.want3)
		}
		slices.Sort(x.actual4)
		if slices.Equal(x.actual4, x.want4) == false {
			t.Errorf("%s.Items = %v, want %v", x.name, x.actual4, x.want4)
		}
	}
	actual, want := set1.String(), "{}"
	if actual != want {
		t.Errorf("Set.String = %q, want %q", actual, want)
	}
}

func TestSetFunctions(t *testing.T) {
	m := NewSet[int]()
	m.Add(1)
	m.Add(1, 2, 3)
	m.Add(2, 3, 4)
	// Has and HasNo
	testCases := []Tuple2[int, bool]{
		{1, true},
		{0, false},
		{3, true},
		{5, false},
		{4, true},
	}
	for _, x := range testCases {
		item, want := x.Unpack()
		actual := m.Has(item)
		if actual != want {
			t.Errorf("Set.Has(%d) = %v, want %v", item, actual, want)
		}
		want = !want
		actual = m.HasNo(item)
		if actual != want {
			t.Errorf("Set.HasNo(%d) = %v, want %v", item, actual, want)
		}
	}
	// Add and Delete
	m.Add(5)
	if m.Has(5) != true {
		t.Errorf("Set.Add.Has(%d) = %v, want true", 5, m.Has(5))
	}
	m.Delete(5)
	if m.HasNo(5) != true {
		t.Errorf("Set.Delete.HasNo(%d) = %v, want true", 5, m.HasNo(5))
	}

	// Copy and Clear
	mc := m.Copy()
	copyItems := mc.Items()
	slices.Sort(copyItems)
	want := []int{1, 2, 3, 4}
	if slices.Equal(copyItems, want) == false {
		t.Errorf("Set.Copy.Items = %v, want %v", copyItems, want)
	}
	mc.Clear()
	if mc.Len() != 0 {
		t.Errorf("Set.Clear.Len = %d, want 0", mc.Len())
	}
	if m.Len() != 4 {
		t.Errorf("Set.Len = %d, want 4", m.Len())
	}
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
		if slices.Equal(actual, want) == false {
			t.Errorf("Set.%s = %v, want %v", name, actual, want)
		}
	}

	actual := s1.HasIntersection(s2)
	if actual != true {
		t.Errorf("Set.HasIntersection = %v, want true", actual)
	}
	actual = s1.HasNoIntersection(s3)
	if actual != true {
		t.Errorf("Set.HasNoIntersection = %v, want true", actual)
	}
	actual = s1.HasDifference(s2)
	if actual != true {
		t.Errorf("Set.HasDifference = %v, want true", actual)
	}
	actual = s1.HasNoDifference(s4)
	if actual != true {
		t.Errorf("Set.HasNoDifference = %v, want true", actual)
	}
}

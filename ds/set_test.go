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

func TestSetOperations(t *testing.T) {
	// TODO: Test Copy
	// TODO: Test Clear
	// TODO: Test Add
	// TODO: Test Has
	// TODO: Test HasNo
	// TODO: Test Delete
	// TODO: Test Union
	// TODO: Test Intersection
	// TODO: Test Difference
	// TODO: Test HasIntersection, HasNoIntersection
	// TODO: Test HasDifference, HasNoDifference
}

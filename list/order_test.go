package list

import "testing"

func TestOrderFunctions(t *testing.T) {
	type testCase struct {
		name   string
		want   bool
		actual bool
	}
	items := InclusiveRange(1, 10)
	testCases := []testCase{
		{"AllGreater", true, AllGreater(items, 0)},
		{"AllGreater", false, AllGreater(items, 5)},
		{"AllGreaterEqual", true, AllGreaterEqual(items, 1)},
		{"AllGreaterEqual", false, AllGreaterEqual(items, 20)},
		{"AllLesser", true, AllLesser(items, 20)},
		{"AllLesser", false, AllLesser(items, 10)},
		{"AllLesserEqual", true, AllLesserEqual(items, 10)},
		{"AllLesserEqual", false, AllLesserEqual(items, 7)},
	}
	for _, x := range testCases {
		if x.want != x.actual {
			t.Errorf("%s: want %t, got %t", x.name, x.want, x.actual)
		}
	}
}

func TestArgMinMax(t *testing.T) {
	// Empty List
	var empty []int
	actual := ArgMin(empty)
	if actual != -1 {
		t.Errorf("ArgMin() = %d, want -1", actual)
	}
	actual = ArgMax(empty)
	if actual != -1 {
		t.Errorf("ArgMax() = %d, want -1", actual)
	}
	// ArgMin, ArgMax
	items := []int{2, 1, 3, 5, 4, 4, 3}
	actual, want := ArgMin(items), 1
	if actual != want {
		t.Errorf("ArgMin() = %d, want %d", actual, want)
	}
	actual, want = ArgMax(items), 3
	if actual != want {
		t.Errorf("ArgMax() = %d, want %d", actual, want)
	}
}

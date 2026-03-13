package list

import "testing"

func TestOrderFunctions(t *testing.T) {
	// TODO: AllGreater, AllGreaterEqual
	// TODO: AllLesser, AllLesserEqual
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

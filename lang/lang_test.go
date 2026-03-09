package lang

import (
	"testing"
)

func TestTernary(t *testing.T) {
	type testCase struct {
		condition  bool
		valueTrue  int
		valueFalse int
		want       int
	}
	testCases := []testCase{
		{true, 1, 0, 1},
		{false, 1, 0, 0},
		{true, 0, 1, 0},
		{false, 0, 1, 1},
	}
	for _, x := range testCases {
		actual := Ternary(x.condition, x.valueTrue, x.valueFalse)
		if actual != x.want {
			t.Errorf("Ternary(%t, %d, %d) = %d, want %d", x.condition, x.valueTrue, x.valueFalse, actual, x.want)
		}
	}
}

func TestRef(t *testing.T) {
	a, b := 1, 2
	c, d := "c", "d"
	testCases1 := []int{a, b}
	testCases2 := []string{c, d}
	for _, x := range testCases1 {
		ref := Ref(x)
		actual := *ref
		if actual != x {
			t.Errorf("Ref(%d) = %d, want %d", x, actual, x)
		}
	}
	for _, x := range testCases2 {
		ref := Ref(x)
		actual := *ref
		if actual != x {
			t.Errorf("Ref(%s) = %s, want %s", x, actual, x)
		}
	}
}

func TestDeref(t *testing.T) {
	type testCase[T any] struct {
		ref  *T
		want T
	}
	a, b := 1, 2
	testCases1 := []testCase[int]{
		{&a, a},
		{&b, b},
		{nil, 0},
	}
	c, d := "c", "d"
	testCases2 := []testCase[string]{
		{&c, c},
		{&d, d},
		{nil, ""},
	}
	for _, x := range testCases1 {
		actual := Deref(x.ref)
		if actual != x.want {
			t.Errorf("Deref(%v) = %d, want %d", x.ref, actual, x.want)
		}
	}
	for _, x := range testCases2 {
		actual := Deref(x.ref)
		if actual != x.want {
			t.Errorf("Deref(%v) = %s, want %s", x.ref, actual, x.want)
		}
	}
}

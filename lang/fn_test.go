package lang

import "testing"

func TestIdentity(t *testing.T) {
	fn1 := Identity[int]
	actual, want := fn1(5), 5
	if actual != want {
		t.Errorf("Identity(%d) = %d, want %d", want, actual, want)
	}
	fn2 := Identity[string]
	actual2, want2 := fn2("a"), "a"
	if actual2 != want2 {
		t.Errorf("Identity(%q) = %q, want %q", want2, actual2, want2)
	}
}

func TestIsEqual(t *testing.T) {
	fn1 := IsEqual(5)
	fn2 := IsEqual("a")
	testCases := [][2]bool{
		{fn1(5), true},
		{fn2("a"), true},
		{fn1(6), false},
		{fn2("b"), false},
	}
	for _, x := range testCases {
		actual, want := x[0], x[1]
		if actual != want {
			t.Errorf("IsEqual = %v, want %v", actual, want)
		}
	}
}

func TestNotEqual(t *testing.T) {
	fn1 := NotEqual(5)
	fn2 := NotEqual("a")
	testCases := [][2]bool{
		{fn1(5), false},
		{fn2("a"), false},
		{fn1(6), true},
		{fn2("b"), true},
	}
	for _, x := range testCases {
		actual, want := x[0], x[1]
		if actual != want {
			t.Errorf("NotEqual = %v, want %v", actual, want)
		}
	}
}

func TestIsGreater(t *testing.T) {
	fn := IsGreater(5)
	testCases := [][2]bool{
		{fn(4), false},
		{fn(5), false},
		{fn(6), true},
	}
	for _, x := range testCases {
		actual, want := x[0], x[1]
		if actual != want {
			t.Errorf("IsGreater = %v, want %v", actual, want)
		}
	}
}

func TestIsGreaterEqual(t *testing.T) {
	fn := IsGreaterEqual(5)
	testCases := [][2]bool{
		{fn(4), false},
		{fn(5), true},
		{fn(6), true},
	}
	for _, x := range testCases {
		actual, want := x[0], x[1]
		if actual != want {
			t.Errorf("IsGreaterEqual = %v, want %v", actual, want)
		}
	}
}

func TestIsLesser(t *testing.T) {
	fn := IsLesser(5)
	testCases := [][2]bool{
		{fn(4), true},
		{fn(5), false},
		{fn(6), false},
	}
	for _, x := range testCases {
		actual, want := x[0], x[1]
		if actual != want {
			t.Errorf("IsGreater = %v, want %v", actual, want)
		}
	}
}

func TestIsLesserEqual(t *testing.T) {
	fn := IsLesserEqual(5)
	testCases := [][2]bool{
		{fn(4), true},
		{fn(5), true},
		{fn(6), false},
	}
	for _, x := range testCases {
		actual, want := x[0], x[1]
		if actual != want {
			t.Errorf("IsGreater = %v, want %v", actual, want)
		}
	}
}

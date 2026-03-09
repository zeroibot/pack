package number

import "testing"

func TestAbs(t *testing.T) {
	testCases1 := []testCase[int, int]{
		{-5, 5},
		{0, 0},
		{3, 3},
	}
	testCases2 := []testCase[float64, float64]{
		{-5.5, 5.5},
		{0, 0},
		{3.25, 3.25},
	}
	for _, x := range testCases1 {
		actual := Abs(x.input)
		if actual != x.want {
			t.Errorf("Abs(%v) = %d, want %d", x.input, actual, x.want)
		}
	}
	for _, x := range testCases2 {
		actual := Abs(x.input)
		if actual != x.want {
			t.Errorf("Abs(%v) = %f, want %f", x.input, actual, x.want)
		}
	}
}

func TestCeilInt(t *testing.T) {
	testCases := []testCase[float64, int]{
		{6.95, 7},
		{3.14, 4},
		{1.0001, 2},
		{5, 5},
		{0, 0},
		{-1.92, -1},
		{-2.05, -2},
	}
	for _, x := range testCases {
		actual := CeilInt(x.input)
		if actual != x.want {
			t.Errorf("CeilInt(%v) = %d, want %d", x.input, actual, x.want)
		}
	}
}

func TestFloorInt(t *testing.T) {
	testCases := []testCase[float64, int]{
		{6.95, 6},
		{3.14, 3},
		{1.0001, 1},
		{5, 5},
		{0, 0},
		{-1.92, -2},
		{-2.05, -3},
	}
	for _, x := range testCases {
		actual := FloorInt(x.input)
		if actual != x.want {
			t.Errorf("FloorInt(%v) = %d, want %d", x.input, actual, x.want)
		}
	}
}

func TestRoundInt(t *testing.T) {
	testCases := []testCase[float64, int]{
		{6.95, 7},
		{3.14, 3},
		{1.0001, 1},
		{5, 5},
		{0, 0},
		{-1.92, -2},
		{-2.05, -2},
		{3.5, 4},
		{4.5, 5},
		{0.5, 1},
	}
	for _, x := range testCases {
		actual := RoundInt(x.input)
		if actual != x.want {
			t.Errorf("RoundInt(%v) = %d, want %d", x.input, actual, x.want)
		}
	}
}

func TestRoundToEvenInt(t *testing.T) {
	testCases := []testCase[float64, int]{
		{6.95, 7},
		{3.14, 3},
		{1.0001, 1},
		{5, 5},
		{0, 0},
		{-1.92, -2},
		{-2.05, -2},
		{3.5, 4},
		{4.5, 4},
		{0.5, 0},
		{1.5, 2},
		{1.49, 1},
		{1.51, 2},
	}
	for _, x := range testCases {
		actual := RoundToEvenInt(x.input)
		if actual != x.want {
			t.Errorf("RoundToEvenInt(%v) = %d, want %d", x.input, actual, x.want)
		}
	}
}

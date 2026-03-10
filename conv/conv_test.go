package conv

import (
	"slices"
	"testing"
)

type testCase[T, V any] struct {
	input T
	want  V
}

func TestAnyToString(t *testing.T) {
	testCases := []testCase[any, string]{
		{true, "true"},
		{false, "false"},
		{1, "1"},
		{-1, "-1"},
		{0, "0"},
		{1.25, "1.25"},
		{nil, "<nil>"},
	}
	for _, x := range testCases {
		actual := AnyToString(x.input)
		if actual != x.want {
			t.Errorf("AnyToString(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
}

func TestAnyToStringList(t *testing.T) {
	type testCase[T any] struct {
		items []T
		want  []string
	}
	testCases := []testCase[int]{
		{[]int{1, 2, 3}, []string{"1", "2", "3"}},
		{[]int{4, 5, 6, 7}, []string{"4", "5", "6", "7"}},
	}
	for _, x := range testCases {
		actual := AnyToStringList(x.items)
		if !slices.Equal(actual, x.want) {
			t.Errorf("AnyToStringList(%v) = %v; want %v", x.items, actual, x.want)
		}
	}
}

func TestBoolToInt(t *testing.T) {
	testCases := []testCase[bool, int]{
		{true, 1},
		{false, 0},
	}
	for _, x := range testCases {
		actual := BoolToInt(x.input)
		if actual != x.want {
			t.Errorf("BoolToInt(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
}

func TestBoolToFloat(t *testing.T) {
	testCases := []testCase[bool, uint]{
		{true, 1},
		{false, 0},
	}
	for _, x := range testCases {
		actual := BoolToUint(x.input)
		if actual != x.want {
			t.Errorf("BoolToUint(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
}

func TestBoolToString(t *testing.T) {
	testCases := []testCase[bool, string]{
		{true, "true"},
		{false, "false"},
	}
	for _, x := range testCases {
		actual := BoolToString(x.input)
		if actual != x.want {
			t.Errorf("BoolToString(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
}

func TestFloatToInt(t *testing.T) {
	testCases1 := []testCase[float32, int]{
		{-6.77, -6},
		{-1.33, -1},
		{-0.99, 0},
		{0.0, 0},
		{0.25, 0},
		{1.33, 1},
		{2.5, 2},
		{3.67, 3},
		{4.99, 4},
		{5.0, 5},
	}
	testCases2 := []testCase[float64, int]{
		{-6.77, -6},
		{-1.33, -1},
		{-0.99, 0},
		{0.0, 0},
		{0.25, 0},
		{1.33, 1},
		{2.5, 2},
		{3.67, 3},
		{4.99, 4},
		{5.0, 5},
	}
	for _, x := range testCases1 {
		actual := FloatToInt(x.input)
		if actual != x.want {
			t.Errorf("FloatToInt(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
	for _, x := range testCases2 {
		actual := FloatToInt(x.input)
		if actual != x.want {
			t.Errorf("FloatToInt(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
}

func TestFloatToUint(t *testing.T) {
	testCases1 := []testCase[float32, uint]{
		{-6.77, 0},
		{-1.33, 0},
		{-0.99, 0},
		{0.0, 0},
		{0.25, 0},
		{1.33, 1},
		{2.5, 2},
		{3.67, 3},
		{4.99, 4},
		{5.0, 5},
	}
	testCases2 := []testCase[float64, uint]{
		{-6.77, 0},
		{-1.33, 0},
		{-0.99, 0},
		{0.0, 0},
		{0.25, 0},
		{1.33, 1},
		{2.5, 2},
		{3.67, 3},
		{4.99, 4},
		{5.0, 5},
	}
	for _, x := range testCases1 {
		actual := FloatToUint(x.input)
		if actual != x.want {
			t.Errorf("FloatToUint(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
	for _, x := range testCases2 {
		actual := FloatToUint(x.input)
		if actual != x.want {
			t.Errorf("FloatToUint(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
}

func TestFloatToString(t *testing.T) {
	// Note: fmt.Sprintf uses 6 decimal places by default
	testCases1 := []testCase[float32, string]{
		{-6.77, "-6.770000"},
		{-1.33, "-1.330000"},
		{-0.9999, "-0.999900"},
		{0.0, "0.000000"},
		{0.25, "0.250000"},
		{1.33333, "1.333330"},
		{2.5, "2.500000"},
		{3.67, "3.670000"},
		{4.99, "4.990000"},
		{5.0, "5.000000"},
	}
	testCases2 := []testCase[float64, string]{
		{-6.77, "-6.770000"},
		{-1.33, "-1.330000"},
		{-0.9999, "-0.999900"},
		{0.0, "0.000000"},
		{0.25, "0.250000"},
		{1.33333, "1.333330"},
		{2.5, "2.500000"},
		{3.67, "3.670000"},
		{4.99, "4.990000"},
		{5.0, "5.000000"},
	}

	for _, x := range testCases1 {
		actual := FloatToString(x.input)
		if actual != x.want {
			t.Errorf("FloatToString(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
	for _, x := range testCases2 {
		actual := FloatToString(x.input)
		if actual != x.want {
			t.Errorf("FloatToString(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
}

func TestIntToBool(t *testing.T) {
	testCases1 := []testCase[int, bool]{
		{-5, true},
		{0, false},
		{1, true},
		{999, true},
	}
	testCases2 := []testCase[uint, bool]{
		{0, false},
		{1, true},
		{5, true},
		{999, true},
	}
	for _, x := range testCases1 {
		actual := IntToBool(x.input)
		if actual != x.want {
			t.Errorf("IntToBool(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
	for _, x := range testCases2 {
		actual := IntToBool(x.input)
		if actual != x.want {
			t.Errorf("IntToBool(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
}

func TestIntToFloat(t *testing.T) {
	testCases1 := []testCase[int, float64]{
		{-5, -5.0},
		{0, 0.0},
		{1, 1.0},
		{999, 999.0},
	}
	testCases2 := []testCase[uint, float64]{
		{0, 0.0},
		{1, 1.0},
		{5, 5.0},
		{999, 999.0},
	}
	for _, x := range testCases1 {
		actual := IntToFloat(x.input)
		if actual != x.want {
			t.Errorf("IntToFloat(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
	for _, x := range testCases2 {
		actual := IntToFloat(x.input)
		if actual != x.want {
			t.Errorf("IntToFloat(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
}

func TestIntToString(t *testing.T) {
	testCases1 := []testCase[int, string]{
		{-5, "-5"},
		{0, "0"},
		{1, "1"},
		{999, "999"},
	}
	testCases2 := []testCase[uint, string]{
		{0, "0"},
		{1, "1"},
		{5, "5"},
		{999, "999"},
	}
	for _, x := range testCases1 {
		actual := IntToString(x.input)
		if actual != x.want {
			t.Errorf("IntToBool(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
	for _, x := range testCases2 {
		actual := IntToString(x.input)
		if actual != x.want {
			t.Errorf("IntToBool(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
}

func TestIntToUint(t *testing.T) {
	testCases := []testCase[int, uint]{
		{1, 1},
		{67, 67},
		{0, 0},
		{-5, 0},
		{-100, 0},
	}
	for _, x := range testCases {
		actual := IntToUint(x.input)
		if actual != x.want {
			t.Errorf("IntToUint(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
}

func TestUintToInt(t *testing.T) {
	testCases := []testCase[uint, int]{
		{1, 1},
		{67, 67},
		{0, 0},
	}
	for _, x := range testCases {
		actual := UintToInt(x.input)
		if actual != x.want {
			t.Errorf("IntToUint(%v) = %v; want %v", x.input, actual, x.want)
		}
	}
}

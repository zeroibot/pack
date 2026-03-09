package number

import "testing"

type testCase[T, V any] struct {
	input T
	want  V
}

func TestParseInt(t *testing.T) {
	testCases := []testCase[string, int]{
		{"0", 0},
		{"1", 1},
		{"-1", -1},
		{"55", 55},
		{"1A", 0},
		{"number", 0},
	}
	for _, x := range testCases {
		actual := ParseInt(x.input)
		if actual != x.want {
			t.Errorf("ParseInt(%v) = %v, want %v", x.input, actual, x.want)
		}
	}
}

func TestParseUint(t *testing.T) {
	testCases := []testCase[string, uint]{
		{"0", 0},
		{"1", 1},
		{"-1", 0},
		{"-33", 0},
		{"55", 55},
		{"1A", 0},
		{"number", 0},
	}
	for _, x := range testCases {
		actual := ParseUint(x.input)
		if actual != x.want {
			t.Errorf("ParseUint(%v) = %v, want %v", x.input, actual, x.want)
		}
	}
}

func TestParseFloat(t *testing.T) {
	testCases := []testCase[string, float64]{
		{"0", 0},
		{"1", 1},
		{"-1", -1},
		{"55", 55},
		{"1A", 0},
		{"number", 0},
		{"3.1415", 3.1415},
		{"1.69", 1.69},
		{"-3.33", -3.33},
	}
	for _, x := range testCases {
		actual := ParseFloat(x.input)
		if actual != x.want {
			t.Errorf("ParseFloat(%v) = %v, want %v", x.input, actual, x.want)
		}
	}
}

package conv

import (
	"testing"

	"github.com/roidaradal/tst"
)

func TestAnyToString(t *testing.T) {
	testCases := []tst.P1W1[any, string]{
		{true, "true"},
		{false, "false"},
		{1, "1"},
		{-1, "-1"},
		{0, "0"},
		{1.25, "1.25"},
		{nil, "<nil>"},
	}
	tst.AllP1W1(t, testCases, "AnyToString", AnyToString, tst.AssertEqual)
}

func TestAnyToStringList(t *testing.T) {
	testCases := []tst.P1W1[[]int, []string]{
		{[]int{1, 2, 3}, []string{"1", "2", "3"}},
		{[]int{4, 5, 6, 7}, []string{"4", "5", "6", "7"}},
	}
	tst.AllP1W1(t, testCases, "AnyToStringList", AnyToStringList, tst.AssertListEqual)
}

func TestBoolToInt(t *testing.T) {
	testCases := []tst.P1W1[bool, int]{
		{true, 1},
		{false, 0},
	}
	tst.AllP1W1(t, testCases, "BoolToInt", BoolToInt, tst.AssertEqual)
}

func TestBoolToFloat(t *testing.T) {
	testCases := []tst.P1W1[bool, uint]{
		{true, 1},
		{false, 0},
	}
	tst.AllP1W1(t, testCases, "BoolToUint", BoolToUint, tst.AssertEqual)
}

func TestBoolToString(t *testing.T) {
	testCases := []tst.P1W1[bool, string]{
		{true, "true"},
		{false, "false"},
	}
	tst.AllP1W1(t, testCases, "BoolToString", BoolToString, tst.AssertEqual)
}

func TestFloatToInt(t *testing.T) {
	testCases1 := []tst.P1W1[float32, int]{
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
	testCases2 := []tst.P1W1[float64, int]{
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
	tst.AllP1W1(t, testCases1, "FloatToInt", FloatToInt, tst.AssertEqual)
	tst.AllP1W1(t, testCases2, "FloatToInt", FloatToInt, tst.AssertEqual)
}

func TestFloatToUint(t *testing.T) {
	testCases1 := []tst.P1W1[float32, uint]{
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
	testCases2 := []tst.P1W1[float64, uint]{
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
	tst.AllP1W1(t, testCases1, "FloatToUint", FloatToUint, tst.AssertEqual)
	tst.AllP1W1(t, testCases2, "FloatToUint", FloatToUint, tst.AssertEqual)
}

func TestFloatToString(t *testing.T) {
	// Note: fmt.Sprintf uses 6 decimal places by default
	testCases1 := []tst.P1W1[float32, string]{
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
	testCases2 := []tst.P1W1[float64, string]{
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
	tst.AllP1W1(t, testCases1, "FloatToString", FloatToString, tst.AssertEqual)
	tst.AllP1W1(t, testCases2, "FloatToString", FloatToString, tst.AssertEqual)
}

func TestIntToBool(t *testing.T) {
	testCases1 := []tst.P1W1[int, bool]{
		{-5, true},
		{0, false},
		{1, true},
		{999, true},
	}
	testCases2 := []tst.P1W1[uint, bool]{
		{0, false},
		{1, true},
		{5, true},
		{999, true},
	}
	tst.AllP1W1(t, testCases1, "IntToBool", IntToBool, tst.AssertEqual)
	tst.AllP1W1(t, testCases2, "IntToBool", IntToBool, tst.AssertEqual)
}

func TestIntToFloat(t *testing.T) {
	testCases1 := []tst.P1W1[int, float64]{
		{-5, -5.0},
		{0, 0.0},
		{1, 1.0},
		{999, 999.0},
	}
	testCases2 := []tst.P1W1[uint, float64]{
		{0, 0.0},
		{1, 1.0},
		{5, 5.0},
		{999, 999.0},
	}
	tst.AllP1W1(t, testCases1, "IntToFloat", IntToFloat, tst.AssertEqual)
	tst.AllP1W1(t, testCases2, "IntToFloat", IntToFloat, tst.AssertEqual)
}

func TestIntToString(t *testing.T) {
	testCases1 := []tst.P1W1[int, string]{
		{-5, "-5"},
		{0, "0"},
		{1, "1"},
		{999, "999"},
	}
	testCases2 := []tst.P1W1[uint, string]{
		{0, "0"},
		{1, "1"},
		{5, "5"},
		{999, "999"},
	}
	tst.AllP1W1(t, testCases1, "IntToString", IntToString, tst.AssertEqual)
	tst.AllP1W1(t, testCases2, "IntToString", IntToString, tst.AssertEqual)
}

func TestIntToUint(t *testing.T) {
	testCases := []tst.P1W1[int, uint]{
		{1, 1},
		{67, 67},
		{0, 0},
		{-5, 0},
		{-100, 0},
	}
	tst.AllP1W1(t, testCases, "IntToUint", IntToUint, tst.AssertEqual)
}

func TestUintToInt(t *testing.T) {
	testCases := []tst.P1W1[uint, int]{
		{1, 1},
		{67, 67},
		{0, 0},
	}
	tst.AllP1W1(t, testCases, "UintToInt", UintToInt, tst.AssertEqual)
}

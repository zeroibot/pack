package number

import (
	"testing"

	"github.com/roidaradal/tst"
)

func TestCommaString(t *testing.T) {
	testCases := []tst.P1W1[int, string]{
		{123456789, "123,456,789"},
		{1234567890, "1,234,567,890"},
		{23456789, "23,456,789"},
		{6789, "6,789"},
		{789, "789"},
		{89, "89"},
		{9, "9"},
		{-123, "-123"},
		{-23, "-23"},
		{-3, "-3"},
		{-12345, "-12,345"},
		{-1234, "-1,234"},
	}
	tst.AllP1W1(t, testCases, "CommaString", CommaString, tst.AssertEqual)
}

func TestCommaDecimalString(t *testing.T) {
	testCases := []tst.P2W1[float64, uint, string]{
		{1333.14, 1, "1,333.1"},
		{123.14, 2, "123.14"},
		{123456.14, 0, "123,456"},
		{34567.1625, 2, "34,567.16"},
		{77.5, 4, "77.5000"},
		{1235.75, 2, "1,235.75"},
		{12345.75, 1, "12,345.8"},
		{999.95, 0, "1,000"},
		{999.95, 1, "1,000.0"},
		{-34567.23, 3, "-34,567.230"},
	}
	tst.AllP2W1(t, testCases, "CommaDecimalString", CommaDecimalString, tst.AssertEqual)
}

func TestDecimalString(t *testing.T) {
	testCases := []tst.P2W1[float64, uint, string]{
		{3.14, 1, "3.1"},
		{3.14, 2, "3.14"},
		{3.14, 0, "3"},
		{3.1625, 2, "3.16"},
		{7.5, 4, "7.5000"},
		{3.75, 2, "3.75"},
		{3.75, 1, "3.8"},
		{3.95, 0, "4"},
		{3.95, 1, "4.0"},
	}
	tst.AllP2W1(t, testCases, "DecimalString", DecimalString, tst.AssertEqual)
}

func TestParseInt(t *testing.T) {
	testCases := []tst.P1W1[string, int]{
		{"0", 0},
		{"1", 1},
		{"-1", -1},
		{"55", 55},
		{"1A", 0},
		{"number", 0},
	}
	tst.AllP1W1(t, testCases, "ParseInt", ParseInt, tst.AssertEqual)
}

func TestParseUint(t *testing.T) {
	testCases := []tst.P1W1[string, uint]{
		{"0", 0},
		{"1", 1},
		{"-1", 0},
		{"-33", 0},
		{"55", 55},
		{"1A", 0},
		{"number", 0},
	}
	tst.AllP1W1(t, testCases, "ParseUint", ParseUint, tst.AssertEqual)
}

func TestParseFloat(t *testing.T) {
	testCases := []tst.P1W1[string, float64]{
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
	tst.AllP1W1(t, testCases, "ParseFloat", ParseFloat, tst.AssertEqual)
}

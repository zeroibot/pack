package lang

import (
	"testing"

	"github.com/roidaradal/tst"
)

func TestIdentity(t *testing.T) {
	tst.AssertEqual(t, "Identity", Identity[int](5), 5)
	tst.AssertEqual(t, "Identity", Identity[string]("a"), "a")
}

func TestIsEqual(t *testing.T) {
	testCases1 := []tst.P1W1[int, bool]{
		{5, true},
		{6, false},
	}
	testCases2 := []tst.P1W1[string, bool]{
		{"a", true},
		{"b", false},
	}
	tst.AllP1W1(t, testCases1, "IsEqual", IsEqual(5), tst.AssertEqual)
	tst.AllP1W1(t, testCases2, "IsEqual", IsEqual("a"), tst.AssertEqual)
}

func TestNotEqual(t *testing.T) {
	testCases1 := []tst.P1W1[int, bool]{
		{5, false},
		{6, true},
	}
	testCases2 := []tst.P1W1[string, bool]{
		{"a", false},
		{"b", true},
	}
	tst.AllP1W1(t, testCases1, "NotEqual", NotEqual(5), tst.AssertEqual)
	tst.AllP1W1(t, testCases2, "NotEqual", NotEqual("a"), tst.AssertEqual)
}

func TestIsGreater(t *testing.T) {
	testCases := []tst.P1W1[int, bool]{
		{4, false},
		{5, false},
		{6, true},
	}
	tst.AllP1W1(t, testCases, "IsGreater", IsGreater(5), tst.AssertEqual)
}

func TestIsGreaterEqual(t *testing.T) {
	testCases := []tst.P1W1[int, bool]{
		{4, false},
		{5, true},
		{6, true},
	}
	tst.AllP1W1(t, testCases, "IsGreaterEqual", IsGreaterEqual(5), tst.AssertEqual)
}

func TestIsLesser(t *testing.T) {
	testCases := []tst.P1W1[int, bool]{
		{4, true},
		{5, false},
		{6, false},
	}
	tst.AllP1W1(t, testCases, "IsLesser", IsLesser(5), tst.AssertEqual)
}

func TestIsLesserEqual(t *testing.T) {
	testCases := []tst.P1W1[int, bool]{
		{4, true},
		{5, true},
		{6, false},
	}
	tst.AllP1W1(t, testCases, "IsLesserEqual", IsLesserEqual(5), tst.AssertEqual)
}

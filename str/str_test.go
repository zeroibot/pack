package str

import (
	"testing"

	"github.com/roidaradal/tst"
)

func TestLen(t *testing.T) {
	testCases := []tst.P1W1[string, int]{
		{"", 0},
		{"abc", 3},
		{"12345", 5},
		{"X", 1},
		{"XY", 2},
	}
	tst.AllP1W1(t, testCases, "Len", Len, tst.AssertEqual)
}

func TestIsEmpty(t *testing.T) {
	type testCase = tst.P1W1[string, bool]
	testCases1 := []testCase{
		{"", true},
		{"123", false},
		{"a", false},
	}
	testCases2 := tst.Convert(testCases1, func(tc testCase) testCase {
		return testCase{P1: tc.P1, W1: !tc.W1}
	})
	tst.AllP1W1(t, testCases1, "IsEmpty", IsEmpty, tst.AssertEqual)
	tst.AllP1W1(t, testCases2, "NotEmpty", NotEmpty, tst.AssertEqual)
}

func TestGuard(t *testing.T) {
	testCases := []tst.P2W1[string, string, string]{
		{"", "default", "default"},
		{"abc", "default", "abc"},
		{"", "unknown", "unknown"},
		{"def", "unknown", "def"},
	}
	tst.AllP2W1(t, testCases, "Guard", Guard, tst.AssertEqual)
}

func TestWrap(t *testing.T) {
	testCases := []tst.P2W1[string, string, string]{
		{"hello", "()", "(hello)"},
		{"hello", "(", "(hello"},
		{"world", "<>", "<world>"},
		{"column", "``", "`column`"},
		{"default", "", "default"},
	}
	tst.AllP2W1(t, testCases, "Wrap", Wrap, tst.AssertEqual)
}

func TestWrapList(t *testing.T) {
	testCases := []tst.P2W1[[]string, string, string]{
		{[]string{"a", "b", "c"}, "[]", "[a, b, c]"},
		{[]string{"a => 1", "b => 2"}, "{}", "{a => 1, b => 2}"},
		{[]string{"1", "2", "3"}, "", "1, 2, 3"},
	}
	tst.AllP2W1(t, testCases, "WrapList", WrapList, tst.AssertEqual)
}

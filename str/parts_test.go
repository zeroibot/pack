package str

import (
	"testing"

	"github.com/roidaradal/tst"
)

func TestCleanSplit(t *testing.T) {
	testCases := []tst.P2W1[string, string, []string]{
		{"apple | banana | cat", "|", []string{"apple", "banana", "cat"}},
		{" apple|banana|cat ", "|", []string{"apple", "banana", "cat"}},
		{"1, 2, 3, 4", ",", []string{"1", "2", "3", "4"}},
		{"1, 2, 3, 4", "|", []string{"1, 2, 3, 4"}},
		{"", ".", []string{""}},
	}
	tst.AllP2W1(t, testCases, "CleanSplit", CleanSplit, tst.AssertListEqual)
}

func TestCleanSplitN(t *testing.T) {
	testCases := []tst.P3W1[string, string, int, []string]{
		{"apple | banana | cat", "|", 3, []string{"apple", "banana", "cat"}},
		{" apple|banana|cat ", "|", 3, []string{"apple", "banana", "cat"}},
		{"apple | banana | cat", "|", 2, []string{"apple", "banana | cat"}},
		{" apple|banana|cat ", "|", 1, []string{"apple|banana|cat"}},
		{"1, 2, 3, 4", ",", 5, []string{"1", "2", "3", "4"}},
		{"1, 2, 3, 4", ",", 3, []string{"1", "2", "3, 4"}},
		{"1, 2, 3, 4", "|", 4, []string{"1, 2, 3, 4"}},
		{"", ".", 1, []string{""}},
	}
	tst.AllP3W1(t, testCases, "CleanSplitN", CleanSplitN, tst.AssertListEqual)
}

func TestSpaceSplit(t *testing.T) {
	testCases := []tst.P1W1[string, []string]{
		{"a\tb\tc", []string{"a", "b", "c"}},
		{"a\nb\nc", []string{"a", "b", "c"}},
		{"a  b  c", []string{"a", "b", "c"}},
	}
	tst.AllP1W1(t, testCases, "SpaceSplit", SpaceSplit, tst.AssertListEqual)
}

func TestJoin(t *testing.T) {
	testCases := []tst.P2W1[string, []string, string]{
		{",", []string{"a", "b", "c"}, "a,b,c"},
		{"\t", []string{"a", "b", "c"}, "a\tb\tc"},
		{"+", []string{"1", "2", "3"}, "1+2+3"},
		{"x", []string{}, ""},
	}
	joinFn := func(glue string, parts []string) string {
		return Join(glue, parts...)
	}
	tst.AllP2W1(t, testCases, "Join", joinFn, tst.AssertEqual)
}

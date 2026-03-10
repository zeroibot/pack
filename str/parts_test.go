package str

import (
	"slices"
	"testing"
)

func TestCleanSplit(t *testing.T) {
	type testCase struct {
		text string
		sep  string
		want []string
	}
	testCases := []testCase{
		{"apple | banana | cat", "|", []string{"apple", "banana", "cat"}},
		{" apple|banana|cat ", "|", []string{"apple", "banana", "cat"}},
		{"1, 2, 3, 4", ",", []string{"1", "2", "3", "4"}},
		{"1, 2, 3, 4", "|", []string{"1, 2, 3, 4"}},
		{"", ".", []string{""}},
	}
	for _, x := range testCases {
		actual := CleanSplit(x.text, x.sep)
		if !slices.Equal(actual, x.want) {
			t.Errorf("CleanSplit(%q, %q) = %v, want %v", x.text, x.sep, actual, x.want)
		}
	}
}

func TestCleanSplitN(t *testing.T) {
	type testCase struct {
		text string
		sep  string
		n    int
		want []string
	}
	testCases := []testCase{
		{"apple | banana | cat", "|", 3, []string{"apple", "banana", "cat"}},
		{" apple|banana|cat ", "|", 3, []string{"apple", "banana", "cat"}},
		{"apple | banana | cat", "|", 2, []string{"apple", "banana | cat"}},
		{" apple|banana|cat ", "|", 1, []string{"apple|banana|cat"}},
		{"1, 2, 3, 4", ",", 5, []string{"1", "2", "3", "4"}},
		{"1, 2, 3, 4", ",", 3, []string{"1", "2", "3, 4"}},
		{"1, 2, 3, 4", "|", 4, []string{"1, 2, 3, 4"}},
		{"", ".", 1, []string{""}},
	}
	for _, x := range testCases {
		actual := CleanSplitN(x.text, x.sep, x.n)
		if !slices.Equal(actual, x.want) {
			t.Errorf("CleanSplitN(%q, %q, %d) = %v, want %v", x.text, x.sep, x.n, actual, x.want)
		}
	}
}

func TestSpaceSplit(t *testing.T) {
	type testCase struct {
		text string
		want []string
	}
	testCases := []testCase{
		{"a\tb\tc", []string{"a", "b", "c"}},
		{"a\nb\nc", []string{"a", "b", "c"}},
		{"a  b  c", []string{"a", "b", "c"}},
	}
	for _, x := range testCases {
		actual := SpaceSplit(x.text)
		if !slices.Equal(actual, x.want) {
			t.Errorf("SpaceSplit(%q) = %v, want %v", x.text, actual, x.want)
		}
	}
}

func TestJoin(t *testing.T) {
	type testCase struct {
		glue  string
		parts []string
		want  string
	}
	testCases := []testCase{
		{",", []string{"a", "b", "c"}, "a,b,c"},
		{"\t", []string{"a", "b", "c"}, "a\tb\tc"},
		{"+", []string{"1", "2", "3"}, "1+2+3"},
		{"x", []string{}, ""},
	}
	for _, x := range testCases {
		actual := Join(x.glue, x.parts...)
		if actual != x.want {
			t.Errorf("Join(%q, %v) = %q, want %q", x.glue, x.parts, actual, x.want)
		}
	}
}

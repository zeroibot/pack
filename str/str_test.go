package str

import "testing"

type testCase[T any] struct {
	text string
	want T
}

func TestLen(t *testing.T) {
	testCases := []testCase[int]{
		{"", 0},
		{"abc", 3},
		{"12345", 5},
		{"X", 1},
		{"XY", 2},
	}
	for _, x := range testCases {
		actual := Len(x.text)
		if x.want != actual {
			t.Errorf("Len(%q) = %d, want %d", x.text, actual, x.want)
		}
	}
}

func TestIsEmpty(t *testing.T) {
	testCases := []testCase[bool]{
		{"", true},
		{"123", false},
		{"a", false},
	}
	for _, x := range testCases {
		actual := IsEmpty(x.text)
		if x.want != actual {
			t.Errorf("IsEmpty(%q) = %v, want %v", x.text, actual, x.want)
		}
		opposite := NotEmpty(x.text)
		if !x.want != opposite {
			t.Errorf("NotEmpty(%q) = %v, want %v", x.text, opposite, x.want)
		}
	}
}

func TestGuard(t *testing.T) {
	type testCase struct {
		text  string
		guard string
		want  string
	}
	testCases := []testCase{
		{"", "default", "default"},
		{"abc", "default", "abc"},
		{"", "unknown", "unknown"},
		{"def", "unknown", "def"},
	}
	for _, x := range testCases {
		actual := Guard(x.text, x.guard)
		if x.want != actual {
			t.Errorf("Guard(%q, %q) = %q, want %q", x.text, x.guard, actual, x.want)
		}
	}
}

func TestWrap(t *testing.T) {
	type testCase struct {
		text    string
		wrapper string
		want    string
	}
	testCases := []testCase{
		{"hello", "()", "(hello)"},
		{"hello", "(", "(hello"},
		{"world", "<>", "<world>"},
		{"column", "``", "`column`"},
		{"default", "", "default"},
	}
	for _, x := range testCases {
		actual := Wrap(x.text, x.wrapper)
		if x.want != actual {
			t.Errorf("Wrap(%q, %q) = %q, want %q", x.text, x.wrapper, actual, x.want)
		}
	}
}

func TestWrapList(t *testing.T) {
	type testCase struct {
		items   []string
		wrapper string
		want    string
	}
	testCases := []testCase{
		{[]string{"a", "b", "c"}, "[]", "[a, b, c]"},
		{[]string{"a => 1", "b => 2"}, "{}", "{a => 1, b => 2}"},
		{[]string{"1", "2", "3"}, "", "1, 2, 3"},
	}
	for _, x := range testCases {
		actual := WrapList(x.items, x.wrapper)
		if x.want != actual {
			t.Errorf("WrapList(%v, %q) = %q, want %q", x.items, x.wrapper, actual, x.want)
		}
	}
}

package dict

import (
	"cmp"
	"maps"
	"slices"
	"testing"
)

func TestDict(t *testing.T) {
	m := map[string]int{
		"apple":  5,
		"orange": 3,
		"banana": 2,
	}
	m2 := Copy(m)
	if maps.Equal(m, m2) == false {
		t.Errorf("Copy() = %v, want %v", m2, m)
	}
	size := Len(m)
	if size != 3 {
		t.Errorf("Len() = %d, want 3", size)
	}
	notEmpty := NotEmpty(m)
	if notEmpty != true {
		t.Errorf("NotEmpty() = %t, want true", notEmpty)
	}
	expKeys := []string{"apple", "banana", "orange"}
	actualKeys := Keys(m)
	slices.Sort(actualKeys)
	if slices.Equal(actualKeys, expKeys) == false {
		t.Errorf("Keys() = %v, want %v", actualKeys, expKeys)
	}
	expValues := []int{2, 3, 5}
	actualValues := Values(m)
	slices.Sort(actualValues)
	if slices.Equal(actualValues, expValues) == false {
		t.Errorf("Values() = %v, want %v", actualValues, expValues)
	}
	expEntries := []Entry[string, int]{{"apple", 5}, {"banana", 2}, {"orange", 3}}
	actualEntries := Entries(m)
	slices.SortFunc(actualEntries, func(e1, e2 Entry[string, int]) int {
		return cmp.Compare(e1.Key, e2.Key)
	})
	if slices.Equal(actualEntries, expEntries) == false {
		t.Errorf("Entries() = %v, want %v", actualEntries, expEntries)
	}
	expStrings := []string{"<apple: 5>", "<banana: 2>", "<orange: 3>"}
	for i, entry := range actualEntries {
		actualString := entry.String()
		if actualString != expStrings[i] {
			t.Errorf("Entry.String() = %q, want %q", actualString, expStrings[i])
		}
	}
	noKeyCases := []Entry[string, bool]{
		{"apple", false},
		{"grape", true},
	}
	for _, x := range noKeyCases {
		key, want := x.Tuple()
		actual := NoKey(m, key)
		if actual != want {
			t.Errorf("NoKey(%q) = %v, want %v", key, actual, want)
		}
	}
	noValueCases := []Entry[int, bool]{
		{3, false},
		{5, false},
		{1, true},
		{69, true},
	}
	for _, x := range noValueCases {
		value, want := x.Tuple()
		actual := NoValue(m, value)
		if actual != want {
			t.Errorf("NoValue(%d) = %v, want %v", value, actual, want)
		}
	}
	getCases := []Entry[string, int]{
		{"apple", 5},
		{"zebra", 0},
	}
	for _, x := range getCases {
		key, want := x.Tuple()
		actual := m[key]
		if actual != want {
			t.Errorf("map[%q] = %v, want %v", key, actual, want)
		}
	}
	defaultValue := 69
	getCases = []Entry[string, int]{
		{"orange", 3},
		{"cherry", defaultValue},
	}
	for _, x := range getCases {
		key, want := x.Tuple()
		actual := GetOrDefault(m, key, defaultValue)
		if actual != want {
			t.Errorf("GetOrDefault(%q, %v) = %v, want %v", key, defaultValue, actual, want)
		}
		SetDefault(m, key, defaultValue)
		actual = m[key]
		if actual != want {
			t.Errorf("SetDefault(%q, %v) = %v, want %v", key, defaultValue, actual, want)
		}
	}
	keyFnCases := []Entry[bool, func(string) bool]{
		{false, func(key string) bool { return key == "apple" }},
		{true, func(key string) bool { return key == "zebra" }},
	}
	for _, x := range keyFnCases {
		want, test := x.Tuple()
		actual := NoKeyFunc(m, test)
		if actual != want {
			t.Errorf("NoKeyFunc = %v, want %v", actual, want)
		}
	}
	valueFnCases := []Entry[bool, func(int) bool]{
		{true, func(value int) bool { return value > 100 }},
		{false, func(value int) bool { return value == 5 }},
	}
	for _, x := range valueFnCases {
		want, test := x.Tuple()
		actual := NoValueFunc(m, test)
		if actual != want {
			t.Errorf("NoValueFunc = %v, want %v", actual, want)
		}
	}

	Clear(m)
	isEmpty := IsEmpty(m)
	if isEmpty != true {
		t.Errorf("IsEmpty() = %t, want true", isEmpty)
	}
}

package ds

import (
	"cmp"
	"maps"
	"slices"
	"testing"
)

func TestMap(t *testing.T) {
	m := Map[string, int]{
		"apple":  5,
		"orange": 3,
		"banana": 2,
	}
	expString := "{apple: 5, banana: 2, orange: 3}"
	actualString := m.String()
	if actualString != expString {
		t.Errorf("Map.String() = %q, want %q", actualString, expString)
	}
	m2 := m.Copy()
	if maps.Equal(m, m2) == false {
		t.Errorf("Map.Copy() = %v, want %v", m2, m)
	}
	size := m.Len()
	if size != 3 {
		t.Errorf("Map.Len() = %d, want 3", size)
	}
	notEmpty := m.NotEmpty()
	if notEmpty != true {
		t.Errorf("Map.NotEmpty() = %t, want true", notEmpty)
	}
	expKeys := []string{"apple", "banana", "orange"}
	actualKeys := m.Keys()
	slices.Sort(actualKeys)
	if slices.Equal(actualKeys, expKeys) == false {
		t.Errorf("Map.Keys() = %v, want %v", actualKeys, expKeys)
	}
	expValues := []int{2, 3, 5}
	actualValues := m.Values()
	slices.Sort(actualValues)
	if slices.Equal(actualValues, expValues) == false {
		t.Errorf("Map.Values() = %v, want %v", actualValues, expValues)
	}
	expEntries := []Tuple2[string, int]{{"apple", 5}, {"banana", 2}, {"orange", 3}}
	actualEntries := m.Entries()
	slices.SortFunc(actualEntries, func(e1, e2 Tuple2[string, int]) int {
		return cmp.Compare(e1.V1, e2.V1)
	})
	if slices.Equal(actualEntries, expEntries) == false {
		t.Errorf("Map.Entries() = %v, want %v", actualEntries, expEntries)
	}
	keyCases := []Tuple2[string, bool]{
		{"apple", false},
		{"grape", true},
	}
	for _, x := range keyCases {
		key, want := x.Values()
		actual := m.NoKey(key)
		if actual != want {
			t.Errorf("Map.NoKey(%q) = %v, want %v", key, actual, want)
		}
	}
	getCases := []Tuple2[string, int]{
		{"apple", 5},
		{"zebra", 0},
	}
	for _, x := range getCases {
		key, want := x.Values()
		actual := m[key]
		if actual != want {
			t.Errorf("Map[%q] = %v, want %v", key, actual, want)
		}
	}
	defaultValue := 69
	getCases = []Tuple2[string, int]{
		{"orange", 3},
		{"cherry", defaultValue},
	}
	for _, x := range getCases {
		key, want := x.Values()
		actual := m.GetOrDefault(key, defaultValue)
		if actual != want {
			t.Errorf("Map.GetOrDefault(%q, %v) = %v, want %v", key, defaultValue, actual, want)
		}
		m.SetDefault(key, defaultValue)
		actual = m[key]
		if actual != want {
			t.Errorf("Map.SetDefault(%q, %v) = %v, want %v", key, defaultValue, actual, want)
		}
	}
	keyFnCases := []Tuple2[func(string) bool, bool]{
		{func(key string) bool { return key == "apple" }, false},
		{func(key string) bool { return key == "zebra" }, true},
	}
	for _, x := range keyFnCases {
		test, want := x.Values()
		actual := m.NoKeyFunc(test)
		if actual != want {
			t.Errorf("Map.NoKeyFunc = %v, want %v", actual, want)
		}
	}
	valueFnCases := []Tuple2[func(int) bool, bool]{
		{func(value int) bool { return value > 100 }, true},
		{func(value int) bool { return value == 5 }, false},
	}
	for _, x := range valueFnCases {
		test, want := x.Values()
		actual := m.NoValueFunc(test)
		if actual != want {
			t.Errorf("Map.NoValueFunc = %v, want %v", actual, want)
		}
	}

	m.Clear()
	isEmpty := m.IsEmpty()
	if isEmpty != true {
		t.Errorf("Map.IsEmpty() = %t, want true", isEmpty)
	}
}

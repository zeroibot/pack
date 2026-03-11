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
	actualKeys2 := m.SortedKeysFunc(func(k1, k2 string) int { return cmp.Compare(k1, k2) })
	if slices.Equal(actualKeys2, expKeys) == false {
		t.Errorf("Map.SortedKeysFunc() = %v, want %v", actualKeys2, expKeys)
	}
	expValues := []int{2, 3, 5}
	actualValues := m.Values()
	slices.Sort(actualValues)
	if slices.Equal(actualValues, expValues) == false {
		t.Errorf("Map.Values() = %v, want %v", actualValues, expValues)
	}
	actualValues2 := m.SortedValuesFunc(func(v1, v2 int) int { return cmp.Compare(v1, v2) })
	if slices.Equal(actualValues2, expValues) == false {
		t.Errorf("Map.SortedValuesFunc() = %v, want %v", actualValues2, expValues)
	}
	expEntries := []Tuple2[string, int]{{"apple", 5}, {"banana", 2}, {"orange", 3}}
	actualEntries := m.Entries()
	sortFn := func(e1, e2 Tuple2[string, int]) int {
		return cmp.Compare(e1.V1, e2.V1)
	}
	slices.SortFunc(actualEntries, sortFn)
	if slices.Equal(actualEntries, expEntries) == false {
		t.Errorf("Map.Entries() = %v, want %v", actualEntries, expEntries)
	}
	actualEntries2 := m.SortedEntriesFunc(sortFn)
	if slices.Equal(actualEntries2, expEntries) == false {
		t.Errorf("Map.SortedEntriesFunc() = %v, want %v", actualEntries2, expEntries)
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
	mf := m.Filter(func(key string, value int) bool { return key != "zebra" && value <= 50 })
	actualEntries = mf.SortedEntriesFunc(sortFn)
	if slices.Equal(actualEntries, expEntries) == false {
		t.Errorf("Map.Filter.Entries = %v, want %v", actualEntries, expEntries)
	}

	m.Delete("cherry")
	if m.HasKey("cherry") {
		t.Errorf("Map.Delete, HasKey = true, want false")
	}

	m2 = Map[string, int]{
		"cherry": 30,
		"banana": 10,
	}
	expEntries = []Tuple2[string, int]{
		{"apple", 5}, {"banana", 10}, {"cherry", 30}, {"orange", 3},
	}
	m.Update(m2)
	actualEntries = m.SortedEntriesFunc(sortFn)
	if slices.Equal(actualEntries, expEntries) == false {
		t.Errorf("Map.Update = %v, want %v", actualEntries, expEntries)
	}

	m.Clear()
	isEmpty := m.IsEmpty()
	if isEmpty != true {
		t.Errorf("Map.IsEmpty() = %t, want true", isEmpty)
	}
}

func TestZipMap(t *testing.T) {
	keys := []string{"a", "b", "c"}
	values := List[int]{1, 2, 3}
	m := ZipMap(keys, values)
	for i, key := range keys {
		want := values[i]
		actual, ok := m[key]
		if actual != want || !ok {
			t.Errorf("ZipMap[%q] = %d, %v, want %d, true", key, actual, ok, want)
		}
		key = key + "x"
		actual, ok = m[key]
		if actual != 0 || ok {
			t.Errorf("ZipMap[%q] = %d, %v, want 0, false", key, actual, ok)
		}
	}
	keys2, values2 := m.Unzip()
	keySet1, keySet2 := NewSetFrom(keys), NewSetFrom(keys2)
	valSet1, valSet2 := NewSetFrom(values), NewSetFrom(values2)
	if keySet1.HasNoDifference(keySet2) == false {
		t.Errorf("Unzip.Keys = %v, want %v", keys2, keys)
	}
	if valSet1.HasNoDifference(valSet2) == false {
		t.Errorf("Unzip.Values = %v, want %v", values2, values)
	}
	values2 = List[int]{1, 2}
	m2 := ZipMap(keys, values2)
	if m2.Len() != 2 {
		t.Errorf("ZipMap.Len = %d, want 2", m2.Len())
	}
}

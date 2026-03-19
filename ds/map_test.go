package ds

import (
	"cmp"
	"fmt"
	"slices"
	"testing"

	"github.com/roidaradal/tst"
)

func TestMap(t *testing.T) {
	// String
	m := Map[string, int]{"apple": 5, "orange": 3, "banana": 2}
	tst.AssertEqual(t, "Map.String", m.String(), "{apple: 5, banana: 2, orange: 3}")
	// Copy
	m2 := m.Copy()
	tst.AssertMapEqual(t, "Map.Copy", m2, m)
	// Len, NotEmpty
	tst.AssertEqual(t, "Map.Len", m.Len(), 3)
	tst.AssertEqual(t, "Map.NotEmpty", m.NotEmpty(), true)
	// Keys
	wantKeys := []string{"apple", "banana", "orange"}
	keys := m.Keys()
	slices.Sort(keys)
	tst.AssertListEqual(t, "Map.Keys", keys, wantKeys)
	// SortedKeysFunc
	keys = m.SortedKeysFunc(func(k1, k2 string) int { return cmp.Compare(k1, k2) })
	tst.AssertListEqual(t, "Map.SortedKeysFunc", keys, wantKeys)
	// Values
	wantValues := []int{2, 3, 5}
	values := m.Values()
	slices.Sort(values)
	tst.AssertListEqual(t, "Map.Values", values, wantValues)
	values = m.SortedValuesFunc(func(v1, v2 int) int { return cmp.Compare(v1, v2) })
	tst.AssertListEqual(t, "Map.SortedValuesFunc", values, wantValues)
	// Entries
	wantEntries := []Tuple2[string, int]{{"apple", 5}, {"banana", 2}, {"orange", 3}}
	entries := m.Entries()
	sortFn := func(e1, e2 Tuple2[string, int]) int {
		return cmp.Compare(e1.V1, e2.V1)
	}
	slices.SortFunc(entries, sortFn)
	tst.AssertListEqual(t, "Map.Entries", entries, wantEntries)
	entries = m.SortedEntriesFunc(sortFn)
	tst.AssertListEqual(t, "Map.SortedEntriesFunc", entries, wantEntries)
	// NoKey
	keyCases := []tst.P2W1[Map[string, int], string, bool]{
		{m, "apple", false},
		{m, "grape", true},
	}
	tst.AllP2W1(t, keyCases, "NoKey", Map[string, int].NoKey, tst.AssertEqual)
	// Get
	getCases := []Tuple2[string, int]{{"apple", 5}, {"zebra", 0}}
	for _, x := range getCases {
		key, want := x.Unpack()
		option := m.Get(key)
		tst.AssertEqual(t, fmt.Sprintf("Map[%q]", key), m[key], want)
		tst.AssertEqualAnd(t, "Map.Get", option.Value(), want, option.NotNil(), want != 0)
	}
	// GetOrDefault
	defaultValue := 69
	getCases = []Tuple2[string, int]{{"orange", 3}, {"cherry", defaultValue}}
	for _, x := range getCases {
		key, want := x.Unpack()
		tst.AssertEqual(t, "Map.GetOrDefault", m.GetOrDefault(key, defaultValue), want)
		m.SetDefault(key, defaultValue)
		tst.AssertEqual(t, "Map.SetDefault", m[key], want)
	}
	// NoKeyFunc
	keyFnCases := []tst.P2W1[Map[string, int], func(string) bool, bool]{
		{m, func(key string) bool { return key == "apple" }, false},
		{m, func(key string) bool { return key == "zebra" }, true},
	}
	tst.AllP2W1(t, keyFnCases, "Map.NoKeyFunc", Map[string, int].NoKeyFunc, tst.AssertEqual)
	// NoValueFunc
	valueFnCases := []tst.P2W1[Map[string, int], func(int) bool, bool]{
		{m, func(value int) bool { return value > 100 }, true},
		{m, func(value int) bool { return value == 5 }, false},
	}
	tst.AllP2W1(t, valueFnCases, "Map.NoValueFunc", Map[string, int].NoValueFunc, tst.AssertEqual)
	// Filter
	mf := m.Filter(func(key string, value int) bool { return key != "zebra" && value <= 50 })
	tst.AssertListEqual(t, "Map.Filter", mf.SortedEntriesFunc(sortFn), wantEntries)
	// Delete
	m.Delete("cherry")
	tst.AssertFalse(t, "Delete.HasKey", m.HasKey("cherry"))
	// Update
	m2 = Map[string, int]{"cherry": 30, "banana": 10}
	wantEntries = []Tuple2[string, int]{
		{"apple", 5}, {"banana", 10}, {"cherry", 30}, {"orange", 3},
	}
	m.Update(m2)
	tst.AssertListEqual(t, "Map.Update", m.SortedEntriesFunc(sortFn), wantEntries)

	// Clear, IsEmpty
	m.Clear()
	tst.AssertTrue(t, "IsEmpty", m.IsEmpty())
}

func TestZipMap(t *testing.T) {
	keys := []string{"a", "b", "c"}
	values := List[int]{1, 2, 3}
	// ZipMap
	m := ZipMap(keys, values)
	for i, key := range keys {
		actual, ok := m[key]
		tst.AssertEqualAnd(t, "ZipMap", actual, values[i], ok, true)
		actual, ok = m[key+"x"]
		tst.AssertEqualAnd(t, "ZipMap", actual, 0, ok, false)
	}
	// Unzip
	keys2, values2 := m.Unzip()
	keySet1, keySet2 := NewSetFrom(keys), NewSetFrom(keys2)
	valSet1, valSet2 := NewSetFrom(values), NewSetFrom(values2)
	tst.AssertTrue(t, "Unzip.Keys", keySet1.HasNoDifference(keySet2))
	tst.AssertTrue(t, "Unzip.Values", valSet1.HasNoDifference(valSet2))
	// ZipMap
	values2 = List[int]{1, 2}
	m2 := ZipMap(keys, values2)
	tst.AssertMapEqual(t, "ZipMap", m2, map[string]int{"a": 1, "b": 2})
	tst.AssertEqual(t, "ZipMap.Len", m2.Len(), 2)
}

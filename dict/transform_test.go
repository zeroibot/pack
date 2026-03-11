package dict

import (
	"maps"
	"slices"
	"testing"
)

func TestZip(t *testing.T) {
	keys := []string{"a", "b", "c"}
	values := []int{1, 2, 3}
	m := Zip(keys, values)
	for i, key := range keys {
		want := values[i]
		actual, ok := m[key]
		if actual != want || !ok {
			t.Errorf("Zip[%q] = %d, %v, want %d, true", key, actual, ok, want)
		}
		key = key + "x"
		actual, ok = m[key]
		if actual != 0 || ok {
			t.Errorf("Zip[%q] = %d, %v, want 0, false", key, actual, ok)
		}
	}
	keys2, values2 := Unzip(m)
	m2 := Zip(keys2, values2)
	if maps.Equal(m, m2) == false {
		t.Errorf("Unzip.Zip = %v, want %v", m2, m)
	}
	values2 = []int{1, 2}
	m3 := Zip(keys, values2)
	if Len(m3) != 2 {
		t.Errorf("ZipMap.Len = %d, want 2", Len(m3))
	}
}

func TestSwap(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	want := []Entry[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	m2 := Swap(m)
	actual := SortedEntries(m2)
	if slices.Equal(actual, want) == false {
		t.Errorf("Swap.Entries = %v, want %v", actual, want)
	}
}

func TestSwapList(t *testing.T) {
	m := map[string][]int{
		"a": {1, 3, 5},
		"b": {2, 4},
	}
	want := []Entry[int, string]{{1, "a"}, {2, "b"}, {3, "a"}, {4, "b"}, {5, "a"}}
	m2 := SwapList(m)
	actual := SortedEntries(m2)
	if slices.Equal(actual, want) == false {
		t.Errorf("SwapList.Entries = %v, want %v", actual, want)
	}
}

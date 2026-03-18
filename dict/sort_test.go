package dict

import (
	"cmp"
	"testing"

	"github.com/roidaradal/tst"
)

func TestSortedKeys(t *testing.T) {
	m := map[string]int{
		"Alice":   1,
		"Charlie": 2,
		"Bob":     3,
		"Tom":     4,
		"Dan":     5,
	}
	want := []string{"Alice", "Bob", "Charlie", "Dan", "Tom"}
	tst.AssertListEqual(t, "SortedKeys", SortedKeys(m), want)
}

func TestSortedKeysFunc(t *testing.T) {
	m := map[string]int{
		"Alice":    1,
		"Charlie":  2,
		"Bob":      3,
		"Denzel":   4,
		"Emmanuel": 5,
	}
	want := []string{"Bob", "Alice", "Denzel", "Charlie", "Emmanuel"}
	actual := SortedKeysFunc(m, func(k1, k2 string) int { return cmp.Compare(len(k1), len(k2)) })
	tst.AssertListEqual(t, "SortedKeysFunc", actual, want)
}

func TestSortedValues(t *testing.T) {
	m := map[string]int{
		"Alice":   1,
		"Charlie": 2,
		"Bob":     3,
		"Tom":     4,
		"Dan":     5,
	}
	want := []int{1, 2, 3, 4, 5}
	tst.AssertListEqual(t, "SortedValues", SortedValues(m), want)
}

func TestSortedValuesFunc(t *testing.T) {
	m := map[string]int{
		"Alice":   1,
		"Charlie": 2,
		"Bob":     3,
		"Tom":     4,
		"Dan":     5,
	}
	want := []int{5, 4, 3, 2, 1}
	actual := SortedValuesFunc(m, func(v1, v2 int) int { return cmp.Compare(v2, v1) })
	tst.AssertListEqual(t, "SortedValuesFunc", actual, want)
}

func TestSortedEntries(t *testing.T) {
	m := map[string]int{
		"Alice":   1,
		"Charlie": 2,
		"Bob":     3,
		"Tom":     4,
		"Dan":     5,
	}
	want := []Entry[string, int]{
		{"Alice", 1},
		{"Bob", 3},
		{"Charlie", 2},
		{"Dan", 5},
		{"Tom", 4},
	}
	tst.AssertListEqual(t, "SortedEntries", SortedEntries(m), want)
}

func TestSortedEntriesFunc(t *testing.T) {
	m := map[string]int{
		"Alice":   1, // 5 * 1 = 5
		"Charlie": 2, // 7 * 2 = 14
		"Bob":     3, // 3 * 3 = 9
		"Tom":     4, // 3 * 4 = 12
		"Dan":     5, // 3 * 5 = 15
	}
	want := []Entry[string, int]{
		{"Dan", 5},
		{"Charlie", 2},
		{"Tom", 4},
		{"Bob", 3},
		{"Alice", 1},
	}
	actual := SortedEntriesFunc(m, func(e1, e2 Entry[string, int]) int {
		score1 := len(e1.Key) * e1.Value
		score2 := len(e2.Key) * e2.Value
		return cmp.Compare(score2, score1)
	})
	tst.AssertListEqual(t, "SortedEntriesFunc", actual, want)
}

func TestSortValueLists(t *testing.T) {
	m := map[string][]int{
		"section1": {90, 95, 87, 82, 93},
		"section2": {78, 88, 98, 81, 69},
	}
	want1 := []int{82, 87, 90, 93, 95}
	want2 := []int{69, 78, 81, 88, 98}
	SortValueLists(m)
	tst.AssertListEqual(t, "SortValueLists", m["section1"], want1)
	tst.AssertListEqual(t, "SortValueLists", m["section2"], want2)
}

func TestSortValueListsFunc(t *testing.T) {
	m := map[string][]int{
		"section1": {90, 95, 87, 82, 93},
		"section2": {78, 88, 98, 81, 69},
	}
	want1 := []int{95, 93, 90, 87, 82}
	want2 := []int{98, 88, 81, 78, 69}
	SortValueListsFunc(m, func(x1, x2 int) int { return cmp.Compare(x2, x1) })
	tst.AssertListEqual(t, "SortValueListsFunc", m["section1"], want1)
	tst.AssertListEqual(t, "SortValueListsFunc", m["section2"], want2)
}

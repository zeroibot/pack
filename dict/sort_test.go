package dict

import (
	"cmp"
	"slices"
	"testing"
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
	actual := SortedKeys(m)
	if slices.Equal(actual, want) == false {
		t.Errorf("SortedKeys = %v, want %v", actual, want)
	}
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
	if slices.Equal(actual, want) == false {
		t.Errorf("SortedKeysFunc = %v, want %v", actual, want)
	}
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
	actual := SortedValues(m)
	if slices.Equal(actual, want) == false {
		t.Errorf("SortedValues = %v, want %v", actual, want)
	}
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
	if slices.Equal(actual, want) == false {
		t.Errorf("SortedValuesFunc = %v, want %v", actual, want)
	}
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
	actual := SortedEntries(m)
	if slices.Equal(actual, want) == false {
		t.Errorf("SortedEntries = %v, want %v", actual, want)
	}
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
	if slices.Equal(actual, want) == false {
		t.Errorf("SortedEntriesFunc = %v, want %v", actual, want)
	}
}

func TestSortValueLists(t *testing.T) {
	m := map[string][]int{
		"section1": {90, 95, 87, 82, 93},
		"section2": {78, 88, 98, 81, 69},
	}
	want1 := []int{82, 87, 90, 93, 95}
	want2 := []int{69, 78, 81, 88, 98}
	SortValueLists(m)
	actual1 := m["section1"]
	actual2 := m["section2"]
	if slices.Equal(actual1, want1) == false {
		t.Errorf("SortValueLists = %v, want %v", actual1, want1)
	}
	if slices.Equal(actual2, want2) == false {
		t.Errorf("SortValueLists = %v, want %v", actual2, want2)
	}
}

func TestSortValueListsFunc(t *testing.T) {
	m := map[string][]int{
		"section1": {90, 95, 87, 82, 93},
		"section2": {78, 88, 98, 81, 69},
	}
	want1 := []int{95, 93, 90, 87, 82}
	want2 := []int{98, 88, 81, 78, 69}
	SortValueListsFunc(m, func(x1, x2 int) int { return cmp.Compare(x2, x1) })
	actual1 := m["section1"]
	actual2 := m["section2"]
	if slices.Equal(actual1, want1) == false {
		t.Errorf("SortValueListsFunc = %v, want %v", actual1, want1)
	}
	if slices.Equal(actual2, want2) == false {
		t.Errorf("SortValueListsFunc = %v, want %v", actual2, want2)
	}

}

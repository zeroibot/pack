package dict

import (
	"maps"
	"slices"
	"testing"
)

func TestGroupByValue(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 1,
		"d": 2,
	}
	wantEntries := []Entry[int, []string]{
		{1, []string{"a", "c"}},
		{2, []string{"b", "d"}},
	}
	g := GroupByValue(m)
	actualEntries := SortedEntries(g)
	if len(actualEntries) != len(wantEntries) {
		t.Errorf("GroupByValue.Len = %d, want %d", len(actualEntries), len(wantEntries))
	}
	for i, actual := range actualEntries {
		want := wantEntries[i]
		if actual.Key != want.Key {
			t.Errorf("GroupByValue.Key = %d, want %d", actual.Key, want.Key)
		}
		slices.Sort(actual.Value)
		if slices.Equal(actual.Value, want.Value) == false {
			t.Errorf("GroupByValue.Value = %v, want %v", actual.Value, want.Value)
		}
	}
}

func TestGroupByFunc(t *testing.T) {
	type person struct {
		name  string
		score int
	}
	m := map[string]person{
		"alpha": {"Alice", 1},
		"beta":  {"Bob", 2},
		"gamma": {"Gus", 1},
		"delta": {"Dan", 2},
	}
	wantEntries := []Entry[int, []byte]{
		{1, []byte{'a', 'g'}},
		{2, []byte{'b', 'd'}},
	}
	keyFn := func(key string) byte {
		return key[0]
	}
	valueFn := func(p person) int {
		return p.score
	}
	g := GroupByFunc(m, keyFn, valueFn)
	actualEntries := SortedEntries(g)
	if len(actualEntries) != len(wantEntries) {
		t.Errorf("GroupByFunc.Len = %d, want %d", len(actualEntries), len(wantEntries))
	}
	for i, actual := range actualEntries {
		want := wantEntries[i]
		if actual.Key != want.Key {
			t.Errorf("GroupByFunc.Key = %d, want %d", actual.Key, want.Key)
		}
		slices.Sort(actual.Value)
		if slices.Equal(actual.Value, want.Value) == false {
			t.Errorf("GroupByFunc.Value = %v, want %v", actual.Value, want.Value)
		}
	}
}

func TestGroupByValueList(t *testing.T) {
	m := map[string][]int{
		"a": {1, 2, 3},
		"b": {2, 3},
		"c": {1, 4, 5},
	}
	wantEntries := []Entry[int, []string]{
		{1, []string{"a", "c"}},
		{2, []string{"a", "b"}},
		{3, []string{"a", "b"}},
		{4, []string{"c"}},
		{5, []string{"c"}},
	}
	g := GroupByValueList(m)
	actualEntries := SortedEntries(g)
	if len(actualEntries) != len(wantEntries) {
		t.Errorf("GroupByValueList.Len = %d, want %d", len(actualEntries), len(wantEntries))
	}
	for i, actual := range actualEntries {
		want := wantEntries[i]
		if actual.Key != want.Key {
			t.Errorf("GroupByValueList.Key = %d, want %d", actual.Key, want.Key)
		}
		slices.Sort(actual.Value)
		if slices.Equal(actual.Value, want.Value) == false {
			t.Errorf("GroupByValueList.Value = %v, want %v", actual.Value, want.Value)
		}
	}
}

func TestGroupByFuncList(t *testing.T) {
	type person struct {
		score int
	}
	m := map[string][]person{
		"A": {{1}, {2}, {4}, {2}},
		"B": {{1}, {3}, {2}, {1}},
	}
	wantEntries := []Entry[int, []string]{
		{1, []string{"A", "B", "B"}},
		{2, []string{"A", "A", "B"}},
		{3, []string{"B"}},
		{4, []string{"A"}},
	}
	keyFn := func(key string) string { return key }
	valueFn := func(p person) int { return p.score }
	g := GroupByFuncList(m, keyFn, valueFn)
	actualEntries := SortedEntries(g)
	if len(actualEntries) != len(wantEntries) {
		t.Errorf("GroupByFuncList.Len = %d, want %d", len(actualEntries), len(wantEntries))
	}
	for i, actual := range actualEntries {
		want := wantEntries[i]
		if actual.Key != want.Key {
			t.Errorf("GroupByFuncList.Key = %d, want %d", actual.Key, want.Key)
		}
		slices.Sort(actual.Value)
		if slices.Equal(actual.Value, want.Value) == false {
			t.Errorf("GroupByFuncList.Value = %v, want %v", actual.Value, want.Value)
		}
	}
}

func TestTallyValues(t *testing.T) {
	m := map[string]int{
		"a": 95,
		"b": 92,
		"c": 90,
		"d": 94,
		"e": 88,
		"f": 90,
	}
	counter := TallyValues(m, []int{88, 90, 92, 94, 96})
	want := Counter[int]{
		88: 1,
		90: 2,
		92: 1,
		94: 1,
		96: 0,
	}
	if maps.Equal(counter, want) == false {
		t.Errorf("TallyValues = %v, want %v", counter, want)
	}
}

func TestTallyFunc(t *testing.T) {
	m := map[string]int{
		"a": 95,
		"b": 92,
		"c": 90,
		"d": 94,
		"e": 88,
		"f": 84,
		"g": 79,
	}
	valueFn := func(value int) int {
		return 10 * (value / 10)
	}
	counter := TallyFunc(m, valueFn)
	want := Counter[int]{
		70: 1,
		80: 2,
		90: 4,
	}
	if maps.Equal(counter, want) == false {
		t.Errorf("TallyFunc = %v, want %v", counter, want)
	}
}

package list

import (
	"maps"
	"reflect"
	"slices"
	"testing"
)

func TestCompareAllAny(t *testing.T) {
	var ints0 []int
	ints1 := []int{1, 1, 1, 1}
	ints2 := []int{1, 2, 3, 4}
	ints3 := []int{2, 2, 1, 2}
	var bools0 []bool
	bools1 := []bool{true, true, true}
	bools2 := []bool{false, false}
	bools3 := []bool{false, true, false}
	// AllEqual
	actual := AllEqual(ints0, 5)
	if actual != false {
		t.Errorf("AllEqual() = %t, want %t", actual, false)
	}
	actual = AllEqual(ints1, 1)
	if actual != true {
		t.Errorf("AllEqual() = %t, want %t", actual, true)
	}
	actual = AllEqual(ints2, 3)
	if actual != false {
		t.Errorf("AllEqual() = %t, want %t", actual, false)
	}
	actual = AllEqual(ints3, 2)
	if actual != false {
		t.Errorf("AllEqual() = %t, want %t", actual, false)
	}
	// AllTrue, All False
	actual = AllTrue(bools0)
	if actual != false {
		t.Errorf("AllTrue() = %t, want %t", actual, false)
	}
	actual = AllFalse(bools0)
	if actual != false {
		t.Errorf("AllFalse() = %t, want %t", actual, false)
	}
	actual = AllTrue(bools1)
	if actual != true {
		t.Errorf("AllTrue() = %t, want %t", actual, true)
	}
	actual = AllTrue(bools3)
	if actual != false {
		t.Errorf("AllTrue() = %t, want %t", actual, false)
	}
	actual = AllFalse(bools2)
	if actual != true {
		t.Errorf("AllFalse() = %t, want %t", actual, true)
	}
	actual = AllFalse(bools3)
	if actual != false {
		t.Errorf("AllFalse() = %t, want %t", actual, false)
	}
	// Has, HasNo
	actual, actual2 := Has(ints0, 1), HasNo(ints0, 1)
	if actual != false {
		t.Errorf("Has() = %t, want %t", actual, false)
	}
	if actual2 != true {
		t.Errorf("HasNo() = %t, want %t", actual2, true)
	}
	actual, actual2 = Has(ints1, 1), HasNo(ints1, 1)
	if actual != true {
		t.Errorf("Has() = %t, want %t", actual, true)
	}
	if actual2 != false {
		t.Errorf("HasNo() = %t, want %t", actual2, false)
	}
	actual, actual2 = Has(ints3, 1), HasNo(ints3, 1)
	if actual != true {
		t.Errorf("Has() = %t, want %t", actual, true)
	}
	if actual2 != false {
		t.Errorf("HasNo() = %t, want %t", actual2, false)
	}
	actual, actual2 = Has(ints2, 5), HasNo(ints2, 5)
	if actual != false {
		t.Errorf("Has() = %t, want %t", actual, false)
	}
	if actual2 != true {
		t.Errorf("HasNo() = %t, want %t", actual2, true)
	}
	// AnyTrue, AnyFalse
	actual, actual2 = AnyTrue(bools1), AnyFalse(bools1)
	if actual != true {
		t.Errorf("AnyTrue() = %t, want %t", actual, true)
	}
	if actual2 != false {
		t.Errorf("AnyFalse() = %t, want %t", actual2, false)
	}
	actual, actual2 = AnyTrue(bools2), AnyFalse(bools2)
	if actual != false {
		t.Errorf("AnyTrue() = %t, want %t", actual, false)
	}
	if actual2 != true {
		t.Errorf("AnyFalse() = %t, want %t", actual2, true)
	}
	actual, actual2 = AnyTrue(bools3), AnyFalse(bools3)
	if actual != true {
		t.Errorf("AnyTrue() = %t, want %t", actual, true)
	}
	if actual2 != true {
		t.Errorf("AnyFalse() = %t, want %t", actual2, true)
	}
}

func TestIndexFunctions(t *testing.T) {
	// IndexLookup
	items := []string{" ", "A", "B", "C"}
	wantMap := map[string]int{" ": 0, "A": 1, "B": 2, "C": 3}
	actualMap := IndexLookup(items)
	if maps.Equal(wantMap, actualMap) == false {
		t.Errorf("IndexLookup() = %v, want %v", actualMap, wantMap)
	}
	// IndexOf
	actualIdx := IndexOf(items, "A")
	if actualIdx != 1 {
		t.Errorf("IndexOf() = %d, want %d", actualIdx, 1)
	}
	actualIdx = IndexOf(items, "C")
	if actualIdx != 3 {
		t.Errorf("IndexOf() = %d, want %d", actualIdx, 3)
	}
	actualIdx = IndexOf(items, "X")
	if actualIdx != -1 {
		t.Errorf("IndexOf() = %d, want %d", actualIdx, -1)
	}
	// AllIndexOf
	ints := []int{1, 2, 3, 1, 2, 3, 1}
	wantInts := []int{0, 3, 6}
	actualInts := AllIndexOf(ints, 1)
	if slices.Equal(actualInts, wantInts) == false {
		t.Errorf("AllIndexOf() = %d, want %d", actualInts, wantInts)
	}
	wantInts = []int{2, 5}
	actualInts = AllIndexOf(ints, 3)
	if slices.Equal(actualInts, wantInts) == false {
		t.Errorf("AllIndexOf() = %d, want %d", actualInts, wantInts)
	}
	wantInts = []int{}
	actualInts = AllIndexOf(ints, 69)
	if slices.Equal(actualInts, wantInts) == false {
		t.Errorf("AllIndexOf() = %d, want %d", actualInts, wantInts)
	}
	// GetOrDefault
	defaultValue := 69
	actual := GetOrDefault(ints, 3, defaultValue)
	if actual != 3 {
		t.Errorf("GetOrDefault() = %d, want %d", actual, 3)
	}
	actual = GetOrDefault(ints, 4, defaultValue)
	if actual != defaultValue {
		t.Errorf("GetOrDefault() = %d, want %d", actual, defaultValue)
	}
	// Remove
	ints2 := Copy(ints)
	wantInts = []int{1, 2, 1, 2, 3, 1}
	ints2, ok := Remove(ints2, 3)
	if !ok || slices.Equal(ints2, wantInts) == false {
		t.Errorf("Remove() = %v, %t, want %v, true", ints2, ok, wantInts)
	}
	ints2, ok = Remove(ints2, 69)
	if ok || slices.Equal(ints2, wantInts) == false {
		t.Errorf("Remove() = %v, %t, want %v, false", ints2, ok, wantInts)
	}
	// RemoveAll
	ints2 = Copy(ints)
	wantInts = []int{2, 3, 2, 3}
	ints2 = RemoveAll(ints2, 1)
	if slices.Equal(ints2, wantInts) == false {
		t.Errorf("RemoveAll() = %v, want %v", ints2, wantInts)
	}
	ints2 = RemoveAll(ints2, 5)
	if slices.Equal(ints2, wantInts) == false {
		t.Errorf("RemoveAll() = %v, want %v", ints2, wantInts)
	}
}

func TestTallyFunctions(t *testing.T) {
	// Count
	chars := []byte{'a', 'b', 'a', 'a', 'c', 'd', 'b', 'c'}
	countCases := [][2]int{
		{Count(chars, 'a'), 3},
		{Count(chars, 'd'), 1},
		{Count(chars, 'x'), 0},
	}
	for _, x := range countCases {
		actual, want := x[0], x[1]
		if actual != want {
			t.Errorf("Count() = %d, want %d", actual, want)
		}
	}
	// GroupByFunc
	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := map[string][]int{
		"odd":  {1, 3, 5, 7, 9},
		"even": {2, 4, 6, 8},
	}
	oddOrEven := func(x int) string {
		if x%2 == 0 {
			return "even"
		}
		return "odd"
	}
	actual := GroupByFunc(ints, oddOrEven)
	if reflect.DeepEqual(want, actual) == false {
		t.Errorf("GroupByFunc() = %v, want %v", actual, want)
	}
	// Tally
	wantTally := map[byte]int{
		'a': 3,
		'b': 2,
		'c': 2,
		'd': 1,
	}
	actualTally := Tally(chars)
	if maps.Equal(wantTally, actualTally) == false {
		t.Errorf("Tally() = %v, want %v", actualTally, wantTally)
	}
	// TallyFunc
	wantCounts := map[string]int{
		"odd":  5,
		"even": 4,
	}
	actualCounts := TallyFunc(ints, oddOrEven)
	if maps.Equal(wantCounts, actualCounts) == false {
		t.Errorf("TallyFunc() = %v, want %v", actualCounts, wantCounts)
	}
}

func TestUniqueFunctions(t *testing.T) {
	// TODO: CountUnique, CountUniqueFunc
	// TODO: AllSame, AllSameFunc
	// TODO: AllUnique, AllUniqueFunc
	// TODO: Deduplicate, DeduplicateFunc
}

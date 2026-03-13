package list

import (
	"slices"
	"strings"
	"testing"
)

func TestFn(t *testing.T) {
	// Filter, CountFunc
	numbers := []int{1, 2, 3, 4, 5, 6, 7}
	fn := func(x int) bool { return x%2 == 0 }
	want, wantCount := []int{2, 4, 6}, 3
	actual := Filter(numbers, fn)
	if slices.Equal(want, actual) == false {
		t.Errorf("Filter() = %v, want %v", actual, want)
	}
	actualCount := CountFunc(numbers, fn)
	if actualCount != wantCount {
		t.Errorf("CountFunc() = %v, want %v", actualCount, wantCount)
	}
	fn = func(x int) bool { return x > 10 }
	want, wantCount = []int{}, 0
	actual = Filter(numbers, fn)
	if slices.Equal(want, actual) == false {
		t.Errorf("Filter() = %v, want %v", actual, want)
	}
	actualCount = CountFunc(numbers, fn)
	if actualCount != wantCount {
		t.Errorf("CountFunc() = %v, want %v", actualCount, wantCount)
	}
	fn = func(x int) bool { return x <= 10 }
	want, wantCount = numbers, len(numbers)
	actual = Filter(numbers, fn)
	if slices.Equal(want, actual) == false {
		t.Errorf("Filter() = %v, want %v", actual, want)
	}
	actualCount = CountFunc(numbers, fn)
	if actualCount != wantCount {
		t.Errorf("CountFunc() = %v, want %v", actualCount, wantCount)
	}
	// FilterIndexed
	want = []int{1, 2, 4, 6, 7}
	actual = FilterIndexed(numbers, func(i, x int) bool { return x%2 == 0 || i%3 == 0 })
	if slices.Equal(want, actual) == false {
		t.Errorf("FilterIndexed() = %v, want %v", actual, want)
	}
	// Reduce
	wantSum := 28
	actualSum := Reduce(numbers, 0, func(result, item int) int {
		return result + item
	})
	if wantSum != actualSum {
		t.Errorf("Reduce() = %d, want %d", actualSum, wantSum)
	}
	// Apply
	want = []int{2, 4, 6, 8, 10, 12, 14}
	actual = Apply(numbers, func(x int) int { return x * 2 })
	if slices.Equal(want, actual) == false {
		t.Errorf("Apply() = %v, want %v", actual, want)
	}
}

func TestFnMap(t *testing.T) {
	// Map, MapIndexed
	type person struct {
		name string
		age  int
	}
	persons := []person{
		{"Alice", 20},
		{"Bob", 30},
		{"Charlie", 19},
		{"Dave", 15},
	}
	want := []string{"Alice", "Bob", "Charlie", "Dave"}
	actual := Map(persons, func(x person) string { return x.name })
	if slices.Equal(want, actual) == false {
		t.Errorf("Map() = %v, want %v", actual, want)
	}
	wantList := []int{20, 30, 19, 15}
	actualList := Map(persons, func(x person) int { return x.age })
	if slices.Equal(wantList, actualList) == false {
		t.Errorf("Map() = %v, want %v", actualList, wantList)
	}
	wantList = []int{0, 30, 38, 45}
	actualList = MapIndexed(persons, func(i int, x person) int { return i * x.age })
	if slices.Equal(wantList, actualList) == false {
		t.Errorf("MapIndexed() = %v, want %v", actualList, wantList)
	}
	// MapIf, MapIndexedIf
	want = []string{"Alice", "Bob"}
	actual = MapIf(persons, func(x person) (string, bool) { return x.name, x.age >= 20 })
	if slices.Equal(want, actual) == false {
		t.Errorf("MapIf() = %v, want %v", actual, want)
	}
	wantList = []int{30, 15}
	actualList = MapIf(persons, func(x person) (int, bool) { return x.age, x.age%3 == 0 })
	if slices.Equal(wantList, actualList) == false {
		t.Errorf("MapIf() = %v, want %v", actualList, wantList)
	}
	wantList = []int{20, 19, 15}
	actualList = MapIndexedIf(persons, func(i int, x person) (int, bool) { return x.age, i == 0 || x.age < 20 })
	if slices.Equal(wantList, actualList) == false {
		t.Errorf("MapIndexedIf() = %v, want %v", actualList, wantList)
	}
	// MapList
	items2 := []string{" ", "A", "B", "C", "D", "E"}
	indexes := []int{2, 5, 1, 4, 0, 3, 1, 2}
	result := MapList(indexes, items2)
	output := strings.Join(result, "")
	if output != "BEAD CAB" {
		t.Errorf("MapList() = %s, want [BEAD CAB]", output)
	}
	// MapLookup
	keys := []string{"a", "c"}
	weight := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	wantList = []int{1, 3}
	actualList = MapLookup(keys, weight)
	if slices.Equal(wantList, actualList) == false {
		t.Errorf("MapLookup() = %v, want %v", actualList, wantList)
	}
}

func TestSumProduct(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6}
	// Sum
	actual, want := Sum(items), 21
	if actual != want {
		t.Errorf("Sum() = %d, want %d", actual, want)
	}
	// Product
	actual, want = Product(items), 720
	if actual != want {
		t.Errorf("Product() = %d, want %d", actual, want)
	}
}

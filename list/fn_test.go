package list

import (
	"strings"
	"testing"

	"github.com/zeroibot/tst"
)

func TestFn(t *testing.T) {
	// Filter
	numbers := []int{1, 2, 3, 4, 5, 6, 7}
	fn1 := func(x int) bool { return x%2 == 0 }
	fn2 := func(x int) bool { return x > 10 }
	fn3 := func(x int) bool { return x <= 10 }

	testCases := []tst.P2W1[[]int, func(int) bool, []int]{
		{numbers, fn1, []int{2, 4, 6}},
		{numbers, fn2, []int{}},
		{numbers, fn3, numbers},
	}
	tst.AllP2W1(t, testCases, "Filter", Filter, tst.AssertListEqual)

	// CountFunc
	testCases2 := []tst.P2W1[[]int, func(int) bool, int]{
		{numbers, fn1, 3},
		{numbers, fn2, 0},
		{numbers, fn3, len(numbers)},
	}
	tst.AllP2W1(t, testCases2, "CountFunc", CountFunc, tst.AssertEqual)

	// FilterIndexed
	want := []int{1, 2, 4, 6, 7}
	actual := FilterIndexed(numbers, func(i, x int) bool { return x%2 == 0 || i%3 == 0 })
	tst.AssertListEqual(t, "FilterIndexed", actual, want)

	// Reduce
	wantSum := 28
	actualSum := Reduce(numbers, 0, func(result, item int) int {
		return result + item
	})
	tst.AssertEqual(t, "Reduce", actualSum, wantSum)

	// Apply
	want = []int{2, 4, 6, 8, 10, 12, 14}
	actual = Apply(numbers, func(x int) int { return x * 2 })
	tst.AssertListEqual(t, "Apply", actual, want)
}

func TestFnMap(t *testing.T) {
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
	// Map, MapIndexed
	want1 := []string{"Alice", "Bob", "Charlie", "Dave"}
	actual1 := Map(persons, func(x person) string { return x.name })
	tst.AssertListEqual(t, "Map", actual1, want1)

	want2 := []int{20, 30, 19, 15}
	actual2 := Map(persons, func(x person) int { return x.age })
	tst.AssertListEqual(t, "Map", actual2, want2)

	want2 = []int{0, 30, 38, 45}
	actual2 = MapIndexed(persons, func(i int, x person) int { return i * x.age })
	tst.AssertListEqual(t, "MapIndexed", actual2, want2)

	// MapIf, MapIndexedIf
	want1 = []string{"Alice", "Bob"}
	actual1 = MapIf(persons, func(x person) (string, bool) { return x.name, x.age >= 20 })
	tst.AssertListEqual(t, "MapIf", actual1, want1)

	want2 = []int{30, 15}
	actual2 = MapIf(persons, func(x person) (int, bool) { return x.age, x.age%3 == 0 })
	tst.AssertListEqual(t, "MapIf", actual2, want2)

	want2 = []int{20, 19, 15}
	actual2 = MapIndexedIf(persons, func(i int, x person) (int, bool) { return x.age, i == 0 || x.age < 20 })
	tst.AssertListEqual(t, "MapIndexedIf", actual2, want2)

	// MapList
	items2 := []string{" ", "A", "B", "C", "D", "E"}
	indexes := []int{2, 5, 1, 4, 0, 3, 1, 2}
	result := MapList(indexes, items2)
	tst.AssertEqual(t, "MapList", strings.Join(result, ""), "BEAD CAB")

	// MapLookup
	keys := []string{"a", "c"}
	weight := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	want2 = []int{1, 3}
	actual2 = MapLookup(keys, weight)
	tst.AssertListEqual(t, "MapLookup", actual2, want2)
}

func TestSumProduct(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6}
	tst.AssertEqual(t, "Sum", Sum(items), 21)
	tst.AssertEqual(t, "Product", Product(items), 720)
	// SumOf
	type Person struct {
		Age     int
		Balance float64
	}
	getAge := func(p Person) int { return p.Age }
	getBalance := func(p Person) float64 { return p.Balance }
	persons := []Person{
		{20, 15.0},
		{25, 18.0},
		{22, 20.0},
	}
	tst.AssertEqual(t, "SumOf", SumOf(persons, getAge), 67)
	tst.AssertEqual(t, "SumOf", SumOf(persons, getBalance), 53.0)

	// TODO: ProductOf
	// TODO: SumIndex
	// TODO: SumKey
	// TODO: ProductIndex
	// TODO: ProductKey
}

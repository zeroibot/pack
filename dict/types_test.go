package dict

import (
	"fmt"
	"slices"
	"testing"

	"github.com/roidaradal/tst"
)

func TestNewCounterFor(t *testing.T) {
	keys := []string{"a", "b", "c"}
	counter := NewCounterFor(keys)
	counterKeys := Keys(counter)
	slices.Sort(counterKeys)
	tst.AssertListEqual(t, "NewCounterFor.Keys", counterKeys, keys)
	for _, key := range keys {
		tst.AssertEqual(t, fmt.Sprintf("Counter[%s]", key), counter[key], 0)
	}
}

func TestNewCounterFunc(t *testing.T) {
	type person struct {
		name    string
		age     int
		friends []person
	}
	people := []person{
		{"John", 25, nil},
		{"Peter", 26, nil},
		{"Alice", 24, nil},
	}
	keys := []string{"Alice", "John", "Peter"}
	counter := NewCounterFunc(people, func(p person) string { return p.name })
	counterKeys := Keys(counter)
	slices.Sort(counterKeys)
	tst.AssertListEqual(t, "NewCounterFunc.Keys", counterKeys, keys)
	for _, key := range keys {
		tst.AssertEqual(t, fmt.Sprintf("Counter[%s]", key), counter[key], 0)
	}
}

func TestUpdateCounter(t *testing.T) {
	keys := []int{1, 2, 3, 4}
	items := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	counter := NewCounterFor(keys)
	UpdateCounter(counter, items)
	for _, key := range keys {
		tst.AssertEqual(t, fmt.Sprintf("Counter[%d]", key), counter[key], key)
	}
}

func TestUpdateCounterFunc(t *testing.T) {
	type student struct {
		name  string
		grade int
	}
	keys := []int{1, 2, 3, 4, 5}
	students := []student{
		{"John", 1},
		{"Peter", 2},
		{"Alice", 1},
		{"Bob", 3},
		{"Cat", 4},
		{"Dan", 2},
	}
	counter := NewCounterFor(keys)
	UpdateCounterFunc(counter, students, func(s student) int { return s.grade })
	want := Counter[int]{1: 2, 2: 2, 3: 1, 4: 1, 5: 0}
	tst.AssertMapEqual(t, "UpdateCounterFunc", counter, want)
}

func TestNewFlagsFor(t *testing.T) {
	keys := []string{"a", "b", "c"}
	flags := NewFlagsFor(keys, true)
	flagKeys := Keys(flags)
	slices.Sort(flagKeys)
	tst.AssertListEqual(t, "NewFlagsFor.Keys", flagKeys, keys)
	for _, key := range keys {
		tst.AssertEqual(t, fmt.Sprintf("Flags[%s]", key), flags[key], true)
	}
}

func TestNewFlagsFunc(t *testing.T) {
	type person struct {
		name    string
		age     int
		friends []person
	}
	people := []person{
		{"John", 25, nil},
		{"Peter", 26, nil},
		{"Alice", 24, nil},
	}
	keys := []string{"Alice", "John", "Peter"}
	flags := NewFlagsFunc(people, true, func(p person) string { return p.name })
	flagKeys := Keys(flags)
	slices.Sort(flagKeys)
	tst.AssertListEqual(t, "NewFlagsFunc.Keys", flagKeys, keys)
	for _, key := range keys {
		tst.AssertEqual(t, fmt.Sprintf("Flags[%s]", key), flags[key], true)
	}
}

func TestLookupFunc(t *testing.T) {
	lookup := map[string]int{"John": 1, "Peter": 2, "Alice": 3}
	fn := LookupFunc(lookup)
	for key, want := range lookup {
		actual, ok := fn(key)
		tst.AssertEqualAnd(t, fmt.Sprintf("LookupFunc(%s)", key), actual, want, ok, true)
	}
	for _, key := range []string{"Bob", "Cat"} {
		actual, ok := fn(key)
		tst.AssertEqualAnd(t, fmt.Sprintf("LookupFunc(%s)", key), actual, 0, ok, false)
	}
}

func TestMustLookupFunc(t *testing.T) {
	lookup := map[string]int{"John": 1, "Peter": 2, "Alice": 3}
	fn := MustLookupFunc(lookup)
	for key, want := range lookup {
		tst.AssertEqual(t, fmt.Sprintf("MustLookupFunc(%s)", key), fn(key), want)
	}

	defer tst.AssertPanic(t, "MustLookupFunc")
	fn("Bob") // should panic
}

func TestGet(t *testing.T) {
	obj := Object{"name": "John", "age": 25, "isActive": true}
	actual1, ok := Get[string](obj, "name")
	tst.AssertEqualAnd(t, "Get(name)", actual1, "John", ok, true)
	actual2, ok := Get[int](obj, "age")
	tst.AssertEqualAnd(t, "Get(age)", actual2, 25, ok, true)

	// Test non-existent key
	actual3, ok := Get[int](obj, "salary")
	tst.AssertEqualAnd(t, "Get(salary)", actual3, 0, ok, false)

	// Test wrong key type
	actual4, ok := Get[int](obj, "isActive")
	tst.AssertEqualAnd(t, "Get(isActive)", actual4, 0, ok, false)
}

func TestGetRef(t *testing.T) {
	var lastUpdated *string = nil
	ipAddress := new("127.0.0.1")
	obj := Object{
		"name":        "John",
		"age":         25,
		"lastUpdated": lastUpdated,
		"ipAddress":   ipAddress,
	}

	actual1 := GetRef[string](obj, "ipAddress")
	tst.AssertEqual(t, "GetRef(ipAddress)", actual1, ipAddress)
	actual2 := GetRef[string](obj, "lastUpdated")
	tst.AssertEqual(t, "GetRef(lastUpdated)", actual2, lastUpdated)

	// Non-existent key
	actual3 := GetRef[string](obj, "browserInfo")
	tst.AssertEqual(t, "GetRef(browserInfo)", actual3, nil)

	// Wrong type
	actual4 := GetRef[string](obj, "name")
	tst.AssertEqual(t, "GetRef(name)", actual4, nil)
	actual5 := GetRef[int](obj, "ipAddress")
	tst.AssertEqual(t, "GetRef(ipAddress)", actual5, nil)
}

func TestGetList(t *testing.T) {
	names := []string{"John", "Peter"}
	ages := []int{25, 26}
	obj := Object{"names": names, "ages": ages}

	actual1 := GetList[string](obj, "names")
	tst.AssertListEqual(t, "GetList(names)", actual1, names)
	actual2 := GetList[int](obj, "ages")
	tst.AssertListEqual(t, "GetList(ages)", actual2, ages)

	// Non-existent key
	actual3 := GetList[int](obj, "salaries")
	tst.AssertListEqual(t, "GetList(salaries)", actual3, nil)

	// Wrong type
	actual4 := GetList[int](obj, "names")
	tst.AssertListEqual(t, "GetList(names)", actual4, nil)
}

func TestCounterUpdate(t *testing.T) {
	c1 := Counter[string]{"apple": 5, "banana": 3}
	c2 := Counter[string]{"apricot": 8, "banana": 7, "cherry": 2}
	expectedEntries := []Entry[string, int]{
		{"apple", 5}, {"apricot", 8}, {"banana", 10}, {"cherry", 2},
	}
	CounterUpdate(c1, c2)
	tst.AssertListEqual(t, "CounterUpdate.Entries", SortedEntries(c1), expectedEntries)
}

func TestMergeCounters(t *testing.T) {
	c1 := Counter[string]{"apple": 5, "banana": 3}
	c2 := Counter[string]{"apricot": 8, "banana": 7, "cherry": 2}
	expectedEntries := []Entry[string, int]{
		{"apple", 5}, {"apricot", 8}, {"banana", 10}, {"cherry", 2},
	}
	c3 := MergeCounters(c1, c2)
	tst.AssertListEqual(t, "MergeCounters.Entries", SortedEntries(c3), expectedEntries)

	expectedEntries1 := []Entry[string, int]{{"apple", 5}, {"banana", 3}}
	expectedEntries2 := []Entry[string, int]{{"apricot", 8}, {"banana", 7}, {"cherry", 2}}
	tst.AssertListEqual(t, "MergeCounters.Entries", SortedEntries(c1), expectedEntries1)
	tst.AssertListEqual(t, "MergeCounters.Entries", SortedEntries(c2), expectedEntries2)
}

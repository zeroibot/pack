package dict

import (
	"maps"
	"slices"
	"testing"
)

func TestNewCounterFor(t *testing.T) {
	keys := []string{"a", "b", "c"}
	counter := NewCounterFor(keys)
	counterKeys := Keys(counter)
	slices.Sort(counterKeys)
	if slices.Equal(keys, counterKeys) == false {
		t.Errorf("NewCounterFor Keys = %v, want %v", counterKeys, keys)
	}
	for _, key := range keys {
		actual := counter[key]
		if actual != 0 {
			t.Errorf("Counter[%s] = %v, want 0", key, actual)
		}
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
	if slices.Equal(keys, counterKeys) == false {
		t.Errorf("NewCounterFunc Keys = %v, want %v", counterKeys, keys)
	}
	for _, key := range keys {
		actual := counter[key]
		if actual != 0 {
			t.Errorf("Counter[%s] = %v, want 0", key, actual)
		}
	}
}

func TestUpdateCounter(t *testing.T) {
	keys := []int{1, 2, 3, 4}
	counter := NewCounterFor(keys)
	items := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	UpdateCounter(counter, items)
	for _, key := range keys {
		actual := counter[key]
		if actual != key {
			t.Errorf("Counter[%d] = %d, want %d", key, actual, key)
		}
	}
}

func TestUpdateCounterFunc(t *testing.T) {
	type student struct {
		name  string
		grade int
	}
	keys := []int{1, 2, 3, 4, 5}
	counter := NewCounterFor(keys)
	students := []student{
		{"John", 1},
		{"Peter", 2},
		{"Alice", 1},
		{"Bob", 3},
		{"Cat", 4},
		{"Dan", 2},
	}
	expected := Counter[int]{
		1: 2,
		2: 2,
		3: 1,
		4: 1,
		5: 0,
	}
	UpdateCounterFunc(counter, students, func(s student) int { return s.grade })
	if maps.Equal(counter, expected) == false {
		t.Errorf("UpdateCounterFunc = %v, want %v", counter, expected)
	}
}

func TestNewFlagsFor(t *testing.T) {
	keys := []string{"a", "b", "c"}
	flags := NewFlagsFor(keys, true)
	flagKeys := Keys(flags)
	slices.Sort(flagKeys)
	if slices.Equal(flagKeys, keys) == false {
		t.Errorf("NewFlagsFor Keys = %v, want %v", flagKeys, keys)
	}
	for _, key := range keys {
		actual := flags[key]
		if actual != true {
			t.Errorf("Flags[%s] = %v, want true", key, actual)
		}
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
	if slices.Equal(keys, flagKeys) == false {
		t.Errorf("NewFlagsFunc Keys = %v, want %v", flagKeys, keys)
	}
	for _, key := range keys {
		actual := flags[key]
		if actual != true {
			t.Errorf("Flags[%s] = %v, want true", key, actual)
		}
	}
}

func TestLookupFunc(t *testing.T) {
	lookup := map[string]int{
		"John":  1,
		"Peter": 2,
		"Alice": 3,
	}
	fn := LookupFunc(lookup)
	for key, want := range lookup {
		actual, ok := fn(key)
		if actual != want || !ok {
			t.Errorf("LookupFunc(%s) = %v, %v want %v, true", key, actual, ok, want)
		}
	}
	for _, key := range []string{"Bob", "Cat"} {
		actual, ok := fn(key)
		if actual != 0 || ok != false {
			t.Errorf("LookupFunc(%s) = %v, %v want 0, false", key, actual, ok)
		}
	}
}

func TestMustLookupFunc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustLookupFunc did not panic")
		}
	}()
	lookup := map[string]int{
		"John":  1,
		"Peter": 2,
		"Alice": 3,
	}
	fn := MustLookupFunc(lookup)
	for key, want := range lookup {
		actual := fn(key)
		if actual != want {
			t.Errorf("MustLookupFunc(%s) = %v want %v", key, actual, want)
		}
	}
	// Test for panic
	fn("Bob")
}

func TestGet(t *testing.T) {
	obj := Object{
		"name":     "John",
		"age":      25,
		"isActive": true,
	}
	want1 := "John"
	actual1, ok := Get[string](obj, "name")
	if actual1 != want1 || !ok {
		t.Errorf("Get(name) = %q, %v want %q, true", actual1, ok, want1)
	}
	want2 := 25
	actual2, ok := Get[int](obj, "age")
	if actual2 != want2 || !ok {
		t.Errorf("Get(age) = %d, want %d, true", actual2, want2)
	}
	// Test non-existent key
	want3 := 0
	actual3, ok := Get[int](obj, "salary")
	if actual3 != want3 || ok {
		t.Errorf("Get(salary) = %d, %v, want %d, false", actual3, ok, want3)
	}
	// Test wrong key type
	want4 := 0
	actual4, ok := Get[int](obj, "isActive")
	if actual4 != want4 || ok {
		t.Errorf("Get(isActive) = %d, %v, want %d, false", actual4, ok, want4)
	}
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
	if actual1 != ipAddress {
		t.Errorf("GetRef(browserInfo) = %v, want %v", actual1, ipAddress)
	}
	actual2 := GetRef[string](obj, "lastUpdated")
	if actual2 != lastUpdated {
		t.Errorf("GetRef(browserInfo) = %v, want %v", actual2, lastUpdated)
	}
	// Non-existent key
	actual3 := GetRef[string](obj, "browserInfo")
	if actual3 != nil {
		t.Errorf("GetRef(browserInfo) = %v, want nil", actual3)
	}
	// Wrong type
	actual4 := GetRef[string](obj, "name")
	if actual4 != nil {
		t.Errorf("GetRef(name) = %v, want nil", actual4)
	}
	actual5 := GetRef[int](obj, "ipAddress")
	if actual5 != nil {
		t.Errorf("GetRef(ipAddress) = %v, want nil", actual5)
	}
}

func TestGetList(t *testing.T) {
	names := []string{"John", "Peter"}
	ages := []int{25, 26}
	obj := Object{
		"names": names,
		"ages":  ages,
	}
	actual1 := GetList[string](obj, "names")
	if slices.Equal(names, actual1) == false {
		t.Errorf("GetList(names) = %v, want %v", actual1, names)
	}
	actual2 := GetList[int](obj, "ages")
	if slices.Equal(ages, actual2) == false {
		t.Errorf("GetList(ages) = %v, want %v", actual2, ages)
	}
	// Non-existent key
	actual3 := GetList[int](obj, "salaries")
	if actual3 != nil {
		t.Errorf("GetList(salaries) = %v, want nil", actual3)
	}
	// Wrong type
	actual4 := GetList[int](obj, "names")
	if actual4 != nil {
		t.Errorf("GetList(names) = %v, want nil", actual4)
	}
}

func TestCounterUpdate(t *testing.T) {
	c1 := Counter[string]{
		"apple":  5,
		"banana": 3,
	}
	c2 := Counter[string]{
		"apricot": 8,
		"banana":  7,
		"cherry":  2,
	}
	expectedEntries := []Entry[string, int]{
		{"apple", 5}, {"apricot", 8}, {"banana", 10}, {"cherry", 2},
	}
	CounterUpdate(c1, c2)
	actualEntries := SortedEntries(c1)
	if slices.Equal(expectedEntries, actualEntries) == false {
		t.Errorf("CounterUpdate.Entries = %v, want %v", actualEntries, expectedEntries)
	}
}

func TestMergeCounters(t *testing.T) {
	c1 := Counter[string]{
		"apple":  5,
		"banana": 3,
	}
	c2 := Counter[string]{
		"apricot": 8,
		"banana":  7,
		"cherry":  2,
	}
	expectedEntries := []Entry[string, int]{
		{"apple", 5}, {"apricot", 8}, {"banana", 10}, {"cherry", 2},
	}
	c3 := MergeCounters(c1, c2)
	actualEntries := SortedEntries(c3)
	if slices.Equal(expectedEntries, actualEntries) == false {
		t.Errorf("MergeCounters.Entries = %v, want %v", actualEntries, expectedEntries)
	}
	expectedEntries1 := []Entry[string, int]{{"apple", 5}, {"banana", 3}}
	expectedEntries2 := []Entry[string, int]{{"apricot", 8}, {"banana", 7}, {"cherry", 2}}
	actualEntries1 := SortedEntries(c1)
	actualEntries2 := SortedEntries(c2)
	if slices.Equal(expectedEntries1, actualEntries1) == false {
		t.Errorf("MergeCounters.Entries = %v, want %v", actualEntries1, expectedEntries1)
	}
	if slices.Equal(expectedEntries2, actualEntries2) == false {
		t.Errorf("MergeCounters.Entries = %v, want %v", actualEntries2, expectedEntries2)
	}
}

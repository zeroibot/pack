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
			t.Errorf("Counter[%s].Value = %v, want 0", key, actual)
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
			t.Errorf("Counter[%s].Value = %v, want 0", key, actual)
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
			t.Errorf("Counter[%d].Value = %d, want %d", key, actual, key)
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

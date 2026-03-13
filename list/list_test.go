package list

import (
	"fmt"
	"slices"
	"testing"
)

func TestList(t *testing.T) {
	type testCase struct {
		name         string
		actual, want int
	}
	// Range, InclusiveRange
	r1 := Range(1, 5)
	r2 := InclusiveRange(1, 5)
	want1 := []int{1, 2, 3, 4}
	want2 := []int{1, 2, 3, 4, 5}
	if slices.Equal(r1, want1) == false {
		t.Errorf("Range() = %v, want %v", r1, want1)
	}
	if slices.Equal(r2, want2) == false {
		t.Errorf("InclusiveRange() = %v, want %v", r2, want2)
	}
	// RepeatedItem
	r1 = RepeatedItem(5, 3)
	r2 = RepeatedItem(3, 5)
	want1 = []int{5, 5, 5}
	want2 = []int{3, 3, 3, 3, 3}
	if slices.Equal(r1, want1) == false {
		t.Errorf("RepeatedItem() = %v, want %v", r1, want1)
	}
	if slices.Equal(r2, want2) == false {
		t.Errorf("RepeatedItem() = %v, want %v", r2, want2)
	}
	// NewEmpty, Len, Cap, LastIndex
	l1 := NewEmpty[int](5)
	l2 := []string{"a", "b", "c"}
	testCases := []testCase{
		{"Len", Len(l1), 0},
		{"Cap", Cap(l1), 5},
		{"LastIndex", LastIndex(l1), -1},
		{"Len", Len(l2), 3},
		{"Cap", Cap(l2), 3},
		{"LastIndex", LastIndex(l2), 2},
	}
	for _, x := range testCases {
		if x.actual != x.want {
			t.Errorf("%s() = %d, want %d", x.name, x.actual, x.want)
		}
	}
	// IsEmpty, NotEmpty
	if IsEmpty(l1) != true {
		t.Errorf("IsEmpty() = %v, want true", IsEmpty(l1))
	}
	if NotEmpty(l2) != true {
		t.Errorf("NotEmpty() = %v, want true", NotEmpty(l2))
	}
	// Copy
	l3 := Copy(l2)
	if slices.Equal(l2, l3) == false {
		t.Errorf("Copy() = %v, want %v", l3, l2)
	}
	l2[0] = "x"
	l3[1] = "r"
	actual, want := fmt.Sprintf("%v", l2), "[x b c]"
	if actual != want {
		t.Errorf("List.String() = %s, want %s", actual, want)
	}
	actual, want = fmt.Sprintf("%v", l3), "[a r c]"
	if actual != want {
		t.Errorf("List.String() = %s, want %s", actual, want)
	}
}

func TestListRandom(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustGetRandom() did not panic")
		}
	}()
	// GetRandom and MustGetRandom
	l1 := NewEmpty[int](3) // empty
	for range 5 {
		item, ok := GetRandom(l1)
		if ok || item != 0 {
			t.Errorf("EmptyList.GetRandom() = %d, %t, want 0, false", item, ok)
		}
	}
	l := InclusiveRange(1, 100)
	for range 100 {
		value, ok := GetRandom(l)
		if !ok || !(1 <= value && value <= 100) {
			t.Errorf("GetRandom() = %v, want 1..100", value)
		}
		value = MustGetRandom(l)
		if !(1 <= value && value <= 100) {
			t.Errorf("MustGetRandom() = %v, want 1..100", value)
		}
	}
	// Shuffle
	l2 := []int{1, 2, 3, 4, 5, 6, 7}
	l3 := Copy(l2)
	Shuffle(l3)
	if slices.Equal(l2, l3) == true {
		t.Errorf("Shuffle() = %v, want not original %v", l3, l2)
	}

	MustGetRandom(l1) // should panic (empty list)
}

func TestListMethods(t *testing.T) {
	// TODO: ToAny
	// TODO: IndexFunc, AllIndexFunc
	// TODO: RemoveFunc, RemoveAllFunc
	// TODO: GetFuncOrDefault
	// TODO: Last, MustLast
}

func TestListCheck(t *testing.T) {
	// Any, NotAny
	items := []int{1, 2, 3, 4, 5, 6}
	fn1 := func(x int) bool { return x%2 == 0 && x%3 == 0 }
	fn2 := func(x int) bool { return x > 10 }
	fn3 := func(x int) bool { return x <= 10 }
	result := Any(items, fn1)
	if result != true {
		t.Errorf("Any() = %v, want true", result)
	}
	result = NotAny(items, fn1)
	if result != false {
		t.Errorf("NotAny() = %v, want false", result)
	}
	result = Any(items, fn2)
	if result != false {
		t.Errorf("Any() = %v, want false", result)
	}
	result = NotAny(items, fn2)
	if result != true {
		t.Errorf("NotAny() = %v, want true", result)
	}
	// All
	var empty []int
	result = All(empty, fn1)
	if result != false {
		t.Errorf("All() = %v, want false", result)
	}
	result = All(items, fn1)
	if result != false {
		t.Errorf("All() = %v, want false", result)
	}
	result = All(items, fn3)
	if result != true {
		t.Errorf("All() = %v, want true", result)
	}
	// AnyIndexed, NotAnyIndexed
	fn4 := func(i, x int) bool { return i >= 0 && x%2 == 0 && x%3 == 0 }
	fn5 := func(i, x int) bool { return i > 10 && x > 10 }
	fn6 := func(i, x int) bool { return i < 10 && x <= 10 }
	result = AnyIndexed(items, fn4)
	if result != true {
		t.Errorf("AnyIndexed() = %v, want true", result)
	}
	result = NotAnyIndexed(items, fn4)
	if result != false {
		t.Errorf("NotAnyIndexed() = %v, want false", result)
	}
	result = AnyIndexed(items, fn5)
	if result != false {
		t.Errorf("AnyIndexed() = %v, want false", result)
	}
	result = NotAnyIndexed(items, fn5)
	if result != true {
		t.Errorf("NotAnyIndexed() = %v, want true", result)
	}
	// AllIndexed
	result = AllIndexed(empty, fn4)
	if result != false {
		t.Errorf("AllIndexed() = %v, want false", result)
	}
	result = AllIndexed(items, fn4)
	if result != false {
		t.Errorf("AllIndexed() = %v, want false", result)
	}
	result = AllIndexed(items, fn6)
	if result != true {
		t.Errorf("AllIndexed() = %v, want true", result)
	}
}

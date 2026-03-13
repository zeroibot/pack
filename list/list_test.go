package list

import "testing"

func TestList(t *testing.T) {
	// TODO: NewEmpty
	// TODO: Range, InclusiveRange
	// TODO: RepeatedItem
	// TODO: Len, Cap, LastIndex
	// TODO: IsEmpty, NotEmpty
	// TODO: Copy
}

func TestListMethods(t *testing.T) {
	// TODO: ToAny
	// TODO: IndexFunc, AllIndexFunc
	// TODO: RemoveFunc, RemoveAllFunc
	// TODO: GetFuncOrDefault
	// TODO: Last, MustLast
	// TODO: GetRandom, MustGetRandom
	// TODO: Shuffle
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

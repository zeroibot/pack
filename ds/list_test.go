package ds

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	l1 := NewList[int](5)
	l2 := List[string]{"a", "b", "c"}
	testCases := []Tuple3[string, int, int]{
		{"Len", l1.Len(), 0},
		{"Cap", l1.Cap(), 5},
		{"LastIndex", l1.LastIndex(), -1},
		{"Len", l2.Len(), 3},
		{"Cap", l2.Cap(), 3},
		{"LastIndex", l2.LastIndex(), 2},
	}
	for _, x := range testCases {
		name, actual, want := x.Values()
		if actual != want {
			t.Errorf("List.%s = %d, want %d", name, actual, want)
		}
	}
	if l1.IsEmpty() != true {
		t.Errorf("List.Empty = %v, want true", l1.IsEmpty())
	}
	if l2.NotEmpty() != true {
		t.Errorf("List.NotEmpty = %v, want true", l2.NotEmpty())
	}
	l3 := l2.Copy()
	if slices.Equal(l2, l3) == false {
		t.Errorf("List.Copy() = %v, want %v", l3, l2)
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

func TestListMethods(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("List.MustLast() did not panic")
		}
	}()
	// ToAny
	items := List[int]{1, 2, 3}
	anyItems := items.ToAny()
	actualString, wantString := fmt.Sprintf("%v", anyItems), "[1 2 3]"
	if actualString != wantString {
		t.Errorf("List.ToAny() = %s, want %s", actualString, wantString)
	}
	// IndexFunc
	items = List[int]{1, 2, 3, 4, 1, 2, 4, 2, 5, 3}
	wantIndex := 2
	actualIndex := items.IndexFunc(func(x int) bool { return x == 3 })
	if wantIndex != actualIndex {
		t.Errorf("List.IndexFunc() = %d, want %d", actualIndex, wantIndex)
	}
	wantIndex = 8
	actualIndex = items.IndexFunc(func(x int) bool { return x == 5 })
	if wantIndex != actualIndex {
		t.Errorf("List.IndexFunc() = %d, want %d", actualIndex, wantIndex)
	}
	// AllIndexFunc
	wantList := List[int]{1, 5, 7}
	actualList := items.AllIndexFunc(func(x int) bool { return x == 2 })
	if slices.Equal(actualList, wantList) == false {
		t.Errorf("List.AllIndexFunc() = %v, wantList %v", actualList, wantList)
	}
	wantList = List[int]{1, 3, 5, 6, 7}
	actualList = items.AllIndexFunc(func(x int) bool { return x%2 == 0 })
	if slices.Equal(actualList, wantList) == false {
		t.Errorf("List.AllIndexFunc() = %v, wantList %v", actualList, wantList)
	}
	// RemoveFunc
	items = List[int]{1, 2, 1, 2, 3, 2}
	items2 := items.Copy()
	items2, ok := items2.RemoveFunc(func(x int) bool { return x == 2 })
	wantList = List[int]{1, 1, 2, 3, 2}
	if !ok || slices.Equal(wantList, items2) == false {
		t.Errorf("List.RemoveFunc() = %v, %t, wantList %v, true", items2, ok, wantList)
	}
	items2, ok = items2.RemoveFunc(func(x int) bool { return x == 4 })
	if ok || slices.Equal(wantList, items2) == false {
		t.Errorf("List.RemoveFunc() = %v, %t, wantList %v, false", items2, ok, wantList)
	}
	// RemoveAllFunc
	items2 = items.Copy()
	items2 = items2.RemoveAllFunc(func(x int) bool { return x == 2 })
	wantList = List[int]{1, 1, 3}
	if slices.Equal(wantList, items2) == false {
		t.Errorf("List.RemoveAllFunc() = %v, wantList %v", items, wantList)
	}
	// Get
	items = List[int]{1, 2, 3}
	option := items.Get(1)
	if option.IsNil() || option.Value() != 2 {
		t.Errorf("List.Get() = %v, want 2", option)
	}
	option = items.Get(3)
	if option.NotNil() {
		t.Errorf("List.Get() = %v, want nil", option)
	}
	option = items.Get(-1)
	if option.NotNil() {
		t.Errorf("List.Get() = %v, want nil", option)
	}
	// GetFuncOrDefault
	defaultValue := 69
	actual := items.GetFuncOrDefault(func(x int) bool { return x == 3 }, defaultValue)
	if actual != 3 {
		t.Errorf("List.GetFuncOrDefault() = %d, want 3", actual)
	}
	actual = items.GetFuncOrDefault(func(x int) bool { return x == 4 }, defaultValue)
	if actual != defaultValue {
		t.Errorf("List.GetFuncOrDefault() = %d, want %d", actual, defaultValue)
	}
	// Last
	option = items.Last(1)
	if option.IsNil() || option.Value() != 3 {
		t.Errorf("List.Last() = %v, want 3", option)
	}
	option = items.Last(3)
	if option.IsNil() || option.Value() != 1 {
		t.Errorf("List.Last() = %v, want 1", option)
	}
	option = items.Last(0)
	if option.NotNil() {
		t.Errorf("List.Last() = %v, want nil", option)
	}
	option = items.Last(4)
	if option.NotNil() {
		t.Errorf("List.Last() = %v, want nil", option)
	}
	// MustLast
	actual = items.MustLast(1)
	if actual != 3 {
		t.Errorf("List.MustLast() = %d, want 3", actual)
	}
	actual = items.MustLast(3)
	if actual != 1 {
		t.Errorf("List.MustLast() = %d, want 1", actual)
	}
	items.MustLast(4) // should panic
}

func TestListRandom(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("List.MustGetRandom() did not panic")
		}
	}()
	// GetRandom and MustGetRandom
	l1 := NewList[int](3) // empty
	for range 5 {
		item := l1.GetRandom()
		if !item.IsNil() {
			t.Errorf("EmptyList.GetRandom() = %v, want nil", item)
		}
	}
	l := NewInclusiveRange(1, 100).ToList()
	for range 100 {
		item := l.GetRandom()
		value := item.Value()
		if item.IsNil() || !(1 <= value && value <= 100) {
			t.Errorf("List.GetRandom() = %v, want 1..100", item)
		}
		value = l.MustGetRandom()
		if !(1 <= value && value <= 100) {
			t.Errorf("List.MustGetRandom() = %v, want 1..100", item)
		}
	}
	// Shuffle
	l2 := List[int]{1, 2, 3, 4, 5, 6, 7}
	l3 := l2.Copy()
	l3.Shuffle()
	if slices.Equal(l2, l3) == true {
		t.Errorf("List.Shuffle = %v, want not original %v", l3, l2)
	}

	l1.MustGetRandom() // should panic (empty list)
}

func TestListCheck(t *testing.T) {
	// Any, NotAny
	items := List[int]{1, 2, 3, 4, 5, 6}
	fn1 := func(x int) bool { return x%2 == 0 && x%3 == 0 }
	fn2 := func(x int) bool { return x > 10 }
	fn3 := func(x int) bool { return x <= 10 }
	result := items.Any(fn1)
	if result != true {
		t.Errorf("List.Any() = %v, want true", result)
	}
	result = items.NotAny(fn1)
	if result != false {
		t.Errorf("List.NotAny() = %v, want false", result)
	}
	result = items.Any(fn2)
	if result != false {
		t.Errorf("List.Any() = %v, want false", result)
	}
	result = items.NotAny(fn2)
	if result != true {
		t.Errorf("List.NotAny() = %v, want true", result)
	}
	// All
	empty := List[int]{}
	result = empty.All(fn1)
	if result != false {
		t.Errorf("List.All() = %v, want false", result)
	}
	result = items.All(fn1)
	if result != false {
		t.Errorf("List.All() = %v, want false", result)
	}
	result = items.All(fn3)
	if result != true {
		t.Errorf("List.All() = %v, want true", result)
	}
	// AnyIndexed, NotAnyIndexed
	fn4 := func(i, x int) bool { return i >= 0 && x%2 == 0 && x%3 == 0 }
	fn5 := func(i, x int) bool { return i > 10 && x > 10 }
	fn6 := func(i, x int) bool { return i < 10 && x <= 10 }
	result = items.AnyIndexed(fn4)
	if result != true {
		t.Errorf("List.AnyIndexed() = %v, want true", result)
	}
	result = items.NotAnyIndexed(fn4)
	if result != false {
		t.Errorf("List.NotAnyIndexed() = %v, want false", result)
	}
	result = items.AnyIndexed(fn5)
	if result != false {
		t.Errorf("List.AnyIndexed() = %v, want false", result)
	}
	result = items.NotAnyIndexed(fn5)
	if result != true {
		t.Errorf("List.NotAnyIndexed() = %v, want true", result)
	}
	// AllIndexed
	result = empty.AllIndexed(fn4)
	if result != false {
		t.Errorf("List.AllIndexed() = %v, want false", result)
	}
	result = items.AllIndexed(fn4)
	if result != false {
		t.Errorf("List.AllIndexed() = %v, want false", result)
	}
	result = items.AllIndexed(fn6)
	if result != true {
		t.Errorf("List.AllIndexed() = %v, want true", result)
	}
}

func TestListFn(t *testing.T) {
	// MapList
	items := List[string]{" ", "A", "B", "C", "D", "E"}
	indexes := []int{2, 5, 1, 4, 0, 3, 1, 2}
	result := items.MapList(indexes)
	output := strings.Join(result, "")
	if output != "BEAD CAB" {
		t.Errorf("List.MapList() = %s, want [BEAD CAB]", output)
	}
	// Filter && CountFunc
	numbers := List[int]{1, 2, 3, 4, 5, 6, 7}
	fn := func(x int) bool { return x%2 == 0 }
	want, wantCount := List[int]{2, 4, 6}, 3
	actual := numbers.Filter(fn)
	if slices.Equal(want, actual) == false {
		t.Errorf("List.Filter() = %v, want %v", actual, want)
	}
	actualCount := numbers.CountFunc(fn)
	if actualCount != wantCount {
		t.Errorf("List.CountFunc() = %v, want %v", actualCount, wantCount)
	}
	fn = func(x int) bool { return x > 10 }
	want, wantCount = List[int]{}, 0
	actual = numbers.Filter(fn)
	if slices.Equal(want, actual) == false {
		t.Errorf("List.Filter() = %v, want %v", actual, want)
	}
	actualCount = numbers.CountFunc(fn)
	if actualCount != wantCount {
		t.Errorf("List.CountFunc() = %v, want %v", actualCount, wantCount)
	}
	fn = func(x int) bool { return x <= 10 }
	want, wantCount = numbers, len(numbers)
	actual = numbers.Filter(fn)
	if slices.Equal(want, actual) == false {
		t.Errorf("List.Filter() = %v, want %v", actual, want)
	}
	actualCount = numbers.CountFunc(fn)
	if actualCount != wantCount {
		t.Errorf("List.CountFunc() = %v, want %v", actualCount, wantCount)
	}
	// FilterIndexed
	want = List[int]{1, 2, 4, 6, 7}
	actual = numbers.FilterIndexed(func(i, x int) bool { return x%2 == 0 || i%3 == 0 })
	if slices.Equal(want, actual) == false {
		t.Errorf("List.FilterIndexed() = %v, want %v", actual, want)
	}
	// Reduce
	wantSum := 28
	actualSum := numbers.Reduce(0, func(result, item int) int {
		return result + item
	})
	if wantSum != actualSum {
		t.Errorf("List.Reduce() = %d, want %d", actualSum, wantSum)
	}
	// Apply
	want = List[int]{2, 4, 6, 8, 10, 12, 14}
	actual = numbers.Apply(func(x int) int { return x * 2 })
	if slices.Equal(want, actual) == false {
		t.Errorf("List.Apply() = %v, want %v", actual, want)
	}
}

func TestNumList(t *testing.T) {
	// ToList
	n := NumList[int]{1, 2, 3, 4, 5, 6}
	l := n.ToList()
	if l.Len() != 6 {
		t.Errorf("NumList.ToList.Len() = %d, want 5", l.Len())
	}
	// Sum
	actual, want := n.Sum(), 21
	if actual != want {
		t.Errorf("NumList.Sum() = %d, want %d", actual, want)
	}
	// Product
	actual, want = n.Product(), 720
	if actual != want {
		t.Errorf("NumList.Product() = %d, want %d", actual, want)
	}
}

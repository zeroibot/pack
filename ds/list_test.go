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
	// TODO: ToAny
	// TODO: IndexFunc, AllIndexFunc
	// TODO: RemoveFunc, RemoveAllFunc
	// TODO: Get, GetFuncOrDefault
	// TODO: Last, MustLast
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
	// TODO: Any, AnyIndexed
	// TODO: NotAny, NotAnyIndexed
	// TODO: All, AllIndexed
	// TODO: CountFunc
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
	// Filter
	numbers := List[int]{1, 2, 3, 4, 5, 6, 7}
	want := List[int]{2, 4, 6}
	actual := numbers.Filter(func(x int) bool { return x%2 == 0 })
	if slices.Equal(want, actual) == false {
		t.Errorf("List.Filter() = %v, want %v", actual, want)
	}
	want = List[int]{}
	actual = numbers.Filter(func(x int) bool { return x > 10 })
	if slices.Equal(want, actual) == false {
		t.Errorf("List.Filter() = %v, want %v", actual, want)
	}
	want = numbers
	actual = numbers.Filter(func(x int) bool { return x <= 10 })
	if slices.Equal(want, actual) == false {
		t.Errorf("List.Filter() = %v, want %v", actual, want)
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

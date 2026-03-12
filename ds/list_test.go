package ds

import (
	"fmt"
	"slices"
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
	// TODO: GetRandom, MustGetRandom
	// TODO: Shuffle
}

func TestCheckMethods(t *testing.T) {
	// TODO: Any, AnyIndexed
	// TODO: NotAny, NotAnyIndexed
	// TODO: All, AllIndexed
	// TODO: CountFunc
}

func TestFnMethods(t *testing.T) {
	// TODO: MapList
	// TODO: Filter, FilterIndexed
	// TODO: Reduce
	// TODO: Apply
}

func TestNumList(t *testing.T) {
	// TODO: NumList.ToList
	// TODO: Sum, Product
}

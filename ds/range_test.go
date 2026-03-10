package ds

import (
	"reflect"
	"testing"
)

func TestRange(t *testing.T) {
	r1 := NewRange(0, 5)
	r2 := NewInclusiveRange(1, 10)
	textCases := []Pair[string]{
		{r1.String(), "[0, 5)"},
		{r2.String(), "[1, 11)"},
	}
	for _, x := range textCases {
		text, want := x.Values()
		if text != want {
			t.Errorf("Range.String() = %q, want %q", text, want)
		}
	}
	a, b := r1.Limits()
	if a != 0 || b != 5 {
		t.Errorf("Range.Limits() = %d, %d, want = %d, %d", a, b, 0, 5)
	}
	a, b = r2.Limits()
	if a != 1 || b != 11 {
		t.Errorf("Range.Limits() = %d, %d, want = %d, %d", a, b, 11, 0)
	}
	size := r1.Len()
	if size != 5 {
		t.Errorf("Range.Len() = %d, want = %d", size, 5)
	}
	size = r2.Len()
	if size != 10 {
		t.Errorf("Range.Len() = %d, want = %d", size, 10)
	}
	r3 := r1.Copy()
	a, b = r3.Limits()
	if a != 0 || b != 5 {
		t.Errorf("Range.Copy.Limits() = %d, %d, want = %d, %d", a, b, 0, 5)
	}
	hasCases := []Tuple2[int, bool]{
		{3, true},
		{0, true},
		{-1, false},
		{4, true},
		{5, false},
		{99, false},
	}
	for _, x := range hasCases {
		item, want := x.Values()
		actual := r1.Has(item)
		if actual != want {
			t.Errorf("Range.Has(%d) = %v, want = %v", item, actual, want)
		}
	}
	hasCases = []Tuple2[int, bool]{
		{3, true},
		{5, true},
		{1, true},
		{0, false},
		{-1, false},
		{10, true},
		{11, false},
		{67, false},
	}
	for _, x := range hasCases {
		item, want := x.Values()
		actual := r2.Has(item)
		if actual != want {
			t.Errorf("Range.Has(%d) = %v, want = %v", item, actual, want)
		}
	}
	testCases := []Tuple3[int, int, string]{
		{r1.Sum(), 10, "Sum"},
		{r2.Sum(), 55, "Sum"},
		{r1.Product(), 0, "Product"},
		{r2.Product(), 3628800, "Product"},
	}
	for _, x := range testCases {
		actual, want, name := x.Values()
		if actual != want {
			t.Errorf("Range.%s() = %d, want = %d", name, actual, want)
		}
	}
	sliceCases := []Pair[[]int]{
		{r1.ToSlice(), []int{0, 1, 2, 3, 4}},
		{r2.ToSlice(), []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	}
	for _, x := range sliceCases {
		items, want := x.Values()
		if reflect.DeepEqual(items, want) == false {
			t.Errorf("Range.ToSlice() = %v, want = %v", items, want)
		}
	}
}

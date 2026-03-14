package ds

import (
	"slices"
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
		text, want := x.Unpack()
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
	for _, isReversed := range []bool{r1.IsReversed(), r2.IsReversed()} {
		if isReversed != false {
			t.Errorf("Range.IsReversed() = %t, want false", isReversed)
		}
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
		item, want := x.Unpack()
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
		item, want := x.Unpack()
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
		actual, want, name := x.Unpack()
		if actual != want {
			t.Errorf("Range.%s() = %d, want = %d", name, actual, want)
		}
	}
	sliceCases := []Pair[[]int]{
		{r1.ToSlice(), []int{0, 1, 2, 3, 4}},
		{r2.ToSlice(), []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	}
	for _, x := range sliceCases {
		items, want := x.Unpack()
		if slices.Equal(items, want) == false {
			t.Errorf("Range.ToSlice() = %v, want = %v", items, want)
		}
	}
}

func TestReverseRange(t *testing.T) {
	r1 := NewRange(5, 0)
	r2 := NewInclusiveRange(10, 1)
	textCases := []Pair[string]{
		{r1.String(), "[5, 0)"},
		{r2.String(), "[10, 0)"},
	}
	for _, x := range textCases {
		text, want := x.Unpack()
		if text != want {
			t.Errorf("Range.String() = %q, want %q", text, want)
		}
	}
	a, b := r1.Limits()
	if a != 5 || b != 0 {
		t.Errorf("Range.Limits() = %d, %d, want = %d, %d", a, b, 5, 0)
	}
	a, b = r2.Limits()
	if a != 10 || b != 0 {
		t.Errorf("Range.Limits() = %d, %d, want = %d, %d", a, b, 10, 0)
	}
	for _, isReversed := range []bool{r1.IsReversed(), r2.IsReversed()} {
		if isReversed != true {
			t.Errorf("Range.IsReversed() = %t, want true", isReversed)
		}
	}
	size := r1.Len()
	if size != 5 {
		t.Errorf("Range.Len() = %d, want = %d", size, 5)
	}
	size = r2.Len()
	if size != 10 {
		t.Errorf("Range.Len() = %d, want = %d", size, 10)
	}
	r3 := r2.Copy()
	a, b = r3.Limits()
	if a != 10 || b != 0 {
		t.Errorf("Range.Copy.Limits() = %d, %d, want = %d, %d", a, b, 10, 0)
	}
	hasCases := []Tuple2[int, bool]{
		{3, true},
		{0, false},
		{-1, false},
		{1, true},
		{5, true},
		{9, true},
		{10, true},
		{11, false},
		{99, false},
	}
	for _, x := range hasCases {
		item, want := x.Unpack()
		actual := r2.Has(item)
		if actual != want {
			t.Errorf("Range.Has(%d) = %v, want = %v", item, actual, want)
		}
	}
	testCases := []Tuple3[int, int, string]{
		{r1.Sum(), 15, "Sum"},
		{r2.Sum(), 55, "Sum"},
		{r1.Product(), 120, "Product"},
		{r2.Product(), 3628800, "Product"},
	}
	for _, x := range testCases {
		actual, want, name := x.Unpack()
		if actual != want {
			t.Errorf("Range.%s() = %d, want = %d", name, actual, want)
		}
	}
	sliceCases := []Pair[List[int]]{
		{r1.ToList(), List[int]{5, 4, 3, 2, 1}},
		{r2.ToList(), List[int]{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}},
	}
	for _, x := range sliceCases {
		items, want := x.Unpack()
		if slices.Equal(items, want) == false {
			t.Errorf("Range.ToList() = %v, want = %v", items, want)
		}
	}
}

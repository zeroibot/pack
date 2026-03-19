package ds

import (
	"testing"

	"github.com/roidaradal/tst"
)

func TestRange(t *testing.T) {
	r1 := NewRange(0, 5)
	r2 := NewInclusiveRange(1, 10)
	r3 := NewRange(5, 0)
	r4 := NewInclusiveRange(10, 1)
	// Range.String
	textCases := [][2]string{
		{r1.String(), "[0, 5)"},
		{r2.String(), "[1, 11)"},
		{r3.String(), "[5, 0)"},
		{r4.String(), "[10, 0)"},
	}
	tst.All(t, textCases, "Range.String", tst.AssertEqual)
	// Range.Limits
	limitCases := []tst.P1W2[Range[int], int, int]{
		{r1, 0, 5}, {r2, 1, 11},
		{r3, 5, 0}, {r4, 10, 0},
	}
	tst.AllP1W2(t, limitCases, "Range.Limits", Range[int].Limits, tst.AssertEqual[int], tst.AssertEqual[int])
	// IsReversed
	revCases := []tst.P1W1[Range[int], bool]{
		{r1, false}, {r2, false}, {r3, true}, {r4, true},
	}
	tst.AllP1W1(t, revCases, "Range.IsReversed", Range[int].IsReversed, tst.AssertEqual)
	// Len
	lenCases := []tst.P1W1[Range[int], int]{
		{r1, 5}, {r2, 10}, {r3, 5}, {r4, 10},
	}
	tst.AllP1W1(t, lenCases, "Range.Len", Range[int].Len, tst.AssertEqual)
	// Copy
	r5 := r1.Copy()
	a, b := r5.Limits()
	tst.AssertEqual2(t, "Range.Copy.Limits", a, 0, b, 5)
	r6 := r4.Copy()
	a, b = r6.Limits()
	tst.AssertEqual2(t, "Range.Copy.Limits", a, 10, b, 0)
	// Has
	hasCases := []tst.P2W1[Range[int], int, bool]{
		{r1, 3, true}, {r1, 0, true}, {r1, -1, false},
		{r1, 4, true}, {r1, 5, false}, {r1, 99, false},
		{r2, 3, true}, {r2, 5, true}, {r2, 1, true},
		{r2, 0, false}, {r2, -1, false}, {r2, 10, true},
		{r2, 11, false}, {r2, 67, false},
		{r4, 3, true}, {r4, 0, false}, {r4, -1, false},
		{r4, 1, true}, {r4, 5, true}, {r4, 9, true},
		{r4, 10, true}, {r4, 11, false}, {r4, 99, false},
	}
	tst.AllP2W1(t, hasCases, "Range.Has", Range[int].Has, tst.AssertEqual)
	// Sum, Product
	sumCases := []tst.P1W1[Range[int], int]{
		{r1, 10}, {r2, 55}, {r3, 15}, {r4, 55},
	}
	productCases := []tst.P1W1[Range[int], int]{
		{r1, 0}, {r2, 3628800}, {r3, 120}, {r4, 3628800},
	}
	tst.AllP1W1(t, sumCases, "Range.Sum", Range[int].Sum, tst.AssertEqual)
	tst.AllP1W1(t, productCases, "Range.Product", Range[int].Product, tst.AssertEqual)
	// ToSlice
	sliceCases := []tst.P1W1[Range[int], []int]{
		{r1, []int{0, 1, 2, 3, 4}},
		{r2, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	}
	tst.AllP1W1(t, sliceCases, "Range.ToSlice", Range[int].ToSlice, tst.AssertListEqual)
	// ToList
	listCases := []tst.P1W1[Range[int], List[int]]{
		{r3, List[int]{5, 4, 3, 2, 1}},
		{r4, List[int]{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}},
	}
	tst.AllP1W1(t, listCases, "Range.ToList", Range[int].ToList, tst.AssertListEqual)
}

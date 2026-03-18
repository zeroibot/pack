package lang

import (
	"slices"
	"testing"

	"github.com/roidaradal/tst"
)

func TestSortAscending(t *testing.T) {
	items1 := []int{5, 1, 2, 4, 3}
	want1 := []int{1, 2, 3, 4, 5}
	slices.SortFunc(items1, SortAscending[int]())
	tst.AssertListEqual(t, "SortAscending", items1, want1)

	items2 := []byte{'d', 'a', 'c', 'b', 'e'}
	want2 := []byte{'a', 'b', 'c', 'd', 'e'}
	slices.SortFunc(items2, SortAscending[byte]())
	tst.AssertListEqual(t, "SortAscending", items2, want2)
}

func TestSortDescending(t *testing.T) {
	items1 := []int{5, 1, 2, 4, 3}
	want1 := []int{5, 4, 3, 2, 1}
	slices.SortFunc(items1, SortDescending[int]())
	tst.AssertListEqual(t, "SortDescending", items1, want1)

	items2 := []byte{'d', 'a', 'c', 'b', 'e'}
	want2 := []byte{'e', 'd', 'c', 'b', 'a'}
	slices.SortFunc(items2, SortDescending[byte]())
	tst.AssertListEqual(t, "SortDescending", items2, want2)
}

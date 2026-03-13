package list

import "testing"

func TestFn(t *testing.T) {
	// TODO: CountFunc
	// TODO: Map, MapIndexed
	// TODO: MapIf, MapIndexedIf
	// TODO: MapList, MapLookup
	// TODO: Filter, FilterIndexed
	// TODO: Reduce
	// TODO: Apply
}

func TestSumProduct(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6}
	// Sum
	actual, want := Sum(items), 21
	if actual != want {
		t.Errorf("Sum() = %d, want %d", actual, want)
	}
	// Product
	actual, want = Product(items), 720
	if actual != want {
		t.Errorf("Product() = %d, want %d", actual, want)
	}
}

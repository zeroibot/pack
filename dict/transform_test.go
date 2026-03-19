package dict

import (
	"fmt"
	"testing"

	"github.com/roidaradal/tst"
)

func TestZip(t *testing.T) {
	// Zip
	keys := []string{"a", "b", "c"}
	values := []int{1, 2, 3}
	m := Zip(keys, values)
	for i, key := range keys {
		actual, ok := m[key]
		tst.AssertEqualAnd(t, "Zip", actual, values[i], ok, true)
		actual, ok = m[key+"x"]
		tst.AssertEqualAnd(t, "Zip", actual, 0, ok, false)
	}
	// Unzip
	keys2, values2 := Unzip(m)
	m2 := Zip(keys2, values2)
	tst.AssertMapEqual(t, "Unzip.Zip", m2, m)
	// Zip.Len
	values2 = []int{1, 2}
	m3 := Zip(keys, values2)
	want := map[string]int{"a": 1, "b": 2}
	tst.AssertMapEqual(t, "Zip", m3, want)
	tst.AssertEqual(t, "Zip.Len", Len(m3), 2)
}

func TestSwap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	want := []Entry[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	tst.AssertListEqual(t, "Swap.Entries", SortedEntries(Swap(m)), want)
}

func TestSwapList(t *testing.T) {
	m := map[string][]int{
		"a": {1, 3, 5},
		"b": {2, 4},
	}
	want := []Entry[int, string]{{1, "a"}, {2, "b"}, {3, "a"}, {4, "b"}, {5, "a"}}
	tst.AssertListEqual(t, "SwapList.Entries", SortedEntries(SwapList(m)), want)
}

func TestFromStruct(t *testing.T) {
	type config struct {
		A, B, C int
	}
	cfg := &config{1, 2, 3}
	want := map[string]int{"A": 1, "B": 2, "C": 3}
	actual, err := FromStruct[int](cfg)
	tst.AssertMapEqualError(t, "FromStruct", actual, want, err, false)

	// Test nil input
	want = map[string]int{}
	actual, err = FromStruct[int, config](nil)
	tst.AssertMapEqualError(t, "FromStruct", actual, want, err, false)

	// Test unmarshal error
	actual2, err := FromStruct[string, config](cfg)
	var want2 map[string]string = nil
	tst.AssertMapEqualError(t, "FromStruct", actual2, want2, err, true)

	// Test marshal error
	type config2 struct {
		Item any
	}
	cfg2 := &config2{Item: make(chan int)}
	actual3, err := FromStruct[any, config2](cfg2)
	var want3 map[string]any = nil
	tst.AssertMapEqualError(t, "FromStruct", actual3, want3, err, true)
}

func TestToStruct(t *testing.T) {
	type config struct {
		A, B, C int
	}
	obj := Object{"A": 1, "B": 2, "C": 3}
	want := &config{1, 2, 3}
	actual, err := ToStruct[config](obj)
	tst.AssertDeepEqualError(t, "ToStruct", actual, want, err, false)

	// Test nil input
	want = &config{0, 0, 0}
	actual, err = ToStruct[config](nil)
	tst.AssertDeepEqualError(t, "ToStruct", actual, want, err, false)

	// Test unmarshal error
	type config2 struct {
		A, B string
	}
	actual2, err := ToStruct[config2](obj)
	var want2 *config2 = nil
	tst.AssertDeepEqualError(t, "ToStruct", actual2, want2, err, true)

	// Test marshal error
	obj2 := Object{"A": make(chan int), "B": 5}
	actual3, err := ToStruct[config2](obj2)
	tst.AssertDeepEqualError(t, "ToStruct", actual3, want2, err, true)
}

func TestToObject(t *testing.T) {
	type config struct {
		A, B, C int
	}
	cfg := &config{1, 2, 3}
	want := Object{"A": 1, "B": 2, "C": 3}
	actual, err := ToObject(cfg)
	tst.AssertTrue(t, "ToObject.Error", err == nil)
	compareObjects(t, actual, want)
}

func TestPruned(t *testing.T) {
	type config struct {
		A, B, C int
	}
	cfg := &config{1, 2, 3}
	want := Object{"B": 2, "C": 3}
	actual, err := Pruned(cfg, "B", "C")
	tst.AssertTrue(t, "Pruned.Error", err == nil)
	compareObjects(t, actual, want)

	type config2 struct {
		Item any
	}
	cfg2 := &config2{Item: make(chan int)}
	actual2, err := Pruned(cfg2, "Item")
	var want2 Object = nil
	tst.AssertMapEqualError(t, "Pruned", actual2, want2, err, true)
}

func compareObjects(t *testing.T, actual, want Object) {
	// Note: cannot use maps.Equal because map value type is <any>
	// The <any> type is not comparable, so even if the map values are the same, the comparison fails
	actualKeys, wantKeys := SortedKeys(actual), SortedKeys(want)
	tst.AssertListEqual(t, "ToObject.Keys", actualKeys, wantKeys)
	for _, key := range wantKeys {
		wantValue := fmt.Sprintf("%v", want[key])
		actualValue := fmt.Sprintf("%v", actual[key])
		tst.AssertEqual(t, fmt.Sprintf("ToObject[%q]", key), actualValue, wantValue)
	}
}

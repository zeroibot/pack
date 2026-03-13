package list

import "testing"

func TestCompareAllAny(t *testing.T) {
	var ints0 []int
	ints1 := []int{1, 1, 1, 1}
	ints2 := []int{1, 2, 3, 4}
	ints3 := []int{2, 2, 1, 2}
	var bools0 []bool
	bools1 := []bool{true, true, true}
	bools2 := []bool{false, false}
	bools3 := []bool{false, true, false}
	// AllEqual
	actual := AllEqual(ints0, 5)
	if actual != false {
		t.Errorf("AllEqual() = %t, want %t", actual, false)
	}
	actual = AllEqual(ints1, 1)
	if actual != true {
		t.Errorf("AllEqual() = %t, want %t", actual, true)
	}
	actual = AllEqual(ints2, 3)
	if actual != false {
		t.Errorf("AllEqual() = %t, want %t", actual, false)
	}
	actual = AllEqual(ints3, 2)
	if actual != false {
		t.Errorf("AllEqual() = %t, want %t", actual, false)
	}
	// AllTrue, All False
	actual = AllTrue(bools0)
	if actual != false {
		t.Errorf("AllTrue() = %t, want %t", actual, false)
	}
	actual = AllFalse(bools0)
	if actual != false {
		t.Errorf("AllFalse() = %t, want %t", actual, false)
	}
	actual = AllTrue(bools1)
	if actual != true {
		t.Errorf("AllTrue() = %t, want %t", actual, true)
	}
	actual = AllTrue(bools3)
	if actual != false {
		t.Errorf("AllTrue() = %t, want %t", actual, false)
	}
	actual = AllFalse(bools2)
	if actual != true {
		t.Errorf("AllFalse() = %t, want %t", actual, true)
	}
	actual = AllFalse(bools3)
	if actual != false {
		t.Errorf("AllFalse() = %t, want %t", actual, false)
	}
	// Has, HasNo
	actual, actual2 := Has(ints0, 1), HasNo(ints0, 1)
	if actual != false {
		t.Errorf("Has() = %t, want %t", actual, false)
	}
	if actual2 != true {
		t.Errorf("HasNo() = %t, want %t", actual2, true)
	}
	actual, actual2 = Has(ints1, 1), HasNo(ints1, 1)
	if actual != true {
		t.Errorf("Has() = %t, want %t", actual, true)
	}
	if actual2 != false {
		t.Errorf("HasNo() = %t, want %t", actual2, false)
	}
	actual, actual2 = Has(ints3, 1), HasNo(ints3, 1)
	if actual != true {
		t.Errorf("Has() = %t, want %t", actual, true)
	}
	if actual2 != false {
		t.Errorf("HasNo() = %t, want %t", actual2, false)
	}
	actual, actual2 = Has(ints2, 5), HasNo(ints2, 5)
	if actual != false {
		t.Errorf("Has() = %t, want %t", actual, false)
	}
	if actual2 != true {
		t.Errorf("HasNo() = %t, want %t", actual2, true)
	}
	actual, actual2 = AnyTrue(bools1), AnyFalse(bools1)
	if actual != true {
		t.Errorf("AnyTrue() = %t, want %t", actual, true)
	}
	if actual2 != false {
		t.Errorf("AnyFalse() = %t, want %t", actual2, false)
	}
	actual, actual2 = AnyTrue(bools2), AnyFalse(bools2)
	if actual != false {
		t.Errorf("AnyTrue() = %t, want %t", actual, false)
	}
	if actual2 != true {
		t.Errorf("AnyFalse() = %t, want %t", actual2, true)
	}
	actual, actual2 = AnyTrue(bools3), AnyFalse(bools3)
	if actual != true {
		t.Errorf("AnyTrue() = %t, want %t", actual, true)
	}
	if actual2 != true {
		t.Errorf("AnyFalse() = %t, want %t", actual2, true)
	}
}

func TestListCompareQueries(t *testing.T) {
	// TODO: IndexLookup
	// TODO: IndexOf, AllIndexOf
	// TODO: Remove, RemoveAll
	// TODO: GetOrDefault
}

func TestListCompareMethods(t *testing.T) {
	// TODO: Tally, TallyFunc
	// TODO: CountUnique, CountUniqueFunc
	// TODO: AllSame, AllSameFunc
	// TODO: AllUnique, AllUniqueFunc
	// TODO: Deduplicate, DeduplicateFunc
	// TODO: Count
	// TODO: GroupByFunc
}

package dyn

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/roidaradal/tst"
)

func TestPointers(t *testing.T) {
	// AddressOf
	x, y, z := 5, 5, 6
	xp, yp, zp := new(x), new(y), new(z)
	testCases := []tst.P1W1[any, string]{
		{xp, fmt.Sprintf("%p", xp)},
		{yp, fmt.Sprintf("%p", yp)},
		{5, "0x0"},
		{nil, "0x0"},
	}
	tst.AllP1W1(t, testCases, "AddressOf", AddressOf, tst.AssertEqual)

	// Deref
	testCases2 := []tst.P1W2[any, any, bool]{
		{x, x, false},
		{xp, x, true},
		{zp, z, true},
	}
	tst.AllP1W2(t, testCases2, "Deref", Deref, tst.AssertEqualAny, tst.AssertEqual)

	// MustDeref
	testCases3 := []tst.P1W1[any, any]{
		{yp, y},
		{zp, z},
	}
	tst.AllP1W1(t, testCases3, "MustDeref", MustDeref, tst.AssertEqualAny)

	defer tst.AssertPanic(t, "MustDeref")
	MustDeref(z) // should panic
}

func TestIsNil(t *testing.T) {
	type testCase = tst.P1W1[any, bool]
	// IsZero
	var zeroInt int
	var zeroFloat float64
	var zeroBool bool
	var zeroString string
	var zeroList []int
	var zeroMap map[int]string
	myInt, myFloat := 1, 2.5
	myString, myList, myMap := "hello", []int{1, 2, 3}, map[int]string{1: "A", 2: "B"}
	testCases := []testCase{
		{zeroInt, true}, {zeroFloat, true}, {zeroBool, true},
		{zeroString, true}, {zeroList, true}, {zeroMap, true},
		{myInt, false}, {myFloat, false}, {true, false},
		{myString, false}, {myList, false}, {myMap, false},
	}
	tst.AllP1W1(t, testCases, "IsZero", IsZero, tst.AssertEqual)

	// IsNil, NotNil
	var ptrInt *int
	var zeroCh chan int
	myCh := make(chan int)
	var zeroFn func()
	myFn := func() {}
	var zeroErr error
	myErr := fmt.Errorf("my error")
	testCases2 := []testCase{
		{myInt, false}, {myFloat, false}, {true, false},
		{myString, false}, {myList, false}, {myMap, false},
		{new(myInt), false}, {myCh, false}, {myFn, false}, {myErr, false},
		{nil, true}, {zeroList, true}, {zeroMap, true},
		{ptrInt, true}, {zeroCh, true}, {zeroFn, true}, {zeroErr, true},
	}
	tst.AllP1W1(t, testCases2, "IsNil", IsNil, tst.AssertEqual)
	testCases2 = tst.FlipP1W1(testCases2)
	tst.AllP1W1(t, testCases2, "NotNil", NotNil, tst.AssertEqual)
}

func TestIsEqual(t *testing.T) {
	type testCase = tst.P2W1[any, any, bool]
	int1, int2, int3 := 5, 5, 6
	list1 := []int{1, 2, 3}
	list2 := []int{1, 2, 3}
	list3 := []int{3, 2, 1}
	list4 := []string{"a", "b", "c"}
	map1 := map[int]string{1: "a", 2: "b", 3: "c"}
	map2 := map[int]string{2: "b", 3: "c", 1: "a"}
	map3 := map[int]string{3: "a", 2: "x", 1: "z"}
	map4 := map[string]int{"a": 1, "b": 2}
	testCases := []testCase{
		{int1, int2, true}, {int1, int1, true},
		{int2, int3, false}, {int3, list1, false},
		{new(int1), new(int2), true}, {new(int2), new(int3), false},
		{list1, list2, true}, {list1, list1, true},
		{list2, list3, false}, {list3, list4, false},
		{&list1, list2, true}, {&list2, &list3, false},
		{map1, map2, true}, {map2, map2, true},
		{map2, map3, false}, {map3, map4, false},
		{map1, &map2, true}, {&map2, &map3, false},
	}
	tst.AllP2W1(t, testCases, "IsEqual", IsEqual, tst.AssertEqual)
	testCases = tst.FlipP2W1(testCases)
	tst.AllP2W1(t, testCases, "NotEqual", NotEqual, tst.AssertEqual)
}

func TestDerefValue(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	// DerefValue
	p := person{"John", 20}
	wantValue := reflect.ValueOf(p)
	actualValue, ok := DerefValue(&p)
	if !ok || reflect.DeepEqual(actualValue.Interface(), wantValue.Interface()) == false {
		t.Errorf("DerefValue() = %v, %t, want %v, true", actualValue, ok, wantValue)
	}
	actualValue, ok = DerefValue(p)
	if ok || actualValue.IsValid() == true {
		t.Errorf("DerefValue() = %v, %t, want <invalid reflect.Value>, false", actualValue, ok)
	}

	// MustDerefValue
	actualValue = MustDerefValue(&p)
	if reflect.DeepEqual(actualValue.Interface(), wantValue.Interface()) == false {
		t.Errorf("MustDerefValue() = %v, want %v", actualValue, wantValue)
	}

	defer tst.AssertPanic(t, "MustDerefValue")
	MustDerefValue(p) // should panic
}

func TestRefValue(t *testing.T) {
	type person struct {
		Name     string
		password string
	}
	p := person{"John", "123"}
	structValue := MustDerefValue(&p)
	nameField := structValue.FieldByName("Name")
	pwdField := structValue.FieldByName("password")
	testCases := []tst.P1W2[reflect.Value, any, bool]{
		{nameField, &p.Name, true},
		{pwdField, nil, false},
		{structValue.FieldByName("unknown"), nil, false},
	}
	tst.AllP1W2(t, testCases, "RefValue", RefValue, tst.AssertEqualAny, tst.AssertEqual[bool])

	actual := MustRefValue(nameField)
	tst.AssertEqualAny(t, "MustRefValue", actual, &p.Name)

	defer tst.AssertPanic(t, "MustRefValue")
	MustRefValue(pwdField) // should panic
}

func TestAnyValue(t *testing.T) {
	type person struct {
		Name     string
		password string
	}
	p := person{"John", "123"}
	structValue := MustDerefValue(&p)
	nameField := structValue.FieldByName("Name")
	pwdField := structValue.FieldByName("password")
	testCases := []tst.P1W2[reflect.Value, any, bool]{
		{nameField, p.Name, true},
		{pwdField, nil, false},
		{structValue.FieldByName("unknown"), nil, false},
	}
	tst.AllP1W2(t, testCases, "AnyValue", AnyValue, tst.AssertEqualAny, tst.AssertEqual[bool])

	actual := MustAnyValue(nameField)
	tst.AssertEqualAny(t, "MustAnyValue", actual, p.Name)

	defer tst.AssertPanic(t, "MustAnyValue")
	MustAnyValue(pwdField) // should panic
}

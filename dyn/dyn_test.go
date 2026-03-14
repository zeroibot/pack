package dyn

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPointers(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustDeref() did not panic")
		}
	}()
	// AddressOf
	x, y, z := 5, 5, 6
	xp, yp, zp := new(x), new(y), new(z)
	testCases := [][2]string{
		{fmt.Sprintf("%p", xp), AddressOf(xp)},
		{fmt.Sprintf("%p", yp), AddressOf(yp)},
		{"0x0", AddressOf(5)},
		{"0x0", AddressOf(nil)},
	}
	for _, x := range testCases {
		want, actual := x[0], x[1]
		if want != actual {
			t.Errorf("AddressOf = %s, want %s", actual, want)
		}
	}
	// Deref
	actual, ok := Deref(x)
	if ok || actual != x {
		t.Errorf("Deref() = %v, %t, want %v, false", actual, ok, x)
	}
	actual, ok = Deref(xp)
	if !ok || actual != x {
		t.Errorf("Deref() = %v, %t, want %v, false", actual, ok, x)
	}
	actual, ok = Deref(zp)
	if !ok || actual != z {
		t.Errorf("Deref() = %v, %t, want %v, false", actual, ok, z)
	}
	value := MustDeref(yp)
	if value != y {
		t.Errorf("MustDeref() = %v, want %v", value, y)
	}
	value = MustDeref(zp)
	if value != z {
		t.Errorf("MustDeref() = %v, want %v", value, z)
	}
	MustDeref(z) // should panic
}

func TestIsNil(t *testing.T) {
	// IsZero
	var zeroInt int
	var zeroFloat float64
	var zeroBool bool
	var zeroString string
	var zeroList []int
	var zeroMap map[int]string
	myInt, myFloat := 1, 2.5
	myString, myList, myMap := "hello", []int{1, 2, 3}, map[int]string{1: "A", 2: "B"}
	testCases := [][2]bool{
		{true, IsZero(zeroInt)}, {true, IsZero(zeroFloat)}, {true, IsZero(zeroBool)},
		{true, IsZero(zeroString)}, {true, IsZero(zeroList)}, {true, IsZero(zeroMap)},
		{false, IsZero(myInt)}, {false, IsZero(myFloat)}, {false, IsZero(true)},
		{false, IsZero(myString)}, {false, IsZero(myList)}, {false, IsZero(myMap)},
	}
	for _, x := range testCases {
		want, actual := x[0], x[1]
		if want != actual {
			t.Errorf("IsZero = %t, want %t", actual, want)
		}
	}
	// IsNil, NotNil
	type testCase struct {
		want  bool
		input any
	}
	var ptrInt *int
	var zeroCh chan int
	myCh := make(chan int)
	var zeroFn func()
	myFn := func() {}
	var zeroErr error
	myErr := fmt.Errorf("my error")
	testCases2 := []testCase{
		{false, myInt}, {false, myFloat}, {false, true},
		{false, myString}, {false, myList}, {false, myMap},
		{false, new(myInt)}, {false, myCh}, {false, myFn}, {false, myErr},
		{true, nil}, {true, zeroList}, {true, zeroMap},
		{true, ptrInt}, {true, zeroCh}, {true, zeroFn}, {true, zeroErr},
	}
	for _, x := range testCases2 {
		actual := IsNil(x.input)
		if actual != x.want {
			t.Errorf("IsNil = %t, want %t", actual, x.want)
		}
		actual, want := NotNil(x.input), !x.want
		if actual != want {
			t.Errorf("NotNil = %t, want %t", actual, want)
		}
	}
}

func TestIsEqual(t *testing.T) {
	type testCase struct {
		want           bool
		input1, input2 any
	}
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
		{true, int1, int2}, {true, int1, int1},
		{false, int2, int3}, {false, int3, list1},
		{true, new(int1), new(int2)}, {false, new(int2), new(int3)},
		{true, list1, list2}, {true, list1, list1},
		{false, list2, list3}, {false, list3, list4},
		{true, &list1, list2}, {false, &list2, &list3},
		{true, map1, map2}, {true, map2, map2},
		{false, map2, map3}, {false, map3, map4},
		{true, map1, &map2}, {false, &map2, &map3},
	}
	for _, x := range testCases {
		actual := IsEqual(x.input1, x.input2)
		if actual != x.want {
			t.Errorf("IsEqual(%v, %v) = %t, want %t", x.input1, x.input2, actual, x.want)
		}
		actual, want := NotEqual(x.input1, x.input2), !x.want
		if actual != want {
			t.Errorf("NotEqual(%v, %v) = %t, want %t", x.input1, x.input2, actual, want)
		}
	}
}

func TestDerefValue(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustDerefValue() did not panic")
		}
	}()
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
	MustDerefValue(p) // should panic
}

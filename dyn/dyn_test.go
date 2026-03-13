package dyn

import (
	"fmt"
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
	// TODO: IsZero
	// TODO: IsNil
	// TODO: NotNil
}

func TestIsEqual(t *testing.T) {
	// TODO: IsEqual
	// TODO: NotEqual
}

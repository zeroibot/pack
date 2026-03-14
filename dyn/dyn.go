// Package dyn (dynamic) is an extension to the `reflect` package
package dyn

import (
	"fmt"
	"reflect"
)

// AddressOf returns the memory address of given item as string
func AddressOf(x any) string {
	if IsNil(x) || !IsPointer(x) {
		return "0x0" // nil address
	}
	return fmt.Sprintf("%p", x)
}

// Deref dereferences the given pointer, and returns flag if it is valid
func Deref(x any) (any, bool) {
	if !IsPointer(x) {
		return x, false // return item as is
	}
	return MustDeref(x), true
}

// MustDeref dereferences the given pointer, panics if item is not a pointer
func MustDeref(x any) any {
	return reflect.ValueOf(x).Elem().Interface()
}

// DerefValue dereferences the given pointer, and returns reflect.Value and flag if it is valid
func DerefValue(x any) (reflect.Value, bool) {
	if !IsPointer(x) {
		var zero reflect.Value
		return zero, false
	}
	return MustDerefValue(x), true
}

// MustDerefValue dereferences the given pointer and returns reflect.Value
// This panics if item is not a pointer
func MustDerefValue(x any) reflect.Value {
	return reflect.ValueOf(x).Elem()
}

// RefValue returns an `any` object that holds a reference to given reflect.Value, and flag if it is valid
func RefValue(value reflect.Value) (any, bool) {
	if !value.CanAddr() {
		return nil, false
	}
	ref := value.Addr()
	return AnyValue(ref)
}

// MustRefValue returns an `any` object that holds a reference to given reflect.Value
// This panics if value is not addressable or if it is private
func MustRefValue(value reflect.Value) any {
	return value.Addr().Interface()
}

// AnyValue returns an `any` object that holds the value of given reflect.Value, and flag if it is valid
func AnyValue(value reflect.Value) (any, bool) {
	if !value.CanInterface() {
		return nil, false
	}
	return MustAnyValue(value), true
}

// MustAnyValue returns an `any` object that holds the value of given reflect.Value
// This panics if value is private.
func MustAnyValue(value reflect.Value) any {
	return value.Interface()
}

// IsZero checks if given item has zero value
func IsZero(x any) bool {
	return reflect.ValueOf(x).IsZero()
}

// IsNil checks if given item is nil
func IsNil(x any) bool {
	if x == nil {
		return true
	}
	switch reflect.TypeOf(x).Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		return reflect.ValueOf(x).IsNil()
	default:
		return false
	}
}

// NotNil checks if given item is not nil
func NotNil(x any) bool {
	return !IsNil(x)
}

// IsEqual checks if the two `any` values are equal
func IsEqual(x, y any) bool {
	// Dereference item1 if pointer and not null
	if IsPointer(x) && NotNil(x) {
		return IsEqual(MustDeref(x), y)
	}
	// Dereference item2 if pointer and not null
	if IsPointer(y) && NotNil(y) {
		return IsEqual(x, MustDeref(y))
	}
	return reflect.DeepEqual(x, y)
}

// NotEqual checks if the two `any` values are not equal
func NotEqual(x, y any) bool {
	return !IsEqual(x, y)
}

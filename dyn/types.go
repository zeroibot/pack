package dyn

import "reflect"

// IsPointer checks if given item is a pointer
func IsPointer(x any) bool {
	return reflect.TypeOf(x).Kind() == reflect.Pointer
}

// IsStruct checks if given item is a struct
func IsStruct(x any) bool {
	return reflect.TypeOf(x).Kind() == reflect.Struct
}

// IsStructPointer checks if given item is a pointer to a struct
func IsStructPointer(x any) bool {
	if !IsPointer(x) {
		return false
	}
	return IsStruct(MustDeref(x))
}

// TypeName returns the base type name of given item (dereferences pointers)
func TypeName(x any) string {
	if IsPointer(x) {
		return TypeName(MustDeref(x))
	}
	return reflect.TypeOf(x).Name()
}

// FullTypeName returns the full type name of given item (*Type for pointers)
func FullTypeName(x any) string {
	if IsPointer(x) {
		return "*" + FullTypeName(MustDeref(x))
	}
	return reflect.TypeOf(x).Name()
}

package dyn

import (
	"fmt"
	"reflect"
)

// GetStructFieldTag extracts the tag value from the struct field
func GetStructFieldTag(structField reflect.StructField, tagKey string) (string, bool) {
	return structField.Tag.Lookup(tagKey)
}

// SetStructField sets item.field = value for the given struct pointer, returns flag for successful operation
func SetStructField(structRef any, field string, value any) bool {
	// Get struct field
	structField, ok := getStructField(structRef, field)
	if !ok {
		return false
	}
	// Check if struct field can be updated
	if !structField.CanSet() {
		return false
	}
	fieldValue := reflect.ValueOf(value)
	// Check if correct type
	if structField.Type() != fieldValue.Type() {
		return false
	}
	structField.Set(fieldValue)
	return true
}

// MustSetStructField sets item.field = value for the given struct pointer.
// If the given item is not a struct pointer, this is a no-op.
// This panics if the given value cannot be assigned to the struct field (wrong value type or field does not exist)
func MustSetStructField(structRef any, field string, value any) {
	// Check if struct pointer
	if !IsStructPointer(structRef) {
		return
	}
	// Set struct field value
	item := reflect.ValueOf(structRef).Elem()
	item.FieldByName(field).Set(reflect.ValueOf(value))
}

// GetStructField gets item.field from given struct pointer, and returns flag for successful get
func GetStructField(structRef any, field string) (any, bool) {
	// Get struct field
	structField, ok := getStructField(structRef, field)
	if !ok {
		return nil, false
	}
	// Check if safe to get value
	if !structField.CanInterface() {
		return nil, false
	}
	return structField.Interface(), true
}

// MustGetStructField gets item.field from given struct pointer.
// It returns nil if the given item is not a struct pointer.
// This panics if the given field is not accessible (not found or private)
func MustGetStructField(structRef any, field string) any {
	// Check if struct pointer
	if !IsStructPointer(structRef) {
		return nil
	}
	// Get struct field
	item := reflect.ValueOf(structRef).Elem()
	return item.FieldByName(field).Interface()
}

// GetStructFieldAs gets item.field from given struct pointer, and type coerces it into T, and returns flag for successful get
func GetStructFieldAs[T any](structRef any, field string) (T, bool) {
	structField, ok := GetStructField(structRef, field)
	if !ok {
		var zero T
		return zero, false
	}
	value, ok := structField.(T)
	return value, ok
}

// GetStructFieldAsString gets item.field from given struct pointer, and return field value as string
func GetStructFieldAsString(structRef any, field string) (string, bool) {
	fieldValue, ok := GetStructField(structRef, field)
	return fmt.Sprintf("%v", fieldValue), ok
}

// MustGetStructFieldAsString gets item.field from given struct pointer, and returns field value as string
// It returns "<nil>" if the given item is not a struct pointer.
// This panics if the given field is not accessible (not found or private).
func MustGetStructFieldAsString(structRef any, field string) string {
	fieldValue := MustGetStructField(structRef, field)
	return fmt.Sprintf("%v", fieldValue)
}

// getStructField gets the struct field as reflect.Value from given struct pointer
func getStructField(structRef any, field string) (reflect.Value, bool) {
	// Check if struct pointer
	if !IsStructPointer(structRef) {
		var zero reflect.Value
		return zero, false
	}
	// Dereference struct pointer = struct
	item := reflect.ValueOf(structRef).Elem()
	// Check if struct field exists
	structField := item.FieldByName(field)
	return structField, structField.IsValid()
}

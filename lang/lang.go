// Package lang contains quality-of-life language functions
package lang

// Ternary mimics the ternary operator : condition ? valueTrue : valueFalse
func Ternary[T any](condition bool, valueTrue, valueFalse T) T {
	if condition {
		return valueTrue
	}
	return valueFalse
}

// Ref returns a reference to given item
func Ref[T any](value T) *T {
	return new(value)
}

// Deref de-references the given item pointer.
// If null pointer, returns the zero value of the item
func Deref[T any](ref *T) T {
	var item T
	if ref != nil {
		item = *ref
	}
	return item
}

// Identity returns the passed argument, used as function argument
func Identity[T any](item T) T {
	return item
}

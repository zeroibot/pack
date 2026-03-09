// Package number contains number types and functions
package number

// Type interface unifies the number types
type Type interface {
	Integer | Float
}

// Integer interface unifies the integer types
type Integer interface {
	~uint | ~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Float interface unifies the float types
type Float interface {
	~float32 | ~float64
}

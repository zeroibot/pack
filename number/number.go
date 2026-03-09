// Package number contains number types and functions
package number

// Type interface unifies the number types
type Type interface {
	Int | Uint | Float
}

// Integer interface unifies the integer types
type Integer interface {
	Int | Uint
}

// Int interface unifies the Int types
type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Uint interface unifies the Uint types
type Uint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Float interface unifies the float types
type Float interface {
	~float32 | ~float64
}

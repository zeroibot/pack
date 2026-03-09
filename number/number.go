// Package number contains number types and functions
package number

import "math"

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

// Abs gets the absolute value of the Number
func Abs[T Type](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Ceil gets the ceiling value, as float64
func Ceil[F Float](f F) float64 {
	return math.Ceil(float64(f))
}

// Floor gets the floor value, as float64
func Floor[F Float](f F) float64 {
	return math.Floor(float64(f))
}

// Round rounds the value to the nearest whole float64
func Round[F Float](f F) float64 {
	return math.Round(float64(f))
}

// RoundToEven rounds the value to the nearest whole float64, with ties (0.5) rounding to nearest even number
func RoundToEven[F Float](f F) float64 {
	return math.RoundToEven(float64(f))
}

// CeilInt gets the ceiling value, as int
func CeilInt[F Float](f F) int {
	return int(Ceil(f))
}

// FloorInt gets the floor value, as int
func FloorInt[F Float](f F) int {
	return int(Floor(f))
}

// RoundInt rounds the value to the nearest int
func RoundInt[F Float](f F) int {
	return int(Round(f))
}

// RoundToEvenInt rounds the value to the nearest int, with ties (0.5) rounding to nearest even number
func RoundToEvenInt[F Float](f F) int {
	return int(RoundToEven(f))
}

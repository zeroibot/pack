// Package conv contains basic data type conversion functions
package conv

import (
	"fmt"

	"github.com/roidaradal/pack/number"
)

// BoolToInt converts bool to int (true = 1, false = 0)
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// BoolToUint converts bool to uint (true = 1, false = 0)
func BoolToUint(b bool) uint {
	if b {
		return 1
	}
	return 0
}

// BoolToString converts bool to string
func BoolToString(b bool) string {
	return fmt.Sprintf("%t", b)
}

// FloatToInt converts a Float to int
func FloatToInt[F number.Float](f F) int {
	return int(f)
}

// FloatToUint converts a Float to uint, clips to 0 if negative float
func FloatToUint[F number.Float](f F) uint {
	if f <= 0 {
		return 0
	}
	return uint(f)
}

// FloatToString converts a Float to string
func FloatToString[F number.Float](f F) string {
	return fmt.Sprintf("%f", f)
}

// IntToBool converts an Integer to bool (0 = false, else = true)
func IntToBool[I number.Integer](i I) bool {
	return i != 0
}

// IntToFloat converts an Integer to float64
func IntToFloat[I number.Integer](i I) float64 {
	return float64(i)
}

// IntToString converts an Integer to string
func IntToString[I number.Integer](i I) string {
	return fmt.Sprintf("%d", i)
}

// IntToUint converts an Int to uint, clips to 0 if negative int
func IntToUint[I number.Int](i I) uint {
	if i <= 0 {
		return 0
	}
	return uint(i)
}

// UintToInt converts a Uint to int
func UintToInt[I number.Uint](i I) int {
	return int(i)
}

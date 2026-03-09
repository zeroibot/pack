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

// FloatToInt converts a float to int
func FloatToInt[F number.Float](f F) int {
	return int(f)
}

// FloatToUint converts a float to uint, clips to 0 if negative float
func FloatToUint[F number.Float](f F) uint {
	if f <= 0 {
		return 0
	}
	return uint(f)
}

// FloatToString converts a float to string
func FloatToString[F number.Float](f F) string {
	return fmt.Sprintf("%f", f)
}

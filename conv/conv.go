// Package conv contains conversion functions
package conv

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/roidaradal/pack/ds"
)

// Abs gets the absolute value of an Integer
func Abs[I ds.Integer](i I) I {
	if i < 0 {
		return -i
	}
	return i
}

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

// IntToString converts an Integer to string
func IntToString[I ds.Integer](i I) string {
	return fmt.Sprintf("%d", i)
}

// FloatToString converts a Float to string
func FloatToString[F ds.Float](f F) string {
	return fmt.Sprintf("%f", f)
}

// FloatToStringDecimal converts a Float to string, with N decimal places
func FloatToStringDecimal[F ds.Float](f F, decimalPlaces int) string {
	format := fmt.Sprintf("%%.%df", decimalPlaces)
	return fmt.Sprintf(format, f)
}

// IntToUint converts Integer to uint, clips to 0 if negative int
func IntToUint[I ds.Integer](i I) uint {
	if i < 0 {
		return 0
	}
	return uint(i)
}

// IntToFloat converts Integer to float64
func IntToFloat[I ds.Integer](i I) float64 {
	return float64(i)
}

// IntToBool converts Integer to bool (0 is false, else true)
func IntToBool[I ds.Integer](i I) bool {
	return i != 0
}

// UintToInt converts uint to int
func UintToInt(i uint) int {
	return int(i)
}

// FloatToInt converts Float to int
func FloatToInt[F ds.Float](f F) int {
	return int(f)
}

// FloatToUint converts Float to uint, clips to 0 if negative float
func FloatToUint[F ds.Float](f F) uint {
	if f < 0 {
		return 0
	}
	return uint(f)
}

// StringToInt parses the string as int, defaults to 0 if invalid int
func StringToInt(s string) int {
	text := strings.TrimSpace(s)
	number, err := strconv.Atoi(text)
	if err != nil {
		return 0
	}
	return number
}

// StringToUint parses the string as uint, defaults to 0 if invalid uint
func StringToUint(s string) uint {
	text := strings.TrimSpace(s)
	number, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return 0
	}
	return uint(number)
}

// StringToFloat parses the string as float64, defaults to 0 if invalid float
func StringToFloat(s string) float64 {
	text := strings.TrimSpace(s)
	number, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0
	}
	return number
}

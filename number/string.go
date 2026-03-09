package number

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// CommaString returns the Integer string, grouped by commas
func CommaString[I Integer](i I) string {
	isNegative := i < 0
	if isNegative {
		i = -i
	}
	digits := fmt.Sprintf("%d", i)
	return groupDigits(digits, isNegative)
}

// CommaDecimalString returns the Float string, grouped by commas, using N decimal places
func CommaDecimalString[F Float](f F, decimalPlaces uint) string {
	isNegative := f < 0
	if isNegative {
		f = -f
	}
	text := DecimalString(f, decimalPlaces)
	parts := strings.Split(text, ".")
	whole := parts[0]
	fraction := ""
	if len(parts) > 1 {
		fraction = parts[1]
	}
	out := groupDigits(whole, isNegative)
	if fraction != "" {
		out = out + "." + fraction
	}
	return out
}

// groupDigits groups the digits into 3 and joins by comma
func groupDigits(text string, isNegative bool) string {
	digits := []byte(text)
	numDigits := len(digits)
	numChunks := int(math.Ceil(float64(numDigits) / 3))
	chunks := make([]string, numChunks)
	chunkIdx := numChunks - 1
	for limit := numDigits; limit > 0; limit -= 3 {
		start := max(0, limit-3)
		chunks[chunkIdx] = string(digits[start:limit])
		chunkIdx--
	}
	out := strings.Join(chunks, ",")
	if isNegative {
		out = "-" + out
	}
	return out
}

// DecimalString returns the Float string, using N decimal places
func DecimalString[F Float](f F, decimalPlaces uint) string {
	format := fmt.Sprintf("%%.%df", decimalPlaces)
	return fmt.Sprintf(format, f)
}

// ParseInt parses the string as int, default to 0 if invalid int
func ParseInt(text string) int {
	value, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		return 0
	}
	return value
}

// ParseUint parses the string as uint, default to 0 if invalid uint
func ParseUint(text string) uint {
	value, err := strconv.ParseUint(strings.TrimSpace(text), 10, 64)
	if err != nil {
		return 0
	}
	return uint(value)
}

// ParseFloat parses the string as float64, default to 0 if invalid float
func ParseFloat(text string) float64 {
	value, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
	if err != nil {
		return 0
	}
	return value
}

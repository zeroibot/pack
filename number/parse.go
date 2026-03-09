package number

import (
	"strconv"
	"strings"
)

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

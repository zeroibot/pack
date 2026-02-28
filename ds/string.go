package ds

import (
	"strconv"
	"strings"
)

type String string

// Len returns the String length
func (s String) Len() Int {
	return Int(len(s))
}

// IsEmpty checks if the String is empty
func (s String) IsEmpty() Boolean {
	return s == ""
}

// NotEmpty checks if the String is not empty
func (s String) NotEmpty() Boolean {
	return s != ""
}

// ToInt parses the string as an Int, defaults to 0 if invalid Int
func (s String) ToInt() Int {
	text := strings.TrimSpace(string(s))
	number, err := strconv.Atoi(text)
	if err != nil {
		return 0
	}
	return Int(number)
}

// ToUint parses the string as Uint, defaults to 0 if invalid Uint
func (s String) ToUint() Uint {
	text := strings.TrimSpace(string(s))
	number, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return 0
	}
	return Uint(number)
}

// ToFloat parses the string as a Float, defaults to 0 if invalid Float
func (s String) ToFloat() Float {
	text := strings.TrimSpace(string(s))
	number, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0
	}
	return Float(number)
}

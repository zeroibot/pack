// Package str contains string functions
package str

// Len returns the string length
func Len(s string) int {
	return len(s)
}

// IsEmpty checks if string is empty
func IsEmpty(s string) bool {
	return s == ""
}

// NotEmpty checks if string is not empty
func NotEmpty(s string) bool {
	return s != ""
}

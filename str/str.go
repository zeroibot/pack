// Package str contains string functions
package str

import (
	"fmt"
	"strings"
)

// Len returns the string length
func Len(text string) int {
	return len(text)
}

// IsEmpty checks if string is empty
func IsEmpty(text string) bool {
	return text == ""
}

// NotEmpty checks if string is not empty
func NotEmpty(text string) bool {
	return text != ""
}

// Guard returns the guard string if given string is empty, otherwise returns the given string
func Guard(text, guard string) string {
	if text == "" {
		return guard
	}
	return text
}

// Wrap wraps the given text with the given wrapper string (first char = left, second char= right)
func Wrap(text, wrapper string) string {
	var left byte = ' '
	var right byte = ' '
	if len(wrapper) > 0 {
		left = wrapper[0]
	}
	if len(wrapper) > 1 {
		right = wrapper[1]
	}
	text = fmt.Sprintf("%c%s%c", left, text, right)
	return strings.TrimSpace(text)
}

// WrapList joins the list items with a comma, and wraps the resulting string in given text (first char = left, second char = right)
func WrapList(items []string, wrapper string) string {
	return Wrap(strings.Join(items, ", "), wrapper)
}

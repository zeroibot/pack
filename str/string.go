// Package str contains string functions
package str

import (
	"strings"
	"unicode"
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

// CleanSplit splits the text by separator, trims each part's extra whitespace
func CleanSplit(text, sep string) []string {
	parts := strings.Split(text, sep)
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}

// CleanSplitN splits the text by separator, maximum of N parts, trims each part's extra whitespace
func CleanSplitN(text, sep string, count int) []string {
	parts := strings.SplitN(text, sep, count)
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}

// SpaceSplit splits text by whitespace
func SpaceSplit(text string) []string {
	return strings.Fields(strings.TrimSpace(text))
}

// CommaSplit splits text by comma, then trims each part's extra whitespace
func CommaSplit(text string) []string {
	return CleanSplit(text, ",")
}

// Lines splits the text by \n, then trims each line's extra whitespace
func Lines(text string) []string {
	return CleanSplit(text, "\n")
}

// Join joins the string parts by glue
func Join(glue string, parts ...string) string {
	return strings.Join(parts, glue)
}

// StartsWithUpper checks if string starts with uppercase letter
func StartsWithUpper(text string) bool {
	first := text[0]
	return 'A' <= first && first <= 'Z' // A-Z
}

// StartsWithLower checks if string starts with lowercase letter
func StartsWithLower(text string) bool {
	first := text[0]
	return 'a' <= first && first <= 'z' // a-z
}

// StartsWithDigit checks if string starts with digit
func StartsWithDigit(text string) bool {
	first := text[0]
	return '0' <= first && first <= '9'
}

// SpacePrefix gets the leading whitespace
func SpacePrefix(text string) string {
	suffix := TrimLeftSpace(text)
	return strings.TrimSuffix(text, suffix)
}

// SpaceSuffix gets the trailing whitespace
func SpaceSuffix(text string) string {
	prefix := TrimRightSpace(text)
	return strings.TrimPrefix(text, prefix)
}

// TrimLeftSpace trims left whitespace
func TrimLeftSpace(text string) string {
	return strings.TrimLeftFunc(text, unicode.IsSpace)
}

// TrimRightSpace trims right whitespace
func TrimRightSpace(text string) string {
	return strings.TrimRightFunc(text, unicode.IsSpace)
}

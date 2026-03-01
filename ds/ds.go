// Package ds contains data structures
package ds

import "strings"

// Number interface unifies the number types
type Number interface {
	Integer | Float
}

// Integer interface unifies the integer types
type Integer interface {
	~int | ~uint | ~int64
}

// Float interface unifies the float types
type Float interface {
	~float32 | ~float64
}

type StringBuilder struct {
	items []string
}

// NewStringBuilder creates a new StringBuilder
func NewStringBuilder() *StringBuilder {
	return &StringBuilder{items: make([]string, 0)}
}

// Add adds a string to the StringBuilder
func (sb *StringBuilder) Add(item string) {
	sb.items = append(sb.items, item)
}

// Build builds the full string parts, joined by the separator
func (sb *StringBuilder) Build(separator string) string {
	return strings.Join(sb.items, separator)
}

package str

import (
	"fmt"
	"slices"
	"strings"
)

// Builder is a utility type for building strings
type Builder struct {
	items []string
}

// NewBuilder creates a new string Builder
func NewBuilder() *Builder {
	return new(Builder{items: make([]string, 0)})
}

// Add adds a string to the string Builder
func (b *Builder) Add(item string) {
	b.items = append(b.items, item)
}

// AddFmt adds a formatted string to the string Builder
func (b *Builder) AddFmt(template string, args ...any) {
	b.Add(fmt.Sprintf(template, args...))
}

// AddItems adds multiple strings to the string Builder
func (b *Builder) AddItems(items ...string) {
	b.items = append(b.items, items...)
}

// Build builds the string parts, joined by the glue
func (b *Builder) Build(glue string) string {
	return strings.Join(b.items, glue)
}

// Repeat creates a new string by repeating given string and joining by glue
func Repeat(count int, text, glue string) string {
	return strings.Join(slices.Repeat([]string{text}, count), glue)
}

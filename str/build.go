package str

import (
	"fmt"
	"math/rand/v2"
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

const (
	upperLetters string = "ABCDEFGHJKLMNPQRSTUVWXYZ"
	lowerLetters string = "abcdefghjkmnpqrstuvwxyz"
	numbers      string = "23456789"
)

// RandomString creates a random string of given length using uppercase, lowercase letters, and numbers (flags)
func RandomString(length uint, useUpper, useLower, useNumber bool) string {
	charSource := ""
	if useUpper {
		charSource += upperLetters
	}
	if useLower {
		charSource += lowerLetters
	}
	if useNumber {
		charSource += numbers
	}
	numChars := len(charSource)
	if numChars == 0 {
		return ""
	}
	b := make([]byte, length)
	for i := range length {
		idx := rand.IntN(numChars)
		b[i] = charSource[idx]
	}
	return string(b)
}

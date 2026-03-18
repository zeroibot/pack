package tst

import (
	"maps"
	"slices"
	"testing"
)

type AssertFn[T any] = func(*testing.T, string, T, T)

// AssertEqual asserts that the two given values are equal
func AssertEqual[T comparable](t *testing.T, name string, a, b T) {
	if a != b {
		t.Errorf("%s = %v, want %v", name, a, b)
	}
}

// AssertListEqual asserts that the two given lists are equal
func AssertListEqual[S ~[]T, T comparable](t *testing.T, name string, a, b S) {
	if slices.Equal(a, b) == false {
		t.Errorf("%s = %v, want %v", name, a, b)
	}
}

// AssertMapEqual asserts that the two given maps are equal
func AssertMapEqual[M1, M2 ~map[K]V, K, V comparable](t *testing.T, name string, a M1, b M2) {
	if maps.Equal(a, b) == false {
		t.Errorf("%s = %v, want %v", name, a, b)
	}
}

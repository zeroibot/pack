package lang

import "cmp"

// Identity returns the passed argument, used as function argument
func Identity[T any](item T) T {
	return item
}

// IsEqual returns a function that checks if value is equal to given value, used as function argument
func IsEqual[T comparable](value T) func(T) bool {
	return func(item T) bool {
		return item == value
	}
}

// NotEqual returns a function that checks if value is not equal to given value, used as function argument
func NotEqual[T comparable](value T) func(T) bool {
	return func(item T) bool {
		return item != value
	}
}

// IsGreater returns a function that checks if value is greater than given value, used as function argument
func IsGreater[T cmp.Ordered](value T) func(T) bool {
	return func(item T) bool {
		return item > value
	}
}

// IsGreaterEqual returns a function that checks if value is greater than or equal to given value, used as function argument
func IsGreaterEqual[T cmp.Ordered](value T) func(T) bool {
	return func(item T) bool {
		return item >= value
	}
}

// IsLesser returns a function that checks if value is lesser than given value, used as function argument
func IsLesser[T cmp.Ordered](value T) func(T) bool {
	return func(item T) bool {
		return item < value
	}
}

// IsLesserEqual returns a function that checks if value is lesser than or equal to given value, used as function argument
func IsLesserEqual[T cmp.Ordered](value T) func(T) bool {
	return func(item T) bool {
		return item <= value
	}
}

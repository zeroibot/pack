package qb

import (
	"cmp"

	"github.com/roidaradal/pack/ds"
)

// Condition interface unifies all Condition objects:
// BuildCondition() method outputs the condition string and parameter values
type Condition interface {
	BuildCondition() (string, []any) // Return (condition string, parameter values)
}

// EmptyCondition creates a matchAllCondition
func EmptyCondition() Condition {
	return matchAllCondition{}
}

// EqualTo creates an Equal Condition
func EqualTo[T comparable](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opEqual)
}

// NotEqualTo creates a NotEqual Condition
func NotEqualTo[T comparable](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opNotEqual)
}

// HasPrefix creates a Prefix Condition
func HasPrefix(this *Instance, fieldRef *string, value string) Condition {
	return newValueCondition(this, fieldRef, value, opPrefix)
}

// HasSuffix creates a Suffix Condition
func HasSuffix(this *Instance, fieldRef *string, value string) Condition {
	return newValueCondition(this, fieldRef, value, opSuffix)
}

// HasSubstring creates a Substring Condition
func HasSubstring(this *Instance, fieldRef *string, value string) Condition {
	return newValueCondition(this, fieldRef, value, opSubstring)
}

// GreaterThan creates a GreaterThan Condition
func GreaterThan[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opGreater)
}

// GreaterEqualTo creates a GreaterThanOrEqual Condition
func GreaterEqualTo[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opGreaterEqual)
}

// LesserThan creates a LesserThan Condition
func LesserThan[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opLesser)
}

// LesserEqualTo creates a LesserThanOrEqual Condition
func LesserEqualTo[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opLesserEqual)
}

// InValues creates an In Condition
func InValues[T comparable](this *Instance, fieldRef *T, values ds.List[T]) Condition {
	return newListCondition(this, fieldRef, values, opIn, opEqual)
}

// NotInValues creates a NotIn Condition
func NotInValues[T comparable](this *Instance, fieldRef *T, values ds.List[T]) Condition {
	return newListCondition(this, fieldRef, values, opNotIn, opNotEqual)
}

// AndCondition creates an And Condition
func AndCondition(conditions ...Condition) Condition {
	return newMultiCondition(opAnd, conditions...)
}

// OrCondition creates an Or Condition
func OrCondition(conditions ...Condition) Condition {
	return newMultiCondition(opOr, conditions...)
}

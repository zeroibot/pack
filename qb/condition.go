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

// EqualCondition creates an Equal Condition
func EqualCondition[T comparable](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opEqual)
}

// NotEqualCondition creates a NotEqual Condition
func NotEqualCondition[T comparable](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opNotEqual)
}

// PrefixCondition creates a Prefix Condition
func PrefixCondition(this *Instance, fieldRef *string, value string) Condition {
	return newValueCondition(this, fieldRef, value, opPrefix)
}

// SuffixCondition creates a Suffix Condition
func SuffixCondition(this *Instance, fieldRef *string, value string) Condition {
	return newValueCondition(this, fieldRef, value, opSuffix)
}

// SubstringCondition creates a Substring Condition
func SubstringCondition(this *Instance, fieldRef *string, value string) Condition {
	return newValueCondition(this, fieldRef, value, opSubstring)
}

// GreaterCondition creates a GreaterThan Condition
func GreaterCondition[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opGreater)
}

// GreaterEqualCondition creates a GreaterThanOrEqual Condition
func GreaterEqualCondition[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opGreaterEqual)
}

// LesserCondition creates a LesserThan Condition
func LesserCondition[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opLesser)
}

// LesserEqualCondition creates a LesserThanOrEqual Condition
func LesserEqualCondition[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opLesserEqual)
}

// InCondition creates an In Condition
func InCondition[T comparable](this *Instance, fieldRef *T, values ds.List[T]) Condition {
	return newListCondition(this, fieldRef, values, opIn, opEqual)
}

// NotInCondition creates a NotIn Condition
func NotInCondition[T comparable](this *Instance, fieldRef *T, values ds.List[T]) Condition {
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

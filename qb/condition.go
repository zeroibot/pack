package qb

import (
	"cmp"

	"github.com/roidaradal/pack/ds"
)

// Condition interface unifies all Condition objects:
// Build() method outputs the condition string and parameter values
type Condition interface {
	Build() (string, ds.List[any]) // Return (condition string, parameter values)
}

// NoCondition creates a matchAllCondition
func NoCondition() Condition {
	return matchAllCondition{}
}

// Equal creates an Equal Condition
func Equal[T comparable](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opEqual)
}

// NotEqual creates a NotEqual Condition
func NotEqual[T comparable](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opNotEqual)
}

// Prefix creates a Prefix Condition
func Prefix(this *Instance, fieldRef *string, value string) Condition {
	return newValueCondition(this, fieldRef, value, opPrefix)
}

// Suffix creates a Suffix Condition
func Suffix(this *Instance, fieldRef *string, value string) Condition {
	return newValueCondition(this, fieldRef, value, opSuffix)
}

// Substring creates a Substring Condition
func Substring(this *Instance, fieldRef *string, value string) Condition {
	return newValueCondition(this, fieldRef, value, opSubstring)
}

// Greater creates a GreaterThan Condition
func Greater[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opGreater)
}

// GreaterEqual creates a GreaterThanOrEqual Condition
func GreaterEqual[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opGreaterEqual)
}

// Lesser creates a LesserThan Condition
func Lesser[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opLesser)
}

// LesserEqual creates a LesserThanOrEqual Condition
func LesserEqual[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opLesserEqual)
}

// In creates an In Condition
func In[T comparable](this *Instance, fieldRef *T, values ds.List[T]) Condition {
	return newListCondition(this, fieldRef, values, opIn, opEqual)
}

// NotIn creates a NotIn Condition
func NotIn[T comparable](this *Instance, fieldRef *T, values ds.List[T]) Condition {
	return newListCondition(this, fieldRef, values, opNotIn, opNotEqual)
}

// And creates an And Condition
func And(conditions ...Condition) Condition {
	return newMultiCondition(opAnd, conditions...)
}

// Or creates an Or Condition
func Or(conditions ...Condition) Condition {
	return newMultiCondition(opOr, conditions...)
}

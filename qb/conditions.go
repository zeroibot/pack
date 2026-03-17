package qb

import (
	"cmp"
	"slices"
	"strings"

	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/list"
)

// Condition interface unifies all Condition objects:
// BuildCondition() method outputs the condition string and parameter values
type Condition interface {
	BuildCondition() (string, []any) // Return (condition string, parameter values)
}

// DualCondition interface holds a Condition Builder and a struct Tester
type DualCondition[T any] interface {
	Condition
	Test(T) bool
}

// EmptyCondition creates a matchAllCondition
func EmptyCondition() Condition {
	return matchAllCondition{}
}

// NoCondition creates a matchAllCombo
func NoCondition[T any]() DualCondition[T] {
	return matchAllCombo[T]{}
}

// EqualTo creates an Equal Condition
func EqualTo[T comparable](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opEqual)
}

// Equal creates an Equal Combo
func Equal[T any, V comparable](this *Instance, fieldRef *V, value V) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return fieldValue == value
	})
	condition := newValueCondition(this, fieldRef, value, opEqual)
	return newValueCombo(condition, test)
}

// NotEqualTo creates a NotEqual Condition
func NotEqualTo[T comparable](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opNotEqual)
}

// NotEqual creates a NotEqual Combo
func NotEqual[T any, V comparable](this *Instance, fieldRef *V, value V) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return fieldValue != value
	})
	condition := newValueCondition(this, fieldRef, value, opNotEqual)
	return newValueCombo(condition, test)
}

// HasPrefix creates a Prefix Condition
func HasPrefix(this *Instance, fieldRef *string, value string) Condition {
	return newValueCondition(this, fieldRef, value, opPrefix)
}

// Prefix creates a Prefix Combo
func Prefix[T any](this *Instance, fieldRef *string, value string) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue string) bool {
		return strings.HasPrefix(fieldValue, value)
	})
	condition := newValueCondition(this, fieldRef, value, opPrefix)
	return newValueCombo(condition, test)
}

// HasSuffix creates a Suffix Condition
func HasSuffix(this *Instance, fieldRef *string, value string) Condition {
	return newValueCondition(this, fieldRef, value, opSuffix)
}

// Suffix creates a Suffix Combo
func Suffix[T any](this *Instance, fieldRef *string, value string) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue string) bool {
		return strings.HasSuffix(fieldValue, value)
	})
	condition := newValueCondition(this, fieldRef, value, opSuffix)
	return newValueCombo(condition, test)
}

// HasSubstring creates a Substring Condition
func HasSubstring(this *Instance, fieldRef *string, value string) Condition {
	return newValueCondition(this, fieldRef, value, opSubstring)
}

// Substring creates a Substring Combo
func Substring[T any](this *Instance, fieldRef *string, value string) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue string) bool {
		return strings.Contains(fieldValue, value)
	})
	condition := newValueCondition(this, fieldRef, value, opSubstring)
	return newValueCombo(condition, test)
}

// GreaterThan creates a GreaterThan Condition
func GreaterThan[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opGreater)
}

// Greater creates a GreaterThan Combo
func Greater[T any, V cmp.Ordered](this *Instance, fieldRef *V, value V) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return fieldValue > value
	})
	condition := newValueCondition(this, fieldRef, value, opGreater)
	return newValueCombo(condition, test)
}

// GreaterEqualTo creates a GreaterThanOrEqual Condition
func GreaterEqualTo[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opGreaterEqual)
}

// GreaterEqual creates a GreaterThanOrEqual Combo
func GreaterEqual[T any, V cmp.Ordered](this *Instance, fieldRef *V, value V) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return fieldValue >= value
	})
	condition := newValueCondition(this, fieldRef, value, opGreaterEqual)
	return newValueCombo(condition, test)
}

// LesserThan creates a LesserThan Condition
func LesserThan[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opLesser)
}

// Lesser creates a LesserThan Combo
func Lesser[T any, V cmp.Ordered](this *Instance, fieldRef *V, value V) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return fieldValue < value
	})
	condition := newValueCondition(this, fieldRef, value, opLesser)
	return newValueCombo(condition, test)
}

// LesserEqualTo creates a LesserThanOrEqual Condition
func LesserEqualTo[T cmp.Ordered](this *Instance, fieldRef *T, value T) Condition {
	return newValueCondition(this, fieldRef, value, opLesserEqual)
}

// LesserEqual creates a LesserThanOrEqual Combo
func LesserEqual[T any, V cmp.Ordered](this *Instance, fieldRef *V, value V) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return fieldValue <= value
	})
	condition := newValueCondition(this, fieldRef, value, opLesserEqual)
	return newValueCombo(condition, test)
}

// InValues creates an In Condition
func InValues[T comparable](this *Instance, fieldRef *T, values ds.List[T]) Condition {
	return newListCondition(this, fieldRef, values, opIn, opEqual)
}

// In creates an In Combo
func In[T any, V comparable](this *Instance, fieldRef *V, values ds.List[V]) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return slices.Contains(values, fieldValue)
	})
	condition := newListCondition(this, fieldRef, values, opIn, opEqual)
	return newListCombo(condition, test)
}

// NotInValues creates a NotIn Condition
func NotInValues[T comparable](this *Instance, fieldRef *T, values ds.List[T]) Condition {
	return newListCondition(this, fieldRef, values, opNotIn, opNotEqual)
}

// NotIn creates a NotIn Combo
func NotIn[T any, V comparable](this *Instance, fieldRef *V, values ds.List[V]) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return !slices.Contains(values, fieldValue)
	})
	condition := newListCondition(this, fieldRef, values, opNotIn, opNotEqual)
	return newListCombo(condition, test)
}

// AndCondition creates an And Condition
func AndCondition(conditions ...Condition) Condition {
	return newMultiCondition(opAnd, conditions...)
}

// And creates an And Combo
func And[T any](conditions ...DualCondition[T]) DualCondition[T] {
	test := func(item T) bool {
		return list.All(conditions, func(c DualCondition[T]) bool {
			return c.Test(item)
		})
	}
	return newMultiCombo(conditions, opAnd, test)
}

// OrCondition creates an Or Condition
func OrCondition(conditions ...Condition) Condition {
	return newMultiCondition(opOr, conditions...)
}

// Or creates an Or Combo
func Or[T any](conditions ...DualCondition[T]) DualCondition[T] {
	test := func(item T) bool {
		return list.Any(conditions, func(c DualCondition[T]) bool {
			return c.Test(item)
		})
	}
	return newMultiCombo(conditions, opOr, test)
}

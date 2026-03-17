package qb

import (
	"cmp"
	"slices"
	"strings"

	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/list"
)

// DualCondition interface holds a Condition Builder and a struct Tester
type DualCondition[T any] interface {
	Condition
	Test(T) bool
}

// NoCondition creates a matchAllCombo
func NoCondition[T any]() DualCondition[T] {
	return matchAllCombo[T]{}
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

// NotEqual creates a NotEqual Combo
func NotEqual[T any, V comparable](this *Instance, fieldRef *V, value V) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return fieldValue != value
	})
	condition := newValueCondition(this, fieldRef, value, opNotEqual)
	return newValueCombo(condition, test)
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

// Suffix creates a Suffix Combo
func Suffix[T any](this *Instance, fieldRef *string, value string) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue string) bool {
		return strings.HasSuffix(fieldValue, value)
	})
	condition := newValueCondition(this, fieldRef, value, opSuffix)
	return newValueCombo(condition, test)
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

// Greater creates a GreaterThan Combo
func Greater[T any, V cmp.Ordered](this *Instance, fieldRef *V, value V) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return fieldValue > value
	})
	condition := newValueCondition(this, fieldRef, value, opGreater)
	return newValueCombo(condition, test)
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

// Lesser creates a LesserThan Combo
func Lesser[T any, V cmp.Ordered](this *Instance, fieldRef *V, value V) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return fieldValue < value
	})
	condition := newValueCondition(this, fieldRef, value, opLesser)
	return newValueCombo(condition, test)
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

// In creates an In Combo
func In[T any, V comparable](this *Instance, fieldRef *V, values ds.List[V]) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return slices.Contains(values, fieldValue)
	})
	condition := newListCondition(this, fieldRef, values, opIn, opEqual)
	return newListCombo(condition, test)
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

// And creates an And Combo
func And[T any](conditions ...DualCondition[T]) DualCondition[T] {
	test := func(item T) bool {
		return list.All(conditions, func(c DualCondition[T]) bool {
			return c.Test(item)
		})
	}
	return newMultiCombo(conditions, opAnd, test)
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

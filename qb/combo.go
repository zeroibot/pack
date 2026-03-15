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

// NoConditionTest creates a matchAllCombo
func NoConditionTest[T any]() DualCondition[T] {
	return matchAllCombo[T]{}
}

// EqualTest creates an Equal Combo
func EqualTest[T any, V comparable](this *Instance, fieldRef *V, value V) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return fieldValue == value
	})
	condition := newValueCondition(this, fieldRef, value, opEqual)
	return newValueCombo(condition, test)
}

// NotEqualTest creates a NotEqual Combo
func NotEqualTest[T any, V comparable](this *Instance, fieldRef *V, value V) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return fieldValue != value
	})
	condition := newValueCondition(this, fieldRef, value, opNotEqual)
	return newValueCombo(condition, test)
}

// PrefixTest creates a Prefix Combo
func PrefixTest[T any](this *Instance, fieldRef *string, value string) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue string) bool {
		return strings.HasPrefix(fieldValue, value)
	})
	condition := newValueCondition(this, fieldRef, value, opPrefix)
	return newValueCombo(condition, test)
}

// SuffixTest creates a Suffix Combo
func SuffixTest[T any](this *Instance, fieldRef *string, value string) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue string) bool {
		return strings.HasSuffix(fieldValue, value)
	})
	condition := newValueCondition(this, fieldRef, value, opSuffix)
	return newValueCombo(condition, test)
}

// SubstringTest creates a Substring Combo
func SubstringTest[T any](this *Instance, fieldRef *string, value string) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue string) bool {
		return strings.Contains(fieldValue, value)
	})
	condition := newValueCondition(this, fieldRef, value, opSubstring)
	return newValueCombo(condition, test)
}

// GreaterTest creates a GreaterThan Combo
func GreaterTest[T any, V cmp.Ordered](this *Instance, fieldRef *V, value V) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return fieldValue > value
	})
	condition := newValueCondition(this, fieldRef, value, opGreater)
	return newValueCombo(condition, test)
}

// GreaterEqualTest creates a GreaterThanOrEqual Combo
func GreaterEqualTest[T any, V cmp.Ordered](this *Instance, fieldRef *V, value V) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return fieldValue >= value
	})
	condition := newValueCondition(this, fieldRef, value, opGreaterEqual)
	return newValueCombo(condition, test)
}

// LesserTest creates a LesserThan Combo
func LesserTest[T any, V cmp.Ordered](this *Instance, fieldRef *V, value V) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return fieldValue < value
	})
	condition := newValueCondition(this, fieldRef, value, opLesser)
	return newValueCombo(condition, test)
}

// LesserEqualTest creates a LesserThanOrEqual Combo
func LesserEqualTest[T any, V cmp.Ordered](this *Instance, fieldRef *V, value V) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return fieldValue <= value
	})
	condition := newValueCondition(this, fieldRef, value, opLesserEqual)
	return newValueCombo(condition, test)
}

// InTest creates an In Combo
func InTest[T any, V comparable](this *Instance, fieldRef *V, values ds.List[V]) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return slices.Contains(values, fieldValue)
	})
	condition := newListCondition(this, fieldRef, values, opIn, opEqual)
	return newListCombo(condition, test)
}

// NotInTest creates a NotIn Combo
func NotInTest[T any, V comparable](this *Instance, fieldRef *V, values ds.List[V]) DualCondition[T] {
	fieldName := this.getFieldName(fieldRef)
	test := createFieldValueTest[T](fieldName, func(fieldValue V) bool {
		return !slices.Contains(values, fieldValue)
	})
	condition := newListCondition(this, fieldRef, values, opNotIn, opNotEqual)
	return newListCombo(condition, test)
}

// AndTest creates an And Combo
func AndTest[T any](conditions ...DualCondition[T]) DualCondition[T] {
	test := func(item T) bool {
		return list.All(conditions, func(c DualCondition[T]) bool {
			return c.Test(item)
		})
	}
	return newMultiCombo(conditions, opAnd, test)
}

// OrTest creates an Or Combo
func OrTest[T any](conditions ...DualCondition[T]) DualCondition[T] {
	test := func(item T) bool {
		return list.Any(conditions, func(c DualCondition[T]) bool {
			return c.Test(item)
		})
	}
	return newMultiCombo(conditions, opOr, test)
}

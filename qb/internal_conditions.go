package qb

import (
	"fmt"
	"strings"

	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/dyn"
	"github.com/roidaradal/pack/str"
)

type testFn[T any] func(T) bool
type columnValuePair = ds.Tuple2[string, any]
type columnValueListPair = ds.Tuple2[string, []any]

// Internal: create the test for combos
func createFieldValueTest[T, V any](fieldName string, test testFn[V]) testFn[T] {
	return func(item T) bool {
		fieldValue, ok := getStructFieldValue[V](&item, fieldName)
		if !ok {
			return false
		}
		return test(fieldValue)
	}
}

const (
	falseCondition string = "false"
	trueCondition  string = "true"
)

type operator string

const (
	opEqual        operator = "="
	opNotEqual     operator = "<>"
	opGreater      operator = ">"
	opGreaterEqual operator = ">="
	opLesser       operator = "<"
	opLesserEqual  operator = "<="
	opIn           operator = "IN"
	opNotIn        operator = "NOT IN"
	opAnd          operator = "AND"
	opOr           operator = "OR"
	opPrefix       operator = "PREFIX"
	opSuffix       operator = "SUFFIX"
	opSubstring    operator = "SUBSTRING"
)

// Internal: create new Column-Value pair
func newColumnValue[T any](this *Instance, fieldRef *T, value T) ds.Option[columnValuePair] {
	column := this.Column(fieldRef)
	if column == "" {
		return ds.Nil[columnValuePair]()
	}
	column = this.dbType.prepareIdentifier(column)
	return ds.NewOption(&columnValuePair{V1: column, V2: value})
}

// Internal: create new Column-ValueList pair
func newColumnValueList[T any](this *Instance, fieldRef *T, values ds.List[T]) ds.Option[columnValueListPair] {
	column := this.Column(fieldRef)
	if column == "" {
		return ds.Nil[columnValueListPair]()
	}
	column = this.dbType.prepareIdentifier(column)
	return ds.NewOption(&columnValueListPair{V1: column, V2: values.ToAny()})
}

// Internal: Create new Column-Value pair, by getting the column name from type and field name
func newFieldColumnValue(this *Instance, typeName, fieldName string, value any) ds.Option[columnValuePair] {
	column := this.getFieldColumnName(typeName, fieldName)
	if column == "" {
		return ds.Nil[columnValuePair]()
	}
	column = this.dbType.prepareIdentifier(column)
	return ds.NewOption(&columnValuePair{V1: column, V2: value})
}

// Common: return 'false' as condition, empty list of values
func falseConditionValues() (string, []any) {
	return falseCondition, []any{}
}

// Common: return 'true' as condition, empty list of values
func trueConditionValues() (string, []any) {
	return trueCondition, []any{}
}

// Missing Condition: the default for UPDATE, DELETE to ensure condition is set.
type missingCondition struct{}

// BuildCondition to 'WHERE false'
func (c missingCondition) BuildCondition() (string, []any) {
	return falseConditionValues()
}

// Missing Combo: uses missingCondition
type missingCombo[T any] struct {
	missingCondition
}

func (c missingCombo[T]) Test(_ T) bool {
	return false
}

// MatchAll Condition: default for SELECT (no condition).
type matchAllCondition struct{}

// BuildCondition to 'WHERE true'
func (c matchAllCondition) BuildCondition() (string, []any) {
	return trueConditionValues()
}

// MatchAll Combo: uses matchAllCondition
type matchAllCombo[T any] struct {
	matchAllCondition
}

func (c matchAllCombo[T]) Test(_ T) bool {
	return true
}

// Value Condition: uses Column-Value pair (one value)
type valueCondition struct {
	pair ds.Option[columnValuePair]
	operator
}

// newValueCondition creates a new valueCondition
func newValueCondition[T any](this *Instance, fieldRef *T, value T, op operator) valueCondition {
	return valueCondition{
		newColumnValue(this, fieldRef, value),
		op,
	}
}

func (c valueCondition) BuildCondition() (string, []any) {
	if c.pair.IsNil() {
		// no pair = false condition
		return falseConditionValues()
	}
	column, value := c.pair.Value().Unpack()
	if column == "" {
		// no column = false condition
		return falseConditionValues()
	}
	return soloConditionValues(column, c.operator, value)
}

// Value Combo: uses valueCondition
type valueCombo[T any] struct {
	valueCondition
	test testFn[T]
}

// newValueCombo creates a new valueCombo
func newValueCombo[T any](condition valueCondition, test testFn[T]) valueCombo[T] {
	return valueCombo[T]{condition, test}
}

func (c valueCombo[T]) Test(item T) bool {
	return c.test(item)
}

// List Condition: uses Column-ValueList pair (multiple values)
type listCondition struct {
	pair         ds.Option[columnValueListPair]
	listOperator operator
	soloOperator operator
}

// newListCondition creates a new listCondition
func newListCondition[T any](this *Instance, fieldRef *T, values ds.List[T], listOperator, soloOperator operator) listCondition {
	return listCondition{
		newColumnValueList(this, fieldRef, values),
		listOperator,
		soloOperator,
	}
}

func (c listCondition) BuildCondition() (string, []any) {
	if c.pair.IsNil() {
		// no pair = false condition
		return falseConditionValues()
	}
	column, values := c.pair.Value().Unpack()
	numValues := len(values)
	if column == "" || numValues == 0 {
		// no column or no values = false condition
		return falseConditionValues()
	} else if numValues == 1 {
		// only 1 value = solo condition
		return soloConditionValues(column, c.soloOperator, values[0])
	}
	placeholders := str.Repeat(numValues, "?", ", ")
	condition := fmt.Sprintf("%s %s (%s)", column, c.listOperator, placeholders)
	return condition, values
}

// List Combo: uses listCondition
type listCombo[T any] struct {
	listCondition
	test testFn[T]
}

// newListCombo creates a new listCombo
func newListCombo[T any](condition listCondition, test testFn[T]) listCombo[T] {
	return listCombo[T]{condition, test}

}

func (c listCombo[T]) Test(item T) bool {
	return c.test(item)
}

// Multi Condition: used for joining multiple conditions through AND, OR
type multiCondition struct {
	conditions ds.List[Condition]
	operator
}

// newMultiCondition creates a new multiCondition
func newMultiCondition(op operator, conditions ...Condition) multiCondition {
	return multiCondition{
		conditions: conditions,
		operator:   op,
	}
}

func (c multiCondition) BuildCondition() (string, []any) {
	return buildMultiCondition(c.operator, c.conditions...)
}

// Multi Combo: uses multiCondition
type multiCombo[T any] struct {
	conditions ds.List[DualCondition[T]]
	operator
	test testFn[T]
}

// newMultiCombo creates a new multiCombo
func newMultiCombo[T any](conditions ds.List[DualCondition[T]], op operator, test testFn[T]) multiCombo[T] {
	return multiCombo[T]{conditions, op, test}
}

func (c multiCombo[T]) BuildCondition() (string, []any) {
	conditions := make([]Condition, len(c.conditions))
	for _, condition := range c.conditions {
		conditions = append(conditions, condition)
	}
	return buildMultiCondition(c.operator, conditions...)
}

func (c multiCombo[T]) Test(item T) bool {
	return c.test(item)
}

// Internal: build condition string and query parameter values list (corresponds to ? in the query);
// Used for solo value conditions
func soloConditionValues(column string, op operator, value any) (string, []any) {
	isValueNil := dyn.IsNil(value)
	if op == opEqual && isValueNil {
		return fmt.Sprintf("%s IS NULL", column), []any{}
	} else if op == opNotEqual && isValueNil {
		return fmt.Sprintf("%s IS NOT NULL", column), []any{}
	} else if op == opPrefix {
		// <column> LIKE 'prefix%'
		prefix := fmt.Sprintf("%v%%", value)
		return fmt.Sprintf("%s LIKE ?", column), []any{prefix}
	} else if op == opSuffix {
		// <column> LIKE '%suffix'
		suffix := fmt.Sprintf("%%%v", value)
		return fmt.Sprintf("%s LIKE ?", column), []any{suffix}
	} else if op == opSubstring {
		// <column> LIKE '%substring%'
		substring := fmt.Sprintf("%%%v%%", value)
		return fmt.Sprintf("%s LIKE ?", column), []any{substring}
	}
	return fmt.Sprintf("%s %s ?", column, op), []any{value}
}

// Internal: common steps for building the multi-condition
func buildMultiCondition(op operator, conditions ...Condition) (string, []any) {
	numConditions := len(conditions)
	switch numConditions {
	case 0:
		// no conditions = false condition
		return falseConditionValues()
	case 1:
		// one condition = only build that one
		return conditions[0].BuildCondition()
	default:
		// multiple conditions
		conditionStrings := make([]string, 0, numConditions)
		allValues := make([]any, 0)
		for _, condition := range conditions {
			if condition == nil {
				continue // skip nil conditions
			}
			conditionString, values := condition.BuildCondition()
			if conditionString == falseCondition {
				// If any condition fails, return false condition immediately
				return falseConditionValues()
			}
			conditionStrings = append(conditionStrings, conditionString)
			allValues = append(allValues, values...)
		}
		// Join by operator and wrap in parentheses
		glue := fmt.Sprintf(" %s ", op)
		condition := fmt.Sprintf("(%s)", strings.Join(conditionStrings, glue))
		return condition, allValues
	}
}

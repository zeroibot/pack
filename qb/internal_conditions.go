package qb

import (
	"fmt"
	"strings"

	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/str"
)

// Missing Condition: the default for UPDATE, DELETE to ensure condition is set.
type missingCondition struct{}

// Build missingCondition to 'WHERE false'
func (c missingCondition) Build() (string, ds.List[any]) {
	return falseConditionValues()
}

// MatchAll Condition: default for SELECT (no condition).
type matchAllCondition struct{}

// Build matchAllCondition to 'WHERE true'
func (c matchAllCondition) Build() (string, ds.List[any]) {
	return trueConditionValues()
}

// Value Condition: uses Column-Value pair (one value)
type valueCondition struct {
	pair     ds.Option[columnValuePair]
	operator operation
}

// newValueCondition creates a new valueCondition
func newValueCondition[T any](this *Instance, fieldRef *T, value T, operator operation) valueCondition {
	return valueCondition{
		newColumnValue(this, fieldRef, value),
		operator,
	}
}

// Build valueCondition
func (c valueCondition) Build() (string, ds.List[any]) {
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

// List Condition: uses Column-ValueList pair (multiple values)
type listCondition struct {
	pair         ds.Option[columnValueListPair]
	listOperator operation
	soloOperator operation
}

// newListCondition creates a new listCondition
func newListCondition[T any](this *Instance, fieldRef *T, values ds.List[T], listOperator, soloOperator operation) listCondition {
	return listCondition{
		newColumnValueList(this, fieldRef, values),
		listOperator,
		soloOperator,
	}
}

// Build listCondition
func (c listCondition) Build() (string, ds.List[any]) {
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

// Multi Condition: used for joining multiple conditions through AND, OR
type multiCondition struct {
	conditions ds.List[Condition]
	operator   operation
}

// newMultiCondition creates a new multiCondition
func newMultiCondition(operator operation, conditions ...Condition) multiCondition {
	return multiCondition{
		conditions: conditions,
		operator:   operator,
	}
}

// Build multiCondition
func (c multiCondition) Build() (string, ds.List[any]) {
	numConditions := c.conditions.Len()
	switch numConditions {
	case 0:
		// no conditions = false condition
		return falseConditionValues()
	case 1:
		// one condition = only build that one
		return c.conditions[0].Build()
	default:
		// multiple conditions
		conditions := make([]string, 0, numConditions)
		allValues := make(ds.List[any], 0)
		for _, condition := range c.conditions {
			if condition == nil {
				continue // skip nil conditions
			}
			conditionString, values := condition.Build()
			if conditionString == falseCondition {
				// If any condition fails, return false condition immediately
				return falseConditionValues()
			}
			conditions = append(conditions, conditionString)
			allValues = append(allValues, values...)
		}
		// Join by operator and wrap in parentheses
		glue := fmt.Sprintf(" %s ", c.operator)
		condition := fmt.Sprintf("(%s)", strings.Join(conditions, glue))
		return condition, allValues
	}
}

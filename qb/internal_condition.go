package qb

import (
	"fmt"

	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/dyn"
)

type columnValuePair = ds.Tuple2[string, any]
type columnValueListPair = ds.Tuple2[string, ds.List[any]]

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
	return ds.NewOption(&columnValuePair{V1: column, V2: value})
}

// Internal: create new Column-ValueList pair
func newColumnValueList[T any](this *Instance, fieldRef *T, values ds.List[T]) ds.Option[columnValueListPair] {
	column := this.Column(fieldRef)
	if column == "" {
		return ds.Nil[columnValueListPair]()
	}
	return ds.NewOption(&columnValueListPair{V1: column, V2: values.ToAny()})
}

// Internal: Create new Column-Value pair, by getting the column name from type and field name
func newFieldColumnValue(this *Instance, typeName, fieldName string, value any) ds.Option[columnValuePair] {
	column := this.getFieldColumnName(typeName, fieldName)
	if column == "" {
		return ds.Nil[columnValuePair]()
	}
	return ds.NewOption(&columnValuePair{V1: column, V2: value})
}

// Common: return 'false' as condition, empty list of values
func falseConditionValues() (string, ds.List[any]) {
	return falseCondition, ds.List[any]{}
}

// Common: return 'true' as condition, empty list of values
func trueConditionValues() (string, ds.List[any]) {
	return trueCondition, ds.List[any]{}
}

// Internal: build condition string and query parameter values list (corresponds to ? in the query);
// Used for solo value conditions
func soloConditionValues(column string, op operator, value any) (string, ds.List[any]) {
	isValueNil := dyn.IsNil(value)
	if op == opEqual && isValueNil {
		return fmt.Sprintf("%s IS NULL", column), ds.List[any]{}
	} else if op == opNotEqual && isValueNil {
		return fmt.Sprintf("%s IS NOT NULL", column), ds.List[any]{}
	} else if op == opPrefix {
		// <column> LIKE 'prefix%'
		prefix := fmt.Sprintf("%v%%", value)
		return fmt.Sprintf("%s LIKE ?", column), ds.List[any]{prefix}
	} else if op == opSuffix {
		// <column> LIKE '%suffix'
		suffix := fmt.Sprintf("%%%v", value)
		return fmt.Sprintf("%s LIKE ?", column), ds.List[any]{suffix}
	} else if op == opSubstring {
		// <column> LIKE '%substring%'
		substring := fmt.Sprintf("%%%v%%", value)
		return fmt.Sprintf("%s LIKE ?", column), ds.List[any]{substring}
	}
	return fmt.Sprintf("%s %s ?", column, op), ds.List[any]{value}
}

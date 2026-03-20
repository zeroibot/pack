package qb

import (
	"fmt"
	"strings"

	"github.com/roidaradal/pack/dyn"
)

// CountQuery counts the number of rows that satisfy the condition
type CountQuery[T any] struct {
	conditionQuery[T]
}

// ValueQuery selects a single column value from the table from one row that satisfies the condition
// Type T = object type, V = value type
type ValueQuery[T, V any] struct {
	conditionQuery[T]
	typeName   string
	columnName string
	reader     RowReader[T]
}

// SelectRowQuery selects a single row from the table that satisfies the condition
type SelectRowQuery[T any] struct {
	conditionQuery[T]
	columnList
	reader RowReader[T]
}

// TopRowQuery selects the top N rows from the table that satisfy the condition
type TopRowQuery[T any] struct {
	conditionQuery[T]
	orderedLimit
	columnList
	reader RowReader[T]
}

// NewCountQuery creates a new CountQuery
func NewCountQuery[T any](this *Instance, table string) *CountQuery[T] {
	q := new(CountQuery[T])
	q.initializeRequired(this, table)
	return q
}

// NewValueQuery creates a new ValueQuery
func NewValueQuery[T, V any](this *Instance, table string, fieldRef *V) *ValueQuery[T, V] {
	q := new(ValueQuery[T, V])
	q.initializeRequired(this, table)
	var item T
	q.typeName = dyn.TypeName(item)
	columnName := this.Column(fieldRef)
	if columnName != "" {
		// Note: create RowReader before preparing identifier, since RowReader cannot recognized a processed column
		q.reader = NewRowReader[T](this, columnName)
		columnName = this.prepareIdentifier(columnName)
	}
	q.columnName = columnName
	return q
}

// NewSelectRowQuery creates a new SelectRowQuery, which only uses selected columns
func NewSelectRowQuery[T any](this *Instance, table string, reader RowReader[T]) *SelectRowQuery[T] {
	q := new(SelectRowQuery[T])
	q.initializeRequired(this, table)
	q.reader = reader
	q.columns = make([]string, 0)
	return q
}

// NewFullSelectRowQuery creates a new SelectRowQuery, which uses all columns
func NewFullSelectRowQuery[T any](this *Instance, table string, reader RowReader[T]) *SelectRowQuery[T] {
	q := NewSelectRowQuery(this, table, reader)
	var item T
	q.Columns(this, this.allColumns(item)...)
	return q
}

// NewTopRowQuery creates a new TopRowQuery
func NewTopRowQuery[T any](this *Instance, table string, reader RowReader[T]) *TopRowQuery[T] {
	q := new(TopRowQuery[T])
	q.initializeRequired(this, table)
	q.reader = reader
	q.limit = 1
	var item T
	q.Columns(this, this.allColumns(item)...)
	return q
}

// BuildQuery returns the query string and parameter values of CountQuery
func (q *CountQuery[T]) BuildQuery() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil {
		return emptyQueryValues()
	}
	query := "SELECT COUNT(*) FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.table, condition)
	return query, values
}

// BuildQuery returns the query string and parameter values of ValueQuery
func (q *ValueQuery[T, V]) BuildQuery() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || q.columnName == "" {
		return emptyQueryValues()
	}
	query := "SELECT %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.columnName, q.table, condition)
	return query, values
}

// BuildQuery returns the query string and parameter values of SelectRowQuery
func (q *SelectRowQuery[T]) BuildQuery() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || len(q.columns) == 0 {
		return emptyQueryValues()
	}
	columns := strings.Join(q.columns, ", ")
	query := "SELECT %s FROM %s WHERE %s LIMIT 1"
	query = fmt.Sprintf(query, columns, q.table, condition)
	return query, values
}

// BuildQuery returns the query string and parameter values of TopRowQuery
func (q *TopRowQuery[T]) BuildQuery() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || len(q.columns) == 0 || len(q.orders) == 0 {
		return emptyQueryValues()
	}
	columns := strings.Join(q.columns, ", ")
	query := "SELECT %s FROM %s WHERE %s %s"
	query = fmt.Sprintf(query, columns, q.table, condition, q.orderLimitString())
	return query, values
}

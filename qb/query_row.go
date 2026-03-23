package qb

import (
	"fmt"
	"strings"

	"github.com/roidaradal/pack/db"
	"github.com/roidaradal/pack/list"
)

// CountQuery counts the number of rows that satisfy the condition
type CountQuery[T any] struct {
	conditionQuery[T]
}

// ValueQuery selects a single column value from the table from one row that satisfies the condition
// Type T = object type, V = value type
type ValueQuery[T, V any] struct {
	conditionQuery[T]
	columnReader[T]
}

// SelectRowQuery selects a single row from the table that satisfies the condition
type SelectRowQuery[T any] struct {
	conditionQuery[T]
	columnsReader[T]
}

// TopRowQuery selects the top N rows from the table that satisfy the condition
type TopRowQuery[T any] struct {
	conditionQuery[T]
	columnsReader[T]
	orderedLimit
}

// TopValueQuery selects the top N values from the table that satisfy the condition
type TopValueQuery[T, V any] struct {
	conditionQuery[T]
	columnReader[T]
	orderedLimit
}

// SumQuery sums up the selected columns for rows that satisfy the condition
type SumQuery[T any] struct {
	conditionQuery[T]
	columnsReader[T]
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
	q.columnReader.initialize(this, this.Column(fieldRef))
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
	q.useAllColumns(this)
	return q
}

// NewTopRowQuery creates a new TopRowQuery
func NewTopRowQuery[T any](this *Instance, table string, reader RowReader[T]) *TopRowQuery[T] {
	q := new(TopRowQuery[T])
	q.initializeRequired(this, table)
	q.limit = 1
	q.reader = reader
	q.useAllColumns(this)
	return q
}

// NewTopValueQuery creates a new TopValueQuery
func NewTopValueQuery[T, V any](this *Instance, table string, fieldRef *V) *TopValueQuery[T, V] {
	q := new(TopValueQuery[T, V])
	q.initializeRequired(this, table)
	q.columnReader.initialize(this, this.Column(fieldRef))
	q.limit = 1
	return q
}

// NewSumQuery creates a new SumQuery
func NewSumQuery[T any](this *Instance, table string, reader RowReader[T]) *SumQuery[T] {
	q := new(SumQuery[T])
	q.initializeOptional(this, table)
	q.reader = reader
	q.columns = make([]string, 0)
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

// BuildQuery returns the query string and parameter values of TopValueQuery
func (q *TopValueQuery[T, V]) BuildQuery() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || q.columnName == "" || len(q.orders) == 0 {
		return emptyQueryValues()
	}
	query := "SELECT %s FROM %s WHERE %s %s"
	query = fmt.Sprintf(query, q.columnName, q.table, condition, q.orderLimitString())
	return query, values
}

// BuildQuery returns the query string and parameter values of SumQuery
func (q *SumQuery[T]) BuildQuery() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || len(q.columns) == 0 {
		return emptyQueryValues()
	}
	sumColumns := list.Map(q.columns, func(column string) string {
		return fmt.Sprintf("SUM(%s)", column)
	})
	columns := strings.Join(sumColumns, ", ")
	query := "SELECT %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, columns, q.table, condition)
	return query, values
}

// Count returns the number of rows that satisfy the CountQuery
func (q *CountQuery[T]) Count(dbc db.Conn) (int, error) {
	query, values, err := preQueryCheck(q, dbc)
	if err != nil {
		return 0, err
	}
	count := 0
	err = dbc.QueryRow(query, values...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Exists checks if there is at least 1 row that satisfies the CountQuery
func (q *CountQuery[T]) Exists(dbc db.Conn) (bool, error) {
	count, err := q.Count(dbc)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// QueryValue executes the ValueQuery and gets the column value
func (q *ValueQuery[T, V]) QueryValue(this *Instance, dbc db.Conn) (V, error) {
	var zero V
	query, values, err := preReadCheck(q, dbc, q.reader)
	if err != nil {
		return zero, err
	}
	row := dbc.QueryRow(query, values...)
	result := q.reader(row)
	if result.IsError() {
		return zero, result.Error()
	}
	return getStructTypedColumnValue[V](this, new(result.Value()), q.typeName, q.columnName)
}

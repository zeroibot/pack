package qb

import (
	"fmt"

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

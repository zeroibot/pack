package qb

import (
	"fmt"
)

// DistinctValuesQuery selects the distinct values for specified column that satisfies the condition
// T = object type, V = value type
type DistinctValuesQuery[T, V any] struct {
	conditionQuery[T]
	columnReader[T]
}

// LookupQuery selects two columns from a table and creates a lookup map for rows that satisfy the condition
// T = object type, K = key type, V = value type
type LookupQuery[T any, K comparable, V any] struct {
	conditionQuery[T]
	typeName    string
	keyColumn   string
	valueColumn string
	reader      RowReader[T]
}

// NewDistinctValuesQuery creates a new DistinctValuesQuery
func NewDistinctValuesQuery[T, V any](this *Instance, table string, fieldRef *V) *DistinctValuesQuery[T, V] {
	q := new(DistinctValuesQuery[T, V])
	q.initializeOptional(this, table)
	q.columnReader.initialize(this, this.Column(fieldRef))
	return q
}

// NewLookupQuery creates a new LookupQuery
func NewLookupQuery[T any, K comparable, V any](this *Instance, table string, keyFieldRef *K, valueFieldRef *V) *LookupQuery[T, K, V] {
	q := new(LookupQuery[T, K, V])
	q.initializeOptional(this, table)
	columns := this.Columns(keyFieldRef, valueFieldRef)
	if len(columns) == 2 {
		q.keyColumn = this.prepareIdentifier(columns[0])
		q.valueColumn = this.prepareIdentifier(columns[1])
		q.reader = NewRowReader[T](this, columns...)
	}
	return q
}

// BuildQuery returns the query string and parameter values of DistinctValuesQuery
func (q *DistinctValuesQuery[T, V]) BuildQuery() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || q.columnName == "" {
		return emptyQueryValues()
	}
	query := "SELECT DISTINCT %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.columnName, q.table, condition)
	return query, values
}

// BuildQuery returns the query string and parameter values of LookupQuery
func (q *LookupQuery[T, K, V]) BuildQuery() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || q.keyColumn == "" || q.valueColumn == "" {
		return emptyQueryValues()
	}
	query := "SELECT %s, %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.keyColumn, q.valueColumn, q.table, condition)
	return query, values
}

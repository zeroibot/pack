package qb

import "fmt"

// DistinctValuesQuery selects the distinct values for specified column that satisfies the condition
// T = object type, V = value type
type DistinctValuesQuery[T, V any] struct {
	conditionQuery[T]
	columnReader[T]
}

// NewDistinctValuesQuery creates a new DistinctValuesQuery
func NewDistinctValuesQuery[T, V any](this *Instance, table string, fieldRef *V) *DistinctValuesQuery[T, V] {
	q := new(DistinctValuesQuery[T, V])
	q.initializeOptional(this, table)
	q.columnReader.initialize(this, this.Column(fieldRef))
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

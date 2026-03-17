package qb

import (
	"fmt"
)

type DeleteQuery[T any] struct {
	conditionQuery[T]
}

// NewDeleteQuery creates a new DeleteQuery
func NewDeleteQuery[T any](this *Instance, table string) *DeleteQuery[T] {
	q := new(DeleteQuery[T])
	q.initializeRequired(this, table)
	return q
}

// BuildQuery returns the query string and parameter values
func (q DeleteQuery[T]) BuildQuery() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil {
		return emptyQueryValues()
	}
	query := "DELETE FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.table, condition)
	return query, values
}

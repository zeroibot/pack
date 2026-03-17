package qb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/roidaradal/pack/dyn"
)

var (
	errEmptyTable = errors.New("empty table")
)

// Query interface unifies all Query types:
// BuildQuery() method outputs the query string and parameter values
type Query interface {
	BuildQuery() (string, []any) // Return (query string, parameter values)
}

// baseQuery is an abstract Query with a table name.
// It does not implement the BuildQuery() method; it is embedded by concrete Queries for method reuse
type baseQuery struct {
	table string
}

// conditionQuery is an abstract Query with a table name and a DualCondition.
// It does not implement the BuildQuery() method; it is embedded by concrete Queries for method reuse
type conditionQuery[T any] struct {
	baseQuery
	condition DualCondition[T]
}

// initialize the baseQuery
func (q *baseQuery) initialize(this *Instance, table string) {
	q.table = this.dbType.prepareIdentifier(table)
}

// initialize a conditionQuery with required condition
func (q *conditionQuery[T]) initializeRequired(this *Instance, table string) {
	q.baseQuery.initialize(this, table)
	// if Condition is not set later (required), defaults to false condition
	q.condition = missingCombo[T]{}
}

// initialize a conditionQuery with optional condition
func (q *conditionQuery[T]) initializeOptional(this *Instance, table string) {
	q.baseQuery.initialize(this, table)
	// if Condition is not set later (optional), defaults to true condition
	q.condition = matchAllCombo[T]{}
}

// Where sets the Query Condition
func (q *conditionQuery[T]) Where(condition DualCondition[T]) {
	q.condition = condition
}

// preBuildCheck checks if the table is set
func (q *baseQuery) preBuildCheck() error {
	if q.table == "" {
		return errEmptyTable
	}
	return nil
}

// preBuildCheck checks if the table is set and builds the Condition
func (q *conditionQuery[T]) preBuildCheck() (string, []any, error) {
	err := q.baseQuery.preBuildCheck()
	condition, values := q.condition.BuildCondition()
	return condition, values, err
}

// ToString builds the Query string
func ToString(q Query) string {
	query, rawValues := q.BuildQuery()
	values := make([]any, len(rawValues))
	formats := make([]string, len(rawValues))
	for i, value := range rawValues {
		typeName := fmt.Sprintf("%T", value)
		if strings.HasPrefix(typeName, "*") {
			values[i] = dyn.MustDeref(value)
		} else {
			values[i] = value
		}
		if _, ok := values[i].(string); ok {
			formats[i] = "%q"
		} else {
			formats[i] = "%v"
		}
	}
	for _, format := range formats {
		query = strings.Replace(query, "?", format, 1)
	}
	return fmt.Sprintf(query, values...)
}

// emptyQueryValues returns an empty query and empty list of values
func emptyQueryValues() (string, []any) {
	return "", []any{}
}

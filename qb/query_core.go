package qb

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errEmptyTable = errors.New("empty table")
)

// Query interface unifies all Query types:
// BuildQuery() method outputs the query string and parameter values
type Query interface {
	BuildQuery() (string, []any) // Return (query string, parameter values)
}

// emptyQueryValues returns an empty query and empty list of values
func emptyQueryValues() (string, []any) {
	return "", []any{}
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
	if table != "" {
		table = this.dbType.prepareIdentifier(table)
	}
	q.table = table
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

// Test uses the test function of the DualCondition
func (q *conditionQuery[T]) Test(item T) bool {
	return q.condition.Test(item)
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

// orderedLimit is an abstractQuery with order column(s) and a limit.
// It does not implement the BuildQuery method; it is embedded by concrete Queries for method reuse
type orderedLimit struct {
	orders []string
	limit  uint
}

// OrderAsc adds a column with ascending order, returns orderedLimit for chaining
func (q *orderedLimit) OrderAsc(this *Instance, column string) *orderedLimit {
	if column == "" {
		return q
	}
	order := fmt.Sprintf("%s ASC", this.dbType.prepareIdentifier(column))
	q.orders = append(q.orders, order)
	return q
}

// OrderDesc adds a column with descending order, returns orderedLimit for chaining
func (q *orderedLimit) OrderDesc(this *Instance, column string) *orderedLimit {
	if column == "" {
		return q
	}
	order := fmt.Sprintf("%s DESC", this.dbType.prepareIdentifier(column))
	q.orders = append(q.orders, order)
	return q
}

// Limit sets the query limit, returns orderedLimit for chaining.
// Setting to 0 removes the limit
func (q *orderedLimit) Limit(limit uint) *orderedLimit {
	q.limit = limit
	return q
}

// fullString builds the orderString and limitString
func (q *orderedLimit) fullString() string {
	output := make([]string, 0, 2)
	orderString := q.orderString()
	if orderString != "" {
		output = append(output, fmt.Sprintf("ORDER BY %s", orderString))
	}
	if q.limit > 0 {
		output = append(output, fmt.Sprintf("LIMIT %d", q.limit))
	}
	return strings.Join(output, " ")
}

// orderString builds the list of orders into a string
func (q *orderedLimit) orderString() string {
	return strings.Join(q.orders, ", ")
}

// mustLimitString builds the orderString and limitString, but only includes the orderString if limitString is not empty
func (q *orderedLimit) mustLimitString() string {
	if q.limit == 0 {
		return ""
	}
	return q.fullString()
}

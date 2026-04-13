package qb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zeroibot/pack/db"
	"github.com/zeroibot/pack/dyn"
)

var (
	errEmptyQuery          = errors.New("empty query")
	errEmptyTable          = errors.New("empty table")
	errFailedResultCheck   = errors.New("failed result check")
	errFailedTypeAssertion = errors.New("failed type assertion")
	errNoChecker           = errors.New("no result checker")
	errNoDBConnection      = errors.New("no db connection")
	errNoReader            = errors.New("no row reader")
	errNoTx                = errors.New("no db transaction")
	errNotFoundField       = errors.New("field not found")
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
	condition Condition
	combo     DualCondition[T]
	useCombo  bool
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
	q.condition = missingCondition{}
	q.combo = missingCombo[T]{}
}

// initialize a conditionQuery with optional condition
func (q *conditionQuery[T]) initializeOptional(this *Instance, table string) {
	q.baseQuery.initialize(this, table)
	// if Condition is not set later (optional), defaults to true condition
	q.condition = matchAllCondition{}
	q.combo = matchAllCombo[T]{}
}

// Test uses the test function of the DualCondition
func (q *conditionQuery[T]) Test(item T) bool {
	return q.combo.Test(item)
}

// Where sets the Query Condition
func (q *conditionQuery[T]) Where(condition Condition) {
	q.useCombo = false
	q.condition = condition
}

func (q *conditionQuery[T]) Where2(condition DualCondition[T]) {
	q.useCombo = true
	q.combo = condition
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
	var condition string
	var values []any
	if q.useCombo {
		condition, values = q.combo.BuildCondition()
	} else {
		condition, values = q.condition.BuildCondition()
	}
	return condition, values, err
}

// preQueryCheck checks if the db connection is set and builds the Query
func preQueryCheck(q Query, dbc db.Conn) (string, []any, error) {
	var err error = nil
	query, values := q.BuildQuery()
	if dbc == nil {
		err = errNoDBConnection
	} else if query == "" {
		err = errEmptyQuery
	}
	return query, values, err
}

// preReadCheck is performed before a SELECT query to check the db connection and reader, and builds the Query
func preReadCheck[T any](q Query, dbc db.Conn, reader RowReader[T]) (string, []any, error) {
	query, values, err := preQueryCheck(q, dbc)
	if err != nil {
		return query, values, err
	}
	if reader == nil {
		err = errNoReader
	}
	return query, values, err
}

// orderedLimit is a Query part with order column(s) and a limit.
// It does not implement the BuildQuery method; it is embedded by concrete Queries for method reuse
type orderedLimit struct {
	orders []string
	limit  uint
}

// columnsReader is a Query part with a list of columns, and a RowReader
// It does not implement the BuildQuery method; it is embedded by concrete Queries for method reuse
type columnsReader[T any] struct {
	columns []string
	reader  RowReader[T]
}

// columnReader is a Query part with typeName, columnName, and a RowReader
// It does not implement the BuildQuery method; it is embedded by concrete Queries for method reuse
type columnReader[T any] struct {
	typeName   string
	columnName string
	reader     RowReader[T]
}

// initialize sets the typeName and columnName, and the reader if the columnName is not blank, for columnReader
func (q *columnReader[T]) initialize(this *Instance, columnName string) {
	var item T
	q.typeName = dyn.TypeName(item)
	if columnName != "" {
		q.columnName = this.prepareIdentifier(columnName)
		q.reader = NewRowReader[T](this, columnName)
	}
}

// useAllColumns sets the columns of columnsReader to all columns of the given type
func (q *columnsReader[T]) useAllColumns(this *Instance) {
	var item T
	q.Columns(this, this.allColumns(item)...)
}

// Columns sets the columns of columnsReader
func (q *columnsReader[T]) Columns(this *Instance, columns ...string) {
	columns2 := make([]string, 0, len(columns))
	for _, column := range columns {
		if column == "" {
			continue
		}
		columns2 = append(columns2, this.prepareIdentifier(column))
	}
	q.columns = columns2
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

// orderLimitString builds the orderString and limitString
func (q *orderedLimit) orderLimitString() string {
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
	return q.orderLimitString()
}

// tryAppend tries to append the given string at the end if it is not blank
func tryAppend(query string, suffix string) string {
	if suffix != "" {
		return fmt.Sprintf("%s %s", query, suffix)
	}
	return query
}

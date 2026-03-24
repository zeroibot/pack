package qb

import (
	"fmt"
	"strings"

	"github.com/roidaradal/pack/db"
	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/number"
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

// SelectRowsQuery selects the rows from the table that satisfy the condition
type SelectRowsQuery[T any] struct {
	conditionQuery[T]
	columnsReader[T]
	orderedLimit
	offset uint
}

// GroupCountQuery gets the counts of rows grouped by a column
type GroupCountQuery[T any, K comparable] struct {
	conditionQuery[T]
	groupColumn string
}

// GroupSumQuery gets the sum of row columns grouped by a column
type GroupSumQuery[T any, K comparable, V number.Type] struct {
	conditionQuery[T]
	groupColumn string
	sumColumn   string
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

// NewSelectRowsQuery creates a new SelectRowsQuery, which uses only selected columns
func NewSelectRowsQuery[T any](this *Instance, table string, reader RowReader[T]) *SelectRowsQuery[T] {
	q := new(SelectRowsQuery[T])
	q.initializeOptional(this, table)
	q.reader = reader
	q.columns = make([]string, 0)
	return q
}

// NewFullSelectRowsQuery creates a new SelectRowsQuery, which uses all columns
func NewFullSelectRowsQuery[T any](this *Instance, table string, reader RowReader[T]) *SelectRowsQuery[T] {
	q := NewSelectRowsQuery(this, table, reader)
	q.useAllColumns(this)
	return q
}

// NewGroupCountQuery creates a new GroupCountQuery
func NewGroupCountQuery[T any, K comparable](this *Instance, table string, groupFieldRef *K) *GroupCountQuery[T, K] {
	q := new(GroupCountQuery[T, K])
	q.initializeOptional(this, table)
	columnName := this.Column(groupFieldRef)
	if columnName != "" {
		q.groupColumn = this.prepareIdentifier(columnName)
	}
	return q
}

// NewGroupSumQuery creates a new GroupSumQuery
func NewGroupSumQuery[T any, K comparable, V number.Type](this *Instance, table string, groupFieldRef *K, sumFieldRef *V) *GroupSumQuery[T, K, V] {
	q := new(GroupSumQuery[T, K, V])
	q.initializeOptional(this, table)
	columns := this.Columns(groupFieldRef, sumFieldRef)
	if len(columns) == 2 {
		q.groupColumn = this.prepareIdentifier(columns[0])
		q.sumColumn = this.prepareIdentifier(columns[1])
	}
	return q
}

// Page sets the page number and batch size for a paginated SelectRowsQuery
func (q *SelectRowsQuery[T]) Page(number, batchSize uint) {
	q.offset = (number - 1) * batchSize
	q.limit = batchSize
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

// BuildQuery returns the query string and parameter values of SelectRowsQuery
func (q *SelectRowsQuery[T]) BuildQuery() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || len(q.columns) == 0 {
		return emptyQueryValues()
	}
	columns := strings.Join(q.columns, ", ")
	query := "SELECT %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, columns, q.table, condition)

	tail := make([]string, 0, 2)
	orderString := q.orderString()
	if orderString != "" {
		tail = append(tail, fmt.Sprintf("ORDER BY %s", orderString))
	}
	if q.offset > 0 || q.limit > 0 {
		tail = append(tail, fmt.Sprintf("LIMIT %d, %d", q.offset, q.limit))
	}
	suffix := strings.Join(tail, " ")
	query = tryAppend(query, suffix)

	return query, values
}

// BuildQuery returns the query string and parameter values of GroupCountQuery
func (q *GroupCountQuery[T, K]) BuildQuery() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || q.groupColumn == "" {
		return emptyQueryValues()
	}
	query := "SELECT %s, COUNT(*) FROM %s WHERE %s GROUP BY %s"
	query = fmt.Sprintf(query, q.groupColumn, q.table, condition, q.groupColumn)
	return query, values
}

// BuildQuery returns the query string and parameter values of GroupSumQuery
func (q *GroupSumQuery[T, K, V]) BuildQuery() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || q.groupColumn == "" || q.sumColumn == "" {
		return emptyQueryValues()
	}
	query := "SELECT %s, SUM(%s) FROM %s WHERE %s GROUP BY %s"
	query = fmt.Sprintf(query, q.groupColumn, q.sumColumn, q.table, condition, q.groupColumn)
	return query, values
}

// Query executes the DistinctValuesQuery and returns the list of distinct values
func (q *DistinctValuesQuery[T, V]) Query(this *Instance, dbc db.Conn) ds.Result[[]V] {
	query, values, err := preReadCheck(q, dbc, q.reader)
	if err != nil {
		return ds.Error[[]V](err)
	}

	distinct := make([]V, 0)
	err = readRows(dbc, query, values, q.reader, func(item *T) {
		result := getStructTypedColumnValue[V](this, item, q.typeName, q.columnName)
		if result.IsError() {
			return
		}
		distinct = append(distinct, result.Value())
	})
	if err != nil {
		return ds.Error[[]V](err)
	}

	return ds.NewResult(distinct, nil)
}

// Common: read rows after executing the Query
func readRows[T any](dbc db.Conn, query string, values []any, reader RowReader[T], task func(*T)) error {
	rows, err := dbc.Query(query, values...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		result := reader(rows)
		if result.IsError() {
			continue
		}
		item := new(result.Value())
		task(item)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

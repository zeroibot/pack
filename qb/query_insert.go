package qb

import (
	"fmt"
	"strings"

	"github.com/roidaradal/pack/dict"
	"github.com/roidaradal/pack/list"
	"github.com/roidaradal/pack/str"
)

// InsertRowQuery is used to insert one row to the table
type InsertRowQuery struct {
	baseQuery
	row dict.Object
}

// InsertRowsQuery is used to insert multiple rows to the table
type InsertRowsQuery struct {
	baseQuery
	rows []dict.Object
}

// NewInsertRowQuery creates a new InsertRowQuery
func NewInsertRowQuery(this *Instance, table string) *InsertRowQuery {
	q := new(InsertRowQuery)
	q.baseQuery.initialize(this, table)
	q.row = make(dict.Object)
	return q
}

// NewInsertRowsQuery creates a new InsertRowsQuery
func NewInsertRowsQuery(this *Instance, table string) *InsertRowsQuery {
	q := new(InsertRowsQuery)
	q.baseQuery.initialize(this, table)
	q.rows = make([]dict.Object, 0)
	return q
}

// Row sets the InsertRowQuery's row
func (q *InsertRowQuery) Row(this *Instance, row dict.Object) {
	q.row = prepareRow(this, row)
}

// Rows sets the InsertRowsQuery's rows
func (q *InsertRowsQuery) Rows(this *Instance, rows ...dict.Object) {
	q.rows = list.Map(rows, func(row dict.Object) dict.Object {
		return prepareRow(this, row)
	})
}

// BuildQuery returns the query string and parameter values of InsertRowQuery
func (q *InsertRowQuery) BuildQuery() (string, []any) {
	numColumns := len(q.row)
	err := q.baseQuery.preBuildCheck()
	if err != nil || numColumns == 0 {
		return emptyQueryValues()
	}
	columnList, values := dict.SortedUnzip(q.row)
	columns := strings.Join(columnList, ", ")
	placeholders := str.Repeat(numColumns, "?", ", ")
	query := "INSERT INTO %s (%s) VALUES (%s)"
	query = fmt.Sprintf(query, q.table, columns, placeholders)
	return query, values
}

// BuildQuery returns the query string and parameter values of InsertRowsQuery
func (q *InsertRowsQuery) BuildQuery() (string, []any) {
	numRows := len(q.rows)
	err := q.baseQuery.preBuildCheck()
	if err != nil || numRows == 0 {
		return emptyQueryValues()
	}
	// Set the fixed order and number of columns based on first row
	row1 := q.rows[0]
	fixedSignature := columnSignature(row1)
	numColumns := len(row1)
	if numColumns == 0 {
		return emptyQueryValues()
	}
	values := make([]any, 0, numRows*numColumns)
	columnOrder, values1 := dict.Unzip(row1)
	values = append(values, values1...)
	for _, row := range q.rows[1:] {
		// Ensure same column signature as first row
		if columnSignature(row) != fixedSignature {
			return emptyQueryValues()
		}
		// Follow row1's column order
		for _, column := range columnOrder {
			values = append(values, row[column])
		}
	}
	columns := strings.Join(columnOrder, ", ")
	placeholder := fmt.Sprintf("(%s)", str.Repeat(numColumns, "?", ", "))
	placeholders := str.Repeat(numRows, placeholder, ", ")
	query := "INSERT INTO %s (%s) VALUES %s"
	query = fmt.Sprintf(query, q.table, columns, placeholders)
	return query, values
}

// Common: prepares the column keys of the row Object
func prepareRow(this *Instance, row dict.Object) dict.Object {
	row2 := make(dict.Object, len(row))
	for column, value := range row {
		if column == "" {
			continue
		}
		column = this.dbType.prepareIdentifier(column)
		row2[column] = value
	}
	return row2
}

// Internal: join the sorted column names to check for row column signature
func columnSignature(row dict.Object) string {
	return strings.Join(dict.SortedKeys(row), "/")
}

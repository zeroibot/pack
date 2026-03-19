package qb

import (
	"fmt"
	"strings"

	"github.com/roidaradal/pack/dict"
	"github.com/roidaradal/pack/str"
)

type InsertRowQuery struct {
	baseQuery
	row dict.Object
}

// NewInsertRowQuery creates a new InsertRowQuery
func NewInsertRowQuery(this *Instance, table string) *InsertRowQuery {
	q := new(InsertRowQuery)
	q.baseQuery.initialize(this, table)
	q.row = make(dict.Object)
	return q
}

// Row sets the InsertRowQuery's row
func (q *InsertRowQuery) Row(this *Instance, row dict.Object) {
	row2 := make(dict.Object, len(row))
	for column, value := range row {
		column = this.dbType.prepareIdentifier(column)
		row2[column] = value
	}
	q.row = row2
}

// BuildQuery returns the query string and parameter values of InsertRowQuery
func (q *InsertRowQuery) BuildQuery() (string, []any) {
	numColumns := len(q.row)
	err := q.baseQuery.preBuildCheck()
	if err != nil || numColumns == 0 {
		return emptyQueryValues()
	}
	columnList, values := dict.Unzip(q.row)
	columns := strings.Join(columnList, ", ")
	placeholders := str.Repeat(numColumns, "?", ", ")
	query := "INSERT INTO %s (%s) VALUES (%s)"
	query = fmt.Sprintf(query, q.table, columns, placeholders)
	return query, values
}

package qb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeroibot/pack/db"
	"github.com/zeroibot/pack/dict"
	"github.com/zeroibot/pack/ds"
	"github.com/zeroibot/pack/dyn"
	"github.com/zeroibot/pack/list"
	"github.com/zeroibot/pack/str"
)

// DeleteQuery is used to delete rows from the table
type DeleteQuery[T any] struct {
	conditionQuery[T]
	orderedLimit
}

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

// UpdateQuery is used to update rows in the table
type UpdateQuery[T any] struct {
	conditionQuery[T]
	orderedLimit
	typeName string
	updates  []ds.Option[columnValuePair]
}

type FieldUpdate [2]any                  // [OldValue, NewValue]
type FieldUpdates map[string]FieldUpdate // {FieldName => [OldValue, NewValue]}

// Unpack returns the oldValue and newValue of the FieldUpdate
func (f FieldUpdate) Unpack() (any, any) {
	return f[0], f[1]
}

// NewDeleteQuery creates a new DeleteQuery
func NewDeleteQuery[T any](this *Instance, table string) *DeleteQuery[T] {
	q := new(DeleteQuery[T])
	q.initializeRequired(this, table)
	return q
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

// NewUpdateQuery creates a new UpdateQuery
func NewUpdateQuery[T any](this *Instance, table string) *UpdateQuery[T] {
	q := new(UpdateQuery[T])
	q.initializeRequired(this, table)
	var item T
	q.typeName = dyn.TypeName(item)
	q.updates = make([]ds.Option[columnValuePair], 0)
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

// Update adds a column=value update to the UpdateQuery
func Update[T, V any](this *Instance, q *UpdateQuery[T], fieldRef *V, value V) {
	// Note: cannot be a method because generics are not yet supported in methods
	pairOption := newColumnValue(this, fieldRef, value)
	q.updates = append(q.updates, pairOption)
}

// Update adds a column=value update to the UpdateQuery.
// Note: We lose type-checking of value here, so must be sure that field=value are of the same type.
func (q *UpdateQuery[T]) Update(this *Instance, fieldName string, value any) {
	pairOption := newFieldColumnValue(this, q.typeName, fieldName, value)
	q.updates = append(q.updates, pairOption)
}

// Updates adds column=value updates to the UpdateQuery
// Note: We lose type-checking of value here, so must be sure that field=value are of the same type.
func (q *UpdateQuery[T]) Updates(this *Instance, updates FieldUpdates) {
	fieldNames := dict.SortedKeys(updates)
	for _, fieldName := range fieldNames {
		_, newValue := updates[fieldName].Unpack()
		q.Update(this, fieldName, newValue)
	}
}

// BuildQuery returns the query string and parameter values of DeleteQuery
func (q *DeleteQuery[T]) BuildQuery() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil {
		return emptyQueryValues()
	}
	query := "DELETE FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.table, condition)
	query = tryAppend(query, q.mustLimitString())
	return query, values
}

// BuildQuery returns the query string and parameter values of InsertRowQuery
func (q *InsertRowQuery) BuildQuery() (string, []any) {
	numColumns := len(q.row)
	err := q.baseQuery.preBuildCheck()
	if err != nil || numColumns == 0 {
		return emptyQueryValues()
	}
	columnKeys, values := dict.SortedUnzip(q.row)
	columns := strings.Join(columnKeys, ", ")
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
	columnOrder, values1 := dict.SortedUnzip(row1)
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

// BuildQuery returns the query string and parameter values of UpdateQuery
func (q *UpdateQuery[T]) BuildQuery() (string, []any) {
	numUpdates := len(q.updates)
	condition, conditionValues, err := q.conditionQuery.preBuildCheck()
	if err != nil || numUpdates == 0 {
		return emptyQueryValues()
	}
	values := make([]any, 0, numUpdates+len(conditionValues))
	updates := make([]string, numUpdates)
	for i, pairOption := range q.updates {
		if pairOption.IsNil() {
			// One column=value pair failed = return empty query
			return emptyQueryValues()
		}
		column, value := pairOption.Value().Unpack()
		if column == "" {
			// Blank column = return empty query
			return emptyQueryValues()
		}
		updates[i] = fmt.Sprintf("%s = ?", column)
		values = append(values, value)
	}
	// Add the condition values to the values list
	values = append(values, conditionValues...)
	update := strings.Join(updates, ", ")
	query := "UPDATE %s SET %s WHERE %s"
	query = fmt.Sprintf(query, q.table, update, condition)
	query = tryAppend(query, q.mustLimitString())
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

// ResultChecker is a function that checks the SQL result if a condition is satisfied
type ResultChecker = func(sql.Result) bool

// AssertNothing does not check the SQL result, it always returns true
func AssertNothing(_ sql.Result) bool {
	return true
}

// AssertRowsAffected checks if the SQL result has the expected number of affected rows
func AssertRowsAffected(expected int) ResultChecker {
	return func(result sql.Result) bool {
		return RowsAffected(result) == expected
	}
}

// RowsAffected gets the number of rows affected by the SQL result, returns 0 on error
func RowsAffected(result sql.Result) int {
	count := 0
	if result != nil {
		rowsAffected, err := result.RowsAffected()
		if err == nil {
			count = int(rowsAffected)
		}
	}
	return count
}

// LastInsertID gets the last inserted ID (uint) from SQL result, returns 0 on error
func LastInsertID(result sql.Result) (uint, bool) {
	var insertID uint = 0
	ok := false
	if result != nil {
		id, err := result.LastInsertId()
		if err == nil {
			insertID = uint(id)
			ok = true
		}
	}
	return insertID, ok
}

// Exec executes the given Query, and returns the SQL Result
func Exec(q Query, dbc db.Conn) ds.Result[sql.Result] {
	query, values, err := preQueryCheck(q, dbc)
	if err != nil {
		return ds.Error[sql.Result](err)
	}
	result, err := dbc.Exec(query, values...)
	if err != nil {
		return ds.Error[sql.Result](err)
	}
	return ds.NewResult(result, nil)
}

// ExecTx executes a given Query as part of a database transaction, and rolls back the transaction if any error occurs
func ExecTx(q Query, tx db.Tx, checker ResultChecker) ds.Result[sql.Result] {
	var err error = nil
	query, values := q.BuildQuery()
	if tx == nil {
		err = errNoTx
	} else if query == "" {
		err = errEmptyQuery
	} else if checker == nil {
		err = errNoChecker
	}
	if err != nil {
		err = Rollback(tx, err)
		return ds.Error[sql.Result](err)
	}

	result, err := tx.Exec(query, values...)
	if err != nil {
		err = Rollback(tx, err)
		return ds.Error[sql.Result](err)
	}

	if ok := checker(result); !ok {
		err = Rollback(tx, errFailedResultCheck)
		return ds.Error[sql.Result](err)
	}

	return ds.NewResult(result, nil)
}

// Rollback rolls back the given database transaction
func Rollback(tx db.Tx, err error) error {
	if tx == nil {
		return err
	}
	err2 := tx.Rollback()
	if err2 != nil {
		// Combine original error and rollback error
		return fmt.Errorf("error: %w, rollback error: %w", err, err2)
	}
	// Return original error if successful rollback
	return err
}

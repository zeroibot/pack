package qb

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/pack/db"
)

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
func Exec(q Query, dbc db.Conn) (sql.Result, error) {
	query, values, err := preQueryCheck(q, dbc)
	if err != nil {
		return nil, err
	}
	return dbc.Exec(query, values...)
}

// ExecTx executes a given Query as part of a database transaction, and rolls back the transaction if any error occurs
func ExecTx(q Query, tx db.Tx, checker ResultChecker) (sql.Result, error) {
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
		return nil, Rollback(tx, err)
	}

	result, err := tx.Exec(query, values...)
	if err != nil {
		return nil, Rollback(tx, err)
	}

	if ok := checker(result); !ok {
		return nil, Rollback(tx, errFailedResultCheck)
	}

	return result, nil
}

// Rollback rolls back the given database transaction
func Rollback(tx db.Tx, err error) error {
	err2 := tx.Rollback()
	if err2 != nil {
		// Combine original error and rollback error
		return fmt.Errorf("error: %w, rollback error: %w", err, err2)
	}
	// Return original error if successful rollback
	return err
}

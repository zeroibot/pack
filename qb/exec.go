package qb

import "database/sql"

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

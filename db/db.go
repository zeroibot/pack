// Package db contains database interfaces and functions
package db

import "database/sql"

// Note: sql.Result is an interface with two methods:
// LastInsertId() (int64, error)
// RowsAffected() (int64, error)

// Conn generalizes an sql.Conn object
type Conn interface {
	QueryRow(query string, args ...any) Row
	Query(query string, args ...any) (Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
	Begin() (Tx, error)
}

// Tx generalizes an sql.Tx object
type Tx interface {
	Exec(query string, args ...any) (sql.Result, error)
	Commit() error
	Rollback() error
}

// Row generalizes an sql.Row object
type Row interface {
	Scan(dest ...any) error
}

// Rows generalizes an sql.Rows object
type Rows interface {
	Scan(dest ...any) error
	Next() bool
	Err() error
	Close() error
}

// RowScanner unifies the Row and Rows interface
type RowScanner interface {
	Scan(dest ...any) error
}

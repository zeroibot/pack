package db

import "database/sql"

// Adapter is an adapter for sql.DB so it follows the Conn interface
type Adapter struct {
	db *sql.DB
}

// NewAdapter creates a new Adapter
func NewAdapter(db *sql.DB) *Adapter {
	return &Adapter{db: db}
}

// QueryRow executes a query and returns a Row object
func (a *Adapter) QueryRow(query string, args ...any) Row {
	return a.db.QueryRow(query, args...)
}

// Query executes a query and returns a Rows object
func (a *Adapter) Query(query string, args ...any) (Rows, error) {
	return a.db.Query(query, args...)
}

// Exec executes a query and returns a Result object
func (a *Adapter) Exec(query string, args ...any) (sql.Result, error) {
	return a.db.Exec(query, args...)
}

// Begin starts a transaction
func (a *Adapter) Begin() (Tx, error) {
	return a.db.Begin()
}

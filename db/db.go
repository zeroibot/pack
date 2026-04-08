// Package db contains database interfaces and functions
package db

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/zeroibot/pack/ds"
	"github.com/zeroibot/tst"
)

// Note: sql.Result is an interface with two methods:
// LastInsertId() (int64, error)
// RowsAffected() (int64, error)

// Conn generalizes an sql.Conn object
type Conn interface {
	Begin() (Tx, error)
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (Rows, error)
	QueryRow(query string, args ...any) Row
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

// ConnParams holds parameters for a MySQL DB connection
type ConnParams struct {
	Type     string
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// NewSQLConnection creates a new MySQL DB connection pool
func NewSQLConnection(p *ConnParams) ds.Result[*sql.DB] {
	if p == nil {
		return ds.Error[*sql.DB](fmt.Errorf("sql conn params are not set"))
	}
	if p.Type == "" || p.Host == "" || p.User == "" || p.Password == "" || p.Database == "" || p.Port == "" {
		return ds.Error[*sql.DB](fmt.Errorf("missing sql conn params"))
	}
	address := fmt.Sprintf("%s:%s", p.Host, p.Port)
	cfg := mysql.Config{
		User:                 p.User,
		Passwd:               p.Password,
		Net:                  "tcp",
		Addr:                 address,
		DBName:               p.Database,
		AllowNativePasswords: true,
	}
	dbc, err := sql.Open(p.Type, cfg.FormatDSN())
	if err != nil {
		return ds.Error[*sql.DB](fmt.Errorf("cannot open db: %w", err))
	}
	err = dbc.Ping()
	if err != nil {
		return ds.Error[*sql.DB](fmt.Errorf("cannot ping db: %w", err))
	}
	return ds.NewResult(dbc, nil)
}

// Adapter is an adapter for sql.DB so it follows the Conn interface
type Adapter struct {
	db *sql.DB
}

// NewAdapter creates a new Adapter
func NewAdapter(db *sql.DB) *Adapter {
	return new(Adapter{db: db})
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

// MockAdapter is an adapter for tst.Conn so it follows the Conn interface
type MockAdapter[T any] struct {
	Conn *tst.Conn[T]
}

// NewMockAdapter creates a new MockAdapter
func NewMockAdapter[T any](conn *tst.Conn[T]) *MockAdapter[T] {
	return new(MockAdapter[T]{Conn: conn})
}

// QueryRow executes a query and returns a Row object
func (a *MockAdapter[T]) QueryRow(query string, args ...any) Row {
	return a.Conn.QueryRow(query, args...)
}

// Query executes a query and returns a Rows object
func (a *MockAdapter[T]) Query(query string, args ...any) (Rows, error) {
	return a.Conn.Query(query, args...)
}

// Exec executes a query and returns a Result object
func (a *MockAdapter[T]) Exec(query string, args ...any) (sql.Result, error) {
	return a.Conn.Exec(query, args...)
}

// Begin starts a transaction
func (a *MockAdapter[T]) Begin() (Tx, error) {
	return a.Conn.Begin()
}

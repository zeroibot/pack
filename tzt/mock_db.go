package tzt

import (
	"database/sql"
)

type Conn[T any] struct {
	items  []T
	err    error
	testFn func(T) bool
	rowFn  func([]T) ([]any, error)
}

func NewConn[T any](items ...T) *Conn[T] {
	return new(Conn[T]{items: items, err: nil})
}

func (c *Conn[T]) Begin() (*Tx, error) {
	// TODO: Update for Tx
	return nil, c.err
}

func (c *Conn[T]) Exec(query string, args ...any) (sql.Result, error) {
	// TODO: Implement
	return nil, c.err
}

func (c *Conn[T]) Query(query string, args ...any) (*Rows, error) {
	// TODO: Implement
	return nil, c.err
}

func (c *Conn[T]) QueryRow(query string, args ...any) *Row {
	if c.testFn == nil || c.rowFn == nil {
		return NewRow()
	}
	validItems := make([]T, 0, len(c.items))
	for _, item := range c.items {
		if c.testFn(item) {
			validItems = append(validItems, item)
		}
	}
	items, err := c.rowFn(validItems)
	if err != nil {
		return NewRow()
	}
	return NewRow(items...)
}

func (c *Conn[T]) SetError(err error) {
	c.err = err
}

func (c *Conn[T]) PrepareRow(testFn func(T) bool, rowFn func([]T) ([]any, error)) {
	c.testFn = testFn
	c.rowFn = rowFn
}

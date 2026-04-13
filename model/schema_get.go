package model

import (
	"github.com/zeroibot/pack/dict"
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/qb"
)

// Get performs a SelectRowQuery at schema table
func (s *Schema[T]) Get(rq *my.Request, condition qb.Condition) (T, error) {
	return s.getRowAt(rq, condition, s.Table)
}

// GetAt performs a SelectRowQuery at given table
func (s *Schema[T]) GetAt(rq *my.Request, condition qb.Condition, table string) (T, error) {
	return s.getRowAt(rq, condition, table)
}

// GetOnly performs a SelectRowQuery at schema table and prunes the item with given field names
func (s *Schema[T]) GetOnly(rq *my.Request, condition qb.Condition, fieldNames ...string) (dict.Object, error) {
	item, err := s.getRowAt(rq, condition, s.Table)
	return pruneRow(item, err, fieldNames...)
}

// GetOnlyAt performs a SelectRowQuery at given table and prunes the item with given field names
func (s *Schema[T]) GetOnlyAt(rq *my.Request, condition qb.Condition, table string, fieldNames ...string) (dict.Object, error) {
	item, err := s.getRowAt(rq, condition, table)
	return pruneRow(item, err, fieldNames...)
}

// GetRows performs a SelectRowsQuery at schema table
func (s *Schema[T]) GetRows(rq *my.Request, condition qb.Condition) ([]T, error) {
	return s.getRowsAt(rq, condition, s.Table)
}

// GetRowsAt performs a SelectRowsQuery at given table
func (s *Schema[T]) GetRowsAt(rq *my.Request, condition qb.Condition, table string) ([]T, error) {
	return s.getRowsAt(rq, condition, table)
}

// GetRowsOnly performs a SelectRowsQuery at schema table and prunes the items with given field names
func (s *Schema[T]) GetRowsOnly(rq *my.Request, condition qb.Condition, fieldNames ...string) ([]dict.Object, error) {
	items, err := s.getRowsAt(rq, condition, s.Table)
	return pruneRows(items, err, fieldNames...)
}

// GetRowsOnlyAt performs a SelectRowsQuery at given table and prunes the items with given field names
func (s *Schema[T]) GetRowsOnlyAt(rq *my.Request, condition qb.Condition, table string, fieldNames ...string) ([]dict.Object, error) {
	items, err := s.getRowsAt(rq, condition, table)
	return pruneRows(items, err, fieldNames...)
}

// GetAllRows performs a SelectRowsQuery without condition at schema table
func (s *Schema[T]) GetAllRows(rq *my.Request) ([]T, error) {
	return s.getRowsAt(rq, nil, s.Table)
}

// GetAllRowsAt performs a SelectRowsQuery without condition at given table
func (s *Schema[T]) GetAllRowsAt(rq *my.Request, table string) ([]T, error) {
	return s.getRowsAt(rq, nil, table)
}

// GetAllRowsOnly performs a SelectRowsQuery without condition at schema table and prunes the items with given field names
func (s *Schema[T]) GetAllRowsOnly(rq *my.Request, fieldNames ...string) ([]dict.Object, error) {
	items, err := s.getRowsAt(rq, nil, s.Table)
	return pruneRows(items, err, fieldNames...)
}

// GetAllRowsOnlyAt performs a SelectRowsQuery without condition at given table and prunes the items with given field names
func (s *Schema[T]) GetAllRowsOnlyAt(rq *my.Request, table string, fieldNames ...string) ([]dict.Object, error) {
	items, err := s.getRowsAt(rq, nil, table)
	return pruneRows(items, err, fieldNames...)
}

// Common: create and execute SelectRowQuery at given table
func (s *Schema[T]) getRowAt(rq *my.Request, condition qb.Condition, table string) (T, error) {
	var zero T
	// Check that condition is set
	if condition == nil {
		rq.Fail(my.Err500, "Get condition is not set")
		return zero, fail.MissingParams
	}

	// Build SelectRowQuery and execute
	q := qb.NewFullSelectRowQuery[T](s.Instance, table, s.Reader)
	q.Where(condition)

	result, err := q.QueryRow(rq.DB)
	if err != nil {
		rq.Status = my.Err500
		return zero, err
	}
	return result, nil
}

// Common: create and execute SelectRowsQuery at given table
func (s *Schema[T]) getRowsAt(rq *my.Request, condition qb.Condition, table string) ([]T, error) {
	// Build SelectRowsQuery and execute
	q := qb.NewFullSelectRowsQuery[T](s.Instance, table, s.Reader)
	if condition != nil {
		q.Where(condition)
	}

	items, err := q.Query(rq.DB)
	if err != nil {
		rq.Status = my.Err500
		return nil, err
	}

	return items, nil
}

// Common: prune item with given field names
func pruneRow[T any](item T, err error, fieldNames ...string) (dict.Object, error) {
	if err != nil {
		return nil, err
	}
	object, err := dict.Pruned(&item, fieldNames...)
	if err != nil {
		return nil, err
	}
	return object, nil
}

// Common: prune items with given field names
func pruneRows[T any](items []T, err error, fieldNames ...string) ([]dict.Object, error) {
	if err != nil {
		return nil, err
	}
	objects := make([]dict.Object, 0)
	for _, row := range items {
		object, err := dict.Pruned(&row, fieldNames...)
		if err == nil {
			objects = append(objects, object)
		}
	}
	return objects, nil
}

package model

import (
	"github.com/zeroibot/pack/dict"
	"github.com/zeroibot/pack/ds"
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/qb"
)

// Get performs a SelectRowQuery at schema table
func (s *Schema[T]) Get(rq *my.Request, condition qb.DualCondition[T]) ds.Result[T] {
	return s.getRowAt(rq, condition, s.Table)
}

// GetAt performs a SelectRowQuery at given table
func (s *Schema[T]) GetAt(rq *my.Request, condition qb.DualCondition[T], table string) ds.Result[T] {
	return s.getRowAt(rq, condition, table)
}

// GetOnly performs a SelectRowQuery at schema table and prunes the item with given field names
func (s *Schema[T]) GetOnly(rq *my.Request, condition qb.DualCondition[T], fieldNames ...string) ds.Result[dict.Object] {
	result := s.getRowAt(rq, condition, s.Table)
	return pruneRow(result, fieldNames...)
}

// GetOnlyAt performs a SelectRowQuery at given table and prunes the item with given field names
func (s *Schema[T]) GetOnlyAt(rq *my.Request, condition qb.DualCondition[T], table string, fieldNames ...string) ds.Result[dict.Object] {
	result := s.getRowAt(rq, condition, table)
	return pruneRow(result, fieldNames...)
}

// GetRows performs a SelectRowsQuery at schema table
func (s *Schema[T]) GetRows(rq *my.Request, condition qb.DualCondition[T]) ds.Result[[]T] {
	return s.getRowsAt(rq, condition, s.Table)
}

// GetRowsAt performs a SelectRowsQuery at given table
func (s *Schema[T]) GetRowsAt(rq *my.Request, condition qb.DualCondition[T], table string) ds.Result[[]T] {
	return s.getRowsAt(rq, condition, table)
}

// GetRowsOnly performs a SelectRowsQuery at schema table and prunes the items with given field names
func (s *Schema[T]) GetRowsOnly(rq *my.Request, condition qb.DualCondition[T], fieldNames ...string) ds.Result[[]dict.Object] {
	result := s.getRowsAt(rq, condition, s.Table)
	return pruneRows(result, fieldNames...)
}

// GetRowsOnlyAt performs a SelectRowsQuery at given table and prunes the items with given field names
func (s *Schema[T]) GetRowsOnlyAt(rq *my.Request, condition qb.DualCondition[T], table string, fieldNames ...string) ds.Result[[]dict.Object] {
	result := s.getRowsAt(rq, condition, table)
	return pruneRows(result, fieldNames...)
}

// GetAllRows performs a SelectRowsQuery without condition at schema table
func (s *Schema[T]) GetAllRows(rq *my.Request) ds.Result[[]T] {
	return s.getRowsAt(rq, nil, s.Table)
}

// GetAllRowsAt performs a SelectRowsQuery without condition at given table
func (s *Schema[T]) GetAllRowsAt(rq *my.Request, table string) ds.Result[[]T] {
	return s.getRowsAt(rq, nil, table)
}

// GetAllRowsOnly performs a SelectRowsQuery without condition at schema table and prunes the items with given field names
func (s *Schema[T]) GetAllRowsOnly(rq *my.Request, fieldNames ...string) ds.Result[[]dict.Object] {
	result := s.getRowsAt(rq, nil, s.Table)
	return pruneRows(result, fieldNames...)
}

// GetAllRowsOnlyAt performs a SelectRowsQuery without condition at given table and prunes the items with given field names
func (s *Schema[T]) GetAllRowsOnlyAt(rq *my.Request, table string, fieldNames ...string) ds.Result[[]dict.Object] {
	result := s.getRowsAt(rq, nil, table)
	return pruneRows(result, fieldNames...)
}

// Common: create and execute SelectRowQuery at given table
func (s *Schema[T]) getRowAt(rq *my.Request, condition qb.DualCondition[T], table string) ds.Result[T] {
	// Check that condition is set
	if condition == nil {
		rq.Fail(my.Err500, "Get condition is not set")
		return ds.Error[T](fail.MissingParams)
	}

	// Build SelectRowQuery and execute
	q := qb.NewFullSelectRowQuery[T](s.Instance, table, s.Reader)
	q.Where(condition)

	result := q.QueryRow(rq.DB)
	if result.IsError() {
		rq.Status = my.Err500
	}
	return result
}

// Common: create and execute SelectRowsQuery at given table
func (s *Schema[T]) getRowsAt(rq *my.Request, condition qb.DualCondition[T], table string) ds.Result[[]T] {
	// Build SelectRowsQuery and execute
	q := qb.NewFullSelectRowsQuery[T](s.Instance, table, s.Reader)
	if condition != nil {
		q.Where(condition)
	}

	result := q.Query(rq.DB)
	if result.IsError() {
		rq.Status = my.Err500
	}
	return result
}

// Common: prune item with given field names
func pruneRow[T any](result ds.Result[T], fieldNames ...string) ds.Result[dict.Object] {
	if result.IsError() {
		return ds.Error[dict.Object](result.Error())
	}
	row := new(result.Value())
	object, err := dict.Pruned(row, fieldNames...)
	if err != nil {
		return ds.Error[dict.Object](err)
	}
	return ds.NewResult(object, nil)
}

// Common: prune items with given field names
func pruneRows[T any](result ds.Result[[]T], fieldNames ...string) ds.Result[[]dict.Object] {
	if result.IsError() {
		return ds.Error[[]dict.Object](result.Error())
	}
	objects := make([]dict.Object, 0)
	for _, row := range result.Value() {
		object, err := dict.Pruned(&row, fieldNames...)
		if err == nil {
			objects = append(objects, object)
		}
	}
	return ds.NewResult(objects, nil)
}

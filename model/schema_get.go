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

// Common: create and execute SelectRowQuery at given table
func (s *Schema[T]) getRowAt(rq *my.Request, condition qb.DualCondition[T], table string) ds.Result[T] {
	// Check that condition is set
	if condition == nil {
		rq.Fail(my.Err500, "Get condition is not set")
		return ds.Error[T](fail.MissingParams)
	}

	// Build SelectRowQuery and execute
	q := qb.NewFullSelectRowQuery[T](s.instance, table, s.Reader)
	q.Where(condition)

	result := q.QueryRow(rq.DB)
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

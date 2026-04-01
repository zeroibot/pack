package model

import (
	"github.com/zeroibot/pack/ds"
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/qb"
)

// Count performs a CountQuery at schema table
func (s *Schema[T]) Count(rq *my.Request, condition qb.DualCondition[T]) ds.Result[int] {
	return s.countAt(rq, condition, s.Table)
}

// CountAt performs a CountQuery at the given table
func (s *Schema[T]) CountAt(rq *my.Request, condition qb.DualCondition[T], table string) ds.Result[int] {
	return s.countAt(rq, condition, table)
}

// Sum performs a SumQuery at schema table
func (s *Schema[T]) Sum(rq *my.Request, columns []string, reader qb.RowReader[T], condition qb.DualCondition[T]) ds.Result[T] {
	return s.sumAt(rq, columns, reader, condition, s.Table)
}

// SumAt performs a SumQuery at the given table
func (s *Schema[T]) SumAt(rq *my.Request, columns []string, reader qb.RowReader[T], condition qb.DualCondition[T], table string) ds.Result[T] {
	return s.sumAt(rq, columns, reader, condition, table)
}

// Common: create and execute CountQuery at given table
func (s *Schema[T]) countAt(rq *my.Request, condition qb.DualCondition[T], table string) ds.Result[int] {
	// Check that condition is set
	if condition == nil {
		rq.Fail(my.Err500, "Count condition is not set")
		return ds.Error[int](fail.MissingParams)
	}

	// Build CountQuery and execute
	q := qb.NewCountQuery[T](s.instance, table)
	q.Where(condition)

	result := q.Count(rq.DB)
	if result.IsError() {
		rq.Status = my.Err500
	}
	return result
}

// Common: create and execute SumQuery at given table
func (s *Schema[T]) sumAt(rq *my.Request, columns []string, reader qb.RowReader[T], condition qb.DualCondition[T], table string) ds.Result[T] {
	// Check that condition is set
	if condition == nil {
		rq.Fail(my.Err500, "Sum condition is not set")
		return ds.Error[T](fail.MissingParams)
	}

	// Build SumQuery and execute
	this := s.instance
	q := qb.NewSumQuery[T](this, table, reader)
	q.Columns(this, columns...)
	q.Where(condition)

	result := q.Sum(rq.DB)
	if result.IsError() {
		rq.Status = my.Err500
	}
	return result
}

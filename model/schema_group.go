package model

import (
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/qb"
)

// Count performs a CountQuery at schema table
func (s *Schema[T]) Count(rq *my.Request, condition qb.Condition) (int, error) {
	return s.countAt(rq, condition, s.Table)
}

// CountAt performs a CountQuery at the given table
func (s *Schema[T]) CountAt(rq *my.Request, condition qb.Condition, table string) (int, error) {
	return s.countAt(rq, condition, table)
}

// Sum performs a SumQuery at schema table
func (s *Schema[T]) Sum(rq *my.Request, columns []string, reader qb.RowReader[T], condition qb.Condition) (T, error) {
	return s.sumAt(rq, columns, reader, condition, s.Table)
}

// SumAt performs a SumQuery at the given table
func (s *Schema[T]) SumAt(rq *my.Request, columns []string, reader qb.RowReader[T], condition qb.Condition, table string) (T, error) {
	return s.sumAt(rq, columns, reader, condition, table)
}

// Common: create and execute CountQuery at given table
func (s *Schema[T]) countAt(rq *my.Request, condition qb.Condition, table string) (int, error) {
	// Check that condition is set
	if condition == nil {
		rq.Fail(my.Err500, "Count condition is not set")
		return 0, fail.MissingParams
	}

	// Build CountQuery and execute
	q := qb.NewCountQuery[T](s.Instance, table)
	q.Where(condition)

	count, err := q.Count(rq.DB)
	if err != nil {
		rq.Status = my.Err500
		return 0, err
	}
	return count, nil
}

// Common: create and execute SumQuery at given table
func (s *Schema[T]) sumAt(rq *my.Request, columns []string, reader qb.RowReader[T], condition qb.Condition, table string) (T, error) {
	var zero T
	// Check that condition is set
	if condition == nil {
		rq.Fail(my.Err500, "Sum condition is not set")
		return zero, fail.MissingParams
	}

	// Build SumQuery and execute
	this := s.Instance
	q := qb.NewSumQuery[T](this, table, reader)
	q.Columns(this, columns...)
	q.Where(condition)

	result, err := q.Sum(rq.DB)
	if err != nil {
		rq.Status = my.Err500
		return zero, err
	}
	return result, nil
}

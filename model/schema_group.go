package model

import (
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/qb"
)

// Count performs a CountQuery at schema table
func (s *Schema[T]) Count(rq *my.Request, condition qb.DualCondition[T]) (int, error) {
	return s.countAt(rq, condition, s.Table)
}

// CountAt performs a CountQuery at the given table
func (s *Schema[T]) CountAt(rq *my.Request, condition qb.DualCondition[T], table string) (int, error) {
	return s.countAt(rq, condition, table)
}

// Common: create and execute CountQuery at given table
func (s *Schema[T]) countAt(rq *my.Request, condition qb.DualCondition[T], table string) (int, error) {
	// Check that condition is set
	if condition == nil {
		rq.Fail(my.Err500, "Count condition is not set")
		return 0, fail.MissingParams
	}

	// Build CountQuery and execute
	q := qb.NewCountQuery[T](s.instance, table)
	q.Where(condition)
	result := q.Count(rq.DB)
	if result.IsError() {
		rq.Fail(my.Err500, "Failed to count %s", s.Name)
		return 0, result.Error()
	}
	return result.Value(), nil
}

package model

import (
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/qb"
)

// GetAscRows performs a SelectRowsQuery with Ascending order at schema table
func (s *Schema[T]) GetAscRows(rq *my.Request, condition qb.Condition, orderColumn string) ([]T, error) {
	return s.topRowsAt(rq, condition, orderColumn, s.Table, true, 0)
}

// GetAscRowsAt performs a SelectRowsQuery with Ascending order at given table
func (s *Schema[T]) GetAscRowsAt(rq *my.Request, condition qb.Condition, orderColumn string, table string) ([]T, error) {
	return s.topRowsAt(rq, condition, orderColumn, table, true, 0)
}

// GetDescRows performs a SelectRowsQuery with Descending order at schema table
func (s *Schema[T]) GetDescRows(rq *my.Request, condition qb.Condition, orderColumn string) ([]T, error) {
	return s.topRowsAt(rq, condition, orderColumn, s.Table, false, 0)
}

// GetDescRowsAt performs a SelectRowsQuery with Descending order at given table
func (s *Schema[T]) GetDescRowsAt(rq *my.Request, condition qb.Condition, orderColumn string, table string) ([]T, error) {
	return s.topRowsAt(rq, condition, orderColumn, table, false, 0)
}

// Common: create and execute SelectRowsQuery, with order and limit set, at given table
func (s *Schema[T]) topRowsAt(rq *my.Request, condition qb.Condition, orderColumn string, table string, isAscending bool, limit int) ([]T, error) {
	// Build SelectRowsQuery and execute
	this := s.Instance
	q := qb.NewFullSelectRowsQuery[T](this, table, s.Reader)
	if condition != nil {
		q.Where(condition)
	}
	if isAscending {
		q.OrderAsc(this, orderColumn)
	} else {
		q.OrderDesc(this, orderColumn)
	}

	items, err := q.Query(rq.DB)
	if err != nil {
		rq.Status = my.Err500
		return nil, err
	}

	return items, nil
}

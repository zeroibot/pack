package model

import (
	"database/sql"

	"github.com/zeroibot/pack/ds"
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/qb"
)

// UpdateFields performs an UpdateQuery using FieldUpdates at schema table
func (s *Schema[T]) UpdateFields(rq *my.Request, updates qb.FieldUpdates, condition qb.DualCondition[T]) error {
	return s.updateFieldsAt(rq, updates, condition, s.Table, false)
}

// UpdateFieldsAt performs an UpdateQuery using FieldUpdates at given table
func (s *Schema[T]) UpdateFieldsAt(rq *my.Request, updates qb.FieldUpdates, condition qb.DualCondition[T], table string) error {
	return s.updateFieldsAt(rq, updates, condition, table, false)
}

// UpdateTxFields performs an UpdateQuery as part of a transaction using FieldUpdates at schema table
func (s *Schema[T]) UpdateTxFields(rqtx *my.Request, updates qb.FieldUpdates, condition qb.DualCondition[T]) error {
	return s.updateFieldsAt(rqtx, updates, condition, s.Table, true)
}

// UpdateTxFieldsAt performs an UpdateQuery as part of a transaction using FieldUpdates at given table
func (s *Schema[T]) UpdateTxFieldsAt(rqtx *my.Request, updates qb.FieldUpdates, condition qb.DualCondition[T], table string) error {
	return s.updateFieldsAt(rqtx, updates, condition, table, true)
}

// Common: create and execute UpdateQuery using FieldUpdates at given table
func (s *Schema[T]) updateFieldsAt(rq *my.Request, updates qb.FieldUpdates, condition qb.DualCondition[T], table string, isTx bool) error {
	// Check that condition and updates are set
	if condition == nil || len(updates) == 0 {
		rq.Fail(my.Err500, "Update / condition is not set")
		return fail.MissingParams
	}

	// Build UpdateQuery
	this := s.instance
	q := qb.NewUpdateQuery[T](this, table)
	q.Where(condition)
	q.Updates(this, updates)

	// Execute UpdateQuery
	var result ds.Result[sql.Result]
	if isTx {
		rq.AddTxStep(q)
		result = qb.ExecTx(q, rq.Tx, rq.Checker)
	} else {
		result = qb.Exec(q, rq.DB)
	}
	if result.IsError() {
		rq.Fail(my.Err500, "Failed to update %s", s.Name)
		return result.Error()
	}

	rowsUpdated := qb.RowsAffected(result.Value())
	if rowsUpdated != 1 {
		rq.AddFmtLog("Updated: %d %s", rowsUpdated, s.Name)
	}
	return nil
}

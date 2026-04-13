package model

import (
	"database/sql"

	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/qb"
)

// UpdateFn is a function that decorates the UpdateQuery with necessary updates
type UpdateFn[T any] = func(*qb.Instance, *qb.UpdateQuery[T])

// Update performs an UpdateQuery using UpdateFn at schema table
func (s *Schema[T]) Update(rq *my.Request, updateFn UpdateFn[T], condition qb.Condition) error {
	return s.updateAt(rq, updateFn, condition, s.Table, false)
}

// UpdateAt performs an UpdateQuery using UpdateFn at given table
func (s *Schema[T]) UpdateAt(rq *my.Request, updateFn UpdateFn[T], condition qb.Condition, table string) error {
	return s.updateAt(rq, updateFn, condition, table, false)
}

// UpdateTx performs an UpdateQuery using UpdateFn as part of a transaction at schema table
func (s *Schema[T]) UpdateTx(rqtx *my.Request, updateFn UpdateFn[T], condition qb.Condition) error {
	return s.updateAt(rqtx, updateFn, condition, s.Table, true)
}

// UpdateTxAt performs an UpdateQuery using UpdateFn as part of a transaction at given table
func (s *Schema[T]) UpdateTxAt(rqtx *my.Request, updateFn UpdateFn[T], condition qb.Condition, table string) error {
	return s.updateAt(rqtx, updateFn, condition, table, true)
}

// UpdateFields performs an UpdateQuery using FieldUpdates at schema table
func (s *Schema[T]) UpdateFields(rq *my.Request, updates qb.FieldUpdates, condition qb.Condition) error {
	return s.updateFieldsAt(rq, updates, condition, s.Table, false)
}

// UpdateFieldsAt performs an UpdateQuery using FieldUpdates at given table
func (s *Schema[T]) UpdateFieldsAt(rq *my.Request, updates qb.FieldUpdates, condition qb.Condition, table string) error {
	return s.updateFieldsAt(rq, updates, condition, table, false)
}

// UpdateTxFields performs an UpdateQuery as part of a transaction using FieldUpdates at schema table
func (s *Schema[T]) UpdateTxFields(rqtx *my.Request, updates qb.FieldUpdates, condition qb.Condition) error {
	return s.updateFieldsAt(rqtx, updates, condition, s.Table, true)
}

// UpdateTxFieldsAt performs an UpdateQuery as part of a transaction using FieldUpdates at given table
func (s *Schema[T]) UpdateTxFieldsAt(rqtx *my.Request, updates qb.FieldUpdates, condition qb.Condition, table string) error {
	return s.updateFieldsAt(rqtx, updates, condition, table, true)
}

// Common: create and execute UpdateQuery using UpdateFn at given table
func (s *Schema[T]) updateAt(rq *my.Request, updateFn UpdateFn[T], condition qb.Condition, table string, isTx bool) error {
	// Check that condition and updateFn are set
	if condition == nil || updateFn == nil {
		rq.Fail(my.Err500, "Update / condition is not set")
		return fail.MissingParams
	}

	// Build UpdateQuery
	this := s.Instance
	q := qb.NewUpdateQuery[T](this, table)
	q.Where(condition)
	updateFn(this, q) // Call updateFn to add updates

	return s.update(rq, q, isTx)
}

// Common: create and execute UpdateQuery using FieldUpdates at given table
func (s *Schema[T]) updateFieldsAt(rq *my.Request, updates qb.FieldUpdates, condition qb.Condition, table string, isTx bool) error {
	// Check that condition and updates are set
	if condition == nil || len(updates) == 0 {
		rq.Fail(my.Err500, "Update / condition is not set")
		return fail.MissingParams
	}

	// Build UpdateQuery
	this := s.Instance
	q := qb.NewUpdateQuery[T](this, table)
	q.Where(condition)
	q.Updates(this, updates)

	return s.update(rq, q, isTx)
}

// Common: execute UpdateQuery
func (s *Schema[T]) update(rq *my.Request, q *qb.UpdateQuery[T], isTx bool) error {
	// Execute UpdateQuery
	var result sql.Result
	var err error
	if isTx {
		rq.AddTxStep(q)
		result, err = qb.ExecTx(q, rq.Tx, rq.Checker)
	} else {
		result, err = qb.Exec(q, rq.DB)
	}
	if err != nil {
		rq.Fail(my.Err500, "Failed to update %s", s.Name)
		return err
	}

	rowsUpdated := qb.RowsAffected(result)
	if rowsUpdated != 1 {
		rq.AddFmtLog("Updated: %d %s", rowsUpdated, s.Name)
	}
	return nil
}

package model

import (
	"database/sql"

	"github.com/zeroibot/pack/ds"
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/qb"
)

// SetFlag updates the booleanField = flag at schema table
func (s *Schema[T]) SetFlag(rq *my.Request, condition qb.DualCondition[T], field *bool, flag bool) error {
	return s.setFlagsAt(rq, condition, field, flag, 1, s.Table, false)
}

// SetFlagAt updates the booleanField = flag at given table
func (s *Schema[T]) SetFlagAt(rq *my.Request, condition qb.DualCondition[T], field *bool, flag bool, table string) error {
	return s.setFlagsAt(rq, condition, field, flag, 1, table, false)
}

// SetTxFlag updates the booleanField = flag as part of transaction at schema table
func (s *Schema[T]) SetTxFlag(rqtx *my.Request, condition qb.DualCondition[T], field *bool, flag bool) error {
	return s.setFlagsAt(rqtx, condition, field, flag, 1, s.Table, true)
}

// SetTxFlagAt updates the booleanField = flag as part of transaction at given table
func (s *Schema[T]) SetTxFlagAt(rqtx *my.Request, condition qb.DualCondition[T], field *bool, flag bool, table string) error {
	return s.setFlagsAt(rqtx, condition, field, flag, 1, table, true)
}

// SetFlags updates the booleanField = flag at schema table affecting multiple rows
func (s *Schema[T]) SetFlags(rq *my.Request, condition qb.DualCondition[T], field *bool, flag bool, numItems int) error {
	return s.setFlagsAt(rq, condition, field, flag, numItems, s.Table, false)
}

// SetFlagsAt updates the booleanField = flag at given table affecting multiple rows
func (s *Schema[T]) SetFlagsAt(rq *my.Request, condition qb.DualCondition[T], field *bool, flag bool, numItems int, table string) error {
	return s.setFlagsAt(rq, condition, field, flag, numItems, table, false)
}

// SetTxFlags updates the booleanField = flag as part of transaction at schema table affecting multiple rows
func (s *Schema[T]) SetTxFlags(rqtx *my.Request, condition qb.DualCondition[T], field *bool, flag bool, numItems int) error {
	return s.setFlagsAt(rqtx, condition, field, flag, numItems, s.Table, true)
}

// SetTxFlagsAt updates the booleanField = flag as part of transaction at given table affecting multiple rows
func (s *Schema[T]) SetTxFlagsAt(rqtx *my.Request, condition qb.DualCondition[T], field *bool, flag bool, numItems int, table string) error {
	return s.setFlagsAt(rqtx, condition, field, flag, numItems, table, true)
}

// Common: create and execute UpdateQuery, which sets booleanField = flag at given table
func (s *Schema[T]) setFlagsAt(rq *my.Request, condition qb.DualCondition[T], field *bool, flag bool, numItems int, table string, isTx bool) error {
	// Check that condition is set
	if condition == nil {
		rq.Fail(my.Err500, "Update condition is not set")
		return fail.MissingParams
	}

	// Build UpdateQuery
	this := s.Instance
	q := qb.NewUpdateQuery[T](this, table)
	q.Where(condition)
	if numItems == 1 {
		q.Limit(1)
	}
	qb.Update(this, q, field, flag)

	// Execute UpdateQuery
	var result ds.Result[sql.Result]
	if isTx {
		rq.AddTxStep(q)
		checker := qb.AssertRowsAffected(numItems)
		result = qb.ExecTx(q, rq.Tx, checker)
	} else {
		result = qb.Exec(q, rq.DB)
	}
	if result.IsError() {
		rq.Fail(my.Err500, "Failed to update %s flag", s.Name)
		return result.Error()
	}
	rowsUpdated := qb.RowsAffected(result.Value())

	// If not transaction, check if rowsUpdated == numItems
	if !isTx && rowsUpdated != numItems {
		rq.Fail(my.Err500, "Update count mismatch: items = %d, updated = %d", numItems, rowsUpdated)
		return fail.MismatchCount
	}

	if rowsUpdated != 1 {
		rq.AddFmtLog("Updated: %d %s", rowsUpdated, s.Name)
	}
	return nil
}

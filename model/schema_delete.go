package model

import (
	"database/sql"

	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/qb"
)

// Delete performs a DeleteQuery on the schema table using the given condition
func (s *Schema[T]) Delete(rq *my.Request, condition qb.DualCondition[T]) (int, error) {
	return s.deleteAt(rq, condition, s.Table, false)
}

// DeleteAt performs a DeleteQuery on the given table using the given condition
func (s *Schema[T]) DeleteAt(rq *my.Request, condition qb.DualCondition[T], table string) (int, error) {
	return s.deleteAt(rq, condition, table, false)
}

// DeleteTx performs a DeleteQuery as part of a transaction on the schema table using the given condition
func (s *Schema[T]) DeleteTx(rqtx *my.Request, condition qb.DualCondition[T]) (int, error) {
	return s.deleteAt(rqtx, condition, s.Table, true)
}

// DeleteTxAt performs a DeleteQuery as part of a transaction on the given table using the given condition
func (s *Schema[T]) DeleteTxAt(rqtx *my.Request, condition qb.DualCondition[T], table string) (int, error) {
	return s.deleteAt(rqtx, condition, table, true)
}

// Common: create and execute DeleteQuery at the given table using the given condition
func (s *Schema[T]) deleteAt(rq *my.Request, condition qb.DualCondition[T], table string, isTx bool) (int, error) {
	// Check that condition is set
	if condition == nil {
		rq.Fail(my.Err500, "Delete condition is not set")
		return 0, fail.MissingParams
	}

	// Build DeleteQuery
	q := qb.NewDeleteQuery[T](s.instance, table)
	q.Where(condition)

	// Execute DeleteQuery
	var result sql.Result
	var err error
	if isTx {
		rq.AddTxStep(q)
		result, err = qb.ExecTx(q, rq.Tx, rq.Checker)
	} else {
		result, err = qb.Exec(q, rq.DB)
	}
	if err != nil {
		rq.Fail(my.Err500, "Failed to delete %s", s.Name)
		return 0, err
	}

	rowsDeleted := qb.RowsAffected(result)
	if rowsDeleted != 1 {
		rq.AddFmtLog("Deleted: %d %s", rowsDeleted, s.Name)
	}
	return rowsDeleted, nil
}

package model

import (
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/qb"
)

// SetFlag updates the booleanField = flag at schema table
func (s *Schema[T]) SetFlag(rq *my.Request, condition qb.DualCondition[T], field *bool, flag bool) error {
	return s.setFlag(rq, condition, field, flag, s.Table, false)
}

// SetFlagAt updates the booleanField = flag at given table
func (s *Schema[T]) SetFlagAt(rq *my.Request, condition qb.DualCondition[T], field *bool, flag bool, table string) error {
	return s.setFlag(rq, condition, field, flag, table, false)
}

// SetTxFlag updates the booleanField = flag as part of transaction at schema table
func (s *Schema[T]) SetTxFlag(rqtx *my.Request, condition qb.DualCondition[T], field *bool, flag bool) error {
	return s.setFlag(rqtx, condition, field, flag, s.Table, true)
}

// SetTxFlagAt updates the booleanField = flag as part of transaction at given table
func (s *Schema[T]) SetTxFlagAt(rqtx *my.Request, condition qb.DualCondition[T], field *bool, flag bool, table string) error {
	return s.setFlag(rqtx, condition, field, flag, table, true)
}

// Common: create and execute UpdateQuery, which sets booleanField = flag at given table
func (s *Schema[T]) setFlag(rq *my.Request, condition qb.DualCondition[T], field *bool, flag bool, table string, isTx bool) error {
	// Check that condition is set
	if condition == nil {
		rq.Fail(my.Err500, "Update condition is not set")
		return fail.MissingParams
	}

	// Build UpdateQuery
	q := qb.NewUpdateQuery[T](s.instance, table)
	q.Where(condition)
	q.Limit(1)
	qb.Update(s.instance, q, field, flag)

	// Execute UpdateQuery
	var err error
	if isTx {
		rq.AddTxStep(q)
		_, err = qb.ExecTx(q, rq.Tx, rq.Checker)
	} else {
		_, err = qb.Exec(q, rq.DB)
	}
	if err != nil {
		rq.Fail(my.Err500, "Failed to update %s flag", s.Name)
		return err
	}

	return nil
}

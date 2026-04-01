package model

import (
	"database/sql"
	"errors"

	"github.com/zeroibot/pack/ds"
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/qb"
)

// Insert performs an InsertRowQuery at schema table
func (s *Schema[T]) Insert(rq *my.Request, item *T) ds.Result[ID] {
	return s.insertAt(rq, item, s.Table, false)
}

// InsertAt performs an InsertRowQuery at given table
func (s *Schema[T]) InsertAt(rq *my.Request, item *T, table string) ds.Result[ID] {
	return s.insertAt(rq, item, table, false)
}

// InsertTx performs an InsertRowQuery as part of a transaction at schema table
func (s *Schema[T]) InsertTx(rqtx *my.Request, item *T) ds.Result[ID] {
	return s.insertAt(rqtx, item, s.Table, true)
}

// InsertTxAt performs an InsertRowQuery as part of a transaction at given table
func (s *Schema[T]) InsertTxAt(rqtx *my.Request, item *T, table string) ds.Result[ID] {
	return s.insertAt(rqtx, item, table, true)
}

// Common: create and execute InsertRowQuery at given table
func (s *Schema[T]) insertAt(rq *my.Request, item *T, table string, isTx bool) ds.Result[ID] {
	// Check that item is not nil
	if item == nil {
		rq.Fail(my.Err500, "Insert item is nil")
		return ds.Error[ID](fail.MissingParams)
	}

	// Build InsertRowQuery
	this := s.instance
	q := qb.NewInsertRowQuery(this, table)
	q.Row(this, qb.ToRow(this, item))

	// Execute InsertRowQuery
	var result ds.Result[sql.Result]
	if isTx {
		rq.AddTxStep(q)
		result = qb.ExecTx(q, rq.Tx, rq.Checker)
	} else {
		result = qb.Exec(q, rq.DB)
	}
	if result.IsError() {
		rq.Fail(my.Err500, "Failed to insert %s", s.Name)
		return ds.Error[ID](result.Error())
	}
	rowsInserted := qb.RowsAffected(result.Value())

	// If not transaction, check if rowsInserted == 1
	if !isTx && rowsInserted != 1 {
		rq.Fail(my.Err500, "No %s rows inserted", s.Name)
		return ds.Error[ID](errors.New("no rows inserted"))
	}

	// Get last inserted ID
	id, ok := qb.LastInsertID(result.Value())
	if !ok {
		rq.Fail(my.Err500, "Failed to get last inserted ID")
		return ds.Error[ID](errors.New("failed to get last inserted ID"))
	}

	rq.Status = my.OK201
	return ds.NewResult[ID](id, nil)
}

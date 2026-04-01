package model

import (
	"database/sql"
	"errors"

	"github.com/zeroibot/pack/dict"
	"github.com/zeroibot/pack/ds"
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/list"
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

// InsertRows performs an InsertRowsQuery at schema table
func (s *Schema[T]) InsertRows(rq *my.Request, items []T) error {
	return s.insertRowsAt(rq, items, s.Table, false)
}

// InsertRowsAt performs an InsertRowsQuery at given table
func (s *Schema[T]) InsertRowsAt(rq *my.Request, items []T, table string) error {
	return s.insertRowsAt(rq, items, table, false)
}

// InsertTxRows performs an InsertRowsQuery as part of a transaction at schema table
func (s *Schema[T]) InsertTxRows(rqtx *my.Request, items []T) error {
	return s.insertRowsAt(rqtx, items, s.Table, true)
}

// InsertTxRowsAt performs an InsertRowsQuery as part of a transaction at given table
func (s *Schema[T]) InsertTxRowsAt(rqtx *my.Request, items []T, table string) error {
	return s.insertRowsAt(rqtx, items, table, true)
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

// Common: create and execute InsertRowsQuery at given table
func (s *Schema[T]) insertRowsAt(rq *my.Request, items []T, table string, isTx bool) error {
	// Check that items are set
	numItems := len(items)
	if numItems == 0 {
		rq.Fail(my.Err500, "Insert items are empty")
		return fail.MissingParams
	}

	// Build InsertRowsQuery
	this := s.instance
	rows := list.Map(items, func(item T) dict.Object {
		return qb.ToRow(this, &item)
	})
	q := qb.NewInsertRowsQuery(this, table)
	q.Rows(this, rows...)

	// Execute InsertRowsQuery
	var result ds.Result[sql.Result]
	if isTx {
		rq.AddTxStep(q)
		checker := qb.AssertRowsAffected(numItems)
		result = qb.ExecTx(q, rq.Tx, checker)
	} else {
		result = qb.Exec(q, rq.DB)
	}
	if result.IsError() {
		rq.Fail(my.Err500, "Failed to insert %d %s rows", numItems, s.Name)
		return result.Error()
	}
	rowsInserted := qb.RowsAffected(result.Value())

	// If not transaction, check if rowsInserted == numItems
	if !isTx && rowsInserted != numItems {
		rq.Fail(my.Err500, "Insert count mismatch: items = %d, inserted = %d", numItems, rowsInserted)
		return errors.New("count mismatch")
	}

	rq.AddFmtLog("Inserted: %d %s", rowsInserted, s.Name)
	rq.Status = my.OK201
	return nil
}

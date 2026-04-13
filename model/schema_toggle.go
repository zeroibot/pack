package model

import (
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/qb"
)

type toggleParams struct {
	isActive bool   // required
	byID     bool   // required
	id       ID     // required for ToggleID
	code     string // required for ToggleCode
}

// ToggleID toggles a row on/off using the given ID at the schema table
func (s *Schema[T]) ToggleID(rq *my.Request, id ID, isActive bool, Items *Schema[Item]) error {
	p := toggleParams{isActive: isActive, byID: true, id: id}
	return s.toggleAt(rq, Items, p, s.Table, false)
}

// ToggleIDAt toggles a row on/off using the given ID at the given table
func (s *Schema[T]) ToggleIDAt(rq *my.Request, id ID, isActive bool, table string, Items *Schema[Item]) error {
	p := toggleParams{isActive: isActive, byID: true, id: id}
	return s.toggleAt(rq, Items, p, table, false)
}

// ToggleCode toggles a row on/off using the given Code at the schema table
func (s *Schema[T]) ToggleCode(rq *my.Request, code string, isActive bool, Items *Schema[Item]) error {
	p := toggleParams{isActive: isActive, byID: false, code: code}
	return s.toggleAt(rq, Items, p, s.Table, false)
}

// ToggleCodeAt toggles a row on/off using the given Code at the given table
func (s *Schema[T]) ToggleCodeAt(rq *my.Request, code string, isActive bool, table string, Items *Schema[Item]) error {
	p := toggleParams{isActive: isActive, byID: false, code: code}
	return s.toggleAt(rq, Items, p, table, false)
}

// ToggleTxID toggles a row on/off as part of a transaction using the given ID at the schema table
func (s *Schema[T]) ToggleTxID(rqtx *my.Request, id ID, isActive bool, Items *Schema[Item]) error {
	p := toggleParams{isActive: isActive, byID: true, id: id}
	return s.toggleAt(rqtx, Items, p, s.Table, true)
}

// ToggleTxIDAt toggles a row on/off as part of a transaction using the given ID at the given table
func (s *Schema[T]) ToggleTxIDAt(rqtx *my.Request, id ID, isActive bool, table string, Items *Schema[Item]) error {
	p := toggleParams{isActive: isActive, byID: true, id: id}
	return s.toggleAt(rqtx, Items, p, table, true)
}

// ToggleTxCode toggles a row on/off as part of a transaction using the given Code at the schema table
func (s *Schema[T]) ToggleTxCode(rqtx *my.Request, code string, isActive bool, Items *Schema[Item]) error {
	p := toggleParams{isActive: isActive, byID: false, code: code}
	return s.toggleAt(rqtx, Items, p, s.Table, true)
}

// ToggleTxCodeAt toggles a row on/off as part of a transaction using the given Code at the given table
func (s *Schema[T]) ToggleTxCodeAt(rqtx *my.Request, code string, isActive bool, table string, Items *Schema[Item]) error {
	p := toggleParams{isActive: isActive, byID: false, code: code}
	return s.toggleAt(rqtx, Items, p, table, true)
}

// Common: create and execute UpdateQuery, which toggles IsActive true/false of item with ID/Code at given table
func (s *Schema[T]) toggleAt(rq *my.Request, Items *Schema[Item], p toggleParams, table string, isTx bool) error {
	// Check that params has ID or Code
	hasIdentity := false
	if p.byID && p.id != 0 {
		hasIdentity = true
	} else if !p.byID && p.code != "" {
		hasIdentity = true
	}
	if !hasIdentity {
		rq.Fail(my.Err500, "ID/Code to toggle is missing")
		return fail.MissingParams
	}
	// Check that Items schema is not nil
	if Items == nil {
		rq.Fail(my.Err500, "Items schema is nil")
		return fail.MissingParams
	}

	// Build UpdateQuery using Items schema
	item := Items.Ref
	this := s.Instance
	q := qb.NewUpdateQuery[T](this, table)
	var condition1 qb.Condition
	if p.byID {
		condition1 = qb.Equal(this, &item.ID, p.id)
	} else {
		condition1 = qb.Equal(this, &item.Code, p.code)
	}
	condition2 := qb.Equal(this, &item.IsActive, !p.isActive)
	q.Where(qb.And(condition1, condition2))
	qb.Update(this, q, &item.IsActive, p.isActive)

	// Execute UpdateQuery
	var err error
	if isTx {
		rq.AddTxStep(q)
		_, err = qb.ExecTx(q, rq.Tx, rq.Checker)
	} else {
		_, err = qb.Exec(q, rq.DB)
	}
	if err != nil {
		rq.Fail(my.Err500, "Failed to toggle %s", s.Name)
		return err
	}

	return nil
}

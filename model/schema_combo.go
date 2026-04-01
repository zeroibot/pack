package model

import (
	"fmt"

	"github.com/zeroibot/pack/ds"
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/qb"
)

type GetOrCreateParams[T any] struct {
	Name          string
	Owner         string
	PreCondition  qb.DualCondition[T]
	PostCondition qb.DualCondition[T]
	NewFn         func() T
}

type GetOrCreateAndLockParams[T any] struct {
	GetOrCreateParams[T]
	LockField       *bool
	LockConditionFn func(T) qb.DualCondition[T]
}

// GetOrCreate gets the item if it exists, otherwise creates and returns it
func (s *Schema[T]) GetOrCreate(rq *my.Request, p *GetOrCreateParams[T]) ds.Result[T] {
	return s.getOrCreate(rq, p, false)
}

// GetOrCreateTx gets the item if it exists, otherwise creates it as part of a transaction, and return the item
func (s *Schema[T]) GetOrCreateTx(rqtx *my.Request, p *GetOrCreateParams[T]) ds.Result[T] {
	return s.getOrCreate(rqtx, p, true)
}

// Common: get the item if it exists, otherwise create and return it
func (s *Schema[T]) getOrCreate(rq *my.Request, p *GetOrCreateParams[T], isTx bool) ds.Result[T] {
	// Fetch item to check if it exists
	rowsResult := s.GetRows(rq, p.PreCondition)
	if rowsResult.IsError() {
		rq.Fail(my.Err500, "Failed to check if %s exists", p.Name)
		err := rowsResult.Error()
		if isTx {
			err = qb.Rollback(rq.Tx, err) // manual rollback
		}
		return ds.Error[T](err)
	}

	rows := rowsResult.Value()
	numRows := len(rows)
	if numRows > 1 {
		rq.Fail(my.Err500, "Failed to get one %s", p.Name)
		err := fmt.Errorf("public: Multiple %s found", p.Name)
		if isTx {
			err = qb.Rollback(rq.Tx, err) // manual rollback
		}
		return ds.Error[T](err)
	} else if numRows == 0 {
		// Not found = create item
		newItem := p.NewFn()
		var result ds.Result[ID]
		if isTx {
			result = s.InsertTx(rq, &newItem)
		} else {
			result = s.Insert(rq, &newItem)
		}
		if result.IsError() {
			return ds.Error[T](result.Error())
		}

		rq.AddFmtLog("Created %s for %s", p.Name, p.Owner)
		return ds.NewResult(newItem, nil)
	}

	// At this point, numRows == 1
	item := rows[0]
	// Check if item passes PostCondition
	if p.PostCondition != nil && p.PostCondition.Test(item) == false {
		rq.Fail(my.Err500, "Failed to get %s", p.Name)
		err := fail.NotFoundItem
		if isTx {
			err = qb.Rollback(rq.Tx, err) // manual rollback
		}
		return ds.Error[T](err)
	}
	return ds.NewResult(item, nil)
}

// GetAndLockTx gets the item and locks it as part of a transaction
// Note: no need to include IsLocked = true/false in conditions, as this function adds it
func (s *Schema[T]) GetAndLockTx(rqtx *my.Request, lockField *bool, selectCondition qb.DualCondition[T], lockConditionFn func(T) qb.DualCondition[T]) ds.Result[T] {
	this := s.instance
	isUnlocked := qb.Equal[T](this, lockField, false)

	// Get unlocked item
	condition := qb.And(selectCondition, isUnlocked)
	result := s.Get(rqtx, condition)
	if result.IsError() {
		rqtx.Fail(my.Err500, "Failed to get unlocked item")
		// Manual rollback on error of Get
		err := qb.Rollback(rqtx.Tx, result.Error())
		return ds.Error[T](err)
	}

	// Lock item
	item := result.Value()
	lockCondition := qb.And(lockConditionFn(item), isUnlocked)
	err := s.SetTxFlag(rqtx, lockCondition, lockField, true)
	if err != nil {
		rqtx.Fail(my.Err500, "Failed to lock item")
		return ds.Error[T](err)
	}

	return ds.NewResult(item, nil)
}

// GetAndLockTxItems gets a list of items and locks all of them as part of a transaction
// Note: no need to include IsLocked = true/false in conditions, as this function adds it
func (s *Schema[T]) GetAndLockTxItems(rqtx *my.Request, lockField *bool, selectCondition qb.DualCondition[T], lockConditionFn func([]T) qb.DualCondition[T], numItems int) ds.Result[[]T] {
	this := s.instance
	isUnlocked := qb.Equal[T](this, lockField, false)

	// Get unlocked items
	condition := qb.And(selectCondition, isUnlocked)
	result := s.GetRows(rqtx, condition)
	if result.IsError() {
		rqtx.Fail(my.Err500, "Failed to get unlocked items")
		// Manual rollback on error of GetRows
		err := qb.Rollback(rqtx.Tx, result.Error())
		return ds.Error[[]T](err)
	}

	// Check that rows have correct number of items
	items := result.Value()
	if len(items) != numItems {
		rqtx.Fail(my.Err500, "Get count mismatch: items = %d, rows = %d", numItems, len(items))
		// Manual rollback if count mismatch
		err := qb.Rollback(rqtx.Tx, fail.MismatchCount)
		return ds.Error[[]T](err)
	}

	// Lock items
	lockCondition := qb.And(lockConditionFn(items), isUnlocked)
	err := s.SetTxFlags(rqtx, lockCondition, lockField, true, numItems)
	if err != nil {
		rqtx.Fail(my.Err500, "Failed to lock items")
		return ds.Error[[]T](err)
	}

	return ds.NewResult(items, nil)
}

// GetOrCreateAndLockTx runs GetOrCreate and locks the item as part of a transaction
// Note: no need to include IsLocked = true/false in conditions, as this function adds it
func (s *Schema[T]) GetOrCreateAndLockTx(rqtx *my.Request, p *GetOrCreateAndLockParams[T]) ds.Result[T] {
	this := s.instance
	isUnlocked := qb.Equal[T](this, p.LockField, false)

	// Get or create unlocked item
	if p.PostCondition == nil {
		p.PostCondition = isUnlocked
	} else {
		p.PostCondition = qb.And(p.PostCondition, isUnlocked)
	}
	result := s.GetOrCreateTx(rqtx, &p.GetOrCreateParams)
	if result.IsError() {
		return ds.Error[T](result.Error())
	}

	// Lock item
	item := result.Value()
	lockCondition := qb.And(p.LockConditionFn(item), isUnlocked)
	err := s.SetTxFlag(rqtx, lockCondition, p.LockField, true)
	if err != nil {
		rqtx.Fail(my.Err500, "Failed to lock item")
		return ds.Error[T](err)
	}

	return ds.NewResult(item, nil)
}

// UpdateAndGetTx updates one time and gets it as part of a transaction
func (s *Schema[T]) UpdateAndGetTx(rqtx *my.Request, updateFn UpdateFn[T], updateCondition, selectCondition qb.DualCondition[T]) ds.Result[T] {
	this := s.instance

	// Update one item
	q := qb.NewUpdateQuery[T](this, s.Table)
	q.Where(updateCondition)
	q.Limit(1)
	updateFn(this, q)
	err := s.update(rqtx, q, true)
	if err != nil {
		return ds.Error[T](err)
	}

	// Select one item
	result := s.Get(rqtx, selectCondition)
	if result.IsError() {
		// Manual rollback on error of Get
		err = qb.Rollback(rqtx.Tx, result.Error())
		return ds.Error[T](err)
	}

	return result
}

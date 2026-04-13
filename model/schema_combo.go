package model

import (
	"fmt"

	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/qb"
)

type GetOrCreateParams[T any] struct {
	Name          string
	Owner         string
	PreCondition  qb.Condition
	PostCondition qb.DualCondition[T]
	NewFn         func() T
	UpdateFn      func(*T, ID)
}

type GetOrCreateAndLockParams[T any] struct {
	GetOrCreateParams[T]
	LockField       *bool
	LockConditionFn func(T) qb.Condition
}

// GetOrCreate gets the item if it exists, otherwise creates and returns it
func (s *Schema[T]) GetOrCreate(rq *my.Request, p *GetOrCreateParams[T]) (T, error) {
	return s.getOrCreate(rq, p, false)
}

// GetOrCreateTx gets the item if it exists, otherwise creates it as part of a transaction, and return the item
func (s *Schema[T]) GetOrCreateTx(rqtx *my.Request, p *GetOrCreateParams[T]) (T, error) {
	return s.getOrCreate(rqtx, p, true)
}

// Common: get the item if it exists, otherwise create and return it
func (s *Schema[T]) getOrCreate(rq *my.Request, p *GetOrCreateParams[T], isTx bool) (T, error) {
	var zero T
	// Fetch item to check if it exists
	rows, err := s.GetRows(rq, p.PreCondition)
	if err != nil {
		rq.Fail(my.Err500, "Failed to check if %s exists", p.Name)
		if isTx {
			err = qb.Rollback(rq.Tx, err) // manual rollback
		}
		return zero, err
	}

	numRows := len(rows)
	if numRows > 1 {
		rq.Fail(my.Err500, "Failed to get one %s", p.Name)
		err := fmt.Errorf("public: Multiple %s found", p.Name)
		if isTx {
			err = qb.Rollback(rq.Tx, err) // manual rollback
		}
		return zero, err
	} else if numRows == 0 {
		// Not found = create item
		newItem := p.NewFn()
		var id ID
		if isTx {
			id, err = s.InsertTx(rq, &newItem)
		} else {
			id, err = s.Insert(rq, &newItem)
		}
		if err != nil {
			return zero, err
		}
		if p.UpdateFn != nil {
			p.UpdateFn(&newItem, id)
		}

		rq.AddFmtLog("Created %s for %s", p.Name, p.Owner)
		return newItem, nil
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
		return zero, err
	}
	return item, nil
}

// GetAndLockTx gets the item and locks it as part of a transaction
// Note: no need to include IsLocked = true/false in conditions, as this function adds it
func (s *Schema[T]) GetAndLockTx(rqtx *my.Request, lockField *bool, selectCondition qb.Condition, lockConditionFn func(T) qb.Condition) (T, error) {
	var zero T
	this := s.Instance
	isUnlocked := qb.Equal(this, lockField, false)

	// Get unlocked item
	condition := qb.And(selectCondition, isUnlocked)
	item, err := s.Get(rqtx, condition)
	if err != nil {
		rqtx.Fail(my.Err500, "Failed to get unlocked item")
		// Manual rollback on error of Get
		return zero, qb.Rollback(rqtx.Tx, err)
	}

	// Lock item
	lockCondition := qb.And(lockConditionFn(item), isUnlocked)
	err = s.SetTxFlag(rqtx, lockCondition, lockField, true)
	if err != nil {
		rqtx.Fail(my.Err500, "Failed to lock item")
		return zero, err
	}

	return item, nil
}

// GetAndLockTxItems gets a list of items and locks all of them as part of a transaction
// Note: no need to include IsLocked = true/false in conditions, as this function adds it
func (s *Schema[T]) GetAndLockTxItems(rqtx *my.Request, lockField *bool, selectCondition qb.Condition, lockConditionFn func([]T) qb.Condition, numItems int) ([]T, error) {
	this := s.Instance
	isUnlocked := qb.Equal(this, lockField, false)

	// Get unlocked items
	condition := qb.And(selectCondition, isUnlocked)
	items, err := s.GetRows(rqtx, condition)
	if err != nil {
		rqtx.Fail(my.Err500, "Failed to get unlocked items")
		// Manual rollback on error of GetRows
		return nil, qb.Rollback(rqtx.Tx, err)
	}

	// Check that rows have correct number of items
	if len(items) != numItems {
		rqtx.Fail(my.Err500, "Get count mismatch: items = %d, rows = %d", numItems, len(items))
		// Manual rollback if count mismatch
		return nil, qb.Rollback(rqtx.Tx, fail.MismatchCount)
	}

	// Lock items
	lockCondition := qb.And(lockConditionFn(items), isUnlocked)
	err = s.SetTxFlags(rqtx, lockCondition, lockField, true, numItems)
	if err != nil {
		rqtx.Fail(my.Err500, "Failed to lock items")
		return nil, err
	}

	return items, nil
}

// GetOrCreateAndLockTx runs GetOrCreate and locks the item as part of a transaction
// Note: no need to include IsLocked = true/false in conditions, as this function adds it
func (s *Schema[T]) GetOrCreateAndLockTx(rqtx *my.Request, p *GetOrCreateAndLockParams[T]) (T, error) {
	var zero T
	this := s.Instance
	isUnlocked := qb.Equal2[T](this, p.LockField, false)

	// Get or create unlocked item
	if p.PostCondition == nil {
		p.PostCondition = isUnlocked
	} else {
		p.PostCondition = qb.And2(p.PostCondition, isUnlocked)
	}
	item, err := s.GetOrCreateTx(rqtx, &p.GetOrCreateParams)
	if err != nil {
		return zero, err
	}

	// Lock item
	lockCondition := qb.And(p.LockConditionFn(item), isUnlocked)
	err = s.SetTxFlag(rqtx, lockCondition, p.LockField, true)
	if err != nil {
		rqtx.Fail(my.Err500, "Failed to lock item")
		return zero, err
	}

	return item, nil
}

// UpdateAndGetTx updates one time and gets it as part of a transaction
func (s *Schema[T]) UpdateAndGetTx(rqtx *my.Request, updateFn UpdateFn[T], updateCondition, selectCondition qb.Condition) (T, error) {
	var zero T
	this := s.Instance

	// Update one item
	q := qb.NewUpdateQuery[T](this, s.Table)
	q.Where(updateCondition)
	q.Limit(1)
	updateFn(this, q)
	err := s.update(rqtx, q, true)
	if err != nil {
		return zero, err
	}

	// Select one item
	item, err := s.Get(rqtx, selectCondition)
	if err != nil {
		// Manual rollback on error of Get
		return zero, qb.Rollback(rqtx.Tx, err)
	}

	return item, nil
}

// MoveItemTx inserts an item to the insertSchema and deletes the corresponding item from the deleteSchema, as part of a transaction
func MoveItemTx[I, D any](rqtx *my.Request, insertSchema *Schema[I], item *I, deleteSchema *Schema[D], deleteCondition qb.Condition) error {
	// 1) Insert item to insertSchema
	_, err := insertSchema.InsertTx(rqtx, item)
	if err != nil {
		return err
	}
	// 2) Delete from deleteSchema
	_, err = deleteSchema.DeleteTx(rqtx, deleteCondition)
	if err != nil {
		return err
	}
	return nil
}

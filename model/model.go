// Package model contains Schema type and methods
package model

import (
	"github.com/zeroibot/pack/ds"
	"github.com/zeroibot/pack/dyn"
	"github.com/zeroibot/pack/qb"
)

type Schema[T any] struct {
	Name     string
	Ref      *T
	Table    string
	Reader   qb.RowReader[T]
	instance *qb.Instance
}

// NewSchema creates a new Schema and registers its type to qb
func NewSchema[T any](this *qb.Instance, structRef *T, table string) ds.Result[*Schema[T]] {
	err := qb.AddType(this, structRef)
	if err != nil {
		return ds.Error[*Schema[T]](err)
	}
	schema := new(Schema[T]{
		Name:     dyn.TypeName(structRef),
		Ref:      structRef,
		Table:    table,
		Reader:   qb.FullRowReader(this, structRef),
		instance: this,
	})
	return ds.NewResult(schema, nil)
}

// NewSharedSchema creates a new shared Schema (no table)
func NewSharedSchema[T any](this *qb.Instance, structRef *T) ds.Result[*Schema[T]] {
	return NewSchema(this, structRef, "")
}

// AddSchema adds a new Schema and adds error to the error list if any
func AddSchema[T any](this *qb.Instance, structRef *T, table string, errs []error) *Schema[T] {
	result := NewSchema(this, structRef, table)
	if result.IsError() {
		errs = append(errs, result.Error())
	}
	return result.Value()
}

// AddSharedSchema adds a new shared Schema (no table) and adds error to the error list if any
func AddSharedSchema[T any](this *qb.Instance, structRef *T, errs []error) *Schema[T] {
	result := NewSharedSchema(this, structRef)
	if result.IsError() {
		errs = append(errs, result.Error())
	}
	return result.Value()
}

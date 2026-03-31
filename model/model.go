// Package model contains Schema type and methods
package model

import (
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
func NewSchema[T any](this *qb.Instance, structRef *T, table string) (*Schema[T], error) {
	err := qb.AddType(this, structRef)
	if err != nil {
		return nil, err
	}
	schema := new(Schema[T]{
		Name:     dyn.TypeName(structRef),
		Ref:      structRef,
		Table:    table,
		Reader:   qb.FullRowReader(this, structRef),
		instance: this,
	})
	return schema, nil
}

// NewSharedSchema creates a new shared Schema (no table)
func NewSharedSchema[T any](this *qb.Instance, structRef *T) (*Schema[T], error) {
	return NewSchema(this, structRef, "")
}

// AddSchema adds a new Schema and adds error to the error list if any
func AddSchema[T any](this *qb.Instance, structRef *T, table string, errs []error) *Schema[T] {
	schema, err := NewSchema(this, structRef, table)
	if err != nil {
		errs = append(errs, err)
	}
	return schema
}

// AddSharedSchema adds a new shared Schema (no table) and adds error to the error list if any
func AddSharedSchema[T any](this *qb.Instance, structRef *T, errs []error) *Schema[T] {
	schema, err := NewSharedSchema(this, structRef)
	if err != nil {
		errs = append(errs, err)
	}
	return schema
}

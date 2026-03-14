package qb

import (
	"fmt"

	"github.com/roidaradal/pack/dict"
	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/dyn"
)

// rowScanner is an Interface that unifies *sql.Row and *sql.Rows
type rowScanner interface {
	Scan(...any) error
}

// RowReader is a function that reads row values into a struct
type RowReader[T any] = func(rowScanner) (ds.Option[T], error)

// ToRow converts a given struct to map[string]any for row insertion
func ToRow[T any](this *Instance, structRef *T) dict.Object {
	typeName := dyn.TypeName(structRef)
	rowFn, ok := this.typeRowCreators[typeName]
	if !ok {
		return dict.Object{}
	}
	return rowFn(structRef)
}

// FullRowReader creates a RowReader for type T, using all columns
func FullRowReader[T any](this *Instance, structRef *T) RowReader[T] {
	columns := this.allColumns(structRef)
	return NewRowReader[T](this, columns...)
}

// NewRowReader creates a RowReader for type T, with the given columns
func NewRowReader[T any](this *Instance, columns ...string) RowReader[T] {
	return func(row rowScanner) (ds.Option[T], error) {
		var item T
		nilOption := ds.Nil[T]()
		if !dyn.IsStruct(item) {
			return nilOption, fmt.Errorf("not a struct type")
		}
		typeName := dyn.TypeName(item)
		numColumns := len(columns)
		fieldRefs := make([]any, 0, numColumns)
		for _, column := range columns {
			if column == "" {
				continue // skip blank columns
			}
			fieldRef, ok := this.getStructColumnFieldRef(&item, typeName, column)
			if !ok {
				continue // skip if column's field not found
			}
			fieldRefs = append(fieldRefs, fieldRef)
		}
		if len(fieldRefs) != numColumns {
			// Return nil if some columns failed
			return nilOption, fmt.Errorf("incomplete fields")
		}
		err := row.Scan(fieldRefs...)
		return ds.NewOption(&item), err
	}
}

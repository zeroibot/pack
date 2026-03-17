package qb

import (
	"github.com/roidaradal/pack/dict"
	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/dyn"
)

// createRowFn creates a map[string]any from given struct reference, used for row insertion
type createRowFn func(structRef any) dict.Object

// Internal: newRowCreator creates a createRowFn for the given columns.
// The produced function returns a map[string]any from given struct reference, using the given columns as keys.
func (i *Instance) newRowCreator(typeName string, columns ds.List[string]) createRowFn {
	return func(structRef any) dict.Object {
		emptyRow := dict.Object{}
		if !dyn.IsStructPointer(structRef) {
			return emptyRow
		}
		numColumns := columns.Len()
		row := make(dict.Object, numColumns)
		for _, column := range columns {
			value, ok := i.getStructColumnValue(structRef, typeName, column)
			if !ok {
				continue // skip if column value not found
			}
			row[column] = value
		}
		if len(row) != numColumns {
			// Return empty row if some columns failed
			return emptyRow
		}
		return row
	}
}

// Internal: get reference to corresponding field of given struct reference, type name, and column name
func (i *Instance) getStructColumnFieldRef(structRef any, typeName, columnName string) (any, bool) {
	structField, ok := i.getStructColumnField(structRef, typeName, columnName)
	if !ok {
		return nil, false
	}
	return dyn.RefValue(structField)
}

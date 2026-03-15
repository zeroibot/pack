package qb

import (
	"reflect"

	"github.com/roidaradal/pack/conv"
	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/dyn"
)

const (
	columnTag       string = "col" // struct tag used to define column name
	skipColumnValue string = "-"   // set `col:"-" to skip column
)

// Holds information for one struct type's columns and field names
type columnsInfo struct {
	columns        ds.List[string]        // list of column names
	columnFields   ds.Map[string, string] // {ColumnName => FieldName}
	fieldColumns   ds.Map[string, string] // {FieldName => ColumnName}
	addressColumns ds.Map[string, string] // {FieldAddress => ColumnName}
	addressFields  ds.Map[string, string] // {FieldAddress => FieldName}
}

// Internal: given the struct pointer, extract all column and field names
// Uses recursion for embedded struct fields.
func (i *Instance) readStructColumns(structRef any) *columnsInfo {
	info := new(columnsInfo{
		columns:        make(ds.List[string], 0),
		columnFields:   make(ds.Map[string, string]),
		fieldColumns:   make(ds.Map[string, string]),
		addressColumns: make(ds.Map[string, string]),
		addressFields:  make(ds.Map[string, string]),
	})

	// Ensure struct ref is a struct pointer, otherwise the methods below could fail
	if !dyn.IsStructPointer(structRef) {
		return info
	}

	structValue := dyn.MustDerefValue(structRef)
	structType := structValue.Type()
	for idx := range structType.NumField() {
		structField := structType.Field(idx)
		fieldName := structField.Name
		if structField.Anonymous {
			// Embedded struct: get columns using recursion
			// Get reference to inner struct
			innerStructRef, ok := dyn.RefValue(structValue.FieldByName(fieldName))
			if !ok {
				continue // skip if not a valid inner struct ref
			}
			inner := i.readStructColumns(innerStructRef)
			info.columns = append(info.columns, inner.columns...)
			info.columnFields.Update(inner.columnFields)
			info.fieldColumns.Update(inner.fieldColumns)
			info.addressColumns.Update(inner.addressColumns)
			info.addressFields.Update(inner.addressFields)
		} else {
			// Normal field
			columnName, ok := dyn.GetStructFieldTag(structField, columnTag)
			if !ok || columnName == "" {
				// Column name defaults to field name if column tag is not set or blank
				columnName = fieldName
			} else if columnName == skipColumnValue {
				continue // skip if column is explicitly set to skip
			}
			columnName = i.prepareColumn(columnName)
			structFieldRef, ok := dyn.RefValue(structValue.Field(idx))
			if !ok {
				continue // skip if struct field cannot be referenced
			}
			fieldAddress := conv.AnyToString(structFieldRef)
			info.columns = append(info.columns, columnName)
			info.columnFields[columnName] = fieldName
			info.fieldColumns[fieldName] = columnName
			info.addressColumns[fieldAddress] = columnName
			info.addressFields[fieldAddress] = fieldName
		}
	}
	return info
}

// Internal: allColumns returns the column names of given item's type
func (i *Instance) allColumns(item any) ds.List[string] {
	typeName := dyn.TypeName(item)
	return i.typeColumns.GetOrDefault(typeName, ds.List[string]{})
}

// Internal: get corresponding field name from given type's column name
func (i *Instance) getColumnFieldName(typeName, columnName string) string {
	if i.typeColumnFields.NoKey(typeName) {
		return ""
	}
	return i.typeColumnFields[typeName][columnName]
}

// Internal: get corresponding column name from given type's field name
func (i *Instance) getFieldColumnName(typeName, fieldName string) string {
	if i.typeFieldColumns.NoKey(typeName) {
		return ""
	}
	return i.typeFieldColumns[typeName][fieldName]
}

// Internal: common steps for getting the struct field reflect.Value from given struct reference, type name, and column name
func (i *Instance) getStructColumnField(structRef any, typeName, columnName string) (reflect.Value, bool) {
	var zero reflect.Value
	structValue, ok := dyn.DerefValue(structRef)
	if !ok {
		return zero, false // invalid struct pointer
	}
	fieldName := i.getColumnFieldName(typeName, columnName)
	if fieldName == "" {
		return zero, false // field name not found
	}
	return structValue.FieldByName(fieldName), true
}

// Internal: get field value from given struct reference, type name, and column name
func (i *Instance) getStructColumnValue(structRef any, typeName, columnName string) (any, bool) {
	structField, ok := i.getStructColumnField(structRef, typeName, columnName)
	if !ok {
		return nil, false
	}
	return dyn.AnyValue(structField)
}

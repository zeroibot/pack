// Package qb contains the SQL QueryBuilder types and functions
package qb

import (
	"fmt"
	"strings"

	"github.com/roidaradal/pack/db"
	"github.com/roidaradal/pack/dict"
	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/dyn"
	"github.com/roidaradal/pack/list"
	"github.com/roidaradal/pack/str"
)

type dbType struct {
	name string
}

var (
	MySQL = dbType{"mysql"}
)

// prepareIdentifier wraps the identifier depending on the database type
func (db dbType) prepareIdentifier(identifier string) string {
	switch db.name {
	case MySQL.name:
		return str.Wrap(identifier, "``")
	default:
		return identifier
	}
}

// rawIdentifier removes the identifier wrappers depending on the database type
func (db dbType) rawIdentifier(identifier string) string {
	switch db.name {
	case MySQL.name:
		return strings.Trim(identifier, "`")
	default:
		return identifier
	}
}

// Instance type stores the needed QueryBuilder data for registered types' columns and fields
type Instance struct {
	dbType
	addressColumns   ds.Map[string, string]                 // {FieldAddress => ColumnName}
	addressFields    ds.Map[string, string]                 // {FieldAddress => FieldName}
	typeColumns      ds.Map[string, ds.List[string]]        // {TypeName => []ColumnNames}
	typeColumnFields ds.Map[string, ds.Map[string, string]] // {TypeName => {ColumnName => FieldName}}
	typeFieldColumns ds.Map[string, ds.Map[string, string]] // {TypeName => {FieldName => ColumnName}}
	typeRowCreators  ds.Map[string, createRowFn]            // {TypeName => createRowFn}
}

// NewInstance creates a new QueryBuilder Instance
func NewInstance(db dbType) *Instance {
	return &Instance{
		dbType:           db,
		addressColumns:   make(ds.Map[string, string]),
		addressFields:    make(ds.Map[string, string]),
		typeColumns:      make(ds.Map[string, ds.List[string]]),
		typeColumnFields: make(ds.Map[string, ds.Map[string, string]]),
		typeFieldColumns: make(ds.Map[string, ds.Map[string, string]]),
		typeRowCreators:  make(ds.Map[string, createRowFn]),
	}
}

// AddType registers a new struct type, and stores its column and field names into the Instance.
func AddType[T any](this *Instance, structRef *T) error {
	if !dyn.IsStructPointer(structRef) {
		return fmt.Errorf("parameter is not a struct pointer")
	}
	typeName := dyn.TypeName(structRef)

	info := this.readStructColumns(structRef)
	this.addressColumns.Update(info.addressColumns)
	this.addressFields.Update(info.addressFields)
	this.typeColumns[typeName] = info.columns
	this.typeColumnFields[typeName] = info.columnFields
	this.typeFieldColumns[typeName] = info.fieldColumns
	this.typeRowCreators[typeName] = this.newRowCreator(typeName, info.columns)
	return nil
}

// LookupColumnName looks up the associated column name for the given struct field reference.
// The struct field reference must come from the Type singleton object used during registration.
func (i *Instance) LookupColumnName(fieldRef any) (string, bool) {
	fieldAddress := dyn.AddressOf(fieldRef)
	columnName, ok := i.addressColumns[fieldAddress]
	return columnName, ok
}

// Column gets the associated column name for the given struct field reference.
// The struct field reference must come from the Type singleton object used during registration.
// Returns empty string if the column name is not found.
func (i *Instance) Column(fieldRef any) string {
	fieldAddress := dyn.AddressOf(fieldRef)
	return i.addressColumns[fieldAddress]
}

// Columns gets the associated column names of given struct field references.
// The struct field references must come from the Type singleton object used during registration.
// Returns empty list if at least one column name is not found.
func (i *Instance) Columns(fieldRefs ...any) ds.List[string] {
	columns := list.MapIf(fieldRefs, i.LookupColumnName)
	if len(columns) != len(fieldRefs) {
		// Return empty list if not all columns found
		return ds.List[string]{}
	}
	return columns
}

// Field gets the associated field name for the given struct field reference.
// The struct field reference must come from the Type singleton object used during registration.
func (i *Instance) Field(typeName string, fieldRef any) string {
	return i.getColumnFieldName(typeName, i.Column(fieldRef))
}

// Fields gets the associated field names for the given struct field references.
// The struct field references must come from the Type singleton object used during registration.
// Returns empty list if at least one field name is not found.
func (i *Instance) Fields(typeName string, fieldRefs ...any) ds.List[string] {
	fieldNames := list.MapIf(fieldRefs, func(fieldRef any) (string, bool) {
		fieldName := i.Field(typeName, fieldRef)
		return fieldName, fieldName != ""
	})
	if len(fieldNames) != len(fieldRefs) {
		// Return empty list if not all fields found
		return ds.List[string]{}
	}
	return fieldNames
}

// RowReader is a function that reads row values into a struct
type RowReader[T any] = func(db.RowScanner) ds.Result[T]

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
	return func(row db.RowScanner) ds.Result[T] {
		var item T
		if !dyn.IsStruct(item) {
			return ds.Error[T](fmt.Errorf("not a struct type"))
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
			return ds.Error[T](fmt.Errorf("incomplete fields"))
		}
		err := row.Scan(fieldRefs...)
		if err != nil {
			return ds.Error[T](err)
		}
		return ds.NewResult(item, nil)
	}
}

// ToString builds the Query string
func ToString(q Query) string {
	query, rawValues := q.BuildQuery()
	values := make([]any, len(rawValues))
	formats := make([]string, len(rawValues))
	for i, value := range rawValues {
		typeName := fmt.Sprintf("%T", value)
		if strings.HasPrefix(typeName, "*") {
			values[i] = dyn.MustDeref(value)
		} else {
			values[i] = value
		}
		if _, ok := values[i].(string); ok {
			formats[i] = "%q"
		} else {
			formats[i] = "%v"
		}
	}
	for _, format := range formats {
		query = strings.Replace(query, "?", format, 1)
	}
	return fmt.Sprintf(query, values...)
}

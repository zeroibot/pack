package qb

import (
	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/dyn"
	"github.com/roidaradal/pack/list"
)

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

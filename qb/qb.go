// Package qb contains the SQL QueryBuilder types and functions
package qb

import (
	"fmt"

	"github.com/roidaradal/pack/dyn"
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

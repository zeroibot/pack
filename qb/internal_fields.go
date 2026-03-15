package qb

import "github.com/roidaradal/pack/dyn"

// Internal: get corresponding field name of given field reference
func (i *Instance) getFieldName(fieldRef any) string {
	fieldAddress := dyn.AddressOf(fieldRef)
	return i.addressFields[fieldAddress]
}

// Internal: get the value of field name from given struct pointer, and type coerce into type T
func getStructFieldValue[V any](structRef any, fieldName string) (V, bool) {
	var zero V
	structValue, ok := dyn.DerefValue(structRef)
	if !ok {
		return zero, false
	}
	rawValue, ok := dyn.AnyValue(structValue.FieldByName(fieldName))
	if !ok {
		return zero, false
	}
	value, ok := rawValue.(V)
	if !ok {
		return zero, false
	}
	return value, true
}

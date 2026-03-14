package qb

import "testing"

func TestNewInstance(t *testing.T) {
	type testCase struct {
		name string
		flag bool
	}
	// Manual Instance
	this1 := &Instance{}
	this2 := NewInstance(MySQL)
	tests := []testCase{
		{"addressColumns", this1.addressColumns == nil},
		{"typeColumns", this1.typeColumns == nil},
		{"typeColumnFields", this1.typeColumnFields == nil},
		{"typeFieldColumns", this1.typeFieldColumns == nil},
		{"addressColumns", this2.addressColumns != nil},
		{"typeColumns", this2.typeColumns != nil},
		{"typeColumnFields", this2.typeColumnFields != nil},
		{"typeFieldColumns", this2.typeFieldColumns != nil},
		{"addressColumns", this2.typeColumns.IsEmpty()},
		{"typeColumns", this2.typeColumns.IsEmpty()},
		{"typeColumnFields", this2.typeColumns.IsEmpty()},
		{"typeFieldColumns", this2.typeFieldColumns.IsEmpty()},
	}
	for _, x := range tests {
		if !x.flag {
			t.Errorf("Instance.%s test = %v, want true", x.name, x.flag)
		}
	}
}

func TestInstanceMethods(t *testing.T) {
	// TODO: LookupColumnName
	// TODO: Column
	// TODO: Columns
	// TODO: Field
	// TODO: Fields
}

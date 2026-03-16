package qb

import (
	"slices"
	"testing"
)

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
		{"addressFields", this1.addressFields == nil},
		{"typeColumns", this1.typeColumns == nil},
		{"typeColumnFields", this1.typeColumnFields == nil},
		{"typeFieldColumns", this1.typeFieldColumns == nil},
		{"typeRowCreators", this1.typeRowCreators == nil},
		{"addressColumns", this2.addressColumns != nil},
		{"addressFields", this2.addressFields != nil},
		{"typeColumns", this2.typeColumns != nil},
		{"typeColumnFields", this2.typeColumnFields != nil},
		{"typeFieldColumns", this2.typeFieldColumns != nil},
		{"typeRowCreators", this2.typeRowCreators != nil},
		{"addressColumns", this2.typeColumns.IsEmpty()},
		{"addressFields", this2.typeColumns.IsEmpty()},
		{"typeColumns", this2.typeColumns.IsEmpty()},
		{"typeColumnFields", this2.typeColumns.IsEmpty()},
		{"typeFieldColumns", this2.typeFieldColumns.IsEmpty()},
		{"typeRowCreators", this2.typeRowCreators.IsEmpty()},
	}
	for _, x := range tests {
		if !x.flag {
			t.Errorf("Instance.%s test = %v, want true", x.name, x.flag)
		}
	}
}

func TestInstanceMethods(t *testing.T) {
	type Person struct {
		Name    string
		Address string
		Job     string
		Age     int
	}
	type School struct {
		Name  string
		Logo  string
		Level int `col:"Lvl"`
	}
	this := NewInstance(MySQL)
	p := &Person{}
	s := &School{}
	err := AddType(this, p)
	if err != nil {
		t.Errorf("AddType() error: %v", err)
	}
	err = AddType(this, s)
	if err != nil {
		t.Errorf("AddType() error: %v", err)
	}
	john := &Person{"John", "PH", "dev", 20}
	// LookupColumnName
	type testCase1 struct {
		want     string
		fieldRef any
	}
	testCases1 := []testCase1{
		{"`Job`", &p.Job},
		{"`Lvl`", &s.Level},
		{"", &john.Name}, // Not from type singleton
		{"", p.Name},     // Not a struct reference
		{"", s.Logo},
	}
	for _, x := range testCases1 {
		column, ok := this.LookupColumnName(x.fieldRef)
		wantOk := x.want != ""
		if column != x.want || ok != wantOk {
			t.Errorf("LookupColumnName() = %q, %t, want %q, %t", column, ok, x.want, wantOk)
		}
		column = this.Column(x.fieldRef)
		if column != x.want {
			t.Errorf("Column() = %q, want %q", column, x.want)
		}
	}
	// Field
	type testCase2 struct {
		want     string
		typeName string
		fieldRef any
	}
	testCases2 := []testCase2{
		{"Job", "Person", &p.Job},
		{"Level", "School", &s.Level},
		{"Age", "Person", &p.Age},
		{"", "Person", &john.Name}, // Not from type singleton
		{"", "School", s.Name},     // Not a struct reference
		{"", "People", &p.Age},     // Wrong typeName
	}
	for _, x := range testCases2 {
		field := this.Field(x.typeName, x.fieldRef)
		if field != x.want {
			t.Errorf("Field() = %q, want %q", field, x.want)
		}
	}
	// Columns
	type testCase3 struct {
		want      []string
		fieldRefs []any
	}
	testCases3 := []testCase3{
		{[]string{"`Job`", "`Age`", "`Name`"}, []any{&p.Job, &p.Age, &p.Name}},
		{[]string{"`Lvl`", "`Logo`", "`Name`"}, []any{&s.Level, &s.Logo, &s.Name}},
		{[]string{}, []any{&p.Name, &p.Age, &john.Job}}, // has one invalid fieldRef
		{[]string{}, []any{&s.Name, &s.Logo, s.Level}},  // has one non-struct ref
	}
	for _, x := range testCases3 {
		actual := this.Columns(x.fieldRefs...)
		if slices.Equal(actual, x.want) == false {
			t.Errorf("Columns() = %v, want %v", actual, x.want)
		}
	}
	// Fields
	type testCase4 struct {
		want      []string
		typeName  string
		fieldRefs []any
	}
	testCases4 := []testCase4{
		{[]string{"Job", "Age", "Name"}, "Person", []any{&p.Job, &p.Age, &p.Name}},
		{[]string{"Level", "Logo", "Name"}, "School", []any{&s.Level, &s.Logo, &s.Name}},
		{[]string{}, "Person", []any{&p.Name, &p.Age, &john.Job}}, // has one invalid fieldRef
		{[]string{}, "School", []any{&s.Name, s.Level}},           // has one non-struct reference
		{[]string{}, "People", []any{&p.Name, &p.Address}},        // wrong typeName
	}
	for _, x := range testCases4 {
		actual := this.Fields(x.typeName, x.fieldRefs...)
		if slices.Equal(actual, x.want) == false {
			t.Errorf("Fields() = %v, want %v", actual, x.want)
		}
	}
}

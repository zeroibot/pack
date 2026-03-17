package qb

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"testing"

	"github.com/roidaradal/pack/dict"
	"github.com/roidaradal/pack/dyn"
)

type mockScanner struct {
	items []any
}

func (m mockScanner) Scan(fieldRefs ...any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic encountered: %v", r)
		}
	}()
	if len(fieldRefs) != len(m.items) {
		return fmt.Errorf("expected %d fieldRefs, got %d", len(m.items), len(fieldRefs))
	}
	for i, fieldRef := range fieldRefs {
		fieldValue := dyn.MustDerefValue(fieldRef)
		fieldValue.Set(reflect.ValueOf(m.items[i]))
	}
	return err
}

func TestAddType(t *testing.T) {
	type User struct {
		Name     string
		Password string
		Age      int
	}
	type School struct {
		Name    string
		Address string
		Level   int
	}
	this := NewInstance(MySQL)
	err := AddType(this, new(5))
	if err == nil {
		t.Errorf("AddType(5) should return error")
	}
	err = AddType(this, &User{})
	if err != nil {
		t.Errorf("AddType error: %s", err.Error())
	}
	err = AddType(this, &School{})
	if err != nil {
		t.Errorf("AddType error: %s", err.Error())
	}
	// Note: checking of result details are checked in the individual methods
}

func TestPrepareIdentifier(t *testing.T) {
	db1 := MySQL
	db2 := dbType{"other"}
	testCases := [][2]string{
		{"`name`", db1.prepareIdentifier("name")},
		{"`age`", db1.prepareIdentifier("age")},
		{"name", db2.prepareIdentifier("name")},
		{"age", db2.prepareIdentifier("age")},
	}
	for _, x := range testCases {
		want, actual := x[0], x[1]
		if want != actual {
			t.Errorf("dbType.prepareIdentifier() = %q, want %q", actual, want)
		}
	}
}

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
		{"Job", &p.Job},
		{"Lvl", &s.Level},
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
		{[]string{"Job", "Age", "Name"}, []any{&p.Job, &p.Age, &p.Name}},
		{[]string{"Lvl", "Logo", "Name"}, []any{&s.Level, &s.Logo, &s.Name}},
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

func TestRowFunctions(t *testing.T) {
	type User struct {
		Name     string
		Password string
		Age      int
		secret   string
	}
	type School struct {
		Name    string
		Address string
	}
	this := NewInstance(MySQL)
	user := &User{"john", "123456", 25, "secret"}
	school := &School{"UP", "Lahug"}
	userRef := new(User)
	err := AddType(this, userRef)
	if err != nil {
		t.Errorf("AddType() error = %v", err)
	}
	// ToRow
	empty := dict.Object{}
	userObj := dict.Object{"Name": "john", "Password": "123456", "Age": 25}
	testCases := [][2]dict.Object{
		{userObj, ToRow(this, user)},
		{empty, ToRow(this, school)},
	}
	for _, x := range testCases {
		want, actual := x[0], x[1]
		if maps.Equal(want, actual) == false {
			t.Errorf("ToRow() = %v, want %v", actual, want)
		}
	}
	// Not a struct type
	intReader := NewRowReader[int](this, "Value", "Decimal")
	intResult := intReader(mockScanner{})
	if intResult.IsError() == false || intResult.Value() != 0 {
		t.Errorf("NewRowReader[int] should return an error")
	}
	// Valid full reader
	fullReader := FullRowReader(this, userRef)
	if fullReader == nil {
		t.Errorf("FullRowReader() should return a rowReader, got nil")
	}
	// Successful read
	result := fullReader(mockScanner{items: []any{"John", "111", 20}})
	if result.NotError() == false {
		t.Errorf("FullRowReader() read = %v, want non-error", result)
	}
	// Check that struct has been filled after fullReader read
	want := User{"John", "111", 20, ""}
	if want != result.Value() {
		t.Errorf("FullRowReader() read = %v, want %v", result.Value(), want)
	}
	// Valid row reader, with specified columns
	nameCol, pwdCol := this.Column(&userRef.Name), this.Column(&userRef.Password)
	rowReader := NewRowReader[User](this, nameCol, pwdCol)
	result = rowReader(mockScanner{items: []any{"Jane", "222"}})
	if result.NotError() == false {
		t.Errorf("RowReader() read = %v, want non-error", result)
	}
	// Check that struct has been filled after rowReader read
	want = User{"Jane", "222", 0, ""}
	if want != result.Value() {
		t.Errorf("RowReader() read = %v, want %v", result.Value(), want)
	}
	emptyUser := User{}
	// Valid row reader, but error in scanning (mocked by incomplete items / invalid type)
	result = rowReader(mockScanner{items: []any{"Jane", 333}})
	if result.IsError() == false || result.Value() != emptyUser {
		t.Errorf("RowReader() read = %v, want Result<%v, error>", result, emptyUser)
	}
	result = rowReader(mockScanner{items: []any{"Jane"}})
	if result.IsError() == false || result.Value() != emptyUser {
		t.Errorf("RowReader() read = %v, want Result<%v, error>", result, emptyUser)
	}
	// Error because of blank columns
	userReader := NewRowReader[User](this, nameCol, pwdCol, "")
	result = userReader(mockScanner{})
	if result.IsError() == false || result.Value() != emptyUser {
		t.Errorf("NewRowReader() read = %v, want Result<%v, error>", result, emptyUser)
	}
	// Error because of unknown column field
	userReader = NewRowReader[User](this, nameCol, pwdCol, "secret")
	result = userReader(mockScanner{})
	if result.IsError() == false || result.Value() != emptyUser {
		t.Errorf("NewRowReader() read = %v, want Result<%v, error>", result, emptyUser)
	}
}

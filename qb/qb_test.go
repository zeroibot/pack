package qb

import (
	"fmt"
	"testing"

	"github.com/roidaradal/pack/dict"
	"github.com/roidaradal/tst"
)

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
	tst.AssertError(t, "AddType", AddType(this, new(5)))
	tst.AssertNoError(t, "AddType", AddType(this, new(User)))
	tst.AssertNoError(t, "AddType", AddType(this, new(School)))
	// Note: checking of result details are checked in the individual methods
}

func TestPrepareIdentifier(t *testing.T) {
	db1 := MySQL
	db2 := dbType{"other"}
	testCases := []tst.P2W1[dbType, string, string]{
		{db1, "name", "`name`"}, {db1, "age", "`age`"},
		{db2, "name", "name"}, {db2, "age", "age"},
	}
	tst.AllP2W1(t, testCases, "dbType.prepareIdentifier", dbType.prepareIdentifier, tst.AssertEqual)
	testCases2 := []tst.P2W1[dbType, string, string]{
		{db1, "`name`", "name"}, {db1, "`age`", "age"},
		{db2, "name", "name"}, {db2, "age", "age"},
	}
	tst.AllP2W1(t, testCases2, "dbType.rawIdentifier", dbType.rawIdentifier, tst.AssertEqual)
}

func TestNewInstance(t *testing.T) {
	this1 := &Instance{}
	this2 := NewInstance(MySQL)
	testCases := map[string][3]bool{
		"addressColumns":   {this1.addressColumns == nil, this2.addressColumns != nil, this2.typeColumns.IsEmpty()},
		"addressFields":    {this1.addressFields == nil, this2.addressFields != nil, this2.addressFields.IsEmpty()},
		"typeColumns":      {this1.typeColumns == nil, this2.typeColumns != nil, this2.typeColumns.IsEmpty()},
		"typeColumnFields": {this1.typeColumnFields == nil, this2.typeColumns != nil, this2.typeColumns.IsEmpty()},
		"typeFieldColumns": {this1.typeFieldColumns == nil, this2.typeFieldColumns != nil, this2.typeFieldColumns.IsEmpty()},
		"typeRowCreators":  {this1.typeRowCreators == nil, this2.typeRowCreators != nil, this2.typeRowCreators.IsEmpty()},
	}
	for name, flags := range testCases {
		f1, f2, f3 := flags[0], flags[1], flags[2]
		tst.AssertTrue(t, fmt.Sprintf("Instance.%s", name), f1 && f2 && f3)
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
	p, s := new(Person), new(School)
	john := &Person{"John", "PH", "dev", 20}
	this := testPrelude2(t, p, s)
	// LookupColumnName
	testCases1 := []tst.P1W2[any, string, bool]{
		{&p.Job, "Job", true}, {&s.Level, "Lvl", true},
		{&john.Name, "", false},                  // not from type singleton
		{p.Name, "", false}, {s.Logo, "", false}, // not a struct reference
	}
	tst.AllP1W2(t, testCases1, "LookupColumnName", this.LookupColumnName, tst.AssertEqual[string], tst.AssertEqual[bool])
	testCases2 := tst.Convert(testCases1, func(tc tst.P1W2[any, string, bool]) tst.P1W1[any, string] {
		return tst.P1W1[any, string]{P1: tc.P1, W1: tc.W1}
	})
	tst.AllP1W1(t, testCases2, "Column", this.Column, tst.AssertEqual)
	// Field
	testCases3 := []tst.P2W1[string, any, string]{
		{"Person", &p.Job, "Job"},
		{"School", &s.Level, "Level"},
		{"Person", &p.Age, "Age"},
		{"Person", &john.Name, ""}, // Not from type singleton
		{"School", s.Name, ""},     // Not a struct reference
		{"People", &p.Age, ""},     // Wrong typeName
	}
	tst.AllP2W1(t, testCases3, "Field", this.Field, tst.AssertEqual)
	// Columns
	testCases4 := []tst.P1W1[[]any, []string]{
		{[]any{&p.Job, &p.Age, &p.Name}, []string{"Job", "Age", "Name"}},
		{[]any{&s.Level, &s.Logo, &s.Name}, []string{"Lvl", "Logo", "Name"}},
		{[]any{&p.Name, &p.Age, &john.Job}, []string{}}, // has one invalid fieldRef
		{[]any{&s.Name, &s.Logo, s.Level}, []string{}},  // has one non-struct ref
	}
	getColumns := func(fieldRefs []any) []string { return this.Columns(fieldRefs...) }
	tst.AllP1W1(t, testCases4, "Columns", getColumns, tst.AssertListEqual)
	// Fields
	testCases5 := []tst.P2W1[string, []any, []string]{
		{"Person", []any{&p.Job, &p.Age, &p.Name}, []string{"Job", "Age", "Name"}},
		{"School", []any{&s.Level, &s.Logo, &s.Name}, []string{"Level", "Logo", "Name"}},
		{"Person", []any{&p.Name, &p.Age, &john.Job}, []string{}}, // has one invalid fieldRef
		{"School", []any{&s.Name, s.Level}, []string{}},           // has one non-struct reference
		{"People", []any{&p.Name, &p.Address}, []string{}},        // wrong typeName
	}
	getFields := func(typeName string, fieldRefs []any) []string { return this.Fields(typeName, fieldRefs...) }
	tst.AllP2W1(t, testCases5, "Fields", getFields, tst.AssertListEqual)
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
	user := new(User{"john", "123456", 25, "secret"})
	school := new(School{"UP", "Lahug"})
	userRef := new(User)
	this := testPrelude(t, userRef)
	// ToRow
	empty := dict.Object{}
	userObj := dict.Object{"Name": "john", "Password": "123456", "Age": 25}
	testCases := [][2]dict.Object{
		{userObj, ToRow(this, user)}, {empty, ToRow(this, school)},
	}
	for _, x := range testCases {
		want, actual := x[0], x[1]
		tst.AssertMapEqual(t, "ToRow", actual, want)
	}

	// Not a struct type
	intReader := NewRowReader[int](this, "Value", "Decimal")
	intResult := intReader(tst.NewRow())
	tst.AssertEqualAnd(t, "NewRowReader[int]", intResult.Value(), 0, intResult.IsError(), true)
	// Valid full reader
	fullReader := FullRowReader(this, userRef)
	tst.AssertTrue(t, "FullRowReader", fullReader != nil)
	// Successful read
	result := fullReader(tst.NewRow("John", "111", 20))
	tst.AssertTrue(t, "FullRowReader", result.NotError())
	// Check that struct has been filled after fullReader read
	want := User{"John", "111", 20, ""}
	tst.AssertEqual(t, "FullRowReader.Read", result.Value(), want)
	// Valid row reader, with specified columns
	nameCol, pwdCol := this.Column(&userRef.Name), this.Column(&userRef.Password)
	rowReader := NewRowReader[User](this, nameCol, pwdCol)
	result = rowReader(tst.NewRow("Jane", "222"))
	tst.AssertTrue(t, "RowReader.Read", result.NotError())
	// Check that struct has been filled after rowReader read
	want = User{"Jane", "222", 0, ""}
	tst.AssertEqual(t, "RowReader.Read", result.Value(), want)
	// Valid row reader, but error in scanning (invalid type)
	emptyUser := User{}
	result = rowReader(tst.NewRow("Jane", 333))
	tst.AssertEqualAnd(t, "RowReader.Read", result.Value(), emptyUser, result.IsError(), true)
	// Valid row reader, but error in scanning (incomplete items)
	result = rowReader(tst.NewRow("Jane"))
	tst.AssertEqualAnd(t, "RowReader.Read", result.Value(), emptyUser, result.IsError(), true)
	// Error because of blank columns
	userReader := NewRowReader[User](this, nameCol, pwdCol, "")
	result = userReader(tst.NewRow())
	tst.AssertEqualAnd(t, "NewRowReader.Read", result.Value(), emptyUser, result.IsError(), true)
	// Error because of unknown column field
	userReader = NewRowReader[User](this, nameCol, pwdCol, "secret")
	result = userReader(tst.NewRow())
	tst.AssertEqualAnd(t, "NewRowReader.Read", result.Value(), emptyUser, result.IsError(), true)
}

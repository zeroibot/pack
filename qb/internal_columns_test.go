package qb

import (
	"testing"

	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/tst"
)

func TestReadStructColumns(t *testing.T) {
	type secret struct {
		password string `col:"password"`
	}
	type Person struct {
		Name    string
		Address string
		secret
		Job string
		Age int
	}
	type Config struct {
		Key   string `col:"AppKey"`
		Value string `col:"AppValue"`
	}
	type School struct {
		Name string
		Logo string
		Config
		Secret   string `col:"-"`
		Level    int    `col:"Lvl"`
		capacity int
	}
	this := NewInstance(MySQL)
	name := "readStructColumns"
	// Not a struct
	info := this.readStructColumns(5)
	tst.AssertTrue(t, name, info.IsEmpty())
	// Not a struct pointer
	info = this.readStructColumns(Person{})
	tst.AssertTrue(t, name, info.IsEmpty())
	// Successful
	info = this.readStructColumns(&Config{})
	tst.AssertListEqual(t, name, info.columns, ds.List[string]{"AppKey", "AppValue"})
	// With embedded
	info = this.readStructColumns(&School{})
	tst.AssertListEqual(t, name, info.columns, ds.List[string]{"Name", "Logo", "AppKey", "AppValue", "Lvl"})
	// With private embedded
	info = this.readStructColumns(&Person{})
	tst.AssertListEqual(t, name, info.columns, ds.List[string]{"Name", "Address", "Job", "Age"})
}

func TestInternalColumnQueries(t *testing.T) {
	type Person struct {
		Name    string
		Address string
		Job     string
		Age     int
	}
	type School struct {
		Name     string
		Logo     string
		Secret   string `col:"-"`
		Level    int    `col:"Lvl"`
		capacity int
	}
	type Config struct {
		Key   string `col:"key"`
		Value string `col:"value"`
	}
	p, s, cfg := new(Person), new(School), new(Config)
	john := &Person{"John", "PH", "dev", 20}
	upc := &School{"UP", "Oblation", "12345", 3, 1000}
	this := testPrelude2(t, p, s)
	// allColumns
	want1 := ds.List[string]{"Name", "Address", "Job", "Age"}
	want2 := ds.List[string]{"Name", "Logo", "Lvl"}
	testCases := []tst.P1W1[any, ds.List[string]]{
		{p, want1}, {john, want1},
		{s, want2}, {upc, want2},
		{cfg, ds.List[string]{}}, // non-registered type
	}
	tst.AllP1W1(t, testCases, "allColumns", this.allColumns, tst.AssertListEqual)
	// getColumnFieldName
	testCases2 := []tst.P2W1[string, string, string]{
		{"Person", "Age", "Age"},
		{"School", "Lvl", "Level"},
		{"Config", "key", ""},
		{"School", "capacity", ""},
		{"School", "Secret", ""},
	}
	tst.AllP2W1(t, testCases2, "getColumnFieldName", this.getColumnFieldName, tst.AssertEqual)
	// getFieldColumnName
	testCases2 = []tst.P2W1[string, string, string]{
		{"Person", "Age", "Age"},
		{"School", "Level", "Lvl"},
		{"Config", "key", ""},
		{"School", "capacity", ""},
		{"School", "Secret", ""},
	}
	tst.AllP2W1(t, testCases2, "getFieldColumnName", this.getFieldColumnName, tst.AssertEqual)
	// getStructColumnValue
	nameCol, lvlCol, ageCol := "Name", "Lvl", "Age"
	testCases3 := []tst.P3W2[any, string, string, any, bool]{
		{john, "Person", nameCol, john.Name, true},
		{upc, "School", lvlCol, upc.Level, true},
		{*john, "Person", nameCol, nil, false},
		{nil, "School", lvlCol, nil, false},
		{john, "School", ageCol, nil, false},
		{john, "Person", lvlCol, nil, false},
	}
	tst.AllP3W2(t, testCases3, "getStructColumnValue", this.getStructColumnValue, tst.AssertEqualAny, tst.AssertEqual)
}

func TestInternalFields(t *testing.T) {
	type Person struct {
		Name     string
		Age      int
		password string
	}
	personRef := &Person{}
	this := testPrelude(t, personRef)
	// getFieldName
	p1 := &Person{"John", 18, "john18"}
	p2 := &Person{"Jane", 19, "jan19"}
	testCases := []tst.P1W1[any, string]{
		{&personRef.Name, "Name"}, {&personRef.Age, "Age"},
		{nil, ""},                   // nil fieldRef
		{p1.Name, ""}, {p1.Age, ""}, // non-pointer fieldRef
		{&p1.Name, ""}, {&p2.Age, ""}, // fieldRefs are not from type singleton
	}
	tst.AllP1W1(t, testCases, "getFieldName", this.getFieldName, tst.AssertEqual)
	// getStructFieldValue
	testCases2 := []tst.P2W2[any, string, string, bool]{
		{p1, "Name", p1.Name, true}, // Success
		{*p1, "Name", "", false},    // not a structRef
		{p2, "Job", "", false},      // unknown field
		{p1, "password", "", false}, // private field
		{p2, "Age", "", false},      // wrong field type
	}
	tst.AllP2W2(t, testCases2, "getStructFieldValue", getStructFieldValue[string], tst.AssertEqual[string], tst.AssertEqual[bool])

	age, ok := getStructFieldValue[int](p2, "Age")
	tst.AssertEqualAnd(t, "getStructFieldValue", age, p2.Age, ok, true)
}

func TestGetStructTypedColumnValue(t *testing.T) {
	type User struct {
		Name     string
		Age      int
		Code     string
		Secret   string `col:"-"`
		password string
	}
	u := new(User)
	typeName := "User"
	this := testPrelude(t, u)
	u1 := new(User{"John", 20, "apple", "123", "456"})
	u2 := new(User{"Jane", 21, "banana", "456", "123"})
	testCases := []tst.P3W2[any, string, string, string, bool]{
		{User{}, typeName, "Name", "", false}, // not a struct ref
		{u1, "Person", "Name", "", false},     // unknown type name
		{u1, typeName, "Job", "", false},      // unknown column name
		{u2, typeName, "", "", false},         // blank column name
		{u2, typeName, "Secret", "", false},   // blank column name
		{u1, typeName, "password", "", false}, // private field
		{u1, typeName, "Name", u1.Name, true}, // success
		{u2, typeName, "Code", u2.Code, true}, // success
	}
	bridgeFn := func(structRef any, typeName, columnName string) (string, bool) {
		res := getStructTypedColumnValue[string](this, structRef, typeName, columnName)
		return res.Value(), res.NotError()
	}
	tst.AllP3W2(t, testCases, "getStructTypedColumnValue", bridgeFn, tst.AssertEqual[string], tst.AssertEqual[bool])

	testCases2 := []tst.P3W2[any, string, string, int, bool]{
		{u1, typeName, "Name", 0, false},    // wrong type
		{u2, typeName, "Code", 0, false},    // wrong type
		{u1, typeName, "Age", u1.Age, true}, // success
		{u2, typeName, "Age", u2.Age, true}, // success
	}
	bridgeFn2 := func(structRef any, typeName, columnName string) (int, bool) {
		res := getStructTypedColumnValue[int](this, structRef, typeName, columnName)
		return res.Value(), res.NotError()
	}
	tst.AllP3W2(t, testCases2, "getStructTypedColumnValue", bridgeFn2, tst.AssertEqual[int], tst.AssertEqual[bool])
}

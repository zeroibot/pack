package qb

import (
	"slices"
	"testing"

	"github.com/roidaradal/pack/ds"
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
		Key   string `col:"key"`
		Value string `col:"value"`
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
	// Not a struct
	info := this.readStructColumns(5)
	if !info.IsEmpty() {
		t.Errorf("readStructColumns() = %v, want empty info", info)
	}
	// Not a struct pointer
	info = this.readStructColumns(Person{})
	if !info.IsEmpty() {
		t.Errorf("readStructColumns() = %v, want empty info", info)
	}
	// Successful
	wantCols := ds.List[string]{"`key`", "`value`"}
	info = this.readStructColumns(&Config{})
	if slices.Equal(wantCols, info.columns) == false {
		t.Errorf("readStructColumns() = %v, want %v", info.columns, wantCols)
	}
	// With embedded
	wantCols = ds.List[string]{"`Name`", "`Logo`", "`key`", "`value`", "`Lvl`"}
	info = this.readStructColumns(&School{})
	if slices.Equal(wantCols, info.columns) == false {
		t.Errorf("readStructColumns() = %v, want %v", info.columns, wantCols)
	}
	// With private embedded
	wantCols = ds.List[string]{"`Name`", "`Address`", "`Job`", "`Age`"}
	info = this.readStructColumns(&Person{})
	if slices.Equal(wantCols, info.columns) == false {
		t.Errorf("readStructColumns() = %v, want %v", info.columns, wantCols)
	}
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
	this := NewInstance(MySQL)
	p := &Person{}
	s := &School{}
	cfg := &Config{}
	err := AddType(this, p)
	if err != nil {
		t.Errorf("AddType() error: %v", err)
	}
	err = AddType(this, s)
	if err != nil {
		t.Errorf("AddType() error: %v", err)
	}
	john := &Person{"John", "PH", "dev", 20}
	upc := &School{"UP", "Oblation", "12345", 3, 1000}
	// allColumns
	wantCols := ds.List[string]{"`Name`", "`Address`", "`Job`", "`Age`"}
	actualCols := this.allColumns(p)
	if slices.Equal(actualCols, wantCols) == false {
		t.Errorf("allColumns() = %v, want %v", actualCols, wantCols)
	}
	actualCols = this.allColumns(john)
	if slices.Equal(actualCols, wantCols) == false {
		t.Errorf("allColumns() = %v, want %v", actualCols, wantCols)
	}
	wantCols = ds.List[string]{"`Name`", "`Logo`", "`Lvl`"}
	actualCols = this.allColumns(s)
	if slices.Equal(actualCols, wantCols) == false {
		t.Errorf("allColumns() = %v, want %v", actualCols, wantCols)
	}
	actualCols = this.allColumns(upc)
	if slices.Equal(actualCols, wantCols) == false {
		t.Errorf("allColumns() = %v, want %v", actualCols, wantCols)
	}
	// non-registered type
	wantCols = ds.List[string]{}
	actualCols = this.allColumns(cfg)
	if slices.Equal(actualCols, wantCols) == false {
		t.Errorf("allColumns() = %v, want %v", actualCols, wantCols)
	}
	// getColumnFieldName
	testCases := [][2]string{
		{"Age", this.getColumnFieldName("Person", "`Age`")},
		{"Level", this.getColumnFieldName("School", "`Lvl`")},
		{"", this.getColumnFieldName("Config", "`key`")},
		{"", this.getColumnFieldName("School", "`capacity`")},
		{"", this.getColumnFieldName("School", "`Secret`")},
	}
	for _, x := range testCases {
		want, actual := x[0], x[1]
		if actual != want {
			t.Errorf("getColumnFieldName() = %q, want %q", actual, want)
		}
	}
	// getFieldColumnName
	testCases = [][2]string{
		{"`Age`", this.getFieldColumnName("Person", "Age")},
		{"`Lvl`", this.getFieldColumnName("School", "Level")},
		{"", this.getFieldColumnName("Config", "key")},
		{"", this.getFieldColumnName("School", "capacity")},
		{"", this.getFieldColumnName("School", "Secret")},
	}
	for _, x := range testCases {
		want, actual := x[0], x[1]
		if actual != want {
			t.Errorf("getFieldColumnName() = %q, want %q", actual, want)
		}
	}
	// getStructColumnValue
	type testCase struct {
		want1                any
		want2                bool
		structRef            any
		typeName, columnName string
	}
	nameCol, lvlCol, ageCol := "`Name`", "`Lvl`", "`Age`"
	testCases2 := []testCase{
		{john.Name, true, john, "Person", nameCol},
		{upc.Level, true, upc, "School", lvlCol},
		{nil, false, *john, "Person", nameCol},
		{nil, false, nil, "School", lvlCol},
		{nil, false, john, "School", ageCol},
		{nil, false, john, "Person", lvlCol},
	}
	for _, x := range testCases2 {
		actual1, actual2 := this.getStructColumnValue(x.structRef, x.typeName, x.columnName)
		if actual1 != x.want1 || actual2 != x.want2 {
			t.Errorf("getStructColumnValue() = %v, %t, want %v, %t", actual1, actual2, x.want1, x.want2)
		}
	}
}

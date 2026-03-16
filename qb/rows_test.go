package qb

import (
	"maps"
	"testing"

	"github.com/roidaradal/pack/dict"
)

type mockScanner struct {
	//this      *Instance
	//structRef any
	//items     []any
}

func (m mockScanner) Scan(fieldRefs ...any) error {
	return nil
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
	empty := dict.Object{}
	userObj := dict.Object{"`Name`": "john", "`Password`": "123456", "`Age`": 25}
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
	intOption, err := intReader(mockScanner{})
	if err == nil || intOption.NotNil() {
		t.Errorf("NewRowReader[int] should return an error")
	}
	// Valid full reader
	fullReader := FullRowReader(this, userRef)
	if fullReader == nil {
		t.Errorf("FullRowReader() should return a rowReader, got nil")
	}
	// Successful read
	option, err := fullReader(mockScanner{})
	if err != nil || option.IsNil() {
		t.Errorf("FullRowReader() read = %v, %v, want <User>, nil", option, err)
	}
	// Valid row reader, with specified columns
	nameCol, pwdCol := this.Column(&userRef.Name), this.Column(&userRef.Password)
	rowReader := NewRowReader[User](this, nameCol, pwdCol)
	option, err = rowReader(mockScanner{})
	if err != nil || option.IsNil() {
		t.Errorf("RowReader() read = %v, %v, want <User>, nil", option, err)
	}
	// Error because of blank columns
	userReader := NewRowReader[User](this, nameCol, pwdCol, "")
	option, err = userReader(mockScanner{})
	if err == nil || option.NotNil() {
		t.Errorf("NewRowReader() read = %v, %v, want nil, err", option, err)
	}
	// Error because of unknown column field
	userReader = NewRowReader[User](this, nameCol, pwdCol, "secret")
	option, err = userReader(mockScanner{})
	if err == nil || option.NotNil() {
		t.Errorf("NewRowReader() read = %v, %v, want nil, err", option, err)
	}
	// TODO: Check that item was read to fieldRefs for FullRowReader
	// TODO: Check that item was read to fieldRefs for RowReader with specified columns
}

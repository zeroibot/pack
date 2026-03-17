package qb

import (
	"testing"
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

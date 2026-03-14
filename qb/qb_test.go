package qb

import (
	"fmt"
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
	fmt.Println(this.addressColumns)
	fmt.Println(this.typeColumns)
	fmt.Println(this.typeColumnFields)
	fmt.Println(this.typeFieldColumns)
	// TODO: Finish this test
}

func TestPrepareColumn(t *testing.T) {
	// TODO: DbType.PrepareColumn
}

package qb

import (
	"errors"
	"testing"

	"github.com/roidaradal/pack/dict"
	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/tst"
)

var (
	errMock     = errors.New("mock error")
	errNotFound = errors.New("not found")
)

// Common steps for creating Instance and adding 1 type
func testPrelude[T any](t *testing.T, typeRef *T) *Instance {
	this := NewInstance(MySQL)
	err := AddType(this, typeRef)
	if err != nil {
		t.Errorf("AddType() error = %v", err)
	}
	return this
}

// Common steps for creating Instance and adding 2 types
func testPrelude2[T1, T2 any](t *testing.T, typeRef1 *T1, typeRef2 *T2) *Instance {
	this := NewInstance(MySQL)
	err := AddType(this, typeRef1)
	if err != nil {
		t.Errorf("AddType() error = %v", err)
	}
	err = AddType(this, typeRef2)
	if err != nil {
		t.Errorf("AddType() error = %v", err)
	}
	return this
}

func TestInternalRows(t *testing.T) {
	type User struct {
		Name     string
		Code     string `col:"Username"`
		password string
	}
	nameCol, codeCol, pwdCol, ageCol := "Name", "Username", "password", "age"
	this := testPrelude(t, &User{})
	// newRowCreator
	user := User{"john", "john67", "123456"}
	rowFn1 := this.newRowCreator("User", ds.List[string]{nameCol})
	rowFn2 := this.newRowCreator("User", ds.List[string]{nameCol, pwdCol})
	rowFn3 := this.newRowCreator("User", ds.List[string]{nameCol, ageCol})
	emptyObj := dict.Object{}
	testCases := [][2]dict.Object{
		{rowFn1(&user), dict.Object{"Name": "john"}},
		{rowFn1(user), emptyObj},  // not a struct pointer
		{rowFn2(&user), emptyObj}, // private field
		{rowFn3(&user), emptyObj}, // non-existent field
	}
	tst.All(t, testCases, "newRowCreator", tst.AssertMapEqual)
	// getStructColumnFieldRef
	testCases2 := []tst.P3W2[any, string, string, any, bool]{
		{&user, "User", codeCol, &user.Code, true},
		{&user, "User", nameCol, &user.Name, true},
		{&user, "User", pwdCol, nil, false},
		{&user, "User", ageCol, nil, false},
	}
	tst.AllP3W2(t, testCases2, "getStructColumnFieldRef", this.getStructColumnFieldRef, tst.AssertEqualAny, tst.AssertEqual)
}

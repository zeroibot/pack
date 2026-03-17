package qb

import (
	"maps"
	"testing"

	"github.com/roidaradal/pack/dict"
	"github.com/roidaradal/pack/ds"
)

func TestInternalRows(t *testing.T) {
	type User struct {
		Name     string
		Code     string `col:"Username"`
		password string
	}
	this := NewInstance(MySQL)
	nameCol := "Name"
	codeCol := "Username"
	pwdCol := "password"
	ageCol := "age"
	err := AddType(this, &User{})
	if err != nil {
		t.Errorf("AddType() error = %v", err)
	}
	// newRowCreator
	user := User{"john", "john67", "123456"}
	rowFn1 := this.newRowCreator("User", ds.List[string]{nameCol})
	rowFn2 := this.newRowCreator("User", ds.List[string]{nameCol, pwdCol})
	rowFn3 := this.newRowCreator("User", ds.List[string]{nameCol, ageCol})
	emptyObj := dict.Object{}
	testCases := [][2]dict.Object{
		{dict.Object{"Name": "john"}, rowFn1(&user)},
		{emptyObj, rowFn1(user)},  // not a struct pointer
		{emptyObj, rowFn2(&user)}, // private field
		{emptyObj, rowFn3(&user)}, // non-existent field
	}
	for _, x := range testCases {
		wantObj, actualObj := x[0], x[1]
		if maps.Equal(wantObj, actualObj) == false {
			t.Errorf("newRowCreator() = %v, want %v", actualObj, wantObj)
		}
	}
	// getStructColumnFieldRef
	type testCase struct {
		column string
		want1  any
		want2  bool
	}
	testCases2 := []testCase{
		{codeCol, &user.Code, true},
		{nameCol, &user.Name, true},
		{pwdCol, nil, false},
		{ageCol, nil, false},
	}
	for _, x := range testCases2 {
		actual1, actual2 := this.getStructColumnFieldRef(&user, "User", x.column)
		if x.want1 != actual1 || x.want2 != actual2 {
			t.Errorf("getStructColumnFieldRef() = %v, %t, want %v, %t", actual1, actual2, x.want1, x.want2)
		}
	}
}

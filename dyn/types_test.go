package dyn

import (
	"testing"

	"github.com/roidaradal/tst"
)

func TestTypeCheck(t *testing.T) {
	type person struct {
		name string
		age  int
	}
	x, y, z := 5, "abc", person{"john", 20}
	xp, yp, zp := new(x), new(y), new(z)

	// IsPointer
	testCases := []tst.P1W1[any, bool]{
		{x, false}, {y, false}, {z, false},
		{xp, true}, {yp, true}, {zp, true},
	}
	tst.AllP1W1(t, testCases, "IsPointer", IsPointer, tst.AssertEqual)

	// IsStruct
	testCases = []tst.P1W1[any, bool]{
		{x, false}, {y, false}, {z, true},
	}
	tst.AllP1W1(t, testCases, "IsStruct", IsStruct, tst.AssertEqual)

	// IsStructPointer
	testCases = []tst.P1W1[any, bool]{
		{x, false}, {y, false}, {z, false},
		{xp, false}, {yp, false}, {zp, true},
	}
	tst.AllP1W1(t, testCases, "IsStructPointer", IsStructPointer, tst.AssertEqual)
}

func TestTypeName(t *testing.T) {
	type person struct {
		name string
		age  int
	}
	type company struct {
		name    string
		address string
	}
	p := person{"john", 20}
	c := company{"google", "usa"}

	// TypeName
	testCases := []tst.P1W1[any, string]{
		{5, "int"}, {2.5, "float64"}, {true, "bool"}, {"abc", "string"},
		{p, "person"}, {c, "company"}, {&p, "person"}, {&c, "company"},
	}
	tst.AllP1W1(t, testCases, "TypeName", TypeName, tst.AssertEqual)

	// FullTypeName
	testCases = []tst.P1W1[any, string]{
		{5, "int"}, {2.5, "float64"}, {false, "bool"}, {"abc", "string"},
		{new(5), "*int"}, {new(2.5), "*float64"}, {new("abc"), "*string"},
		{p, "person"}, {c, "company"}, {&p, "*person"}, {&c, "*company"},
	}
	tst.AllP1W1(t, testCases, "FullTypeName", FullTypeName, tst.AssertEqual)
}

package dyn

import (
	"testing"
)

func TestTypeCheck(t *testing.T) {
	type testCase struct {
		name         string
		want, actual bool
	}
	type person struct {
		name string
		age  int
	}
	x, y, z := 5, "abc", person{"john", 20}
	xp, yp, zp := new(x), new(y), new(z)
	testCases := []testCase{
		{"IsPointer", false, IsPointer(x)},
		{"IsPointer", false, IsPointer(y)},
		{"IsPointer", false, IsPointer(z)},
		{"IsPointer", true, IsPointer(xp)},
		{"IsPointer", true, IsPointer(yp)},
		{"IsPointer", true, IsPointer(zp)},
		{"IsStruct", false, IsStruct(x)},
		{"IsStruct", false, IsStruct(y)},
		{"IsStruct", true, IsStruct(z)},
		{"IsStructPointer", false, IsStructPointer(x)},
		{"IsStructPointer", false, IsStructPointer(y)},
		{"IsStructPointer", false, IsStructPointer(z)},
		{"IsStructPointer", false, IsStructPointer(xp)},
		{"IsStructPointer", false, IsStructPointer(yp)},
		{"IsStructPointer", true, IsStructPointer(zp)},
	}
	for _, c := range testCases {
		if c.want != c.actual {
			t.Errorf("%s() = %v, want %v", c.name, c.actual, c.want)
		}
	}
}

func TestTypeName(t *testing.T) {
	type testCase struct {
		name         string
		want, actual string
	}
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
	testCases := []testCase{
		{"TypeName", "int", TypeName(5)},
		{"FullTypeName", "int", FullTypeName(5)},
		{"TypeName", "float64", TypeName(2.5)},
		{"FullTypeName", "float64", FullTypeName(2.5)},
		{"TypeName", "bool", TypeName(true)},
		{"FullTypeName", "bool", FullTypeName(false)},
		{"TypeName", "string", TypeName("abc")},
		{"FullTypeName", "string", FullTypeName("abc")},
		{"FullTypeName", "*string", FullTypeName(new("abc"))},
		{"TypeName", "person", TypeName(p)},
		{"FullTypeName", "person", FullTypeName(p)},
		{"TypeName", "company", TypeName(c)},
		{"FullTypeName", "company", FullTypeName(c)},
		{"TypeName", "person", TypeName(&p)},
		{"FullTypeName", "*person", FullTypeName(&p)},
		{"TypeName", "company", TypeName(&c)},
		{"FullTypeName", "*company", FullTypeName(&c)},
	}
	for _, x := range testCases {
		if x.want != x.actual {
			t.Errorf("%s() = %v, want %v", x.name, x.actual, x.want)
		}
	}
}

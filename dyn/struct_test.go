package dyn

import (
	"reflect"
	"testing"

	"github.com/roidaradal/tst"
)

type testCase struct {
	structRef any
	field     string
	value     any
	test      func() bool
}

func (tc testCase) MustSetStructField() {
	MustSetStructField(tc.structRef, tc.field, tc.value)
}

func (tc testCase) PostTest() bool {
	if tc.test == nil {
		return true
	}
	return tc.test()
}

func TestSetStructField(t *testing.T) {
	type person struct {
		Name     string
		Age      int
		Weight   float64
		password string
	}
	p := person{"John", 20, 67.9, "secret"}
	items := []int{1, 2, 3}

	// SetStructField
	testCases := []tst.P3W1Post[any, string, any, bool]{
		// Successful
		{&p, "Name", "Johnny", true, func() bool { return p.Name == "Johnny" }},
		{&p, "Age", 25, true, func() bool { return p.Age == 25 }},
		{&p, "Weight", 69.7, true, func() bool { return p.Weight == 69.7 }},
		// Not a struct pointer
		{p, "Name", "Jar", false, nil},
		{25, "Age", 25, false, nil},
		{&items, "Weight", 60.5, false, nil},
		// Non-existent fields
		{&p, "Job", "Dev", false, nil},
		{&p, "Password", "secret", false, nil},
		// Private field property
		{&p, "password", "secret123", false, func() bool { return p.password == "secret" }},
		// Wrong field type
		{&p, "Name", 50, false, func() bool { return p.Name == "Johnny" }},
		{&p, "Weight", "Skinny", false, func() bool { return p.Weight == 69.7 }},
	}
	tst.AllP3W1Post(t, testCases, "SetStructField", SetStructField, tst.AssertEqual)

	// MustSetStructField
	testCases2 := []testCase{
		// Not a struct pointer
		{5, "Age", 5, nil},
		{"Name", "Name", "MyName", nil},
		// Success
		{&p, "Name", "Jane", func() bool { return p.Name == "Jane" }},
		{&p, "Age", 30, func() bool { return p.Age == 30 }},
		{&p, "Weight", 55.5, func() bool { return p.Weight == 55.5 }},
	}
	tst.AllActionPost(t, testCases2, "MustSetStructField", testCase.MustSetStructField)

	defer tst.AssertPanic(t, "MustSetStructField")
	MustSetStructField(&p, "password", "unknown") // should panic
}

func TestGetStructField(t *testing.T) {
	type person struct {
		Name     string
		Age      int
		Weight   float64
		password string
	}
	p := person{"John", 20, 67.9, "secret"}
	// GetStructFieldAs
	testCases1 := []tst.P2W2[any, string, string, bool]{
		{&p, "Name", p.Name, true},
		{p, "Name", "", false},
		{&p, "Job", "", false},
		{&p, "password", "", false},
	}
	tst.AllP2W2(t, testCases1, "GetStructFieldAs[string]", GetStructFieldAs[string], tst.AssertEqual[string], tst.AssertEqual[bool])
	testCases2 := []tst.P2W2[any, string, int, bool]{
		{&p, "Age", p.Age, true},
		{&p, "Name", 0, false},
		{&p, "Weight", 0, false},
	}
	tst.AllP2W2(t, testCases2, "GetStructFieldAs[int]", GetStructFieldAs[int], tst.AssertEqual[int], tst.AssertEqual[bool])

	// GetStructFieldAsString
	testCases3 := []tst.P2W2[any, string, string, bool]{
		{&p, "Name", "John", true},
		{&p, "Age", "20", true},
		{&p, "Weight", "67.9", true},
		{p, "Name", "<nil>", false},
		{&p, "Job", "<nil>", false},
		{&p, "password", "<nil>", false},
	}
	tst.AllP2W2(t, testCases3, "GetStructFieldAsString", GetStructFieldAsString, tst.AssertEqual[string], tst.AssertEqual[bool])
}

func TestMustGetStructField(t *testing.T) {
	type person struct {
		Name   string
		Age    int
		Weight float64
	}
	p := person{"John", 20, 55.5}

	// MustGetStructField
	testCases := []tst.P2W1[any, string, any]{
		{p, "Name", nil},
		{&p, "Name", p.Name},
		{&p, "Age", p.Age},
	}
	tst.AllP2W1(t, testCases, "MustGetStructField", MustGetStructField, tst.AssertEqual)

	// MustGetStructFieldAsString
	testCases2 := []tst.P2W1[any, string, string]{
		{p, "Name", "<nil>"},
		{&p, "Age", "20"},
		{&p, "Weight", "55.5"},
	}
	tst.AllP2W1(t, testCases2, "MustGetStructFieldAsString", MustGetStructFieldAsString, tst.AssertEqual)

	// MustGetStructField panic
	defer tst.AssertPanic(t, "MustGetStructField")
	MustGetStructField(&p, "Job") // should panic
}

func TestMustGetStructField2(t *testing.T) {
	type person struct {
		password string
	}
	p := person{"abc"}
	defer tst.AssertPanic(t, "MustGetStructField")
	MustGetStructField(&p, "password") // should panic
}

func TestGetStructFieldTag(t *testing.T) {
	type person struct {
		Name     string `col:"name" json:"username"`
		Age      int    `col:"age"`
		Password string `col:"pass" json:"-"`
		Secret   string
	}
	p := &person{"John", 20, "abc", "secret"}
	colKey, jsonKey := "col", "json"
	structValue := MustDerefValue(p)
	structType := structValue.Type()
	nameField := structType.Field(0)
	ageField := structType.Field(1)
	passwordField := structType.Field(2)
	secretField := structType.Field(3)
	testCases := []tst.P2W2[reflect.StructField, string, string, bool]{
		{nameField, colKey, "name", true},
		{nameField, jsonKey, "username", true},
		{ageField, colKey, "age", true},
		{ageField, jsonKey, "", false},
		{passwordField, colKey, "pass", true},
		{passwordField, jsonKey, "-", true},
		{secretField, colKey, "", false},
		{secretField, jsonKey, "", false},
	}
	tst.AllP2W2[reflect.StructField, string, string, bool](t, testCases, "GetStructFieldTag", GetStructFieldTag, tst.AssertEqual[string], tst.AssertEqual[bool])
}

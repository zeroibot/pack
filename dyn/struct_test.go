package dyn

import (
	"reflect"
	"testing"

	"github.com/roidaradal/tst"
)

func TestSetStructField(t *testing.T) {
	type person struct {
		Name     string
		Age      int
		Weight   float64
		password string
	}
	type testCase struct {
		want, actual bool
		postTest     func() bool
	}
	p := person{"John", 20, 67.9, "secret"}
	items := []int{1, 2, 3}
	testCases := []testCase{
		// Successful
		{true, SetStructField(&p, "Name", "Johnny"), func() bool {
			return p.Name == "Johnny"
		}},
		{true, SetStructField(&p, "Age", 25), func() bool {
			return p.Age == 25
		}},
		{true, SetStructField(&p, "Weight", 69.7), func() bool {
			return p.Weight == 69.7
		}},
		// Not a struct pointer
		{false, SetStructField(p, "Name", "Jar"), nil},
		{false, SetStructField(25, "Age", 25), nil},
		{false, SetStructField(&items, "Weight", 60.5), nil},
		// Non-existent fields
		{false, SetStructField(&p, "Job", "Dev"), nil},
		{false, SetStructField(&p, "Password", "secret"), nil},
		// Private field property
		{false, SetStructField(&p, "password", "secret123"), func() bool { return p.password == "secret" }},
		// Wrong field type
		{false, SetStructField(&p, "Name", 50), func() bool { return p.Name == "Johnny" }},
		{false, SetStructField(&p, "Weight", "Skinny"), func() bool { return p.Weight == 69.7 }},
	}
	for _, x := range testCases {
		if x.actual != x.want {
			t.Errorf("SetStructField() = %t, want %t", x.actual, x.want)
		}
		if x.postTest != nil && x.postTest() == false {
			t.Errorf("SetStructField() PostTest failed")
		}
	}
	type testCase2 struct {
		structRef any
		field     string
		value     any
		postTest  func() bool
	}
	// MustSetStructField
	testCases2 := []testCase2{
		// Not a struct pointer
		{5, "Age", 5, nil},
		{"Name", "Name", "MyName", nil},
		// Success
		{&p, "Name", "Jane", func() bool { return p.Name == "Jane" }},
		{&p, "Age", 30, func() bool { return p.Age == 30 }},
		{&p, "Weight", 55.5, func() bool { return p.Weight == 55.5 }},
	}
	for _, x := range testCases2 {
		MustSetStructField(x.structRef, x.field, x.value)
		if x.postTest != nil && x.postTest() == false {
			t.Errorf("MustSetStructField PostTest failed")
		}
	}

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
	name, ok := GetStructFieldAs[string](&p, "Name")
	if !ok || name != p.Name {
		t.Errorf("GetStructFieldAs() = %v, %t, want %v, true", name, ok, p.Name)
	}
	age, ok := GetStructFieldAs[int](&p, "Age")
	if !ok || age != p.Age {
		t.Errorf("GetStructFieldAs() = %v, %t, want %v, true", age, ok, p.Age)
	}
	name, ok = GetStructFieldAs[string](p, "Name")
	if ok || name != "" {
		t.Errorf("GetStructFieldAs() = %q, %t, want %q, false", name, ok, "")
	}
	job, ok := GetStructFieldAs[string](&p, "Job")
	if ok || job != "" {
		t.Errorf("GetStructFieldAs() = %q, %t, want %q, false", job, ok, "")
	}
	pwd, ok := GetStructFieldAs[string](&p, "password")
	if ok || pwd != "" {
		t.Errorf("GetStructFieldAs() = %q, %t, want %q, false", pwd, ok, "")
	}
	// GetStructFieldAsString
	type testCase struct {
		structRef any
		field     string
		want1     string
		want2     bool
	}
	testCases := []testCase{
		{&p, "Name", "John", true},
		{&p, "Age", "20", true},
		{&p, "Weight", "67.9", true},
		{p, "Name", "<nil>", false},
		{&p, "Job", "<nil>", false},
		{&p, "password", "<nil>", false},
	}
	for _, x := range testCases {
		actual1, actual2 := GetStructFieldAsString(x.structRef, x.field)
		if actual1 != x.want1 || actual2 != x.want2 {
			t.Errorf("GetStructFieldAsString() = %s, %t, want %s, %t", actual1, actual2, x.want1, x.want2)
		}
	}
}

func TestMustGetStructField(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustGetStructField() did not panic")
		}
	}()
	type person struct {
		Name   string
		Age    int
		Weight float64
	}
	p := person{"John", 20, 55.5}
	none := MustGetStructField(p, "Name")
	if none != nil {
		t.Errorf("MustGetStructField() = %v, want nil", none)
	}
	noneString := MustGetStructFieldAsString(p, "Name")
	if noneString != "<nil>" {
		t.Errorf("MustGetStructFieldAsString() = %q, want <nil>", noneString)
	}
	name := MustGetStructField(&p, "Name")
	if name != p.Name {
		t.Errorf("MustGetStructField() = %s, want %s", name, p.Name)
	}
	age := MustGetStructField(&p, "Age")
	if age != p.Age {
		t.Errorf("MustGetStructField() = %d, want %d", age, p.Age)
	}
	ageString := MustGetStructFieldAsString(&p, "Age")
	if ageString != "20" {
		t.Errorf("MustGetStructFieldAsString() = %q, want 20", ageString)
	}
	weightString := MustGetStructFieldAsString(&p, "Weight")
	if weightString != "55.5" {
		t.Errorf("MustGetStructFieldAsString() = %q, want 55.5", weightString)
	}
	MustGetStructField(&p, "Job") // should panic
}

func TestMustGetStructField2(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustGetStructField() did not panic")
		}
	}()
	type person struct {
		password string
	}
	p := person{"abc"}
	MustGetStructField(&p, "password") // should panic
}

func TestGetStructFieldTag(t *testing.T) {
	type person struct {
		Name     string `col:"name" json:"username"`
		Age      int    `col:"age"`
		Password string `col:"pass" json:"-"`
		Secret   string
	}
	type testCase struct {
		structField reflect.StructField
		key         string
		want1       string
		want2       bool
	}
	p := &person{"John", 20, "abc", "secret"}
	colKey, jsonKey := "col", "json"
	structValue := MustDerefValue(p)
	structType := structValue.Type()
	nameField := structType.Field(0)
	ageField := structType.Field(1)
	passwordField := structType.Field(2)
	secretField := structType.Field(3)
	testCases := []testCase{
		{nameField, colKey, "name", true},
		{nameField, jsonKey, "username", true},
		{ageField, colKey, "age", true},
		{ageField, jsonKey, "", false},
		{passwordField, colKey, "pass", true},
		{passwordField, jsonKey, "-", true},
		{secretField, colKey, "", false},
		{secretField, jsonKey, "", false},
	}
	for _, x := range testCases {
		actual1, actual2 := GetStructFieldTag(x.structField, x.key)
		if actual1 != x.want1 || actual2 != x.want2 {
			t.Errorf("GetStructFieldTag() = %s, %t, want %s, %t", actual1, actual2, x.want1, x.want2)
		}
	}
}

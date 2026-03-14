package dyn

import (
	"testing"
)

func TestSetStructField(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("MustSetStructField() did not panic")
		}
	}()
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

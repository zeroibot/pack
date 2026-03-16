package qb

import "testing"

func TestInternalFields(t *testing.T) {
	type Person struct {
		Name     string
		Age      int
		password string
	}
	personRef := &Person{}
	this := NewInstance(MySQL)
	err := AddType(this, personRef)
	if err != nil {
		t.Errorf("AddType() = %v", err)
	}
	// getFieldName
	p1 := &Person{"John", 18, "john18"}
	p2 := &Person{"Jane", 19, "jan19"}
	testCases := [][2]string{
		{"Name", this.getFieldName(&personRef.Name)},
		{"Age", this.getFieldName(&personRef.Age)},
		{"", this.getFieldName(nil)},     // nil fieldRef
		{"", this.getFieldName(p1.Name)}, // non-pointer fieldRefs
		{"", this.getFieldName(p1.Age)},
		{"", this.getFieldName(&p1.Name)}, // fieldRefs are not from type singleton
		{"", this.getFieldName(&p2.Age)},
	}
	for _, x := range testCases {
		want, actual := x[0], x[1]
		if want != actual {
			t.Errorf("getFieldName() = %q, want %q", actual, want)
		}
	}
	// getStructFieldValue
	name, ok := getStructFieldValue[string](p1, "Name")
	if name != p1.Name || !ok {
		t.Errorf("getStructFieldValue() = %q, %t, want %q, true", name, ok, p1.Name)
	}
	age, ok := getStructFieldValue[int](p2, "Age")
	if age != p2.Age || !ok {
		t.Errorf("getStructFieldValue() = %d, %t, want %d, true", age, ok, p2.Age)
	}
	// Fail to Deref
	name, ok = getStructFieldValue[string](*p1, "Name")
	if name != "" || ok {
		t.Errorf("getStructFieldValue() = %q, %t, want %q, false", name, ok, "")
	}
	// Unknown field
	job, ok := getStructFieldValue[string](p2, "Job")
	if job != "" || ok {
		t.Errorf("getStructFieldValue() = %q, %t, want %q, false", job, ok, "")
	}
	// Private field
	pwd, ok := getStructFieldValue[string](p1, "password")
	if pwd != "" || ok {
		t.Errorf("getStructFieldValue() = %q, %t, want %q, false", pwd, ok, "")
	}
	// Wrong field type
	ageStr, ok := getStructFieldValue[string](p2, "Age")
	if ageStr != "" || ok {
		t.Errorf("getStructFieldValue() = %q, %t, want %q, false", ageStr, ok, "")
	}
}

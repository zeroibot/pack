package qb

import (
	"reflect"
	"slices"
	"testing"

	"github.com/roidaradal/pack/ds"
)

func TestKVCondition(t *testing.T) {
	type Person struct {
		Name    string
		Age     int
		Job     string
		Level   int    `col:"Lvl"`
		Details string `col:"-"`
		secret  string
		IP      *string
	}
	this := NewInstance(MySQL)
	p := &Person{}
	err := AddType(this, p)
	if err != nil {
		t.Errorf("AddType error: %v", err)
	}
	// newColumnValue
	type testCase1 struct {
		isNil  bool
		want   columnValuePair
		actual ds.Option[columnValuePair]
	}
	var empty1 columnValuePair
	testCases1 := []testCase1{
		{false, columnValuePair{V1: "`Lvl`", V2: 5}, newColumnValue(this, &p.Level, 5)},
		{false, columnValuePair{V1: "`Job`", V2: "dev"}, newColumnValue(this, &p.Job, "dev")},
		{true, empty1, newColumnValue(this, &p.secret, "123")},
		{true, empty1, newColumnValue(this, &p.Details, "abc")},
	}
	for _, x := range testCases1 {
		if x.actual.IsNil() != x.isNil || x.actual.Value() != x.want {
			t.Errorf("newColumnValue() = %v, want %v", x.actual, x.want)
		}
	}

	// newFieldColumnValue
	testCases1 = []testCase1{
		{false, columnValuePair{V1: "`Name`", V2: "John"}, newFieldColumnValue(this, "Person", "Name", "John")},
		{false, columnValuePair{V1: "`Age`", V2: 20}, newFieldColumnValue(this, "Person", "Age", 20)},
		// Note: still works even if Name (string) is paired with 20 (int), because the columnValuePair accepts `any` value
		{false, columnValuePair{V1: "`Name`", V2: 20}, newFieldColumnValue(this, "Person", "Name", 20)},
		{true, empty1, newFieldColumnValue(this, "Account", "Name", "John")},   // wrong typeName
		{true, empty1, newFieldColumnValue(this, "Person", "secret", "hello")}, // private field
		{true, empty1, newFieldColumnValue(this, "Person", "Address", "PH")},   // unknown field
	}
	for _, x := range testCases1 {
		if x.actual.IsNil() != x.isNil || x.actual.Value() != x.want {
			t.Errorf("newFieldColumnValue() = %v, want %v", x.actual, x.want)
		}
	}

	// newColumnValueList
	var empty2 columnValueListPair
	type testCase2 struct {
		isNil  bool
		want   columnValueListPair
		actual ds.Option[columnValueListPair]
	}
	testCases2 := []testCase2{
		{false, columnValueListPair{V1: "`Name`", V2: []any{"John", "Jane"}}, newColumnValueList(this, &p.Name, ds.List[string]{"John", "Jane"})},
		{false, columnValueListPair{V1: "`Job`", V2: []any{"dev", "qa"}}, newColumnValueList(this, &p.Job, ds.List[string]{"dev", "qa"})},
		{true, empty2, newColumnValueList(this, &p.secret, ds.List[string]{"123", "456"})},
		{true, empty2, newColumnValueList(this, &p.Details, ds.List[string]{"abc", "def"})},
	}
	for _, x := range testCases2 {
		if x.actual.IsNil() != x.isNil || reflect.DeepEqual(x.actual.Value(), x.want) == false {
			t.Errorf("newColumnValueList() = %v, want %v", x.actual, x.want)
		}
	}

	// falseConditionValues
	emptyValues := make([]any, 0)
	wantValues := emptyValues
	actualCond, actualValues := falseConditionValues()
	if actualCond != "false" || slices.Equal(actualValues, wantValues) == false {
		t.Errorf("falseConditionValues() = %v, %v, want false, []", actualCond, actualValues)
	}

	// trueConditionValues
	actualCond, actualValues = trueConditionValues()
	if actualCond != "true" || slices.Equal(actualValues, wantValues) == false {
		t.Errorf("trueConditionValues() = %v, %v, want true, []", actualCond, actualValues)
	}

	// soloConditionValues

	// normal operators
	kv1 := newColumnValue(this, &p.Name, "John").Value()
	kv2 := newColumnValue(this, &p.Age, 20).Value()
	kv3 := newColumnValue(this, &p.Level, 5).Value()
	kv4 := newColumnValue(this, &p.IP, nil).Value()
	kv5 := newColumnValue(this, &p.Job, "Dev").Value()

	type testCase3 struct {
		wantCond   string
		wantValues []any
		kv         columnValuePair
		op         operator
	}
	testCases3 := []testCase3{
		{"`Name` = ?", []any{"John"}, kv1, opEqual},
		{"`Age` > ?", []any{20}, kv2, opGreater},
		{"`Lvl` <> ?", []any{5}, kv3, opNotEqual},
		{"`IP` IS NULL", emptyValues, kv4, opEqual},
		{"`IP` IS NOT NULL", emptyValues, kv4, opNotEqual},
		{"`Name` LIKE ?", []any{"John%"}, kv1, opPrefix},
		{"`Name` LIKE ?", []any{"%John"}, kv1, opSuffix},
		{"`Job` LIKE ?", []any{"%Dev%"}, kv5, opSubstring},
	}
	for _, x := range testCases3 {
		cond, values := soloConditionValues(x.kv.V1, x.op, x.kv.V2)
		if cond != x.wantCond || slices.Equal(values, x.wantValues) == false {
			t.Errorf("soloConditionValues() = %q, %v, want %q, %v", cond, values, x.wantCond, x.wantValues)
		}
	}
}

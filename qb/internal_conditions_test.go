package qb

import (
	"reflect"
	"slices"
	"testing"

	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/list"
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

func TestInternalConditions(t *testing.T) {
	type Person struct {
		Name    string
		Age     int `col:"age"`
		Job     string
		Details string `col:"-"`
	}
	this := NewInstance(MySQL)
	p := &Person{}
	err := AddType(this, p)
	if err != nil {
		t.Errorf("AddType error: %v", err)
	}
	emptyValues := make([]any, 0)
	// missingCondition.BuildCondition()
	cond1 := missingCondition{}
	wantCond := "false"
	actualCond, actualValues := cond1.BuildCondition()
	if actualCond != wantCond || slices.Equal(actualValues, emptyValues) == false {
		t.Errorf("missingCondition.BuildCondition() = %q, %v, want %q, %v", actualCond, actualValues, wantCond, emptyValues)
	}
	// matchAllCondition.BuildCondition()
	cond2 := matchAllCondition{}
	wantCond = "true"
	actualCond, actualValues = cond2.BuildCondition()
	if actualCond != wantCond || slices.Equal(actualValues, emptyValues) == false {
		t.Errorf("matchAllCondition.BuildCondition() = %q, %v, want %q, %v", actualCond, actualValues, wantCond, emptyValues)
	}
	// newValueCondition
	valueCond1 := newValueCondition(this, &p.Name, "John", opNotEqual)
	valueCond2 := newValueCondition(this, &p.Age, 20, opGreater)
	valueCond3 := newValueCondition(this, &p.Details, "Info", opSubstring)
	isNils := []bool{false, false, true}
	for i, valueCond := range []valueCondition{valueCond1, valueCond2, valueCond3} {
		if valueCond.pair.IsNil() != isNils[i] {
			t.Errorf("newValueCondition() = %v, want nil = %t", valueCond.pair, isNils[i])
		}
	}
	// valueCondition.BuildCondition()
	valueCond4 := valueCondition{
		pair:     ds.NewOption(new(columnValuePair{V1: "", V2: 0})),
		operator: opEqual,
	}
	type testCase1 struct {
		wantCond   string
		wantValues []any
		valueCond  valueCondition
	}
	testCases1 := []testCase1{
		{"`Name` <> ?", []any{"John"}, valueCond1},
		{"`age` > ?", []any{20}, valueCond2},
		{"false", emptyValues, valueCond3},
		{"false", emptyValues, valueCond4},
	}
	for _, x := range testCases1 {
		actualCond, actualValues = x.valueCond.BuildCondition()
		if actualCond != x.wantCond || slices.Equal(actualValues, x.wantValues) == false {
			t.Errorf("valueCondition.BuildCondition() = %q, %v, want %q, %v", actualCond, actualValues, x.wantCond, x.wantValues)
		}
	}
	// newListCondition
	listCond1 := newListCondition(this, &p.Age, ds.List[int]{18, 19, 20}, opIn, opEqual)
	listCond2 := newListCondition(this, &p.Job, ds.List[string]{"dev", "qa"}, opNotIn, opNotEqual)
	listCond3 := newListCondition(this, &p.Name, ds.List[string]{"John"}, opIn, opEqual)
	listCond4 := newListCondition(this, &p.Details, ds.List[string]{"a", "b", "c"}, opNotIn, opNotEqual)
	listCond5 := newListCondition(this, &p.Name, ds.List[string]{}, opIn, opEqual)
	isNils = []bool{false, false, false, true, false}
	for i, listCond := range []listCondition{listCond1, listCond2, listCond3, listCond4, listCond5} {
		if listCond.pair.IsNil() != isNils[i] {
			t.Errorf("newListCondition() = %v, want nil = %t", listCond.pair, isNils[i])
		}
	}
	// listCondition.BuildCondition()
	type testCase2 struct {
		wantCond   string
		wantValues []any
		listCond   listCondition
	}
	listCond6 := listCondition{
		pair:         ds.NewOption(new(columnValueListPair{V1: "", V2: []any{1, 2, 3}})),
		listOperator: opIn,
		soloOperator: opEqual,
	}
	testCases2 := []testCase2{
		{"`age` IN (?, ?, ?)", []any{18, 19, 20}, listCond1},
		{"`Job` NOT IN (?, ?)", []any{"dev", "qa"}, listCond2},
		{"`Name` = ?", []any{"John"}, listCond3},
		{"false", emptyValues, listCond4},
		{"false", emptyValues, listCond5},
		{"false", emptyValues, listCond6},
	}
	for _, x := range testCases2 {
		actualCond, actualValues = x.listCond.BuildCondition()
		if actualCond != x.wantCond || slices.Equal(actualValues, x.wantValues) == false {
			t.Errorf("listCondition.BuildCondition() = %q, %v, want %q, %v", actualCond, actualValues, x.wantCond, x.wantValues)
		}
	}
	// newMultiCondition
	multiCond1 := newMultiCondition(opAnd, valueCond1, valueCond2)
	multiCond2 := newMultiCondition(opOr, listCond1, listCond2)
	multiCond3 := newMultiCondition(opOr, valueCond1)
	multiCond4 := newMultiCondition(opAnd)
	multiCond5 := newMultiCondition(opOr, multiCond1, multiCond2, nil)
	multiCond6 := newMultiCondition(opOr, valueCond2, multiCond4)
	// multiCondition.BuildCondition()
	type testCase3 struct {
		wantCond   string
		wantValues []any
		multiCond  multiCondition
	}
	testCases3 := []testCase3{
		{"(`Name` <> ? AND `age` > ?)", []any{"John", 20}, multiCond1},
		{"(`age` IN (?, ?, ?) OR `Job` NOT IN (?, ?))", []any{18, 19, 20, "dev", "qa"}, multiCond2},
		{"`Name` <> ?", []any{"John"}, multiCond3},
		{"false", emptyValues, multiCond4},
		{"((`Name` <> ? AND `age` > ?) OR (`age` IN (?, ?, ?) OR `Job` NOT IN (?, ?)))", []any{"John", 20, 18, 19, 20, "dev", "qa"}, multiCond5},
		{"false", emptyValues, multiCond6},
	}
	for _, x := range testCases3 {
		actualCond, actualValues = x.multiCond.BuildCondition()
		if actualCond != x.wantCond || slices.Equal(actualValues, x.wantValues) == false {
			t.Errorf("multiCondition.BuildCondition() = %q, %v, want %q, %v", actualCond, actualValues, x.wantCond, x.wantValues)
		}
	}
}

func TestInternalCombos(t *testing.T) {
	type Person struct {
		Name    string
		Age     int `col:"age"`
		Job     string
		Details string `col:"-"`
	}
	this := NewInstance(MySQL)
	p := &Person{}
	err := AddType(this, p)
	if err != nil {
		t.Errorf("AddType error: %v", err)
	}
	p1 := Person{"John", 18, "dev", "regular"}
	p2 := Person{"James", 20, "qa", "prob"}
	p3 := Person{"Jill", 25, "admin", "regular"}
	p4 := Person{"Juno", 23, "dev", "prob"}
	p5 := Person{"Jack", 21, "sales", "intern"}
	people := ds.List[Person]{p1, p2, p3, p4, p5}
	// missingCombo.Test
	combo1 := missingCombo[Person]{}
	allFalse := people.All(func(person Person) bool {
		return combo1.Test(person) == false
	})
	if !allFalse {
		t.Errorf("Not all missingCombo.Test() returned false")
	}
	// matchAllCombo.Test
	combo2 := matchAllCombo[Person]{}
	allTrue := people.All(func(person Person) bool {
		return combo2.Test(person) == true
	})
	if !allTrue {
		t.Errorf("Not all matchAllCombo.Test() returned true")
	}
	// newValueCombo, valueCombo.Test
	test1 := createFieldValueTest[Person]("Name", func(name string) bool { return name == "Jill" })
	valueCond1 := newValueCondition(this, &p.Name, "Jill", opEqual)
	valueCombo1 := newValueCombo(valueCond1, test1)
	wantBools := []bool{false, false, true, false, false}
	actualBools := list.Map(people, valueCombo1.Test)
	if slices.Equal(wantBools, actualBools) == false {
		t.Errorf("valueCombo.Test() = %v, want %v", actualBools, wantBools)
	}
	test0 := createFieldValueTest[Person]("Unknown", func(unknown string) bool { return unknown != "A" })
	valueCombo2 := newValueCombo(valueCond1, test0)
	allFalse = people.All(func(person Person) bool {
		return valueCombo2.Test(person) == false
	})
	if !allFalse {
		t.Errorf("Not all valueCombo.Test() returned false")
	}
	// newListCombo, listCombo.Test
	jobs := []string{"dev", "qa"}
	test2 := createFieldValueTest[Person]("Job", func(job string) bool { return list.Has(jobs, job) })
	listCond1 := newListCondition(this, &p.Job, jobs, opIn, opEqual)
	listCombo1 := newListCombo(listCond1, test2)
	wantBools = []bool{true, true, false, true, false}
	actualBools = list.Map(people, listCombo1.Test)
	if slices.Equal(wantBools, actualBools) == false {
		t.Errorf("listCombo.Test() = %v, want %v", actualBools, wantBools)
	}
	// newMultiCombo, multiCombo.Test
	multiTest := func(person Person) bool { return test1(person) || test2(person) }
	multiCombo1 := newMultiCombo(ds.List[DualCondition[Person]]{valueCombo1, listCombo1}, opOr, multiTest)
	wantBools = []bool{true, true, true, true, false}
	actualBools = list.Map(people, multiCombo1.Test)
	if slices.Equal(wantBools, actualBools) == false {
		t.Errorf("multiCombo.Test() = %v, want %v", actualBools, wantBools)
	}
	// multiCombo.BuildCondition
	wantCond := "(`Name` = ? OR `Job` IN (?, ?))"
	wantValues := []any{"Jill", "dev", "qa"}
	actualCond, actualValues := multiCombo1.BuildCondition()
	if actualCond != wantCond || slices.Equal(wantValues, actualValues) == false {
		t.Errorf("multiCombo.BuildCondition() = %q, %v, want %q, %v", actualCond, actualValues, wantCond, wantValues)
	}
}

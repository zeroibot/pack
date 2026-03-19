package qb

import (
	"testing"

	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/list"
	"github.com/roidaradal/tst"
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
	p := &Person{}
	emptyValues := make([]any, 0)
	this := testPrelude(t, p)
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
		tst.AssertEqual2(t, "newColumnValue", x.actual.Value(), x.want, x.actual.IsNil(), x.isNil)
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
		tst.AssertEqual2(t, "newFieldColumnValue", x.actual.Value(), x.want, x.actual.IsNil(), x.isNil)
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
		tst.AssertEqual(t, "newColumnValueList", x.actual.IsNil(), x.isNil)
		tst.AssertDeepEqual(t, "newColumnValueList", x.actual.Value(), x.want)
	}

	// falseConditionValues
	actualCond, actualValues := falseConditionValues()
	tst.AssertEqual(t, "falseConditionValues", actualCond, "false")
	tst.AssertListEqual(t, "falseConditionValues", actualValues, emptyValues)

	// trueConditionValues
	actualCond, actualValues = trueConditionValues()
	tst.AssertEqual(t, "trueConditionValues", actualCond, "true")
	tst.AssertListEqual(t, "trueConditionValues", actualValues, emptyValues)

	// soloConditionValues

	// normal operators
	kv1 := newColumnValue(this, &p.Name, "John").Value()
	kv2 := newColumnValue(this, &p.Age, 20).Value()
	kv3 := newColumnValue(this, &p.Level, 5).Value()
	kv4 := newColumnValue(this, &p.IP, nil).Value()
	kv5 := newColumnValue(this, &p.Job, "Dev").Value()

	testCases3 := []tst.P3W2[string, operator, any, string, []any]{
		{kv1.V1, opEqual, kv1.V2, "`Name` = ?", []any{"John"}},
		{kv2.V1, opGreater, kv2.V2, "`Age` > ?", []any{20}},
		{kv3.V1, opNotEqual, kv3.V2, "`Lvl` <> ?", []any{5}},
		{kv4.V1, opEqual, kv4.V2, "`IP` IS NULL", emptyValues},
		{kv4.V1, opNotEqual, kv4.V2, "`IP` IS NOT NULL", emptyValues},
		{kv1.V1, opPrefix, kv1.V2, "`Name` LIKE ?", []any{"John%"}},
		{kv1.V1, opSuffix, kv1.V2, "`Name` LIKE ?", []any{"%John"}},
		{kv5.V1, opSubstring, kv5.V2, "`Job` LIKE ?", []any{"%Dev%"}},
	}
	tst.AllP3W2(t, testCases3, "soloConditionValues", soloConditionValues, tst.AssertEqual, tst.AssertListEqual)
}

func TestInternalConditions(t *testing.T) {
	type Person struct {
		Name    string
		Age     int `col:"age"`
		Job     string
		Details string `col:"-"`
	}
	p := &Person{}
	emptyValues := make([]any, 0)
	this := testPrelude(t, p)
	// missingCondition.BuildCondition()
	cond1 := missingCondition{}
	actualCond, actualValues := cond1.BuildCondition()
	tst.AssertEqual(t, "missingCondition.BuildCondition", actualCond, "false")
	tst.AssertListEqual(t, "missingCondition.BuildCondition", actualValues, emptyValues)
	// matchAllCondition.BuildCondition()
	cond2 := matchAllCondition{}
	actualCond, actualValues = cond2.BuildCondition()
	tst.AssertEqual(t, "matchAllCondition.BuildCondition", actualCond, "true")
	tst.AssertListEqual(t, "matchAllCondition.BuildCondition", actualValues, emptyValues)
	// newValueCondition
	valueCond1 := newValueCondition(this, &p.Name, "John", opNotEqual)
	valueCond2 := newValueCondition(this, &p.Age, 20, opGreater)
	valueCond3 := newValueCondition(this, &p.Details, "Info", opSubstring)
	isNils := []bool{false, false, true}
	for i, valueCond := range []valueCondition{valueCond1, valueCond2, valueCond3} {
		tst.AssertEqual(t, "newValueCondition", valueCond.pair.IsNil(), isNils[i])
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
	testCases1 := []tst.P1W2[valueCondition, string, []any]{
		{valueCond1, "`Name` <> ?", []any{"John"}},
		{valueCond2, "`age` > ?", []any{20}},
		{valueCond3, "false", emptyValues},
		{valueCond4, "false", emptyValues},
	}
	tst.AllP1W2(t, testCases1, "valueCondition.BuildCondition", valueCondition.BuildCondition, tst.AssertEqual, tst.AssertListEqual)
	// newListCondition
	listCond1 := newListCondition(this, &p.Age, ds.List[int]{18, 19, 20}, opIn, opEqual)
	listCond2 := newListCondition(this, &p.Job, ds.List[string]{"dev", "qa"}, opNotIn, opNotEqual)
	listCond3 := newListCondition(this, &p.Name, ds.List[string]{"John"}, opIn, opEqual)
	listCond4 := newListCondition(this, &p.Details, ds.List[string]{"a", "b", "c"}, opNotIn, opNotEqual)
	listCond5 := newListCondition(this, &p.Name, ds.List[string]{}, opIn, opEqual)
	isNils = []bool{false, false, false, true, false}
	for i, listCond := range []listCondition{listCond1, listCond2, listCond3, listCond4, listCond5} {
		tst.AssertEqual(t, "newListCondition", listCond.pair.IsNil(), isNils[i])
	}
	// listCondition.BuildCondition()
	listCond6 := listCondition{
		pair:         ds.NewOption(new(columnValueListPair{V1: "", V2: []any{1, 2, 3}})),
		listOperator: opIn,
		soloOperator: opEqual,
	}
	testCases2 := []tst.P1W2[listCondition, string, []any]{
		{listCond1, "`age` IN (?, ?, ?)", []any{18, 19, 20}},
		{listCond2, "`Job` NOT IN (?, ?)", []any{"dev", "qa"}},
		{listCond3, "`Name` = ?", []any{"John"}},
		{listCond4, "false", emptyValues},
		{listCond5, "false", emptyValues},
		{listCond6, "false", emptyValues},
	}
	tst.AllP1W2(t, testCases2, "listCondition.BuildCondition", listCondition.BuildCondition, tst.AssertEqual, tst.AssertListEqual)
	// newMultiCondition
	multiCond1 := newMultiCondition(opAnd, valueCond1, valueCond2)
	multiCond2 := newMultiCondition(opOr, listCond1, listCond2)
	multiCond3 := newMultiCondition(opOr, valueCond1)
	multiCond4 := newMultiCondition(opAnd)
	multiCond5 := newMultiCondition(opOr, multiCond1, multiCond2, nil)
	multiCond6 := newMultiCondition(opOr, valueCond2, multiCond4)
	// multiCondition.BuildCondition()
	testCases3 := []tst.P1W2[multiCondition, string, []any]{
		{multiCond1, "(`Name` <> ? AND `age` > ?)", []any{"John", 20}},
		{multiCond2, "(`age` IN (?, ?, ?) OR `Job` NOT IN (?, ?))", []any{18, 19, 20, "dev", "qa"}},
		{multiCond3, "`Name` <> ?", []any{"John"}},
		{multiCond4, "false", emptyValues},
		{multiCond5, "((`Name` <> ? AND `age` > ?) OR (`age` IN (?, ?, ?) OR `Job` NOT IN (?, ?)))", []any{"John", 20, 18, 19, 20, "dev", "qa"}},
		{multiCond6, "false", emptyValues},
	}
	tst.AllP1W2(t, testCases3, "multiCondition.BuildCondition", multiCondition.BuildCondition, tst.AssertEqual, tst.AssertListEqual)
}

func TestInternalCombos(t *testing.T) {
	type Person struct {
		Name    string
		Age     int `col:"age"`
		Job     string
		Details string `col:"-"`
	}
	p := &Person{}
	p1 := Person{"John", 18, "dev", "regular"}
	p2 := Person{"James", 20, "qa", "prob"}
	p3 := Person{"Jill", 25, "admin", "regular"}
	p4 := Person{"Juno", 23, "dev", "prob"}
	p5 := Person{"Jack", 21, "sales", "intern"}
	people := ds.List[Person]{p1, p2, p3, p4, p5}
	this := testPrelude(t, p)
	// missingCombo.Test
	combo1 := missingCombo[Person]{}
	allFalse := people.All(func(person Person) bool {
		return combo1.Test(person) == false
	})
	tst.AssertTrue(t, "missingCombo.Test", allFalse)
	// matchAllCombo.Test
	combo2 := matchAllCombo[Person]{}
	allTrue := people.All(func(person Person) bool {
		return combo2.Test(person) == true
	})
	tst.AssertTrue(t, "matchAllCombo.Test", allTrue)
	// newValueCombo, valueCombo.Test
	test1 := createFieldValueTest[Person]("Name", func(name string) bool { return name == "Jill" })
	valueCond1 := newValueCondition(this, &p.Name, "Jill", opEqual)
	valueCombo1 := newValueCombo(valueCond1, test1)
	actualBools := list.Map(people, valueCombo1.Test)
	tst.AssertListEqual(t, "valueCombo.Test", actualBools, []bool{false, false, true, false, false})
	test0 := createFieldValueTest[Person]("Unknown", func(unknown string) bool { return unknown != "A" })
	valueCombo2 := newValueCombo(valueCond1, test0)
	allFalse = people.All(func(person Person) bool {
		return valueCombo2.Test(person) == false
	})
	tst.AssertTrue(t, "valueCombo.Test", allFalse)
	// newListCombo, listCombo.Test
	jobs := []string{"dev", "qa"}
	test2 := createFieldValueTest[Person]("Job", func(job string) bool { return list.Has(jobs, job) })
	listCond1 := newListCondition(this, &p.Job, jobs, opIn, opEqual)
	listCombo1 := newListCombo(listCond1, test2)
	actualBools = list.Map(people, listCombo1.Test)
	tst.AssertListEqual(t, "listCombo.Test", actualBools, []bool{true, true, false, true, false})
	// newMultiCombo, multiCombo.Test
	multiTest := func(person Person) bool { return test1(person) || test2(person) }
	multiCombo1 := newMultiCombo(ds.List[DualCondition[Person]]{valueCombo1, listCombo1}, opOr, multiTest)
	actualBools = list.Map(people, multiCombo1.Test)
	tst.AssertListEqual(t, "multiCombo.Test", actualBools, []bool{true, true, true, true, false})
	// multiCombo.BuildCondition
	wantCond := "(`Name` = ? OR `Job` IN (?, ?))"
	wantValues := []any{"Jill", "dev", "qa"}
	actualCond, actualValues := multiCombo1.BuildCondition()
	tst.AssertEqual(t, "multiCombo.BuildCondition", actualCond, wantCond)
	tst.AssertListEqual(t, "multiCombo.Test", actualValues, wantValues)
}

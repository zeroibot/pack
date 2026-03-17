package qb

import (
	"slices"
	"testing"

	"github.com/roidaradal/pack/ds"
)

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

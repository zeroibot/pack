package qb

import (
	"slices"
	"testing"

	"github.com/roidaradal/pack/ds"
)

func TestConditions(t *testing.T) {
	type Person struct {
		Name    string
		Address string
		Age     int
		Job     string
		Score   int
	}
	type testCase struct {
		cond       Condition
		wantCond   string
		wantValues ds.List[any]
	}
	this := NewInstance(MySQL)
	p := &Person{}
	err := AddType(this, p)
	if err != nil {
		t.Errorf("AddType error: %v", err)
	}
	condNone := NoCondition()
	condEqual := Equal(this, &p.Name, "John")
	condNotEqual := NotEqual(this, &p.Job, "manager")
	condPrefix := Prefix(this, &p.Job, "assistant")
	condSuffix := Suffix(this, &p.Address, "City")
	condSubstring := Substring(this, &p.Address, "Tower")
	condGreater := Greater(this, &p.Score, 75)
	condGreaterEqual := GreaterEqual(this, &p.Age, 20)
	condLesser := Lesser(this, &p.Age, 60)
	condLesserEqual := LesserEqual(this, &p.Score, 50)
	condIn := In(this, &p.Job, []string{"dev", "qa", "intern"})
	condNotIn := NotIn(this, &p.Score, []int{67, 69})
	condAnd := And(condEqual, condGreater)
	condOr := Or(condLesserEqual, condIn)

	testCases := []testCase{
		{condNone, "true", ds.List[any]{}},
		{condEqual, "`Name` = ?", ds.List[any]{"John"}},
		{condNotEqual, "`Job` <> ?", ds.List[any]{"manager"}},
		{condPrefix, "`Job` LIKE ?", ds.List[any]{"assistant%"}},
		{condSuffix, "`Address` LIKE ?", ds.List[any]{"%City"}},
		{condSubstring, "`Address` LIKE ?", ds.List[any]{"%Tower%"}},
		{condGreater, "`Score` > ?", ds.List[any]{75}},
		{condGreaterEqual, "`Age` >= ?", ds.List[any]{20}},
		{condLesser, "`Age` < ?", ds.List[any]{60}},
		{condLesserEqual, "`Score` <= ?", ds.List[any]{50}},
		{condIn, "`Job` IN (?, ?, ?)", ds.List[any]{"dev", "qa", "intern"}},
		{condNotIn, "`Score` NOT IN (?, ?)", ds.List[any]{67, 69}},
		{condAnd, "(`Name` = ? AND `Score` > ?)", ds.List[any]{"John", 75}},
		{condOr, "(`Score` <= ? OR `Job` IN (?, ?, ?))", ds.List[any]{50, "dev", "qa", "intern"}},
	}
	for _, x := range testCases {
		actualCond, actualValues := x.cond.BuildCondition()
		if actualCond != x.wantCond || slices.Equal(actualValues, x.wantValues) == false {
			t.Errorf("Condition.BuildCondition() = %q, %v, want %q, %v", actualCond, actualValues, x.wantCond, x.wantValues)
		}
	}
}

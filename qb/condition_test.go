package qb

import (
	"slices"
	"testing"
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
		wantValues []any
	}
	this := NewInstance(MySQL)
	p := &Person{}
	err := AddType(this, p)
	if err != nil {
		t.Errorf("AddType error: %v", err)
	}
	condNone := EmptyCondition()
	condEqual := EqualCondition(this, &p.Name, "John")
	condNotEqual := NotEqualCondition(this, &p.Job, "manager")
	condPrefix := PrefixCondition(this, &p.Job, "assistant")
	condSuffix := SuffixCondition(this, &p.Address, "City")
	condSubstring := SubstringCondition(this, &p.Address, "Tower")
	condGreater := GreaterCondition(this, &p.Score, 75)
	condGreaterEqual := GreaterEqualCondition(this, &p.Age, 20)
	condLesser := LesserCondition(this, &p.Age, 60)
	condLesserEqual := LesserEqualCondition(this, &p.Score, 50)
	condIn := InCondition(this, &p.Job, []string{"dev", "qa", "intern"})
	condNotIn := NotInCondition(this, &p.Score, []int{67, 69})
	condAnd := AndCondition(condEqual, condGreater)
	condOr := OrCondition(condLesserEqual, condIn)

	testCases := []testCase{
		{condNone, "true", []any{}},
		{condEqual, "`Name` = ?", []any{"John"}},
		{condNotEqual, "`Job` <> ?", []any{"manager"}},
		{condPrefix, "`Job` LIKE ?", []any{"assistant%"}},
		{condSuffix, "`Address` LIKE ?", []any{"%City"}},
		{condSubstring, "`Address` LIKE ?", []any{"%Tower%"}},
		{condGreater, "`Score` > ?", []any{75}},
		{condGreaterEqual, "`Age` >= ?", []any{20}},
		{condLesser, "`Age` < ?", []any{60}},
		{condLesserEqual, "`Score` <= ?", []any{50}},
		{condIn, "`Job` IN (?, ?, ?)", []any{"dev", "qa", "intern"}},
		{condNotIn, "`Score` NOT IN (?, ?)", []any{67, 69}},
		{condAnd, "(`Name` = ? AND `Score` > ?)", []any{"John", 75}},
		{condOr, "(`Score` <= ? OR `Job` IN (?, ?, ?))", []any{50, "dev", "qa", "intern"}},
	}
	for _, x := range testCases {
		actualCond, actualValues := x.cond.BuildCondition()
		if actualCond != x.wantCond || slices.Equal(actualValues, x.wantValues) == false {
			t.Errorf("Condition.BuildCondition() = %q, %v, want %q, %v", actualCond, actualValues, x.wantCond, x.wantValues)
		}
	}
}

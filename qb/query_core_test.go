package qb

import (
	"errors"
	"slices"
	"testing"
)

func TestQueryCore(t *testing.T) {
	type User struct {
		Username string
		Password string
		Count    int
	}
	table := "users"
	this := NewInstance(MySQL)
	u := new(User)
	err := AddType(this, u)
	if err != nil {
		t.Errorf("AddType() error = %v", err)
	}
	type testCase struct {
		q          *conditionQuery[User]
		wantCond   string
		wantValues []any
		wantErr    error
	}
	// initializeRequired
	q1 := new(conditionQuery[User])
	q1.initializeRequired(this, table)
	wantCond1, wantValues1 := falseConditionValues()
	// initializeOptional
	q2 := new(conditionQuery[User])
	q2.initializeOptional(this, table)
	wantCond2, wantValues2 := trueConditionValues()
	// Where
	q3 := new(conditionQuery[User])
	q3.initializeRequired(this, table)
	q3.Where(Greater[User](this, &u.Count, 5))
	// Empty table
	q4 := new(conditionQuery[User])
	q4.initializeOptional(this, "")
	testCases := []testCase{
		{q1, wantCond1, wantValues1, nil},
		{q2, wantCond2, wantValues2, nil},
		{q3, "`Count` > ?", []any{5}, nil},
		{q4, wantCond2, wantValues2, errEmptyTable},
	}
	for _, x := range testCases {
		actualCond, actualValues, actualErr := x.q.preBuildCheck()
		if actualCond != x.wantCond || slices.Equal(actualValues, x.wantValues) == false {
			t.Errorf("conditionQuery = %q, %v, want %q, %v", actualCond, actualValues, x.wantCond, x.wantValues)
		}
		if (x.wantErr == nil && actualErr != nil) || (x.wantErr != nil && !errors.Is(actualErr, x.wantErr)) {
			t.Errorf("conditionQuery.preBuildCheck error = %v, want %v", actualErr, x.wantErr)
		}
	}
}

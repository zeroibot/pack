package qb

import (
	"errors"
	"slices"
	"testing"

	"github.com/roidaradal/pack/dyn"
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

func TestOrderedLimit(t *testing.T) {
	this := NewInstance(MySQL)
	emptyQuery := new(orderedLimit)
	q0 := new(orderedLimit)
	q1 := new(orderedLimit)
	q2 := new(orderedLimit)
	q3 := new(orderedLimit)
	q4 := new(orderedLimit)
	q5 := new(orderedLimit)
	q6 := new(orderedLimit)
	q7 := new(orderedLimit)

	// OrderAsc, OrderDesc
	q1.OrderAsc(this, "")  // no column = no effect
	q2.OrderDesc(this, "") // no column = no effect
	if dyn.NotEqual(q1, emptyQuery) {
		t.Errorf("OrderAsc() = %v, want %v", q1, emptyQuery)
	}
	if dyn.NotEqual(q2, emptyQuery) {
		t.Errorf("OrderDesc() = %v, want %v", q2, emptyQuery)
	}

	q1.OrderAsc(this, "Name")
	q2.OrderDesc(this, "CreatedAt")
	q3.Limit(5)
	q4.OrderAsc(this, "Code").Limit(10)
	q5.OrderDesc(this, "UpdatedAt").Limit(5)
	q6.OrderAsc(this, "Code").OrderDesc(this, "UpdatedAt").Limit(5)
	q7.OrderDesc(this, "UpdatedAt").OrderAsc(this, "Code").Limit(10)

	// orderString
	type testCase struct {
		q         *orderedLimit
		wantOrder string
		wantLimit uint
	}
	testCases := []testCase{
		{q0, "", 0},
		{q1, "`Name` ASC", 0},
		{q2, "`CreatedAt` DESC", 0},
		{q3, "", 5},
		{q4, "`Code` ASC", 10},
		{q5, "`UpdatedAt` DESC", 5},
		{q6, "`Code` ASC, `UpdatedAt` DESC", 5},
		{q7, "`UpdatedAt` DESC, `Code` ASC", 10},
	}
	for _, x := range testCases {
		actualOrder := x.q.orderString()
		actualLimit := x.q.limit
		if actualOrder != x.wantOrder || actualLimit != x.wantLimit {
			t.Errorf("orderedLimit() = %q, %d, want %q, %d", actualOrder, actualLimit, x.wantOrder, x.wantLimit)
		}
	}

	// fullString
	type testCase2 struct {
		q          *orderedLimit
		wantString string
	}
	c0 := testCase2{q0, ""}
	c3 := testCase2{q3, "LIMIT 5"}
	c4 := testCase2{q4, "ORDER BY `Code` ASC LIMIT 10"}
	c5 := testCase2{q5, "ORDER BY `UpdatedAt` DESC LIMIT 5"}
	c6 := testCase2{q6, "ORDER BY `Code` ASC, `UpdatedAt` DESC LIMIT 5"}
	c7 := testCase2{q7, "ORDER BY `UpdatedAt` DESC, `Code` ASC LIMIT 10"}
	testCases2 := []testCase2{
		c0, c3, c4, c5, c6, c7,
		{q1, "ORDER BY `Name` ASC"},
		{q2, "ORDER BY `CreatedAt` DESC"},
	}
	for _, x := range testCases2 {
		actualString := x.q.fullString()
		if actualString != x.wantString {
			t.Errorf("orderedLimit.String() = %q, want %q", actualString, x.wantString)
		}
	}
	// mustLimitString
	testCases2 = []testCase2{
		c0, c3, c4, c5, c6, c7,
		{q1, ""},
		{q2, ""},
	}
	for _, x := range testCases2 {
		actualString := x.q.mustLimitString()
		if actualString != x.wantString {
			t.Errorf("orderedLimit.String() = %q, want %q", actualString, x.wantString)
		}
	}
}

package qb

import (
	"testing"

	"github.com/roidaradal/pack/dyn"
	"github.com/roidaradal/tst"
)

func TestQueryCore(t *testing.T) {
	type User struct {
		Username string
		Password string
		Count    int
	}
	table := "users"
	u := new(User)
	this := testPrelude(t, u)
	// Note: use testCase struct here because of the error checking
	type testCase struct {
		q          *conditionQuery[User]
		wantCond   string
		wantValues []any
		errNotNil  bool
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
		{q1, wantCond1, wantValues1, false},
		{q2, wantCond2, wantValues2, false},
		{q3, "`Count` > ?", []any{5}, false},
		{q4, wantCond2, wantValues2, true},
	}
	name := "conditionQuery.preBuildCheck"
	for _, x := range testCases {
		actualCond, actualValues, actualErr := x.q.preBuildCheck()
		tst.AssertEqual(t, name, actualCond, x.wantCond)
		tst.AssertListEqualError(t, name, actualValues, x.wantValues, actualErr, x.errNotNil)
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
	tst.AssertTrue(t, "OrderAsc", dyn.IsEqual(q1, emptyQuery))
	tst.AssertTrue(t, "OrderDesc", dyn.IsEqual(q2, emptyQuery))

	q1.OrderAsc(this, "Name")
	q2.OrderDesc(this, "CreatedAt")
	q3.Limit(5)
	q4.OrderAsc(this, "Code").Limit(10)
	q5.OrderDesc(this, "UpdatedAt").Limit(5)
	q6.OrderAsc(this, "Code").OrderDesc(this, "UpdatedAt").Limit(5)
	q7.OrderDesc(this, "UpdatedAt").OrderAsc(this, "Code").Limit(10)

	// orderString
	testCases := []tst.P1W2[*orderedLimit, string, uint]{
		{q0, "", 0},
		{q1, "`Name` ASC", 0},
		{q2, "`CreatedAt` DESC", 0},
		{q3, "", 5},
		{q4, "`Code` ASC", 10},
		{q5, "`UpdatedAt` DESC", 5},
		{q6, "`Code` ASC, `UpdatedAt` DESC", 5},
		{q7, "`UpdatedAt` DESC, `Code` ASC", 10},
	}
	orderStringLimit := func(q *orderedLimit) (string, uint) {
		return q.orderString(), q.limit
	}
	tst.AllP1W2(t, testCases, "OrderString,Limit", orderStringLimit, tst.AssertEqual[string], tst.AssertEqual[uint])

	// fullString
	type testCase2 = tst.P1W1[*orderedLimit, string]
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
	tst.AllP1W1(t, testCases2, "orderedLimit.fullString", (*orderedLimit).fullString, tst.AssertEqual)
	// mustLimitString
	testCases2 = []testCase2{
		c0, c3, c4, c5, c6, c7, {q1, ""}, {q2, ""},
	}
	tst.AllP1W1(t, testCases2, "orderedLimit.mustLimitString", (*orderedLimit).mustLimitString, tst.AssertEqual)
}

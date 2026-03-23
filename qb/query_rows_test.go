package qb

import (
	"fmt"
	"testing"

	"github.com/roidaradal/tst"
)

func TestDistinctValuesQuery(t *testing.T) {
	type User struct {
		Username string
		Age      int
		Extra    string `col:"-"`
		secret   string
	}
	u := new(User)
	table := "users"
	this := testPrelude(t, u)

	// NewDistinctValuesQuery
	q1 := NewDistinctValuesQuery[User](this, table, &u.Username)
	q1.Where(Equal[User](this, &u.Age, 18))
	q2 := NewDistinctValuesQuery[User](this, table, &u.Username) // no condition
	q3 := NewDistinctValuesQuery[User](this, "", &u.Username)    // no table
	q4 := NewDistinctValuesQuery[User](this, table, new(string)) // invalid field (not in struct)
	q5 := NewDistinctValuesQuery[User](this, table, &u.secret)   // private field
	q6 := NewDistinctValuesQuery[User](this, table, &u.Extra)    // blank column

	// DistinctValuesQuery.BuildQuery
	emptyValues := make([]any, 0)
	testCases1 := []tst.P1W2[*DistinctValuesQuery[User, string], string, []any]{
		{q1, "SELECT DISTINCT `Username` FROM `users` WHERE `Age` = ?", []any{18}},
		{q2, "SELECT DISTINCT `Username` FROM `users` WHERE true", emptyValues},
		{q3, "", emptyValues},
		{q4, "", emptyValues},
		{q5, "", emptyValues},
		{q6, "", emptyValues},
	}
	tst.AllP1W2(t, testCases1, "DistinctValuesQuery.BuildQuery", (*DistinctValuesQuery[User, string]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)

	// DistinctValuesQuery.Test
	u1 := User{"Alice", 18, "", ""}
	u2 := User{"Bob", 20, "", ""}
	testCases2 := []tst.P2W1[*DistinctValuesQuery[User, string], User, bool]{
		{q1, u1, true}, {q1, u2, false},
		{q2, u1, true}, {q2, u2, true},
	}
	tst.AllP2W1(t, testCases2, "DistinctValuesQuery.Test", (*DistinctValuesQuery[User, string]).Test, tst.AssertEqual)

	// ToString(DistinctValuesQuery)
	testCases3 := []tst.P1W1[Query, string]{
		{q1, fmt.Sprintf("SELECT DISTINCT `Username` FROM `users` WHERE `Age` = %d", 18)},
		{q2, "SELECT DISTINCT `Username` FROM `users` WHERE true"},
	}
	tst.AllP1W1(t, testCases3, "ToString(DistinctValuesQuery)", ToString, tst.AssertEqual)
}

func TestLookupQuery(t *testing.T) {
	// TODO: NewLookupQuery
	// TODO: LookupQuery.Where
	// TODO: LookupQuery without condition (valid)
	// TODO: LookupQuery.Test
	// TODO: LookupQuery.BuildQuery
}

func TestSelectRowsQuery(t *testing.T) {
	// TODO: NewSelectRowsQuery
	// TODO: NewFullSelectRowsQuery
	// TODO: SelectRowsQuery.Columns
	// TODO: SelectRowsQuery.Where
	// TODO: SelectRowsQuery without condition (valid)
	// TODO: SelectRowsQuery.Test
	// TODO: SelectRowsQuery.OrderAsc, OrderDesc
	// TODO: SelectRowsQuery.Limit
	// TODO: SelectRowsQuery.Page
	// TODO: SelectRowsQuery.BuildQuery
}

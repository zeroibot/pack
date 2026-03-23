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
	type User struct {
		Username string
		Age      int
		Extra    string `col:"-"`
		secret   string
	}
	u := new(User)
	table := "users"
	this := testPrelude(t, u)

	// NewLookupQuery
	q1 := NewLookupQuery[User](this, table, &u.Username, &u.Age)
	q1.Where(Greater[User](this, &u.Age, 18))
	q2 := NewLookupQuery[User](this, table, &u.Username, &u.Age)    // no condition
	q3 := NewLookupQuery[User](this, "", &u.Username, &u.Age)       // no table
	q4 := NewLookupQuery[User](this, table, new(string), &u.Age)    // invalid key field (not in struct)
	q5 := NewLookupQuery[User](this, table, &u.Username, new(int))  // invalid value field (not in struct)
	q6 := NewLookupQuery[User](this, table, &u.secret, &u.Age)      // private key field
	q7 := NewLookupQuery[User](this, table, &u.Username, &u.secret) // private value field
	q8 := NewLookupQuery[User](this, table, &u.Extra, &u.Age)       // blank key column
	q9 := NewLookupQuery[User](this, table, &u.Username, &u.Extra)  // blank value column

	// LookupQuery.BuildQuery
	emptyValues := make([]any, 0)
	testCases1 := []tst.P1W2[*LookupQuery[User, string, int], string, []any]{
		{q1, "SELECT `Username`, `Age` FROM `users` WHERE `Age` > ?", []any{18}},
		{q2, "SELECT `Username`, `Age` FROM `users` WHERE true", emptyValues},
		{q3, "", emptyValues}, {q4, "", emptyValues}, {q5, "", emptyValues},
		{q6, "", emptyValues}, {q8, "", emptyValues},
	}
	tst.AllP1W2(t, testCases1, "LookupQuery.BuildQuery", (*LookupQuery[User, string, int]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)
	testCases4 := []tst.P1W2[*LookupQuery[User, string, string], string, []any]{
		{q7, "", emptyValues}, {q9, "", emptyValues},
	}
	tst.AllP1W2(t, testCases4, "LookupQuery.BuildQuery", (*LookupQuery[User, string, string]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)

	// LookupQuery.Test
	u1 := User{"Alice", 18, "", ""}
	u2 := User{"Bob", 20, "", ""}
	testCases2 := []tst.P2W1[*LookupQuery[User, string, int], User, bool]{
		{q1, u1, false}, {q1, u2, true},
		{q2, u1, true}, {q2, u2, true},
	}
	tst.AllP2W1(t, testCases2, "LookupQuery.Test", (*LookupQuery[User, string, int]).Test, tst.AssertEqual)

	// ToString(LookupQuery)
	testCases3 := []tst.P1W1[Query, string]{
		{q1, fmt.Sprintf("SELECT `Username`, `Age` FROM `users` WHERE `Age` > %d", 18)},
		{q2, "SELECT `Username`, `Age` FROM `users` WHERE true"},
	}
	tst.AllP1W1(t, testCases3, "ToString(LookupQuery)", ToString, tst.AssertEqual)
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

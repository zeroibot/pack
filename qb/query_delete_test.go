package qb

import (
	"fmt"
	"testing"

	"github.com/roidaradal/tst"
)

func TestDeleteQuery(t *testing.T) {
	type User struct {
		Username string
		Password string
		Count    int
		IP       *string
	}
	u := new(User)
	table := "users"
	this := testPrelude(t, u)

	// NewDeleteQuery
	q1 := NewDeleteQuery[User](this, table)
	q1.Where(NotEqual[User](this, &u.Username, "root"))
	q2 := NewDeleteQuery[User](this, table) // no condition
	q3 := NewDeleteQuery[User](this, "")    // blank table

	// DeleteQuery.BuildQuery()
	emptyValues := make([]any, 0)
	testCases1 := []tst.P1W2[*DeleteQuery[User], string, []any]{
		{q1, "DELETE FROM `users` WHERE `Username` <> ?", []any{"root"}},
		{q2, "DELETE FROM `users` WHERE false", emptyValues},
		{q3, "", emptyValues},
	}
	tst.AllP1W2(t, testCases1, "DeleteQuery.BuildQuery", (*DeleteQuery[User]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)

	// DeleteQuery.Test()
	u1 := User{"root", "admin", 1, nil}
	u2 := User{"guest", "", 2, nil}
	testCases2 := []tst.P2W1[*DeleteQuery[User], User, bool]{
		{q1, u1, false},
		{q1, u2, true},
		{q2, u1, false},
		{q2, u2, false},
	}
	tst.AllP2W1(t, testCases2, "DeleteQuery,Test", (*DeleteQuery[User]).Test, tst.AssertEqual)

	// ToString(DeleteQuery)
	homeIP := new("127.0.0.1")
	q4 := NewDeleteQuery[User](this, table)
	q4.Where(Greater[User](this, &u.Count, 5))
	q5 := NewDeleteQuery[User](this, table)
	q5.Where(Equal[User](this, &u.IP, homeIP))
	testCases3 := []tst.P1W1[Query, string]{
		{q1, fmt.Sprintf("DELETE FROM `users` WHERE `Username` <> %q", "root")},
		{q2, "DELETE FROM `users` WHERE false"},
		{q4, "DELETE FROM `users` WHERE `Count` > 5"},
		{q5, fmt.Sprintf("DELETE FROM `users` WHERE `IP` = %q", "127.0.0.1")},
	}
	tst.AllP1W1(t, testCases3, "ToString(DeleteQuery)", ToString, tst.AssertEqual)

	// OrderAsc, OrderDesc, Limit
	q1.OrderAsc(this, "ID").Limit(1)         // OrderAsc + Limit
	q4.OrderDesc(this, "CreatedAt").Limit(5) // OrderDesc + Limit
	q5.Limit(10)                             // Limit only
	q6 := NewDeleteQuery[User](this, table)
	q6.Where(Lesser[User](this, &u.Count, 10))
	q6.OrderDesc(this, "CreatedAt").OrderAsc(this, "ID").Limit(5) // Mixed Orders + Limit
	q7 := NewDeleteQuery[User](this, table)
	q7.Where(Equal[User](this, &u.Username, "admin"))
	q7.OrderAsc(this, "ID") // Order only, no limit
	testCases1 = []tst.P1W2[*DeleteQuery[User], string, []any]{
		{q1, "DELETE FROM `users` WHERE `Username` <> ? ORDER BY `ID` ASC LIMIT 1", []any{"root"}},
		{q4, "DELETE FROM `users` WHERE `Count` > ? ORDER BY `CreatedAt` DESC LIMIT 5", []any{5}},
		{q5, "DELETE FROM `users` WHERE `IP` = ? LIMIT 10", []any{homeIP}},
		{q6, "DELETE FROM `users` WHERE `Count` < ? ORDER BY `CreatedAt` DESC, `ID` ASC LIMIT 5", []any{10}},
		{q7, "DELETE FROM `users` WHERE `Username` = ?", []any{"admin"}},
	}
	tst.AllP1W2(t, testCases1, "DeleteQuery.BuildQuery", (*DeleteQuery[User]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)
}

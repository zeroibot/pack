package qb

import (
	"fmt"
	"slices"
	"testing"
)

func TestDeleteQuery(t *testing.T) {
	type User struct {
		Username string
		Password string
		Count    int
		IP       *string
	}
	table := "users"
	this := NewInstance(MySQL)
	u := new(User)
	err := AddType(this, u)
	if err != nil {
		t.Errorf("AddType() error = %v", err)
	}
	// NewDeleteQuery
	q1 := NewDeleteQuery[User](this, table)
	q1.Where(NotEqual[User](this, &u.Username, "root"))
	q2 := NewDeleteQuery[User](this, table) // no condition
	q3 := NewDeleteQuery[User](this, "")    // blank table
	// DeleteQuery.BuildQuery()
	type testCase1 struct {
		q          *DeleteQuery[User]
		wantQuery  string
		wantValues []any
	}
	emptyValues := make([]any, 0)
	testCases1 := []testCase1{
		{q1, "DELETE FROM `users` WHERE `Username` <> ?", []any{"root"}},
		{q2, "DELETE FROM `users` WHERE false", emptyValues},
		{q3, "", emptyValues},
	}
	for _, x := range testCases1 {
		actualQuery, actualValues := x.q.BuildQuery()
		if actualQuery != x.wantQuery || slices.Equal(actualValues, x.wantValues) == false {
			t.Errorf("DeleteQuery.BuildQuery() = %q, %v, want %q, %v", actualQuery, actualValues, x.wantQuery, x.wantValues)
		}
	}
	// DeleteQuery.Test()
	type testCase2 struct {
		q        *DeleteQuery[User]
		user     User
		wantFlag bool
	}
	u1 := User{"root", "admin", 1, nil}
	u2 := User{"guest", "", 2, nil}
	testCases2 := []testCase2{
		{q1, u1, false},
		{q1, u2, true},
		{q2, u1, false},
		{q2, u2, false},
	}
	for _, x := range testCases2 {
		actualFlag := x.q.Test(x.user)
		if actualFlag != x.wantFlag {
			t.Errorf("DeleteQuery.Test() = %t, want %t", actualFlag, x.wantFlag)
		}
	}
	// ToString(DeleteQuery)
	type testCase3 struct {
		q         *DeleteQuery[User]
		wantQuery string
	}
	homeIP := new("127.0.0.1")
	q4 := NewDeleteQuery[User](this, table)
	q4.Where(Greater[User](this, &u.Count, 5))
	q5 := NewDeleteQuery[User](this, table)
	q5.Where(Equal[User](this, &u.IP, homeIP))
	testCases3 := []testCase3{
		{q1, fmt.Sprintf("DELETE FROM `users` WHERE `Username` <> %q", "root")},
		{q2, "DELETE FROM `users` WHERE false"},
		{q4, "DELETE FROM `users` WHERE `Count` > 5"},
		{q5, fmt.Sprintf("DELETE FROM `users` WHERE `IP` = %q", "127.0.0.1")},
	}
	for _, x := range testCases3 {
		actualQuery := ToString(x.q)
		if actualQuery != x.wantQuery {
			t.Errorf("ToString(DeleteQuery) = %s, want %s", actualQuery, x.wantQuery)
		}
	}
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
	testCases1 = []testCase1{
		{q1, "DELETE FROM `users` WHERE `Username` <> ? ORDER BY `ID` ASC LIMIT 1", []any{"root"}},
		{q4, "DELETE FROM `users` WHERE `Count` > ? ORDER BY `CreatedAt` DESC LIMIT 5", []any{5}},
		{q5, "DELETE FROM `users` WHERE `IP` = ? LIMIT 10", []any{homeIP}},
		{q6, "DELETE FROM `users` WHERE `Count` < ? ORDER BY `CreatedAt` DESC, `ID` ASC LIMIT 5", []any{10}},
		{q7, "DELETE FROM `users` WHERE `Username` = ?", []any{"admin"}},
	}
	for _, x := range testCases1 {
		actualQuery, actualValues := x.q.BuildQuery()
		if actualQuery != x.wantQuery || slices.Equal(actualValues, x.wantValues) == false {
			t.Errorf("DeleteQuery.BuildQuery() = %q, %v, want %q, %v", actualQuery, actualValues, x.wantQuery, x.wantValues)
		}
	}
}

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
	q2 := NewDeleteQuery[User](this, table)
	q3 := NewDeleteQuery[User](this, "")
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
	q4 := NewDeleteQuery[User](this, table)
	q4.Where(Greater[User](this, &u.Count, 5))
	q5 := NewDeleteQuery[User](this, table)
	q5.Where(Equal[User](this, &u.IP, new("127.0.0.1")))
	testCases1 = []testCase1{
		{q1, fmt.Sprintf("DELETE FROM `users` WHERE `Username` <> %q", "root"), emptyValues},
		{q2, "DELETE FROM `users` WHERE false", emptyValues},
		{q4, "DELETE FROM `users` WHERE `Count` > 5", emptyValues},
		{q5, fmt.Sprintf("DELETE FROM `users` WHERE `IP` = %q", "127.0.0.1"), emptyValues},
	}
	for _, x := range testCases1 {
		actualQuery := ToString(x.q)
		if actualQuery != x.wantQuery {
			t.Errorf("ToString(DeleteQuery) = %s, want %s", actualQuery, x.wantQuery)
		}
	}
}

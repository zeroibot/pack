package qb

import (
	"testing"

	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/tst"
)

func TestUpdateQuery(t *testing.T) {
	type User struct {
		Username string
		Password string
		Count    int
		IP       *string
		secret   string
		Code     string `col:"UUID"`
	}
	table := "users"
	u := new(User)
	this := testPrelude(t, u)

	// NewUpdateQuery
	q0 := NewUpdateQuery[User](this, table) // no updates
	q1 := NewUpdateQuery[User](this, "")    // no table

	// Update
	q2 := NewUpdateQuery[User](this, table) // no condition
	Update(this, q2, &u.Username, "admin")
	q3 := NewUpdateQuery[User](this, table) // with multiple updates
	Update(this, q3, &u.Username, "admin")
	Update(this, q3, &u.Password, "123")
	q3.Where(Equal[User](this, &u.Username, "root"))
	q4 := NewUpdateQuery[User](this, table) // has a nil pair
	Update(this, q4, &u.Username, "admin")
	Update(this, q4, &u.secret, "secret")
	q5 := NewUpdateQuery[User](this, table) // pair has a blank column
	q5.updates = append(q5.updates, ds.NewOption(new(columnValuePair{V1: "", V2: 5})))

	// UpdateQuery.Update, UpdateQuery.Updates
	q6 := NewUpdateQuery[User](this, table)
	q6.Update(this, "Count", 5)
	q6.Where(Greater[User](this, &u.Count, 5))
	updates := FieldUpdates{
		"Code":     [2]any{5, 6},
		"Password": [2]any{"hahaha", "horse"},
	}
	q7 := NewUpdateQuery[User](this, table)
	q7.Updates(this, updates)
	q7.Where(Equal[User](this, &u.Username, "groot"))

	// UpdateQuery.BuildQuery
	emptyValues := make([]any, 0)
	testCases := []tst.P1W2[*UpdateQuery[User], string, []any]{
		{q0, "", emptyValues},
		{q1, "", emptyValues},
		{q2, "UPDATE `users` SET `Username` = ? WHERE false", []any{"admin"}},
		{q3, "UPDATE `users` SET `Username` = ?, `Password` = ? WHERE `Username` = ?", []any{"admin", "123", "root"}},
		{q4, "", emptyValues},
		{q5, "", emptyValues},
		{q6, "UPDATE `users` SET `Count` = ? WHERE `Count` > ?", []any{5, 5}},
		{q7, "UPDATE `users` SET `UUID` = ?, `Password` = ? WHERE `Username` = ?", []any{6, "horse", "groot"}},
	}
	tst.AllP1W2(t, testCases, "UpdateQuery.BuildQuery", (*UpdateQuery[User]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)

	// OrderAsc, OrderDesc, Limit
	q3.OrderDesc(this, "CreatedAt").Limit(1)                      // OrderDesc + Limit
	q6.Limit(10)                                                  // Limit only
	q7.OrderDesc(this, "CreatedAt").OrderAsc(this, "ID").Limit(1) // Mixed Orders + Limit
	q2.Where(Equal[User](this, &u.Username, "groot"))
	q2.OrderAsc(this, "Username") // Order only, no limit
	testCases = []tst.P1W2[*UpdateQuery[User], string, []any]{
		{q3, "UPDATE `users` SET `Username` = ?, `Password` = ? WHERE `Username` = ? ORDER BY `CreatedAt` DESC LIMIT 1", []any{"admin", "123", "root"}},
		{q6, "UPDATE `users` SET `Count` = ? WHERE `Count` > ? LIMIT 10", []any{5, 5}},
		{q7, "UPDATE `users` SET `UUID` = ?, `Password` = ? WHERE `Username` = ? ORDER BY `CreatedAt` DESC, `ID` ASC LIMIT 1", []any{6, "horse", "groot"}},
		{q2, "UPDATE `users` SET `Username` = ? WHERE `Username` = ?", []any{"admin", "groot"}},
	}
	tst.AllP1W2(t, testCases, "UpdateQuery.BuildQuery", (*UpdateQuery[User]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)
}

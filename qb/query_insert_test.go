package qb

import (
	"fmt"
	"testing"

	"github.com/roidaradal/pack/dict"
	"github.com/roidaradal/tst"
)

func TestInsertRowQuery(t *testing.T) {
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

	// NewInsertRowQuery
	q0 := NewInsertRowQuery(this, table) // empty row
	q1 := NewInsertRowQuery(this, "")    // no table

	// InsertRowQuery.Row
	q2 := NewInsertRowQuery(this, table)
	q2.Row(this, dict.Object{"Username": "admin", "Password": "123", "": "blank"}) // one blank column
	q3 := NewInsertRowQuery(this, table)
	q3.Row(this, dict.Object{"Username": "john", "Password": "12345", "Count": 5})

	// Using ToRow
	ip := new("127.0.0.1")
	u1 := new(User{Username: "Jane", Password: "6767", Count: 10, Code: "eagle", IP: ip})
	u2 := new(User{Username: "Jack", Password: "6969", Count: 5, Code: "tiger", IP: ip})
	q4 := NewInsertRowQuery(this, table)
	q4.Row(this, ToRow(this, u1))
	q5 := NewInsertRowQuery(this, table)
	q5.Row(this, ToRow(this, u2))
	q6 := NewInsertRowQuery(this, table)
	q6.Row(this, dict.Object{"IP": nil})
	q7 := NewInsertRowQuery(this, table)
	q7.Row(this, dict.Object{"Username": "homer", "IP": nil})
	q8 := NewInsertRowQuery(this, table)
	r := ToRow(this, new(User))
	q8.Row(this, r)
	fmt.Println(r["IP"] == nil, r["IP"], fmt.Sprintf("%T", r["IP"]))

	// InsertRowQuery.BuildQuery
	emptyValues := make([]any, 0)
	testCases := []tst.P1W2[*InsertRowQuery, string, []any]{
		{q0, "", emptyValues},
		{q1, "", emptyValues},
		{q2, "INSERT INTO `users` (`Password`, `Username`) VALUES (?, ?)", []any{"123", "admin"}},
		{q3, "INSERT INTO `users` (`Count`, `Password`, `Username`) VALUES (?, ?, ?)", []any{5, "12345", "john"}},
		{q4, "INSERT INTO `users` (`Count`, `IP`, `Password`, `UUID`, `Username`) VALUES (?, ?, ?, ?, ?)", []any{10, ip, "6767", "eagle", "Jane"}},
		{q5, "INSERT INTO `users` (`Count`, `IP`, `Password`, `UUID`, `Username`) VALUES (?, ?, ?, ?, ?)", []any{5, ip, "6969", "tiger", "Jack"}},
		{q6, "INSERT INTO `users` (`IP`) VALUES (?)", []any{nil}},
		{q7, "INSERT INTO `users` (`IP`, `Username`) VALUES (?, ?)", []any{nil, "homer"}},
		{q8, "INSERT INTO `users` (`Count`, `IP`, `Password`, `UUID`, `Username`) VALUES (?, ?, ?, ?, ?)", []any{0, nil, "", "", ""}},
	}
	// Note: used ListMixedEqual here because of IP (*string) which has nil checking
	tst.AllP1W2(t, testCases, "InsertRowQuery.BuildQuery", (*InsertRowQuery).BuildQuery, tst.AssertEqual, tst.AssertListMixedEqual)
}

func TestInsertRowsQuery(t *testing.T) {
	//type User struct {
	//	Username string
	//	Password string
	//	Count    int
	//	IP       *string
	//	secret   string
	//	Code     string `col:"UUID"`
	//}
	//table := "users"
	//u := new(User)
	//this := testPrelude(t, u)
	//
	//// NewInsertRowsQuery
	//q0 := NewInsertRowsQuery(this, table) // no rows
	//q1 := NewInsertRowsQuery(this, "")    // no table
	//
	//// InsertRowsQuery.Rows
	//q2 := NewInsertRowsQuery(this, table)
	//q2.Rows(this, dict.Object{"Username": "admin", "Password": "123"})
	//q3 := NewInsertRowsQuery(this, table)
	//q3.Rows(this, dict.Object{"Username": "admin", "Password": "123"}, dict.Object{"Username": "root", "Password": "456"})
	//q4 := NewInsertRowsQuery(this, table) // blank column
	//q4.Rows(this, dict.Object{"Username": "admin", "": "blank"})
	//q5 := NewInsertRowsQuery(this, table) // empty row
	//q5.Rows(this, dict.Object{})
	//q6 := NewInsertRowsQuery(this, table) // inconsistent signatures
	//q6.Rows(this, dict.Object{"Username": "admin"}, dict.Object{"Password": "123"})
	//
	//// InsertRowsQuery.BuildQuery
	//emptyValues := make([]any, 0)
	//testCases := []tst.P1W2[*InsertRowsQuery, string, []any]{
	//	{q0, "", emptyValues},
	//	{q1, "", emptyValues},
	//	{q2, "INSERT INTO `users` (`Username`, `Password`) VALUES (?, ?)", []any{"admin", "123"}},
	//	{q3, "INSERT INTO `users` (`Username`, `Password`) VALUES (?, ?), (?, ?)", []any{"admin", "123", "root", "456"}},
	//	{q4, "INSERT INTO `users` (`Username`) VALUES (?)", []any{"admin"}},
	//	{q5, "", emptyValues},
	//	{q6, "", emptyValues},
	//}
	//tst.AllP1W2(t, testCases, "InsertRowsQuery.BuildQuery", (*InsertRowsQuery).BuildQuery, tst.AssertEqual, tst.AssertListEqual)
}

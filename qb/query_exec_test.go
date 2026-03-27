package qb

import (
	"fmt"
	"testing"

	"github.com/roidaradal/pack/conv"
	"github.com/roidaradal/pack/db"
	"github.com/roidaradal/pack/dict"
	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/list"
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

func TestDeleteExec(t *testing.T) {
	type User struct {
		ID   uint
		Name string
		Job  string
	}
	table := "users"
	u := new(User)
	this := testPrelude(t, u)
	u1, u2, u3 := User{1, "Alice", "Dev"}, User{2, "Bob", "Dev"}, User{3, "Charlie", "Dev"}
	u4, u5, u6 := User{4, "Dave", "QA"}, User{5, "Eve", "QA"}, User{6, "Frank", "Sales"}
	u7, u8, u9 := User{7, "Grace", "Sales"}, User{8, "Harry", "UX"}, User{9, "Ivy", "Admin"}
	users := []User{u1, u2, u3, u4, u5, u6, u7, u8, u9}

	q0 := NewDeleteQuery[User](this, "")
	q1 := NewDeleteQuery[User](this, table)
	q1.Where(Equal[User](this, &u.ID, 4))
	q2 := NewDeleteQuery[User](this, table)
	q2.Where(In[User](this, &u.Job, []string{"Sales", "UX"}))
	q3 := NewDeleteQuery[User](this, table)
	q3.Where(Or[User](
		Equal[User](this, &u.Job, "Dev"),
		Equal[User](this, &u.Name, "Ivy"),
	))
	q4 := NewDeleteQuery[User](this, table)
	q4.Where(Greater[User](this, &u.ID, 10))

	execFn := func(test func(User) bool) func([]User) ([]User, error) {
		test = conv.NotFn(test)
		return func(items []User) ([]User, error) {
			return list.Filter(items, test), nil
		}
	}
	newResult := func(test func(User) bool) *tst.Result {
		return tst.NewResult(list.CountFunc(users, test), 0, nil)
	}

	dbc := db.NewMockAdapter(tst.NewConn(users...))
	prep1 := dbc.Conn.PrepExecReset(execFn(q1.Test), newResult(q1.Test), users...)
	prep1b := func() { prep1(); dbc.Conn.SetError(errMock) }
	prep2 := dbc.Conn.PrepExecReset(execFn(q2.Test), newResult(q2.Test), users...)
	prep3 := dbc.Conn.PrepExecReset(execFn(q3.Test), newResult(q3.Test), users...)
	prep4 := dbc.Conn.PrepExecReset(execFn(q4.Test), newResult(q4.Test), users...)

	want1 := []User{u1, u2, u3, u5, u6, u7, u8, u9}
	want2 := []User{u1, u2, u3, u4, u5, u9}
	want3 := []User{u4, u5, u6, u7, u8}
	testCases := []tst.P2W3Pre[*DeleteQuery[User], db.Conn, int, bool, []User]{
		{nil, q0, dbc, 0, false, users},    // empty query
		{nil, q1, nil, 0, false, users},    // no db connection
		{prep1, q1, dbc, 1, true, want1},   // success query1
		{prep1b, q1, dbc, 0, false, users}, // error on query1
		{prep2, q2, dbc, 3, true, want2},   // success query2
		{prep3, q3, dbc, 4, true, want3},   // success query3
		{prep4, q4, dbc, 0, true, users},   // success query4, no rows deleted
	}
	deleteExec := func(q *DeleteQuery[User], dbConn db.Conn) (int, bool, []User) {
		result, err := Exec(q, dbConn)
		if err != nil {
			return 0, false, dbc.Conn.Items()
		}
		return RowsAffected(result), true, dbc.Conn.Items()
	}
	tst.AllP2W3Pre(t, testCases, "DeleteQuery.Exec", deleteExec, tst.AssertEqual[int], tst.AssertEqual[bool], tst.AssertListEqual)
}

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
	q6 := NewInsertRowQuery(this, table)
	q6.Row(this, dict.Object{"IP": nil})
	q7 := NewInsertRowQuery(this, table)
	q7.Row(this, dict.Object{"Username": "homer", "IP": nil})

	// Using ToRow
	ip := new("127.0.0.1")
	u1 := new(User{Username: "Jane", Password: "6767", Count: 10, Code: "eagle", IP: ip})
	u2 := new(User{Username: "Jack", Password: "6969", Count: 5, Code: "tiger", IP: ip})
	q4 := NewInsertRowQuery(this, table)
	q4.Row(this, ToRow(this, u1))
	q5 := NewInsertRowQuery(this, table)
	q5.Row(this, ToRow(this, u2))
	q8 := NewInsertRowQuery(this, table)
	q8.Row(this, ToRow(this, new(User)))

	// InsertRowQuery.BuildQuery
	emptyValues := make([]any, 0)
	testCases := []tst.P1W2[*InsertRowQuery, string, []any]{
		{q0, "", emptyValues},
		{q1, "", emptyValues},
		{q2, "INSERT INTO `users` (`Password`, `Username`) VALUES (?, ?)", []any{"123", "admin"}},
		{q3, "INSERT INTO `users` (`Count`, `Password`, `Username`) VALUES (?, ?, ?)", []any{5, "12345", "john"}},
		{q6, "INSERT INTO `users` (`IP`) VALUES (?)", []any{nil}},
		{q7, "INSERT INTO `users` (`IP`, `Username`) VALUES (?, ?)", []any{nil, "homer"}},
		{q4, "INSERT INTO `users` (`Count`, `IP`, `Password`, `UUID`, `Username`) VALUES (?, ?, ?, ?, ?)", []any{10, ip, "6767", "eagle", "Jane"}},
		{q5, "INSERT INTO `users` (`Count`, `IP`, `Password`, `UUID`, `Username`) VALUES (?, ?, ?, ?, ?)", []any{5, ip, "6969", "tiger", "Jack"}},
		{q8, "INSERT INTO `users` (`Count`, `IP`, `Password`, `UUID`, `Username`) VALUES (?, ?, ?, ?, ?)", []any{0, nil, "", "", ""}},
	}
	// Note: used ListMixedEqual here because of IP (*string) which has nil checking
	tst.AllP1W2(t, testCases, "InsertRowQuery.BuildQuery", (*InsertRowQuery).BuildQuery, tst.AssertEqual, tst.AssertListMixedEqual)
}

func TestInsertRowsQuery(t *testing.T) {
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

	// NewInsertRowsQuery
	q0 := NewInsertRowsQuery(this, table) // no rows
	q1 := NewInsertRowsQuery(this, "")    // no table

	// InsertRowsQuery.Rows
	q2 := NewInsertRowsQuery(this, table)
	q2.Rows(this, dict.Object{"Username": "admin", "Password": "123"})
	q3 := NewInsertRowsQuery(this, table)
	q3.Rows(this, dict.Object{"Username": "admin", "Password": "123"}, dict.Object{"Username": "root", "Password": "456"})
	q4 := NewInsertRowsQuery(this, table) // blank column
	q4.Rows(this, dict.Object{"Username": "admin", "": "blank"})
	q5 := NewInsertRowsQuery(this, table) // empty row
	q5.Rows(this, dict.Object{})
	q6 := NewInsertRowsQuery(this, table) // inconsistent signatures
	q6.Rows(this, dict.Object{"Username": "admin"}, dict.Object{"Password": "123"})

	// Rows with ToRow
	ip1, ip2 := new("127.0.0.1"), new("localhost")
	john := new(User{Username: "John", Password: "1234", Count: 5})
	jack := new(User{Username: "Jack", Password: "6969", Count: 10})
	jane := new(User{Username: "Jane", Password: "6767", Count: 3, Code: "eagle", IP: ip1})
	juno := new(User{Username: "Juno", Password: "3435", Count: 7, Code: "tiger", IP: ip2})
	q7 := NewInsertRowsQuery(this, table)
	q7.Rows(this, ToRow(this, john), ToRow(this, jack))
	q8 := NewInsertRowsQuery(this, table)
	q8.Rows(this, ToRow(this, new(User)))
	q9 := NewInsertRowsQuery(this, table)
	users := []dict.Object{ToRow(this, jack), ToRow(this, jane), ToRow(this, juno)}
	q9.Rows(this, users...)

	// InsertRowsQuery.BuildQuery
	emptyValues := make([]any, 0)
	testCases := []tst.P1W2[*InsertRowsQuery, string, []any]{
		{q0, "", emptyValues},
		{q1, "", emptyValues},
		{q2, "INSERT INTO `users` (`Password`, `Username`) VALUES (?, ?)", []any{"123", "admin"}},
		{q3, "INSERT INTO `users` (`Password`, `Username`) VALUES (?, ?), (?, ?)", []any{"123", "admin", "456", "root"}},
		{q4, "INSERT INTO `users` (`Username`) VALUES (?)", []any{"admin"}},
		{q5, "", emptyValues},
		{q6, "", emptyValues},
		{q7, "INSERT INTO `users` (`Count`, `IP`, `Password`, `UUID`, `Username`) VALUES (?, ?, ?, ?, ?), (?, ?, ?, ?, ?)", []any{5, nil, "1234", "", "John", 10, nil, "6969", "", "Jack"}},
		{q8, "INSERT INTO `users` (`Count`, `IP`, `Password`, `UUID`, `Username`) VALUES (?, ?, ?, ?, ?)", []any{0, nil, "", "", ""}},
		{q9, "INSERT INTO `users` (`Count`, `IP`, `Password`, `UUID`, `Username`) VALUES (?, ?, ?, ?, ?), (?, ?, ?, ?, ?), (?, ?, ?, ?, ?)", []any{10, nil, "6969", "", "Jack", 3, ip1, "6767", "eagle", "Jane", 7, ip2, "3435", "tiger", "Juno"}},
	}
	tst.AllP1W2(t, testCases, "InsertRowsQuery.BuildQuery", (*InsertRowsQuery).BuildQuery, tst.AssertEqual, tst.AssertListMixedEqual)
}

func TestInsertRowExec(t *testing.T) {
	type User struct {
		ID   uint
		Name string
		Job  string
	}
	table := "users"
	u := new(User)
	this := testPrelude(t, u)
	u1, u2, u3 := User{1, "Alice", "Dev"}, User{2, "Bob", "Dev"}, User{3, "Charlie", "Dev"}
	u4 := User{4, "Dave", "QA"}
	//u5, u6 := User{5, "Eve", "QA"}, User{6, "Frank", "Sales"}
	//u7, u8, u9 := User{7, "Grace", "Sales"}, User{8, "Harry", "UX"}, User{9, "Ivy", "Admin"}

	q1 := NewInsertRowQuery(this, "") // no table
	q1.Row(this, ToRow(this, &u1))
	q2 := NewInsertRowQuery(this, table) // no row
	q3 := NewInsertRowQuery(this, table) // insert u3
	q3.Row(this, ToRow(this, &u3))
	q4 := NewInsertRowQuery(this, table) // insert u4
	q4.Row(this, ToRow(this, &u4))

	execFn1 := func(user User) func([]User) ([]User, error) {
		return func(users []User) ([]User, error) {
			users = append(users, user)
			return users, nil
		}
	}
	result1 := func(user User) *tst.Result {
		return tst.NewResult(1, int(user.ID), nil)
	}

	dbc := db.NewMockAdapter(tst.NewConn(u1, u2))
	prep3 := dbc.Conn.PrepExec(execFn1(u3), result1(u3))
	prep3b := func() { dbc.Conn.SetError(errMock) }
	prep4 := dbc.Conn.PrepExec(execFn1(u4), result1(u4))

	testCases := []tst.P2W4Pre[*InsertRowQuery, db.Conn, int, uint, bool, []User]{
		{nil, q1, dbc, 0, 0, false, []User{u1, u2}},          // empty query = no table
		{nil, q2, dbc, 0, 0, false, []User{u1, u2}},          // empty query = no row
		{prep3, q3, dbc, 1, 3, true, []User{u1, u2, u3}},     // success query3
		{prep3b, q3, dbc, 0, 0, false, []User{u1, u2, u3}},   // error on query3
		{prep4, q4, dbc, 1, 4, true, []User{u1, u2, u3, u4}}, // success query4
	}
	insertRow := func(q *InsertRowQuery, dbConn db.Conn) (int, uint, bool, []User) {
		result, err := Exec(q, dbConn)
		if err != nil {
			return 0, 0, false, dbc.Conn.Items()
		}
		var insertID uint = 0
		if id, ok := LastInsertID(result); ok {
			insertID = id
		}
		return RowsAffected(result), insertID, true, dbc.Conn.Items()
	}
	tst.AllP2W4Pre(t, testCases, "InsertRowQuery.Exec", insertRow, tst.AssertEqual[int], tst.AssertEqual[uint], tst.AssertEqual[bool], tst.AssertListEqual)
}
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

func TestResultCheckers(t *testing.T) {
	// TODO: AssertNothing
	// TODO: AssertRowsAffected
	// TODO: RowsAffected
	// TODO: LastInsertID
}

func TestExec(t *testing.T) {
	// TODO: Exec
	// TODO: ExecTx
	// TODO: Rollback
}

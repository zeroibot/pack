package qb

import (
	"fmt"
	"testing"

	"github.com/roidaradal/pack/db"
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
	q7 := NewDistinctValuesQuery[User](this, table, &u.Username) // zero results
	q7.Where(Lesser[User](this, &u.Age, 10))

	// DistinctValuesQuery.BuildQuery
	emptyValues := make([]any, 0)
	testCases1 := []tst.P1W2[*DistinctValuesQuery[User, string], string, []any]{
		{q1, "SELECT DISTINCT `Username` FROM `users` WHERE `Age` = ?", []any{18}},
		{q2, "SELECT DISTINCT `Username` FROM `users` WHERE true", emptyValues},
		{q3, "", emptyValues},
		{q4, "", emptyValues},
		{q5, "", emptyValues},
		{q6, "", emptyValues},
		{q7, "SELECT DISTINCT `Username` FROM `users` WHERE `Age` < ?", []any{10}},
	}
	tst.AllP1W2(t, testCases1, "DistinctValuesQuery.BuildQuery", (*DistinctValuesQuery[User, string]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)

	// DistinctValuesQuery.Test
	u1 := User{"Alice", 18, "", ""}
	u2 := User{"Bob", 20, "", ""}
	testCases2 := []tst.P2W1[*DistinctValuesQuery[User, string], User, bool]{
		{q1, u1, true}, {q1, u2, false},
		{q2, u1, true}, {q2, u2, true},
		{q7, u1, false}, {q7, u2, false},
	}
	tst.AllP2W1(t, testCases2, "DistinctValuesQuery.Test", (*DistinctValuesQuery[User, string]).Test, tst.AssertEqual)

	// ToString(DistinctValuesQuery)
	testCases3 := []tst.P1W1[Query, string]{
		{q1, "SELECT DISTINCT `Username` FROM `users` WHERE `Age` = 18"},
		{q2, "SELECT DISTINCT `Username` FROM `users` WHERE true"},
		{q7, "SELECT DISTINCT `Username` FROM `users` WHERE `Age` < 10"},
	}
	tst.AllP1W1(t, testCases3, "ToString(DistinctValuesQuery)", ToString, tst.AssertEqual)

	// DistinctValuesQuery.Query
	dbc := db.NewMockAdapter(tst.NewConn(u1, u2))
	prep0a := func() { dbc.Conn.SetError(errMock) }
	prep0b := func() { q1.reader = nil }
	getUsername := func(x User) []any { return []any{x.Username} }
	prep1 := dbc.Conn.PrepRows(q1.Test, getUsername)
	prep2 := dbc.Conn.PrepRows(q2.Test, getUsername)
	prep3 := dbc.Conn.PrepRows(q7.Test, getUsername)

	testCases4 := []tst.P2W2Pre[*DistinctValuesQuery[User, string], db.Conn, []string, bool]{
		{nil, q3, dbc, nil, false},                       // empty query
		{nil, q1, nil, nil, false},                       // no db connection
		{prep1, q1, dbc, []string{"Alice"}, true},        // success query1
		{prep2, q2, dbc, []string{"Alice", "Bob"}, true}, // success query2
		{prep3, q7, dbc, []string{}, true},               // empty rows
		{prep0a, q1, dbc, nil, false},                    // error on query
		{prep0b, q1, dbc, nil, false},                    // nil reader
	}
	distinctValuesQuery := func(q *DistinctValuesQuery[User, string], dbc db.Conn) ([]string, bool) {
		res := q.Query(this, dbc)
		return res.Value(), res.NotError()
	}
	tst.AllP2W2Pre(t, testCases4, "DistinctValuesQuery.Query", distinctValuesQuery, tst.AssertListEqual, tst.AssertEqual)
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

	// TODO: LookupQuery.Lookup
}

func TestSelectRowsQuery(t *testing.T) {
	type Product struct {
		ID    int
		Name  string
		Price float64
		Stock int    `col:"Qty"`
		Extra string `col:"-"`
		code  string
	}
	p := new(Product)
	this := testPrelude(t, p)
	table := "products"
	allCols := this.allColumns(p)
	cols2 := this.Columns(&p.Name, &p.Price)
	reader1 := NewRowReader[Product](this, allCols...)
	reader2 := NewRowReader[Product](this, cols2...)

	// NewSelectRowsQuery, NewFullSelectRowsQuery
	q0 := NewSelectRowsQuery[Product](this, "", nil)    // no table
	q1 := NewFullSelectRowsQuery(this, table, reader1)  // all columns
	q2 := NewSelectRowsQuery(this, table, reader2)      // specific columns
	q3 := NewSelectRowsQuery[Product](this, table, nil) // nil reader
	q4 := NewSelectRowsQuery(this, table, reader2)      // no columns set
	q5 := NewSelectRowsQuery(this, table, reader2)      // no condition (optional)
	q6 := NewSelectRowsQuery(this, table, reader2)      // with private column
	q7 := NewSelectRowsQuery(this, table, reader2)      // with skipped column

	// SelectRowsQuery.Columns
	q2.Columns(this, cols2...)
	q5.Columns(this, cols2...)
	q6.Columns(this, append(cols2, this.Column(&p.code))...)
	q7.Columns(this, append(cols2, this.Column(&p.Extra))...)

	// SelectRowsQuery.Where
	q1.Where(Greater[Product](this, &p.Price, 100.0))
	q2.Where(Equal[Product](this, &p.Stock, 50))

	// SelectRowsQuery.OrderAsc, OrderDesc, Limit, Page
	q1.OrderDesc(this, this.Column(&p.Price)).Limit(10)
	q2.OrderAsc(this, this.Column(&p.Name))
	q2.Page(2, 5) // offset 5, limit 5

	// SelectRowsQuery.Test
	p1 := Product{1, "Laptop", 1200.0, 10, "", "p1"}
	p2 := Product{2, "Mouse", 25.0, 50, "", "p2"}
	p3 := Product{3, "Monitor", 300.0, 50, "", "p3"}
	testCases1 := []tst.P2W1[*SelectRowsQuery[Product], Product, bool]{
		{q1, p1, true}, {q1, p2, false}, {q1, p3, true},
		{q2, p1, false}, {q2, p2, true}, {q2, p3, true},
		{q5, p1, true}, {q5, p2, true}, {q5, p3, true},
	}
	tst.AllP2W1(t, testCases1, "SelectRowsQuery.Test", (*SelectRowsQuery[Product]).Test, tst.AssertEqual)

	// SelectRowsQuery.BuildQuery
	emptyValues := make([]any, 0)
	testCases2 := []tst.P1W2[*SelectRowsQuery[Product], string, []any]{
		{q0, "", emptyValues},
		{q1, "SELECT `ID`, `Name`, `Price`, `Qty` FROM `products` WHERE `Price` > ? ORDER BY `Price` DESC LIMIT 0, 10", []any{100.0}},
		{q2, "SELECT `Name`, `Price` FROM `products` WHERE `Qty` = ? ORDER BY `Name` ASC LIMIT 5, 5", []any{50}},
		{q3, "", emptyValues},
		{q4, "", emptyValues},
		{q5, "SELECT `Name`, `Price` FROM `products` WHERE true", emptyValues},
		{q6, "SELECT `Name`, `Price` FROM `products` WHERE true", emptyValues},
		{q7, "SELECT `Name`, `Price` FROM `products` WHERE true", emptyValues},
	}
	tst.AllP1W2(t, testCases2, "SelectRowsQuery.BuildQuery", (*SelectRowsQuery[Product]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)

	// ToString(SelectRowsQuery)
	testCases3 := []tst.P1W1[Query, string]{
		{q1, "SELECT `ID`, `Name`, `Price`, `Qty` FROM `products` WHERE `Price` > 100 ORDER BY `Price` DESC LIMIT 0, 10"},
		{q2, "SELECT `Name`, `Price` FROM `products` WHERE `Qty` = 50 ORDER BY `Name` ASC LIMIT 5, 5"},
		{q5, "SELECT `Name`, `Price` FROM `products` WHERE true"},
	}
	tst.AllP1W1(t, testCases3, "ToString(SelectRowsQuery)", ToString, tst.AssertEqual)

	// TODO: SelectRowsQuery.Query
}

func TestGroupCountQuery(t *testing.T) {
	type User struct {
		ID     int
		Name   string
		Age    int
		Extra  string `col:"-"`
		secret string
	}
	u := new(User)
	this := testPrelude(t, u)
	table := "users"

	// NewGroupCountQuery
	q0 := NewGroupCountQuery[User, string](this, "", &u.Name)        // no table
	q1 := NewGroupCountQuery[User, string](this, table, &u.Name)     // valid
	q2 := NewGroupCountQuery[User, int](this, table, &u.Age)         // another field
	q3 := NewGroupCountQuery[User, string](this, table, &u.Name)     // with condition
	q4 := NewGroupCountQuery[User, string](this, table, new(string)) // invalid field
	q5 := NewGroupCountQuery[User, string](this, table, &u.secret)   // private field
	q6 := NewGroupCountQuery[User, string](this, table, &u.Extra)    // blank column (skipped)

	// GroupCountQuery.Where
	q3.Where(Greater[User](this, &u.Age, 18))

	// GroupCountQuery.BuildQuery
	emptyValues := make([]any, 0)
	testCases1 := []tst.P1W2[*GroupCountQuery[User, string], string, []any]{
		{q0, "", emptyValues},
		{q1, "SELECT `Name`, COUNT(*) FROM `users` WHERE true GROUP BY `Name`", emptyValues},
		{q3, "SELECT `Name`, COUNT(*) FROM `users` WHERE `Age` > ? GROUP BY `Name`", []any{18}},
		{q4, "", emptyValues},
		{q5, "", emptyValues},
		{q6, "", emptyValues},
	}
	tst.AllP1W2(t, testCases1, "GroupCountQuery.BuildQuery", (*GroupCountQuery[User, string]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)

	testCases2 := []tst.P1W2[*GroupCountQuery[User, int], string, []any]{
		{q2, "SELECT `Age`, COUNT(*) FROM `users` WHERE true GROUP BY `Age`", emptyValues},
	}
	tst.AllP1W2(t, testCases2, "GroupCountQuery.BuildQuery", (*GroupCountQuery[User, int]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)

	// GroupCountQuery.Test
	u1 := User{1, "Alice", 20, "", ""}
	u2 := User{2, "Bob", 15, "", ""}
	testCases3 := []tst.P2W1[*GroupCountQuery[User, string], User, bool]{
		{q1, u1, true}, {q1, u2, true},
		{q3, u1, true}, {q3, u2, false},
	}
	tst.AllP2W1(t, testCases3, "GroupCountQuery.Test", (*GroupCountQuery[User, string]).Test, tst.AssertEqual)

	// ToString(GroupCountQuery)
	testCases4 := []tst.P1W1[Query, string]{
		{q1, "SELECT `Name`, COUNT(*) FROM `users` WHERE true GROUP BY `Name`"},
		{q3, fmt.Sprintf("SELECT `Name`, COUNT(*) FROM `users` WHERE `Age` > %d GROUP BY `Name`", 18)},
	}
	tst.AllP1W1(t, testCases4, "ToString(GroupCountQuery)", ToString, tst.AssertEqual)

	// TODO: GroupCountQuery.GroupCount
}

func TestGroupSumQuery(t *testing.T) {
	type Product struct {
		ID      int
		Name    string
		Price   float64
		Qty     int
		Extra   int `col:"-"`
		code    string
		balance float64
	}
	p := new(Product)
	this := testPrelude(t, p)
	table := "products"

	// NewGroupSumQuery
	q0 := NewGroupSumQuery[Product, string, float64](this, "", &p.Name, &p.Price)        // no table
	q1 := NewGroupSumQuery[Product, string, float64](this, table, &p.Name, &p.Price)     // valid
	q2 := NewGroupSumQuery[Product, int, int](this, table, &p.ID, &p.Qty)                // another fields
	q3 := NewGroupSumQuery[Product, string, float64](this, table, &p.Name, &p.Price)     // with condition
	q4 := NewGroupSumQuery[Product, string, float64](this, table, new(string), &p.Price) // invalid group field
	q5 := NewGroupSumQuery[Product, string, float64](this, table, &p.Name, new(float64)) // invalid sum field
	q6 := NewGroupSumQuery[Product, string, float64](this, table, &p.code, &p.Price)     // private group field
	q7 := NewGroupSumQuery[Product, string, float64](this, table, &p.Name, &p.balance)   // private sum field
	q8 := NewGroupSumQuery[Product, int, float64](this, table, &p.Extra, &p.Price)       // skipped group field
	q9 := NewGroupSumQuery[Product, string, int](this, table, &p.Name, &p.Extra)         // skipped sum field

	// GroupSumQuery.Where
	q3.Where(Greater[Product](this, &p.Qty, 10))

	// GroupSumQuery.BuildQuery
	emptyValues := make([]any, 0)
	testCases1 := []tst.P1W2[*GroupSumQuery[Product, string, float64], string, []any]{
		{q0, "", emptyValues},
		{q1, "SELECT `Name`, SUM(`Price`) FROM `products` WHERE true GROUP BY `Name`", emptyValues},
		{q3, "SELECT `Name`, SUM(`Price`) FROM `products` WHERE `Qty` > ? GROUP BY `Name`", []any{10}},
		{q4, "", emptyValues},
		{q5, "", emptyValues},
		{q6, "", emptyValues},
		{q7, "", emptyValues},
	}
	tst.AllP1W2(t, testCases1, "GroupSumQuery.BuildQuery (string, float64)", (*GroupSumQuery[Product, string, float64]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)

	testCases2 := []tst.P1W2[*GroupSumQuery[Product, int, int], string, []any]{
		{q2, "SELECT `ID`, SUM(`Qty`) FROM `products` WHERE true GROUP BY `ID`", emptyValues},
	}
	tst.AllP1W2(t, testCases2, "GroupSumQuery.BuildQuery (int, int)", (*GroupSumQuery[Product, int, int]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)

	testCases6 := []tst.P1W2[*GroupSumQuery[Product, int, float64], string, []any]{
		{q8, "", emptyValues},
	}
	tst.AllP1W2(t, testCases6, "GroupSumQuery.BuildQuery (int, float64)", (*GroupSumQuery[Product, int, float64]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)

	testCases5 := []tst.P1W2[*GroupSumQuery[Product, string, int], string, []any]{
		{q9, "", emptyValues},
	}
	tst.AllP1W2(t, testCases5, "GroupSumQuery.BuildQuery (string, int)", (*GroupSumQuery[Product, string, int]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)

	// GroupSumQuery.Test
	p1 := Product{ID: 1, Name: "Laptop", Price: 1200.0, Qty: 20}
	p2 := Product{ID: 2, Name: "Mouse", Price: 25.0, Qty: 5}
	testCases3 := []tst.P2W1[*GroupSumQuery[Product, string, float64], Product, bool]{
		{q1, p1, true}, {q1, p2, true},
		{q3, p1, true}, {q3, p2, false},
	}
	tst.AllP2W1(t, testCases3, "GroupSumQuery.Test", (*GroupSumQuery[Product, string, float64]).Test, tst.AssertEqual)

	// ToString(GroupSumQuery)
	testCases4 := []tst.P1W1[Query, string]{
		{q1, "SELECT `Name`, SUM(`Price`) FROM `products` WHERE true GROUP BY `Name`"},
		{q3, "SELECT `Name`, SUM(`Price`) FROM `products` WHERE `Qty` > 10 GROUP BY `Name`"},
	}
	tst.AllP1W1(t, testCases4, "ToString(GroupSumQuery)", ToString, tst.AssertEqual)

	// TODO: GroupSumQuery.GroupSum
}

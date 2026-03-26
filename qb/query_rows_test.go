package qb

import (
	"cmp"
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
	q0 := NewLookupQuery[User](this, table, &u.Username, &u.Age)
	q0.Where(Lesser[User](this, &u.Age, 10))
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
		{q0, "SELECT `Username`, `Age` FROM `users` WHERE `Age` < ?", []any{10}},
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
		{q0, u1, false}, {q0, u2, false},
		{q1, u1, false}, {q1, u2, true},
		{q2, u1, true}, {q2, u2, true},
	}
	tst.AllP2W1(t, testCases2, "LookupQuery.Test", (*LookupQuery[User, string, int]).Test, tst.AssertEqual)

	// ToString(LookupQuery)
	testCases3 := []tst.P1W1[Query, string]{
		{q0, "SELECT `Username`, `Age` FROM `users` WHERE `Age` < 10"},
		{q1, "SELECT `Username`, `Age` FROM `users` WHERE `Age` > 18"},
		{q2, "SELECT `Username`, `Age` FROM `users` WHERE true"},
	}
	tst.AllP1W1(t, testCases3, "ToString(LookupQuery)", ToString, tst.AssertEqual)

	// LookupQuery.Lookup
	dbc := db.NewMockAdapter(tst.NewConn(u1, u2))
	prep0a := func() { dbc.Conn.SetError(errMock) }
	prep0b := func() { q1.reader = nil }
	getUsernameAge := func(x User) []any { return []any{x.Username, x.Age} }
	prep1 := dbc.Conn.PrepRows(q1.Test, getUsernameAge)
	prep2 := dbc.Conn.PrepRows(q2.Test, getUsernameAge)
	prep3 := dbc.Conn.PrepRows(q0.Test, getUsernameAge)
	want1 := map[string]int{"Bob": 20}
	want2 := map[string]int{"Alice": 18, "Bob": 20}
	want3 := map[string]int{}

	testCases5 := []tst.P2W2Pre[*LookupQuery[User, string, int], db.Conn, map[string]int, bool]{
		{nil, q3, dbc, nil, false},    // empty query
		{nil, q1, dbc, nil, false},    // no db connection
		{prep1, q1, dbc, want1, true}, // success query1
		{prep2, q2, dbc, want2, true}, // success query2
		{prep3, q0, dbc, want3, true}, // empty results
		{prep0a, q1, dbc, nil, false}, // error on query
		{prep0b, q1, dbc, nil, false}, // nil reader
	}
	lookupQuery := func(q *LookupQuery[User, string, int], dbc db.Conn) (map[string]int, bool) {
		res := q.Lookup(this, dbc)
		return res.Value(), res.NotError()
	}
	tst.AllP2W2Pre(t, testCases5, "LookupQuery.Lookup", lookupQuery, tst.AssertMapEqual, tst.AssertEqual)
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
	q8 := NewSelectRowsQuery(this, table, reader2)      // no results

	// SelectRowsQuery.Columns
	q2.Columns(this, cols2...)
	q5.Columns(this, cols2...)
	q6.Columns(this, append(cols2, this.Column(&p.code))...)
	q7.Columns(this, append(cols2, this.Column(&p.Extra))...)
	q8.Columns(this, cols2...)

	// SelectRowsQuery.Where
	q1.Where(Greater[Product](this, &p.Price, 100.0))
	q2.Where(Equal[Product](this, &p.Stock, 50))
	q8.Where(Equal[Product](this, &p.Name, "Computer"))

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
		{q8, p1, false}, {q8, p2, false}, {q8, p3, false},
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
		{q8, "SELECT `Name`, `Price` FROM `products` WHERE `Name` = ?", []any{"Computer"}},
	}
	tst.AllP1W2(t, testCases2, "SelectRowsQuery.BuildQuery", (*SelectRowsQuery[Product]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)

	// ToString(SelectRowsQuery)
	testCases3 := []tst.P1W1[Query, string]{
		{q1, "SELECT `ID`, `Name`, `Price`, `Qty` FROM `products` WHERE `Price` > 100 ORDER BY `Price` DESC LIMIT 0, 10"},
		{q2, "SELECT `Name`, `Price` FROM `products` WHERE `Qty` = 50 ORDER BY `Name` ASC LIMIT 5, 5"},
		{q5, "SELECT `Name`, `Price` FROM `products` WHERE true"},
		{q8, fmt.Sprintf("SELECT `Name`, `Price` FROM `products` WHERE `Name` = %q", "Computer")},
	}
	tst.AllP1W1(t, testCases3, "ToString(SelectRowsQuery)", ToString, tst.AssertEqual)

	// SelectRowsQuery.Query
	dbc := db.NewMockAdapter(tst.NewConn(p1, p2, p3))
	prep0a := func() { dbc.Conn.SetError(errMock) }
	prep0b := func() { q1.reader = nil }
	getAllColumns := func(x Product) []any { return []any{x.ID, x.Name, x.Price, x.Stock} }
	getNamePrice := func(x Product) []any { return []any{x.Name, x.Price} }
	sortPriceDesc := func(x1, x2 Product) int { return cmp.Compare(x2.Price, x1.Price) }
	prep1 := dbc.Conn.PrepSortRows(q1.Test, getAllColumns, sortPriceDesc, 10)
	prep5 := dbc.Conn.PrepRows(q5.Test, getNamePrice)
	prep8 := dbc.Conn.PrepRows(q8.Test, getNamePrice)
	want1 := []Product{{ID: 1, Name: "Laptop", Price: 1200.0, Stock: 10}, {ID: 3, Name: "Monitor", Price: 300.0, Stock: 50}}
	want5 := []Product{{Name: "Laptop", Price: 1200.0}, {Name: "Mouse", Price: 25.0}, {Name: "Monitor", Price: 300.0}}
	want8 := make([]Product, 0)

	testCases4 := []tst.P2W2Pre[*SelectRowsQuery[Product], db.Conn, []Product, bool]{
		{nil, q0, dbc, nil, false},    // empty query
		{nil, q1, nil, nil, false},    // no DB connection
		{prep1, q1, dbc, want1, true}, // success query1
		{prep5, q5, dbc, want5, true}, // success query5
		{prep8, q8, dbc, want8, true}, // empty results
		{prep0a, q1, dbc, nil, false}, // error on query
		{prep0b, q1, dbc, nil, false}, // nil reader
	}
	selectRowsQuery := func(q *SelectRowsQuery[Product], dbc db.Conn) ([]Product, bool) {
		res := q.Query(dbc)
		return res.Value(), res.NotError()
	}
	tst.AllP2W2Pre(t, testCases4, "SelectRowsQuery.Query", selectRowsQuery, tst.AssertListEqual, tst.AssertEqual)
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
	q7 := NewGroupCountQuery[User, string](this, table, &u.Name)     // no results

	// GroupCountQuery.Where
	q3.Where(Greater[User](this, &u.Age, 18))
	q7.Where(Greater[User](this, &u.ID, 10))

	// GroupCountQuery.BuildQuery
	emptyValues := make([]any, 0)
	testCases1 := []tst.P1W2[*GroupCountQuery[User, string], string, []any]{
		{q0, "", emptyValues},
		{q1, "SELECT `Name`, COUNT(*) FROM `users` WHERE true GROUP BY `Name`", emptyValues},
		{q3, "SELECT `Name`, COUNT(*) FROM `users` WHERE `Age` > ? GROUP BY `Name`", []any{18}},
		{q4, "", emptyValues},
		{q5, "", emptyValues},
		{q6, "", emptyValues},
		{q7, "SELECT `Name`, COUNT(*) FROM `users` WHERE `ID` > ? GROUP BY `Name`", []any{10}},
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
		{q7, u1, false}, {q7, u2, false},
	}
	tst.AllP2W1(t, testCases3, "GroupCountQuery.Test", (*GroupCountQuery[User, string]).Test, tst.AssertEqual)

	// ToString(GroupCountQuery)
	testCases4 := []tst.P1W1[Query, string]{
		{q1, "SELECT `Name`, COUNT(*) FROM `users` WHERE true GROUP BY `Name`"},
		{q3, "SELECT `Name`, COUNT(*) FROM `users` WHERE `Age` > 18 GROUP BY `Name`"},
		{q7, "SELECT `Name`, COUNT(*) FROM `users` WHERE `ID` > 10 GROUP BY `Name`"},
	}
	tst.AllP1W1(t, testCases4, "ToString(GroupCountQuery)", ToString, tst.AssertEqual)

	// GroupCountQuery.GroupCount
	u3 := User{3, "Alice", 5, "", ""}
	u4 := User{4, "Cat", 22, "", ""}
	dbc := db.NewMockAdapter(tst.NewConn(u1, u2, u3, u4))
	groupByName := func(users []User) [][]any {
		counts := make(map[string]int)
		for _, user := range users {
			counts[user.Name] += 1
		}
		values := make([][]any, 0, len(counts))
		for name, count := range counts {
			values = append(values, []any{name, count})
		}
		return values
	}
	prep0 := func() { dbc.Conn.SetError(errMock) }
	prep1 := dbc.Conn.PrepGroup(q1.Test, groupByName)
	prep3 := dbc.Conn.PrepGroup(q3.Test, groupByName)
	prep7 := dbc.Conn.PrepGroup(q7.Test, groupByName)
	want1 := map[string]int{"Alice": 2, "Bob": 1, "Cat": 1}
	want3 := map[string]int{"Alice": 1, "Cat": 1}
	want7 := make(map[string]int)

	testCases5 := []tst.P2W2Pre[*GroupCountQuery[User, string], db.Conn, map[string]int, bool]{
		{nil, q0, dbc, nil, false},    // empty query
		{nil, q1, nil, nil, false},    // no DB connection,
		{prep0, q1, dbc, nil, false},  // error on query
		{prep1, q1, dbc, want1, true}, // success query1
		{prep3, q3, dbc, want3, true}, // success query3
		{prep7, q7, dbc, want7, true}, // empty results
	}
	groupCountQuery := func(q *GroupCountQuery[User, string], dbc db.Conn) (map[string]int, bool) {
		res := q.GroupCount(dbc)
		return res.Value(), res.NotError()
	}
	tst.AllP2W2Pre(t, testCases5, "GroupCountQuery.GroupCount", groupCountQuery, tst.AssertMapEqual, tst.AssertEqual)
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
	q10 := NewGroupSumQuery[Product, string, float64](this, table, &p.Name, &p.Price)    // no results

	// GroupSumQuery.Where
	q3.Where(Greater[Product](this, &p.Qty, 10))
	q10.Where(Greater[Product](this, &p.ID, 10))

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
		{q10, "SELECT `Name`, SUM(`Price`) FROM `products` WHERE `ID` > ? GROUP BY `Name`", []any{10}},
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
		{q10, p1, false}, {q10, p2, false},
	}
	tst.AllP2W1(t, testCases3, "GroupSumQuery.Test", (*GroupSumQuery[Product, string, float64]).Test, tst.AssertEqual)

	// ToString(GroupSumQuery)
	testCases4 := []tst.P1W1[Query, string]{
		{q1, "SELECT `Name`, SUM(`Price`) FROM `products` WHERE true GROUP BY `Name`"},
		{q3, "SELECT `Name`, SUM(`Price`) FROM `products` WHERE `Qty` > 10 GROUP BY `Name`"},
		{q10, "SELECT `Name`, SUM(`Price`) FROM `products` WHERE `ID` > 10 GROUP BY `Name`"},
	}
	tst.AllP1W1(t, testCases4, "ToString(GroupSumQuery)", ToString, tst.AssertEqual)

	// GroupSumQuery.GroupSum
	p3 := Product{ID: 3, Name: "Laptop", Price: 1500.0, Qty: 10}
	p4 := Product{ID: 4, Name: "Monitor", Price: 300.0, Qty: 50}
	dbc := db.NewMockAdapter(tst.NewConn(p1, p2, p3, p4))
	sumPriceByName := func(products []Product) [][]any {
		totalPrice := make(map[string]float64)
		for _, product := range products {
			totalPrice[product.Name] += product.Price
		}
		values := make([][]any, 0, len(totalPrice))
		for name, price := range totalPrice {
			values = append(values, []any{name, price})
		}
		return values
	}
	prep0 := func() { dbc.Conn.SetError(errMock) }
	prep1 := dbc.Conn.PrepGroup(q1.Test, sumPriceByName)
	prep3 := dbc.Conn.PrepGroup(q3.Test, sumPriceByName)
	prep10 := dbc.Conn.PrepGroup(q10.Test, sumPriceByName)
	want1 := map[string]float64{"Laptop": 2700.0, "Mouse": 25.0, "Monitor": 300.0}
	want3 := map[string]float64{"Laptop": 1200.0, "Monitor": 300.0}
	want10 := make(map[string]float64)

	testCases7 := []tst.P2W2Pre[*GroupSumQuery[Product, string, float64], db.Conn, map[string]float64, bool]{
		{nil, q0, dbc, nil, false},       // empty query
		{nil, q1, nil, nil, false},       // no DB connection
		{prep0, q1, dbc, nil, false},     // error on query
		{prep1, q1, dbc, want1, true},    // success query1
		{prep3, q3, dbc, want3, true},    // success query2
		{prep10, q10, dbc, want10, true}, // empty results
	}
	groupSumQuery := func(q *GroupSumQuery[Product, string, float64], dbc db.Conn) (map[string]float64, bool) {
		res := q.GroupSum(dbc)
		return res.Value(), res.NotError()
	}
	tst.AllP2W2Pre(t, testCases7, "GroupSumQuery.GroupSum", groupSumQuery, tst.AssertMapEqual, tst.AssertEqual)
}

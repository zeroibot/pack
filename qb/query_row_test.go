package qb

import (
	"testing"

	"github.com/roidaradal/tst"
)

func TestCountQuery(t *testing.T) {
	type Person struct {
		Name string
		Age  int
		Job  string
	}
	p := new(Person)
	this := testPrelude(t, p)
	table := "persons"
	// NewCountQuery
	q0 := NewCountQuery[Person](this, "")    // blank table
	q1 := NewCountQuery[Person](this, table) // with condition
	q2 := NewCountQuery[Person](this, table) // no condition
	// CountQuery.Where
	q1.Where(GreaterEqual[Person](this, &p.Age, 18))
	// CountQuery.Test
	testCases := []tst.P2W1[*CountQuery[Person], Person, bool]{
		{q1, Person{"John", 20, "dev"}, true},
		{q1, Person{"Jane", 18, "student"}, true},
		{q1, Person{"Alice", 15, "student"}, false},
	}
	tst.AllP2W1(t, testCases, "CountQuery.Test", (*CountQuery[Person]).Test, tst.AssertEqual)
	// CountQuery.BuildQuery
	testCases2 := []tst.P1W2[*CountQuery[Person], string, []any]{
		{q0, "", []any{}},
		{q1, "SELECT COUNT(*) FROM `persons` WHERE `Age` >= ?", []any{18}},
		{q2, "SELECT COUNT(*) FROM `persons` WHERE false", []any{}},
	}
	tst.AllP1W2(t, testCases2, "CountQuery.BuildQuery", (*CountQuery[Person]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)
}

func TestValueQuery(t *testing.T) {
	type User struct {
		ID     int
		Name   string
		Code   string
		Job    string
		Extra  string `col:"-"`
		secret string
	}
	u := new(User)
	this := testPrelude(t, u)
	table := "users"
	// NewValueQuery
	q0 := NewValueQuery[User](this, "", &u.Name)      // no table
	q1 := NewValueQuery[User](this, table, &u.Name)   // with condition
	q2 := NewValueQuery[User](this, table, &u.Code)   // with condition
	q3 := NewValueQuery[User](this, table, &u.Job)    //  no condition
	q4 := NewValueQuery[User](this, table, &u.Extra)  // blank column
	q5 := NewValueQuery[User](this, table, &u.secret) // private field
	// ValueQuery.Where
	q1.Where(Equal[User](this, &u.Code, "admin"))
	q2.Where(Equal[User](this, &u.ID, 2))
	// ValueQuery.Test
	u1 := User{1, "Admin", "admin", "dev", "", "123"}
	u2 := User{2, "Guest", "guest", "dev", "", "456"}
	testCases := []tst.P2W1[*ValueQuery[User, string], User, bool]{
		{q1, u1, true}, {q1, u2, false},
		{q2, u1, false}, {q2, u2, true},
		{q3, u1, false}, {q3, u2, false},
	}
	tst.AllP2W1(t, testCases, "ValueQuery.Test", (*ValueQuery[User, string]).Test, tst.AssertEqual)
	// ValueQuery.BuildQuery
	testCases2 := []tst.P1W2[*ValueQuery[User, string], string, []any]{
		{q0, "", []any{}},
		{q1, "SELECT `Name` FROM `users` WHERE `Code` = ?", []any{"admin"}},
		{q2, "SELECT `Code` FROM `users` WHERE `ID` = ?", []any{2}},
		{q3, "SELECT `Job` FROM `users` WHERE false", []any{}},
		{q4, "", []any{}},
		{q5, "", []any{}},
	}
	tst.AllP1W2(t, testCases2, "ValueQuery.BuildQuery", (*ValueQuery[User, string]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)
}

func TestSelectRowQuery(t *testing.T) {
	type Company struct {
		ID       int
		Name     string
		Code     string
		Age      int
		Type     string `col:"Kind"`
		Extra    string `col:"-"`
		password string
	}
	c := new(Company)
	this := testPrelude(t, c)
	table := "companies"
	cols2 := this.Columns(&c.Name, &c.Type)
	cols3 := this.Columns(&c.Name, &c.Code)
	cols4 := this.Columns(&c.ID)
	cols8 := append([]string{}, cols4...)
	cols8 = append(cols8, this.Column(&c.password))
	cols9 := append([]string{}, cols4...)
	cols9 = append(cols9, this.Column(&c.Extra))
	reader1 := NewRowReader[Company](this, this.allColumns(c)...)
	reader2 := NewRowReader[Company](this, cols2...)
	reader3 := NewRowReader[Company](this, cols3...)
	reader4 := NewRowReader[Company](this, cols4...)
	// NewSelectRowQuery, NewFullSelectRowQuery
	q0 := NewSelectRowQuery[Company](this, "", nil)    // no table
	q1 := NewFullSelectRowQuery(this, table, reader1)  // all columns
	q2 := NewSelectRowQuery(this, table, reader2)      // specific columns
	q3 := NewSelectRowQuery(this, table, reader3)      // specific columns
	q4 := NewSelectRowQuery(this, table, reader4)      // one column
	q5 := NewSelectRowQuery[Company](this, table, nil) // nil reader
	q6 := NewSelectRowQuery(this, table, reader4)      // no columns
	q7 := NewSelectRowQuery(this, table, reader4)      // no condition
	q8 := NewSelectRowQuery(this, table, reader4)      // with private column
	q9 := NewSelectRowQuery(this, table, reader4)      // with blank column
	// SelectRowQuery.Columns
	q2.Columns(this, cols2...)
	q3.Columns(this, cols3...)
	q4.Columns(this, cols4...)
	q7.Columns(this, cols4...)
	q8.Columns(this, cols8...)
	q9.Columns(this, cols9...)
	// SelectRowQuery.Where
	q1.Where(Equal[Company](this, &c.Name, "Google"))
	q2.Where(Equal[Company](this, &c.Code, "XYZ"))
	q3.Where(Greater[Company](this, &c.Age, 10))
	q4.Where(NotIn[Company](this, &c.Type, []string{"IT", "Finance"}))
	// SelectRowQuery.Test
	c1 := Company{1, "Google", "GGL", 25, "IT", "", "search"}
	c2 := Company{2, "Unknown", "XYZ", 5, "Finance", "", "banks"}
	c3 := Company{3, "GoldStar", "GS", 7, "Mining", "", "nuggets"}
	c4 := Company{4, "Apple", "APL", 20, "IT", "", "ios"}
	testCases := []tst.P2W1[*SelectRowQuery[Company], Company, bool]{
		{q1, c1, true}, {q1, c2, false}, {q1, c3, false}, {q1, c4, false},
		{q2, c1, false}, {q2, c2, true}, {q2, c3, false}, {q2, c4, false},
		{q3, c1, true}, {q3, c2, false}, {q3, c3, false}, {q3, c4, true},
		{q4, c1, false}, {q4, c2, false}, {q4, c3, true}, {q4, c4, false},
	}
	tst.AllP2W1(t, testCases, "SelectRowQuery.Test", (*SelectRowQuery[Company]).Test, tst.AssertEqual)
	// SelectRowQuery.BuildQuery
	testCases2 := []tst.P1W2[*SelectRowQuery[Company], string, []any]{
		{q0, "", []any{}},
		{q1, "SELECT `ID`, `Name`, `Code`, `Age`, `Kind` FROM `companies` WHERE `Name` = ? LIMIT 1", []any{"Google"}},
		{q2, "SELECT `Name`, `Kind` FROM `companies` WHERE `Code` = ? LIMIT 1", []any{"XYZ"}},
		{q3, "SELECT `Name`, `Code` FROM `companies` WHERE `Age` > ? LIMIT 1", []any{10}},
		{q4, "SELECT `ID` FROM `companies` WHERE `Kind` NOT IN (?, ?) LIMIT 1", []any{"IT", "Finance"}},
		{q5, "", []any{}},
		{q6, "", []any{}},
		{q7, "SELECT `ID` FROM `companies` WHERE false LIMIT 1", []any{}},
		{q8, "SELECT `ID` FROM `companies` WHERE false LIMIT 1", []any{}},
		{q9, "SELECT `ID` FROM `companies` WHERE false LIMIT 1", []any{}},
	}
	tst.AllP1W2(t, testCases2, "SelectRowQuery.BuildQuery", (*SelectRowQuery[Company]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)
}

func TestTopRowQuery(t *testing.T) {
	type User struct {
		Name    string
		Code    string
		Age     int
		Balance float64
		secret1 string
		secret2 string
	}
	u := new(User)
	this := testPrelude(t, u)
	table := "users"
	cols1 := this.allColumns(u)
	cols2 := this.Columns(&u.Code, &u.Age)
	cols3 := this.Columns(&u.secret1, &u.secret2)
	reader1 := NewRowReader[User](this, cols1...)
	reader2 := NewRowReader[User](this, cols2...)
	reader3 := NewRowReader[User](this, cols3...)
	// NewTopRowQuery
	q0 := NewTopRowQuery[User](this, "", nil)  // no table
	q1 := NewTopRowQuery(this, table, reader1) // all columns
	q2 := NewTopRowQuery(this, table, reader2) // selected columns
	q3 := NewTopRowQuery(this, table, reader2) // no condition
	q4 := NewTopRowQuery(this, table, reader2) // no order
	q5 := NewTopRowQuery(this, table, reader3) // no columns
	// TopRowQuery.Columns
	q2.Columns(this, cols2...)
	q3.Columns(this, cols2...)
	q5.Columns(this, cols3...)
	// TopRowQuery.OrderAsc, OrderDesc, Limit
	q1.OrderDesc(this, this.Column(&u.Age)).OrderAsc(this, this.Column(&u.Name)).Limit(5)
	q2.OrderDesc(this, this.Column(&u.Balance))
	q3.OrderDesc(this, this.Column(&u.Balance))
	q5.OrderAsc(this, this.Column(&u.Age))
	// TopRowQuery.Where
	q1.Where(Greater[User](this, &u.Balance, 0))
	q2.Where(Greater[User](this, &u.Age, 10))
	q5.Where(Greater[User](this, &u.Age, 0))
	// TopRowQuery.Test
	u1 := User{"John", "john", 20, 5.0, "x", "y"}
	u2 := User{"Jean", "jean", 5, 0.0, "z", "z"}
	u3 := User{"Jack", "jack", 10, 15.0, "d", "r"}
	testCases := []tst.P2W1[*TopRowQuery[User], User, bool]{
		{q1, u1, true}, {q1, u2, false}, {q1, u3, true},
		{q2, u1, true}, {q2, u2, false}, {q2, u3, false},
	}
	tst.AllP2W1(t, testCases, "TopRowQuery.Test", (*TopRowQuery[User]).Test, tst.AssertEqual)
	// TopRowQuery.BuildQuery
	testCases2 := []tst.P1W2[*TopRowQuery[User], string, []any]{
		{q0, "", []any{}},
		{q1, "SELECT `Name`, `Code`, `Age`, `Balance` FROM `users` WHERE `Balance` > ? ORDER BY `Age` DESC, `Name` ASC LIMIT 5", []any{0.0}},
		{q2, "SELECT `Code`, `Age` FROM `users` WHERE `Age` > ? ORDER BY `Balance` DESC LIMIT 1", []any{10}},
		{q3, "SELECT `Code`, `Age` FROM `users` WHERE false ORDER BY `Balance` DESC LIMIT 1", []any{}},
		{q4, "", []any{}},
		{q5, "", []any{}},
	}
	tst.AllP1W2(t, testCases2, "TopRowQuery.BuildQuery", (*TopRowQuery[User]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)
}

func TestTopValueQuery(t *testing.T) {
	type User struct {
		ID      int
		Name    string
		Code    string
		Age     int
		Balance float64
		secret  string
	}
	u := new(User)
	this := testPrelude(t, u)
	table := "users"

	// NewTopValueQuery
	q0 := NewTopValueQuery[User, string](this, "", &u.Name)      // no table
	q1 := NewTopValueQuery[User, string](this, table, &u.Name)   // success
	q2 := NewTopValueQuery[User, int](this, table, &u.Age)       // success, another column
	q6 := NewTopValueQuery[User, string](this, table, &u.Code)   // success, another string column
	q3 := NewTopValueQuery[User, string](this, table, &u.Code)   // no condition
	q4 := NewTopValueQuery[User, string](this, table, &u.Code)   // no order
	q5 := NewTopValueQuery[User, string](this, table, &u.secret) // private field

	// TopValueQuery.OrderAsc, OrderDesc, Limit
	q1.OrderAsc(this, this.Column(&u.Name)).OrderDesc(this, this.Column(&u.Age)).Limit(3)
	q3.OrderAsc(this, this.Column(&u.Age))
	q2.OrderDesc(this, this.Column(&u.Balance))
	q6.OrderDesc(this, this.Column(&u.Age))
	q5.OrderAsc(this, this.Column(&u.Age))

	// TopValueQuery.Where
	q1.Where(Greater[User](this, &u.Balance, 0))
	q2.Where(Greater[User](this, &u.Age, 10))
	q5.Where(Greater[User](this, &u.Age, 0))
	q6.Where(Greater[User](this, &u.Age, 10))

	// TopValueQuery.Test
	u1 := User{1, "John", "john", 20, 5.0, "x"}
	u2 := User{2, "Jean", "jean", 5, 0.0, "z"}
	u3 := User{3, "Jack", "jack", 10, 15.0, "d"}
	testCases := []tst.P2W1[*TopValueQuery[User, string], User, bool]{
		{q1, u1, true}, {q1, u2, false}, {q1, u3, true},
		{q6, u1, true}, {q6, u2, false}, {q6, u3, false},
	}
	tst.AllP2W1(t, testCases, "TopValueQuery.Test", (*TopValueQuery[User, string]).Test, tst.AssertEqual)

	// TopValueQuery.BuildQuery
	testCases2 := []tst.P1W2[*TopValueQuery[User, string], string, []any]{
		{q0, "", []any{}},
		{q1, "SELECT `Name` FROM `users` WHERE `Balance` > ? ORDER BY `Name` ASC, `Age` DESC LIMIT 3", []any{0.0}},
		{q3, "SELECT `Code` FROM `users` WHERE false ORDER BY `Age` ASC LIMIT 1", []any{}},
		{q4, "", []any{}},
		{q5, "", []any{}},
		{q6, "SELECT `Code` FROM `users` WHERE `Age` > ? ORDER BY `Age` DESC LIMIT 1", []any{10}},
	}
	tst.AllP1W2(t, testCases2, "TopValueQuery.BuildQuery", (*TopValueQuery[User, string]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)

	// Test TopValueQuery with different value type
	testCases3 := []tst.P1W2[*TopValueQuery[User, int], string, []any]{
		{q2, "SELECT `Age` FROM `users` WHERE `Age` > ? ORDER BY `Balance` DESC LIMIT 1", []any{10}},
	}
	tst.AllP1W2(t, testCases3, "TopValueQuery.BuildQuery (int)", (*TopValueQuery[User, int]).BuildQuery, tst.AssertEqual, tst.AssertListEqual)
}

func TestSumQuery(t *testing.T) {
	// TODO: NewSumQuery
	// TODO: SumQuery.Where
	// TODO: SumQuery without condition (optional)
	// TODO: SumQuery.Columns
	// TODO: SumQuery.BuildQuery
}

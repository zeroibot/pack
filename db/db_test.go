package db

import (
	"database/sql"
	"testing"

	"github.com/zeroibot/tst"
)

func TestDBInterfaces(t *testing.T) {
	// Conn interface
	dbc1 := NewAdapter(new(sql.DB))
	dbc2 := NewMockAdapter(tst.NewConn[any]())
	var conn1 Conn = dbc1
	var conn2 Conn = dbc2
	tst.All(t, [][2]any{{conn1, dbc1}, {conn2, dbc2}}, "Conn", tst.AssertDeepEqual)

	// Tx interface
	tx1 := new(sql.Tx)
	tx2 := tst.NewTx()
	var tx3 Tx = tx1
	var tx4 Tx = tx2
	tst.All(t, [][2]any{{tx3, tx1}, {tx4, tx2}}, "Tx", tst.AssertDeepEqual)

	// Row interface
	row1 := new(sql.Row)
	row2 := tst.NewRow()
	var row3 Row = row1
	var row4 Row = row2
	tst.All(t, [][2]any{{row3, row1}, {row4, row2}}, "Row", tst.AssertDeepEqual)

	// Rows interface
	rows1 := new(sql.Rows)
	rows2 := tst.NewRows()
	var rows3 Rows = rows1
	var rows4 Rows = rows2
	tst.All(t, [][2]any{{rows3, rows1}, {rows4, rows2}}, "Rows", tst.AssertDeepEqual)

	// RowScanner interface
	var scan1 RowScanner = row1
	var scan2 RowScanner = row2
	var scan3 RowScanner = rows1
	var scan4 RowScanner = rows2
	pairs := [][2]any{{scan1, row1}, {scan2, row2}, {scan3, rows1}, {scan4, rows2}}
	tst.All(t, pairs, "RowScanner", tst.AssertDeepEqual)
}

func TestAdapters(t *testing.T) {
	cfg := new(ConnParams{
		Type:     "mysql",
		Host:     "localhost",
		Port:     "3306",
		User:     "test",
		Password: "abcd1234",
		Database: "test",
	})
	conn, err := NewSQLConnection(cfg)
	if err != nil {
		t.Fatal(err)
	}

	query, values := "SELECT `ID` FROM users WHERE `Name` = ?", []any{"John"}
	dbc1 := NewMockAdapter(tst.NewConn[any]())
	dbc2 := NewAdapter(conn)
	for _, dbc := range []Conn{dbc1, dbc2} {
		row := dbc.QueryRow(query, values...)
		tst.AssertTrue(t, "QueryRow", row != nil)
		rows, err := dbc.Query(query, values...)
		tst.AssertTrue(t, "Query", rows != nil && err != nil)
		result, err := dbc.Exec(query, values...)
		tst.AssertTrue(t, "Exec", result == nil && err != nil)
		tx, err := dbc.Begin()
		tst.AssertTrue(t, "Begin", tx != nil && err == nil)
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("cannot close db: %v", err)
	}
}

func TestNewSQLConnection(t *testing.T) {
	clone := func(cfg ConnParams) ConnParams {
		return ConnParams{cfg.Type, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database}
	}
	cfg := ConnParams{
		Type:     "mysql",
		Host:     "localhost",
		Port:     "3306",
		User:     "test",
		Password: "abcd1234",
		Database: "test",
	}
	cfg2 := ConnParams{}
	cfg3 := clone(cfg)
	cfg3.Database = "nonexistent"
	cfg4 := clone(cfg)
	cfg4.Password = "wrong_password"
	cfg5 := clone(cfg)
	cfg5.User = "wrong_user"
	cfg6 := clone(cfg)
	cfg6.Port = "6769"
	cfg7 := clone(cfg)
	cfg7.Port = "xxxx"
	cfg8 := clone(cfg)
	cfg8.Host = "wrong_host"
	cfg9 := clone(cfg)
	cfg9.Type = "mongodb"

	testCases := []tst.P1W2[*ConnParams, bool, bool]{
		{&cfg, true, true},    // success
		{nil, false, false},   // nil params
		{&cfg2, false, false}, // has blank params
		{&cfg3, false, false}, // wrong db
		{&cfg4, false, false}, // wrong password
		{&cfg5, false, false}, // wrong user
		{&cfg6, false, false}, // wrong port number
		{&cfg7, false, false}, // wrong port format
		{&cfg8, false, false}, // wrong host
		{&cfg9, false, false}, // wrong type
	}
	newConn := func(cfg *ConnParams) (bool, bool) {
		conn, err := NewSQLConnection(cfg)
		return conn != nil, err == nil
	}
	tst.AllP1W2(t, testCases, "NewSQLConnection", newConn, tst.AssertEqual[bool], tst.AssertEqual[bool])
}

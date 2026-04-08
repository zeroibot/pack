// Package my contains the Request object
package my

import (
	"net/http"

	"github.com/zeroibot/pack/db"
	"github.com/zeroibot/pack/ds"
)

const (
	OK200  int = http.StatusOK                  // OK
	OK201  int = http.StatusCreated             // Created
	Err400 int = http.StatusBadRequest          // Client-Side Error
	Err401 int = http.StatusUnauthorized        // Unauthenticated
	Err403 int = http.StatusForbidden           // Unauthorized
	Err404 int = http.StatusNotFound            // Not Found
	Err429 int = http.StatusTooManyRequests     // Rate Limiting
	Err500 int = http.StatusInternalServerError // Server-Side Error
)

// Instance type stores the db connection pools for requests
type Instance struct {
	dbConn    db.Conn            // db connection pool
	dbConnMap map[string]db.Conn // map of custom db connection pools
}

// NewInstance creates a new Instance object
func NewInstance(p *db.ConnParams) ds.Result[*Instance] {
	dbc := db.NewSQLConnection(p)
	if dbc.IsError() {
		return ds.Error[*Instance](dbc.Error())
	}
	dbAdapter := db.NewAdapter(dbc.Value())
	instance := new(Instance{
		dbConn:    dbAdapter,
		dbConnMap: make(map[string]db.Conn),
	})
	return ds.NewResult(instance, nil)
}

// AddConnection adds a custom DB connection
func (i *Instance) AddConnection(name string, p *db.ConnParams) error {
	dbc := db.NewSQLConnection(p)
	if dbc.IsError() {
		return dbc.Error()
	}
	i.dbConnMap[name] = db.NewAdapter(dbc.Value())
	return nil
}

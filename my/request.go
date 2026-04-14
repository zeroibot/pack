package my

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/zeroibot/pack/clock"
	"github.com/zeroibot/pack/db"
	"github.com/zeroibot/pack/dict"
	"github.com/zeroibot/pack/qb"
	"github.com/zeroibot/pack/str"
)

var (
	errNoDBConn = errors.New("no db connection")
	errNoDBTx   = errors.New("no db transaction")
)

const MainDB string = "main"

// Task contains the Action and Target of the task
type Task struct {
	Action string
	Target string
}

// String returns the string representation of the task
func (t Task) String() string {
	target := t.Target
	if strings.HasSuffix(target, "-%s") {
		parts := str.CleanSplit(target, "-")
		target = strings.Join(parts[:len(parts)-1], "-")
	}
	return fmt.Sprintf("%s-%s", t.Action, target)
}

// Request is an application request that holds a DB connection, a transaction,
// result checker, transaction queries, request start time, and logs
type Request struct {
	Task
	Name    string
	Params  dict.Object
	DB      db.Conn
	Tx      db.Tx
	Checker qb.ResultChecker
	Status  int
	Now     clock.DateTime
	// Private fields
	this    *Instance
	start   clock.DateTime
	txSteps []qb.Query
	dbMap   map[string]db.Conn
	// Logs
	mu   sync.RWMutex
	logs []string
}

// NewRequest creates a new Request
func (i *Instance) NewRequest(name string, args ...any) (*Request, error) {
	if len(args) > 0 {
		name = fmt.Sprintf(name, args...)
	}
	rq := newRequest(i, name)
	if i.dbConn == nil {
		rq.Status = Err500
		return nil, errNoDBConn
	}
	rq.DB = i.dbConn
	rq.dbMap[MainDB] = i.dbConn
	return rq, nil
}

// NewRequestAt creates a new Request at custom db
func (i *Instance) NewRequestAt(key, name string, args ...any) (*Request, error) {
	if len(args) > 0 {
		name = fmt.Sprintf(name, args...)
	}
	rq := newRequest(i, name)
	conn, ok := i.dbConnMap[key]
	if !ok || conn == nil {
		rq.Status = Err500
		return nil, errNoDBConn
	}
	rq.DB = conn
	rq.dbMap[MainDB] = conn
	return rq, nil
}

// Create a new Request object
func newRequest(this *Instance, name string) *Request {
	return new(Request{
		this:   this,
		Name:   name,
		Params: make(dict.Object),
		Status: OK200,
		start:  clock.DateTimeNow(),
		logs:   make([]string, 0),
		dbMap:  make(map[string]db.Conn),
	})
}

// AddDB adds a database connection to the request
func (rq *Request) AddDB(key string) bool {
	conn, ok := rq.this.dbConnMap[key]
	if !ok {
		return false
	}
	rq.dbMap[key] = conn
	return true
}

// SwitchDB changes the database connection of the request
func (rq *Request) SwitchDB(key string) bool {
	conn, ok := rq.dbMap[key]
	if !ok {
		return false
	}
	rq.DB = conn
	return true
}

// AddLog adds a log message to the request
func (rq *Request) AddLog(message string) {
	message = fmt.Sprintf("%s | %s", clock.DateTimeNow(), message)
	rq.logs = append(rq.logs, message)
}

// AddFmtLog adds a formatted log message to the request
func (rq *Request) AddFmtLog(format string, args ...any) {
	rq.AddLog(fmt.Sprintf(format, args...))
}

// AddDurationLog adds a duration log message to the request
func (rq *Request) AddDurationLog(start time.Time) {
	rq.AddLog(fmt.Sprintf("Time: %v", time.Since(start)))
}

// AddErrorLog adds an error log message to the request
func (rq *Request) AddErrorLog(err error) {
	rq.AddLog(fmt.Sprintf("Error: %s", err.Error()))
}

// SetNow sets the Request.Now field
func (rq *Request) SetNow() clock.DateTime {
	rq.Now = clock.DateTimeNow()
	return rq.Now
}

// Fail adds an log message and sets the status of the request
func (rq *Request) Fail(status int, message string, args ...any) {
	rq.Status = status
	rq.AddFmtLog(message, args...)
}

// Output combines the logs into a single string separated by newline
func (rq *Request) Output() string {
	return strings.Join(rq.logs, "\n")
}

// SubRequest creates a sub-request for concurrent tasks
func (rq *Request) SubRequest() *Request {
	return new(Request{
		Task:   rq.Task,
		Params: rq.Params,
		DB:     rq.DB,
		Status: OK200,
		logs:   make([]string, 0),
	})
}

// MergeLogs performs a concurrent-safe merging of subrequest logs to main request logs
func (rq *Request) MergeLogs(srq *Request) {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	rq.logs = append(rq.logs, srq.logs...)
}

// StartTransaction starts a database transaction
func (rq *Request) StartTransaction(numSteps int) error {
	if rq.DB == nil {
		rq.AddLog("No DB connection")
		rq.Status = Err500
		return errNoDBConn
	}
	tx, err := rq.DB.Begin()
	if err != nil {
		rq.AddLog("Failed to start transaction")
		rq.Status = Err500
		return err
	}
	rq.Tx = tx
	rq.txSteps = make([]qb.Query, 0, numSteps)
	rq.Checker = qb.AssertNothing // default checker
	return nil
}

// CommitTransaction commits the database transaction
func (rq *Request) CommitTransaction() error {
	if rq.DB == nil {
		rq.AddLog("No DB connection")
		rq.Status = Err500
		return errNoDBConn
	}
	if rq.Tx == nil {
		rq.AddLog("No DB transaction")
		rq.Status = Err500
		return errNoDBTx
	}
	err := rq.Tx.Commit()
	if err != nil {
		// Add transaction steps to logs
		for i, q := range rq.txSteps {
			rq.AddFmtLog("Query %d: %s", i+1, qb.ToString(q))
		}
		rq.AddLog("Failed to commit transaction")
		rq.Status = Err500
		return fmt.Errorf("tx commit error: %w", err)
	}
	return nil
}

// AddTxStep adds a transaction step to the request
func (rq *Request) AddTxStep(q qb.Query) {
	rq.txSteps = append(rq.txSteps, q)
}

package tzt

import "database/sql"

type Tx struct{}

func NewTx() *Tx {
	return new(Tx)
}

func (t Tx) Exec(query string, args ...any) (sql.Result, error) {
	return nil, nil
}

func (t Tx) Commit() error {
	return nil
}

func (t Tx) Rollback() error {
	return nil
}

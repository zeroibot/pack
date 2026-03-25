package tzt

import (
	"fmt"
	"reflect"
)

type Row struct {
	items []any
}

func NewRow(items ...any) *Row {
	return new(Row{items: items})
}

func (r *Row) Scan(fieldRefs ...any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic encountered: %v", r)
		}
	}()
	if len(fieldRefs) != len(r.items) {
		return fmt.Errorf("expected %d fieldRefs, got %d", len(r.items), len(fieldRefs))
	}
	for i, fieldRef := range fieldRefs {
		fieldValue := reflect.ValueOf(fieldRef).Elem()
		fieldValue.Set(reflect.ValueOf(r.items[i]))
	}
	return err
}

type Rows struct {
	items [][]any
	index int
	err   error
}

func NewRows(items ...[]any) *Rows {
	return new(Rows{items: items, index: 0, err: nil})
}

func (r *Rows) Scan(fieldRefs ...any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic encountered: %v", r)
		}
	}()
	if !r.Next() {
		return fmt.Errorf("no more rows")
	}
	if len(fieldRefs) != len(r.items[r.index]) {
		return fmt.Errorf("expected %d fieldRefs, got %d", len(r.items[r.index]), len(fieldRefs))
	}
	for i, fieldRef := range fieldRefs {
		fieldValue := reflect.ValueOf(fieldRef).Elem()
		fieldValue.Set(reflect.ValueOf(r.items[r.index][i]))
	}
	r.index += 1
	return err
}

func (r *Rows) Next() bool {
	return r.index < len(r.items)
}

func (r *Rows) Err() error {
	return r.err
}

func (r *Rows) Close() error {
	return nil
}

func (r *Rows) SetError(err error) {
	r.err = err
}

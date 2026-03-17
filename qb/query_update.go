package qb

import (
	"fmt"
	"strings"

	"github.com/roidaradal/pack/dict"
	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/dyn"
)

type UpdateQuery[T any] struct {
	conditionQuery[T]
	orderedLimit
	typeName string
	updates  []ds.Option[columnValuePair]
}

type FieldUpdate [2]any                  // [OldValue, NewValue]
type FieldUpdates map[string]FieldUpdate // {FieldName => [OldValue, NewValue]}

// Unpack returns the oldValue and newValue of the FieldUpdate
func (f FieldUpdate) Unpack() (any, any) {
	return f[0], f[1]
}

// NewUpdateQuery creates a new UpdateQuery
func NewUpdateQuery[T any](this *Instance, table string) *UpdateQuery[T] {
	q := new(UpdateQuery[T])
	q.initializeRequired(this, table)
	var item T
	q.typeName = dyn.TypeName(item)
	q.updates = make([]ds.Option[columnValuePair], 0)
	return q
}

// Update adds a column=value update to the UpdateQuery
func Update[T, V any](this *Instance, q *UpdateQuery[T], fieldRef *V, value V) {
	// Note: cannot be a method because generics are not yet supported in methods
	pairOption := newColumnValue(this, fieldRef, value)
	q.updates = append(q.updates, pairOption)
}

// Update adds a column=value update to the UpdateQuery.
// Note: We lose type-checking of value here, so must be sure that field=value are of the same type.
func (q *UpdateQuery[T]) Update(this *Instance, fieldName string, value any) {
	pairOption := newFieldColumnValue(this, q.typeName, fieldName, value)
	q.updates = append(q.updates, pairOption)
}

// Updates adds column=value updates to the UpdateQuery
// Note: We lose type-checking of value here, so must be sure that field=value are of the same type.
func (q *UpdateQuery[T]) Updates(this *Instance, updates FieldUpdates) {
	fieldNames := dict.SortedKeys(updates)
	for _, fieldName := range fieldNames {
		_, newValue := updates[fieldName].Unpack()
		q.Update(this, fieldName, newValue)
	}
}

// BuildQuery returns the query string and parameter values of UpdateQuery
func (q *UpdateQuery[T]) BuildQuery() (string, []any) {
	numUpdates := len(q.updates)
	condition, conditionValues, err := q.conditionQuery.preBuildCheck()
	if err != nil || numUpdates == 0 {
		return emptyQueryValues()
	}
	values := make([]any, 0, numUpdates+len(conditionValues))
	updates := make([]string, numUpdates)
	for i, pairOption := range q.updates {
		if pairOption.IsNil() {
			// One column=value pair failed = return empty query
			return emptyQueryValues()
		}
		column, value := pairOption.Value().Unpack()
		if column == "" {
			// Blank column = return empty query
			return emptyQueryValues()
		}
		updates[i] = fmt.Sprintf("%s = ?", column)
		values = append(values, value)
	}
	// Add the condition values to the values list
	values = append(values, conditionValues...)
	update := strings.Join(updates, ", ")
	query := "UPDATE %s SET %s WHERE %s"
	query = fmt.Sprintf(query, q.table, update, condition)
	query = tryAppend(query, q.mustLimitString())
	return query, values
}

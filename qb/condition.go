package qb

import "github.com/roidaradal/pack/ds"

// Condition interface unifies all Condition objects:
// Build() method outputs the condition string and parameter values
type Condition interface {
	Build() (string, ds.List[any]) // Return (condition string, parameter values)
}

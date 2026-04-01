// Package fail contains common errors
package fail

import "errors"

var (
	MismatchCount = errors.New("public: Count does not match")
	MissingParams = errors.New("public: Missing required parameters")
	NotFoundItem  = errors.New("public: Item not found")
)

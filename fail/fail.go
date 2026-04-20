// Package fail contains common errors
package fail

import (
	"errors"
	"fmt"
	"strings"
)

var (
	InvalidSession = errors.New("public: Invalid session")
	MismatchCount  = errors.New("public: Count does not match")
	MissingParams  = errors.New("public: Missing required parameters")
	NotAuthorized  = errors.New("public: Not authorized")
	NotFoundItem   = errors.New("public: Item not found")
)

// FromErrors produces a single error from the list of errors
func FromErrors(label string, errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%s: %d errors encountered: %w", label, len(errs), errs[0])
}

// PublicMessage gets the public error message from "public: <message>" error
func PublicMessage(err error) (string, bool) {
	msg := err.Error()
	if strings.HasPrefix(msg, "public:") {
		parts := strings.Split(msg, ":")
		return strings.TrimSpace(parts[1]), true
	}
	return "Unexpected error occurred", false
}

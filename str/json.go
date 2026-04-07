package str

import (
	"encoding/json"
	"strings"
)

// JSON creates a JSON string from struct
func JSON[T any](item T) (string, error) {
	bytes, err := MarshalJSON(item, 0)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// IndentedJSON creates an indented JSON string from struct
func IndentedJSON[T any](item T, tabLength int) (string, error) {
	bytes, err := MarshalJSON(item, tabLength)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// MarshalJSON marshals the struct as JSON, with tabLength indicating whether to indent or not
func MarshalJSON[T any](item T, tabLength int) ([]byte, error) {
	if tabLength <= 0 {
		return json.Marshal(item)
	}
	return json.MarshalIndent(item, "", strings.Repeat(" ", tabLength))
}

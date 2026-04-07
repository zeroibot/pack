package io

import (
	"os"

	"github.com/zeroibot/pack/str"
)

// SaveJSON saves JSON object to given file path
func SaveJSON[T any](item T, path string) error {
	return saveJSON(item, path, 0)
}

// SaveIndentedJSON saves indented JSON object to given file path
func SaveIndentedJSON[T any](item T, path string, tabLength int) error {
	return saveJSON(item, path, tabLength)
}

// Common: save JSONto file path, with given tab length as indentation
func saveJSON[T any](item T, path string, tabLength int) error {
	bytes, err := str.MarshalJSON(item, tabLength)
	if err != nil {
		return err
	}
	err = EnsurePathExists(path)
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, defaultFileMode)
}

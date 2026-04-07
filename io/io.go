package io

import (
	"os"
	"path/filepath"
	"strings"
)

const defaultFileMode = 0o644

// IsDir checks if a path is a directory
func IsDir(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	return f.IsDir()
}

// PathExists checks if a path exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// EnsurePathExists creates all non-existent parent directories of the given path
func EnsurePathExists(path string) error {
	return os.MkdirAll(filepath.Dir(path), defaultFileMode)
}

// BaseFilename returns the base filename of a path
func BaseFilename(path string) string {
	filename := filepath.Base(path)
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

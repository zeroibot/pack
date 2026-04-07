package io

import (
	"encoding/json"
	"os"

	"github.com/zeroibot/pack/ds"
)

// ReadJSON reads a JSON object of given type from given path
func ReadJSON[T any](path string) ds.Result[T] {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return ds.Error[T](err)
	}
	var item T
	err = json.Unmarshal(bytes, &item)
	if err != nil {
		return ds.Error[T](err)
	}
	return ds.NewResult(item, nil)
}

// ReadJSONList reads a JSON list of given type from given path
func ReadJSONList[T any](path string) ds.Result[[]T] {
	return ReadJSON[[]T](path)
}

// ReadJSONMap reads a JSON map of given value type from given path
func ReadJSONMap[V any](path string) ds.Result[map[string]V] {
	return ReadJSON[map[string]V](path)
}

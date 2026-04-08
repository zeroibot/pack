// Package conf contains app configuration types and methods
package conf

import (
	"strings"

	"github.com/zeroibot/pack/dict"
	"github.com/zeroibot/pack/ds"
	"github.com/zeroibot/pack/dyn"
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/lang"
	"github.com/zeroibot/pack/model"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/number"
	"github.com/zeroibot/pack/qb"
	"github.com/zeroibot/pack/str"
)

const (
	keyGlue  string = "."
	listGlue string = "|"
)

type KV struct {
	Key           string `col:"AppKey"`
	Value         string `col:"AppValue"`
	LastUpdatedAt model.DateTime
}

type Defaults struct {
	UintMap       map[string]uint
	IntMap        dict.Ints
	StringMap     dict.Strings
	StringListMap dict.StringLists
}

// Lookup loads a Config map from database
func Lookup(this *qb.Instance, rq *my.Request, appKeys []string, kvSchema *model.Schema[KV]) ds.Result[dict.Strings] {
	if kvSchema == nil {
		rq.Fail(my.Err500, "KV Schema is nil")
		return ds.Error[dict.Strings](fail.MissingParams)
	}
	kv := kvSchema.Ref
	q := qb.NewLookupQuery[KV](this, kvSchema.Table, &kv.Key, &kv.Value)
	q.Where(qb.In[KV](this, &kv.Key, appKeys))
	result := q.Lookup(this, rq.DB)
	if result.IsError() {
		rq.Fail(my.Err500, "Failed to load app config from db")
	}
	return result
}

// Create decorates a Config object with the contents of lookup
func Create[T any](cfg *T, lookup dict.Strings, defaults *Defaults) *T {
	for key := range defaults.UintMap {
		value := uintOrDefault(lookup, defaults.UintMap, key)
		dyn.SetStructField(cfg, getKey(key), value)
	}
	for key := range defaults.IntMap {
		value := intOrDefault(lookup, defaults.IntMap, key)
		dyn.SetStructField(cfg, getKey(key), value)
	}
	for key := range defaults.StringMap {
		value := stringOrDefault(lookup, defaults.StringMap, key)
		dyn.SetStructField(cfg, getKey(key), value)
	}
	for key := range defaults.StringListMap {
		value := stringListOrDefault(lookup, defaults.StringListMap, key)
		dyn.SetStructField(cfg, getKey(key), value)
	}
	return cfg
}

// getKey extracts the second part of <Domain>.<Key>
func getKey(fullKey string) string {
	parts := strings.Split(fullKey, keyGlue)
	if len(parts) != 2 {
		return fullKey
	}
	return parts[1]
}

// uintOrDefault tries to convert lookup[key] to uint, falls back to defaultValue[key]
func uintOrDefault(lookup dict.Strings, defaultValue map[string]uint, key string) uint {
	value := defaultValue[key]
	if lookupValue, ok := lookup[key]; ok {
		value = uint(number.ParseInt(lookupValue))
	}
	return value
}

// intOrDefault t ries to convert lookup[key] to int, falls back to defaultValue[key]
func intOrDefault(lookup dict.Strings, defaultValue dict.Ints, key string) int {
	value := defaultValue[key]
	if lookupValue, ok := lookup[key]; ok {
		value = number.ParseInt(lookupValue)
	}
	return value
}

// stringOrDefault tries to get lookup[key], falls back to defaultValue[key]
func stringOrDefault(lookup dict.Strings, defaultValue dict.Strings, key string) string {
	value, ok := lookup[key]
	return lang.Ternary(ok, value, defaultValue[key])
}

// stringListOrDefault tries to convert lookup[key] to []string, falls back to defaultValue[key]
func stringListOrDefault(lookup dict.Strings, defaultValue dict.StringLists, key string) []string {
	value, ok := lookup[key]
	return lang.Ternary(ok, str.CleanSplit(value, listGlue), defaultValue[key])
}

// Package daemon contains daemon types and functions
package daemon

import (
	"time"

	"github.com/zeroibot/pack/dict"
)

type Instance struct {
	start    *dict.SyncMap[string, time.Time]
	last     *dict.SyncMap[string, time.Time]
	duration *dict.SyncMap[string, time.Duration]
}

// NewInstance creates a new Daemon Instance
func NewInstance() *Instance {
	return new(Instance{
		start:    dict.NewSyncMap[string, time.Time](),
		last:     dict.NewSyncMap[string, time.Time](),
		duration: dict.NewSyncMap[string, time.Duration](),
	})
}

// Package daemon contains daemon types and functions
package daemon

import (
	"fmt"
	"time"

	"github.com/zeroibot/pack/clock"
	"github.com/zeroibot/pack/dict"
	"github.com/zeroibot/pack/ds"
	"github.com/zeroibot/pack/io"
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

/*
Note: DaemonConfig is expected to have this structure:
Use int for Interval and Margin values so we can disable a daemon by setting the value < 0.

type Config struct {
	<Domain> struct {
		<FeatureInterval> int
		...
	}
	...
}
*/

// LoadConfig loads a Daemon Config which is expected to follow the required structure.
// This checks if any Interval values are 0 (invalid).
func LoadConfig[T any](path string) ds.Result[T] {
	result := io.ReadJSON[T](path)
	if result.IsError() {
		return result
	}
	cfg := new(result.Value())
	cfgMap, err := dict.FromStruct[map[string]int](cfg)
	if err != nil {
		return ds.Error[T](err)
	}
	for key := range cfgMap {
		for cfgKey, value := range cfgMap[key] {
			if value == 0 {
				return ds.Error[T](fmt.Errorf("invalid daemon %s.%s: %d", key, cfgKey, value))
			}
		}
	}
	return result
}

// Run runs a task every given interval, where TimeScale = time.Hour, time.Minute, time.Second
func Run[A any](this *Instance, app *A, name string, task func(*A), interval int, timeScale time.Duration) {
	if interval < 0 {
		fmt.Printf("Daemon:%s is disabled\n", name)
		return
	}
	timeInterval := time.Duration(interval) * timeScale
	this.start.Set(name, time.Now())
	this.duration.Set(name, timeInterval)
	go func() {
		for {
			start := time.Now()
			this.last.Set(name, start)
			task(app)
			clock.Sleep(timeInterval, start)
		}
	}()
}

type Info struct {
	Start    string
	Last     string
	Duration string
}

// All returns info on all running daemons
func All(this *Instance) map[string]Info {
	startTimes := this.start.Map()
	lastTimes := this.last.Map()
	durations := this.duration.Map()
	info := make(map[string]Info)
	for name, startTime := range startTimes {
		start := clock.StandardFormat(startTime)
		var last, duration string
		if dict.HasKey(lastTimes, name) {
			last = clock.StandardFormat(lastTimes[name])
		}
		if dict.HasKey(durations, name) {
			duration = fmt.Sprintf("%v", durations[name])
		}
		info[name] = Info{start, last, duration}
	}
	return info
}

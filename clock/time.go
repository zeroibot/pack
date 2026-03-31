package clock

import "time"

// DateTimeNow returns the current datetime in standard format
func DateTimeNow() string {
	return StandardFormat(time.Now())
}

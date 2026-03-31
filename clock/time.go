package clock

import "time"

// DateTimeNow returns the current datetime in standard format
func DateTimeNow() DateTime {
	return StandardFormat(time.Now())
}

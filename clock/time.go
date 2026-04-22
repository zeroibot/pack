package clock

import "time"

// DateTimeNow returns the current datetime in standard format
func DateTimeNow() DateTime {
	return StandardFormat(time.Now())
}

// TimestampNow returns the current datetime in timestamp format
func TimestampNow() DateTime {
	return TimestampFormat(time.Now())
}

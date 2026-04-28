package clock

import (
	"fmt"
	"time"
)

// TimestampNow returns the current datetime in timestamp format
func TimestampNow() DateTime {
	return TimestampFormat(time.Now())
}

// TimeNow returns the current time in hh:mm:ss format
func TimeNow() Time {
	return TimeFormat(time.Now())
}

// HourMinNow returns the current time in hh:mm format
func HourMinNow() string {
	return HourMinFormat(time.Now())
}

// DateNow returns the current date in yyyy-mm-dd format
func DateNow() Date {
	return DateFormat(time.Now())
}

// DateTimeNow returns the current datetime in standard format
func DateTimeNow() DateTime {
	return StandardFormat(time.Now())
}

// DateTimeNowWithExpiry returns current datetime and expiry datetime in standard format
func DateTimeNowWithExpiry(duration time.Duration) (now, expiry DateTime) {
	timeNow := time.Now()
	timeExpiry := timeNow.Add(duration)
	now, expiry = StandardFormat(timeNow), StandardFormat(timeExpiry)
	return now, expiry
}

// DateTimeStart attaches 00:00:00 to given date
func DateTimeStart(date Date) string {
	return fmt.Sprintf("%s 00:00:00", date)
}

// DateTimeEnd attaches 23:59:59 to given date
func DateTimeEnd(date Date) string {
	return fmt.Sprintf("%s 23:59:59", date)
}

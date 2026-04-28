package clock

import (
	"fmt"
	"strings"
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

// Parse the given string as datetime in given format
func Parse(datetime, format string) (time.Time, error) {
	datetime = strings.TrimSpace(datetime)
	return time.Parse(format, datetime)
}

// ParseDateTime parses the given string as datetime in standard format
func ParseDateTime(datetime string) (time.Time, error) {
	datetime = strings.TrimSpace(datetime)
	return time.Parse(standardFmt, datetime)
}

// IsExpired checks if given datetime is already expired (before current time)
func IsExpired(expiry string) bool {
	limit, err := ParseDateTime(expiry)
	if err != nil {
		// Default to expired if invalid datetime
		return true
	}
	return time.Now().After(limit)
}

// ExtendTime returns the extended time in standard format
func ExtendTime(datetime string, duration time.Duration) (string, error) {
	dt, err := ParseDateTime(datetime)
	if err != nil {
		return "", err
	}
	return StandardFormat(dt.Add(duration)), nil
}

// DurationSince returns the duration as string since given datetime, rounded to given duration
func DurationSince(datetime string, round time.Duration) (string, error) {
	dt, err := ParseDateTime(datetime)
	if err != nil {
		return "", err
	}
	duration := time.Since(dt).Round(round)
	return fmt.Sprintf("%v", duration), nil
}

// ElapsedSince returns the duration since given datetime
func ElapsedSince(datetime string) (time.Duration, error) {
	dt, err := ParseDateTime(datetime)
	if err != nil {
		return 0, err
	}
	return time.Since(dt), nil
}

// ExtractTime extracts time part from given datetime string
func ExtractTime(datetime string) string {
	if !IsValidDateTime(datetime) {
		return ""
	}
	parts := strings.Fields(datetime)
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

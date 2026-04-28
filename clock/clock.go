// Package clock contains time and date-related functions
package clock

import (
	"strings"
	"time"
)

const (
	dateFmt      string = "2006-01-02"
	timeFmt      string = "15:04:05"
	hourMinFmt   string = "15:04"
	standardFmt  string = "2006-01-02 15:04:05"
	timestampFmt string = "060102150405"
)

type (
	Date     = string
	Time     = string
	DateTime = string
)

// Sleep pauses for a given duration, given the start time
func Sleep(pause time.Duration, start time.Time) {
	time.Sleep(pause - time.Since(start))
}

// DateFormat formats the given time in date format (yyyy-mm-dd)
func DateFormat(t time.Time) Date {
	return t.Format(dateFmt)
}

// TimeFormat formats the given time in time format (hh:mm:ss)
func TimeFormat(t time.Time) Time {
	return t.Format(timeFmt)
}

// HourMinFormat formats the given time in hh:mm format
func HourMinFormat(t time.Time) string {
	return t.Format(hourMinFmt)
}

// StandardFormat formats the given time in standard format (yyyy-mm-dd hh:mm:ss)
func StandardFormat(t time.Time) DateTime {
	return t.Format(standardFmt)
}

// TimestampFormat formats the given time in timestamp format (yymmddhhmmss)
func TimestampFormat(t time.Time) string {
	return t.Format(timestampFmt)
}

// IsValidDate checks if given yyyy-mm-dd date is valid
func IsValidDate(date string) bool {
	_, err := time.Parse(dateFmt, strings.TrimSpace(date))
	return err == nil
}

// IsValidDateTime checks if given datetime in standard format is valid
func IsValidDateTime(datetime string) bool {
	_, err := time.Parse(standardFmt, strings.TrimSpace(datetime))
	return err == nil
}

// ExtractDate returns the date part of the given datetime string
func ExtractDate(datetime string) string {
	if !IsValidDateTime(datetime) {
		return ""
	}
	return strings.Fields(datetime)[0]
}

// ExtractYearMonth returns the year and month part of the given datetime string
func ExtractYearMonth(datetime string) string {
	date := ExtractDate(datetime)
	if date == "" {
		return ""
	}
	return date[0:7] // yyyy-mm
}

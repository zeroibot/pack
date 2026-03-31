// Package clock contains time and date-related functions
package clock

import "time"

const (
	dateFmt      string = "2006-01-02"
	timeFmt      string = "15:04:05"
	hourMinFmt   string = "15:04"
	standardFmt  string = "2006-01-02 15:04:05"
	timestampFmt string = "060102150405"
)

// DateFormat formats the given time in date format (yyyy-mm-dd)
func DateFormat(t time.Time) string {
	return t.Format(dateFmt)
}

// TimeFormat formats the given time in time format (hh:mm:ss)
func TimeFormat(t time.Time) string {
	return t.Format(timeFmt)
}

// HourMinFormat formats the given time in hh:mm format
func HourMinFormat(t time.Time) string {
	return t.Format(hourMinFmt)
}

// StandardFormat formats the given time in standard format (yyyy-mm-dd hh:mm:ss)
func StandardFormat(t time.Time) string {
	return t.Format(standardFmt)
}

// TimestampFormat formats the given time in timestamp format (yymmddhhmmss)
func TimestampFormat(t time.Time) string {
	return t.Format(timestampFmt)
}

// Sleep pauses for a given duration, given the start time
func Sleep(pause time.Duration, start time.Time) {
	time.Sleep(pause - time.Since(start))
}

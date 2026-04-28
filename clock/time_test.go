package clock

import (
	"testing"
	"time"

	"github.com/zeroibot/tst"
)

func TestTime(t *testing.T) {
	tst.AssertEqual(t, "DateTimeNow", DateTimeNow(), StandardFormat(time.Now()))
	tst.AssertEqual(t, "DateNow", DateNow(), DateFormat(time.Now()))
	tst.AssertEqual(t, "TimestampNow", TimestampNow(), TimestampFormat(time.Now()))
	// TODO: DateTimeStart
	// TODO: DateTimeEnd
	// TODO: DateTimeNowWithExpiry
	// TODO: TimeNow
	// TODO: HourMinNow
	// TODO: Parse, ParseDateTime
	// TODO: IsExpired
	// TODO: ExtendTime
	// TODO: DurationSince, ElapsedSince
	// TODO: ExtractTime
}

package clock

import (
	"testing"
	"time"

	"github.com/zeroibot/tst"
)

func TestClockFormat(t *testing.T) {
	now := time.Now()
	tst.AssertEqual(t, "DateFormat", DateFormat(now), now.Format("2006-01-02"))
	tst.AssertEqual(t, "TimeFormat", TimeFormat(now), now.Format("15:04:05"))
	tst.AssertEqual(t, "HourMinFormat", HourMinFormat(now), now.Format("15:04"))
	tst.AssertEqual(t, "StandardFormat", StandardFormat(now), now.Format("2006-01-02 15:04:05"))
	tst.AssertEqual(t, "TimestampFormat", TimestampFormat(now), now.Format("060102150405"))
	Sleep(1*time.Second, now)
	// TODO: IsValidDate
	// TODO: IsValidDateTime
	// TODO: ExtractDate
	// TODO: ExtractYearMonth
}

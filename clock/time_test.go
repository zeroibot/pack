package clock

import (
	"testing"
	"time"

	"github.com/zeroibot/tst"
)

func TestTime(t *testing.T) {
	tst.AssertEqual(t, "DateTimeNow", DateTimeNow(), StandardFormat(time.Now()))
}

package ds

import (
	"testing"

	"github.com/zeroibot/tst"
)

func TestOption(t *testing.T) {
	value1 := 1
	ref1 := new(value1)
	opt1 := NewOption(ref1)
	// IsNil
	tst.AssertEqual(t, "Option.IsNil", opt1.IsNil(), false)
	// NotNil
	tst.AssertEqual(t, "Option.NotNil", opt1.NotNil(), true)
	// Get
	value, flag := opt1.Get()
	tst.AssertEqualAnd(t, "Option.Get", value, value1, flag, true)
	// Value
	tst.AssertEqual(t, "Option.Value", opt1.Value(), value1)
	// String
	tst.AssertEqual(t, "Option.String", opt1.String(), "1")

	// Nil Options
	var ref2 *int
	opt2 := NewOption(ref2)
	opt3 := Nil[int]()
	for _, opt := range []Option[int]{opt2, opt3} {
		value, flag = opt.Get()
		tst.AssertEqual(t, "Option.IsNil", opt.IsNil(), true)
		tst.AssertEqual(t, "Option.NotNil", opt.NotNil(), false)
		tst.AssertEqualAnd(t, "Option.Get", value, 0, flag, false)
		tst.AssertEqual(t, "Option.Value", opt.Value(), 0)
		tst.AssertEqual(t, "Option.String", opt.String(), "<nil>")
	}
}

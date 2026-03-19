package ds

import (
	"fmt"
	"testing"

	"github.com/roidaradal/tst"
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

func TestResult(t *testing.T) {
	// NewResult, Error
	errNotEven := fmt.Errorf("not an even number")
	res1 := NewResult[int](5, nil)
	res2 := Error[int](errNotEven)
	res3 := NewResult[int](67, nil)

	// IsError, NotError, Error
	testCases1 := []Tuple3[Result[int], bool, bool]{
		{res1, false, false},
		{res2, true, true},
		{res3, false, false},
	}
	for _, x := range testCases1 {
		result, want, notNilErr := x.Unpack()
		tst.AssertEqualError(t, "Result.IsError", result.IsError(), want, result.Error(), notNilErr)
		tst.AssertEqual(t, "Result.NotError", result.NotError(), !want)
	}

	// Get, Value, String
	testCases2 := []Tuple4[Result[int], int, bool, string]{
		{res1, 5, true, "5"},
		{res2, 0, false, fmt.Sprintf("error: %s", errNotEven.Error())},
		{res3, 67, true, "67"},
	}
	for _, x := range testCases2 {
		result, wantValue, wantOk, wantString := x.Unpack()
		actualValue, actualOk := result.Get()
		tst.AssertEqualAnd(t, "Result.Get", actualValue, wantValue, actualOk, wantOk)
		tst.AssertEqual(t, "Result.Value", result.Value(), wantValue)
		tst.AssertEqual(t, "Result.String", result.String(), wantString)
	}
}

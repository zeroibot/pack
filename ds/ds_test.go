package ds

import (
	"errors"
	"fmt"
	"testing"
)

func TestOption(t *testing.T) {
	value1 := 1
	ref1 := new(value1)
	opt1 := NewOption(ref1)

	flag := opt1.IsNil()
	if flag != false {
		t.Errorf("Option.IsNil() = %v, want = false", flag)
	}
	flag = opt1.NotNil()
	if flag != true {
		t.Errorf("Option.NotNil() = %v, want = true", flag)
	}
	value, flag := opt1.Get()
	if value != value1 || flag != true {
		t.Errorf("Option.Get() = (%v, %v), want = (%v, %v)", value, flag, value1, true)
	}
	value = opt1.Value()
	if value != value1 {
		t.Errorf("Option.Value() = %v, want = %v", value, value1)
	}
	text, want := opt1.String(), "1"
	if text != want {
		t.Errorf("Option.String() = %v, want = %v", text, want)
	}

	var ref2 *int
	opt2 := NewOption(ref2)
	opt3 := Nil[int]()
	for _, opt := range []Option[int]{opt2, opt3} {
		flag = opt.IsNil()
		if flag != true {
			t.Errorf("Option.IsNil() = %v, want = true", flag)
		}
		flag = opt.NotNil()
		if flag != false {
			t.Errorf("Option.NotNil() = %v, want = false", flag)
		}
		value, flag = opt.Get()
		if value != 0 || flag != false {
			t.Errorf("Option.Get() = (%v, %v), want = (0, false)", value, flag)
		}
		value = opt.Value()
		if value != 0 {
			t.Errorf("Option.Value() = %v, want = 0", value)
		}
		text, want = opt.String(), "<nil>"
		if text != want {
			t.Errorf("Option.String() = %v, want = %v", text, want)
		}
	}
}

func TestResult(t *testing.T) {
	// NewResult, Error
	errNotEven := fmt.Errorf("not an even number")
	res1 := NewResult[int](5, nil)
	res2 := Error[int](errNotEven)
	res3 := NewResult[int](67, nil)

	// IsError, NotError, Error
	testCases1 := []Tuple3[Result[int], bool, error]{
		{res1, false, nil},
		{res2, true, errNotEven},
		{res3, false, nil},
	}
	for _, x := range testCases1 {
		result, want, wantErr := x.Unpack()
		actual := result.IsError()
		if actual != want {
			t.Errorf("Result.IsError() = %t, want %t", actual, want)
		}
		want = !want
		actual = result.NotError()
		if actual != want {
			t.Errorf("Result.NotError() = %t, want %t", actual, want)
		}
		actualErr := result.Error()
		if errors.Is(actualErr, wantErr) == false {
			t.Errorf("Result.Error() = %v, want = %v", actualErr, wantErr)
		}
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
		if actualValue != wantValue || actualOk != wantOk {
			t.Errorf("Result.Get() = %v, %t, want %v, %t", actualValue, actualOk, wantValue, wantOk)
		}
		actualValue = result.Value()
		if actualValue != wantValue {
			t.Errorf("Result.Value() = %v, want = %v", actualValue, wantValue)
		}
		actualString := result.String()
		if actualString != wantString {
			t.Errorf("Result.String() = %v, want = %v", actualString, wantString)
		}
	}
}

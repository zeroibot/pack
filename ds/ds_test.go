package ds

import "testing"

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

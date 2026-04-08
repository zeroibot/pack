package str

import (
	"testing"

	"github.com/zeroibot/tst"
)

func TestBuilder(t *testing.T) {
	b := NewBuilder()
	b.Add("1")
	b.AddItems("2", "3", "4")
	b.AddFmt("%d,%d", 5, 6)
	actual := b.Build(",")
	want := "1,2,3,4,5,6"
	tst.AssertEqual(t, "Builder", actual, want)
}

func TestRepeat(t *testing.T) {
	testCases := []tst.P3W1[int, string, string, string]{
		{5, "a", "", "aaaaa"},
		{3, "ab", "-", "ab-ab-ab"},
		{0, "a", "x", ""},
		{2, "b", ",", "b,b"},
		{1, "x", "x", "x"},
	}
	tst.AllP3W1(t, testCases, "Repeat", Repeat, tst.AssertEqual)
}

func TestRandomString(t *testing.T) {
	// TODO: RandomString
}

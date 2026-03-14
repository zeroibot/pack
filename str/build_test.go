package str

import "testing"

func TestBuilder(t *testing.T) {
	b := NewBuilder()
	b.Add("1")
	b.AddItems("2", "3", "4")
	b.AddFmt("%d,%d", 5, 6)
	actual := b.Build(",")
	want := "1,2,3,4,5,6"
	if actual != want {
		t.Errorf("Builder: got %q, want %q", actual, want)
	}
}

func TestRepeat(t *testing.T) {
	// TODO: Repeat
}

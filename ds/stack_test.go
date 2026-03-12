package ds

import (
	"slices"
	"testing"
)

func TestStack(t *testing.T) {
	// NewStack
	s := NewStack[int]()
	actualString, wantString := s.String(), "[]"
	if actualString != wantString {
		t.Errorf("Stack.String() = %q, want %q", actualString, wantString)
	}
	actualCount, wantCount := s.Len(), 0
	if actualCount != wantCount {
		t.Errorf("Stack.Len() = %d, want %d", actualCount, wantCount)
	}
	actualFlag := s.IsEmpty()
	if actualFlag != true {
		t.Errorf("Stack.IsEmpty() = %t, want %t", actualFlag, true)
	}
	actualFlag = s.NotEmpty()
	if actualFlag != false {
		t.Errorf("Stack.NotEmpty() = %t, want %t", actualFlag, false)
	}
	// NewStackFrom
	items := []int{1, 2, 3}
	s = NewStackFrom[int](items)
	actualString, wantString = s.String(), "[1 2 3]"
	if actualString != wantString {
		t.Errorf("Stack.String() = %q, want %q", actualString, wantString)
	}
	actualCount, wantCount = s.Len(), len(items)
	if actualCount != wantCount {
		t.Errorf("Stack.Len() = %d, want %d", actualCount, wantCount)
	}
	actualFlag = s.IsEmpty()
	if actualFlag != false {
		t.Errorf("Stack.IsEmpty() = %t, want %t", actualFlag, false)
	}
	actualFlag = s.NotEmpty()
	if actualFlag != true {
		t.Errorf("Stack.NotEmpty() = %t, want %t", actualFlag, true)
	}
	actualItems := s.Items()
	if slices.Equal(items, actualItems) == false {
		t.Errorf("Stack.Items() = %v, want %v", actualItems, items)
	}
	// Copy
	s2 := s.Copy()
	wantItems := s.Items()
	actualItems = s2.Items()
	if slices.Equal(wantItems, actualItems) == false {
		t.Errorf("Stack.Copy.Items() = %v, want %v", actualItems, wantItems)
	}
	// Clear
	s2.Clear()
	actualFlag = s2.IsEmpty()
	if actualFlag != true {
		t.Errorf("Stack.Clear.IsEmpty() = %t, want %t", actualFlag, true)
	}
	// Check that original stack is unchanged
	actualItems = s.Items()
	if slices.Equal(items, actualItems) == false {
		t.Errorf("Stack.Items() = %v, want %v", actualItems, items)
	}
}

func TestStackMethods(t *testing.T) {
	// TODO: Push
	// TODO: Top, MustTop
	// TODO: Pop, MustPop
}

package ds

import (
	"testing"

	"github.com/roidaradal/tst"
)

func TestStack(t *testing.T) {
	// NewStack, Len, IsEmpty, NotEmpty
	s := NewStack[int]()
	tst.AssertEqual(t, "Stack.String", s.String(), "[]")
	tst.AssertEqual(t, "Stack.Len", s.Len(), 0)
	tst.AssertEqual(t, "Stack.IsEmpty", s.IsEmpty(), true)
	tst.AssertEqual(t, "Stack.NotEmpty", s.NotEmpty(), false)
	// NewStackFrom, Stack.Items
	items := []int{1, 2, 3}
	s = NewStackFrom[int](items)
	tst.AssertEqual(t, "Stack.String", s.String(), "[1 2 3]")
	tst.AssertEqual(t, "Stack.Len", s.Len(), len(items))
	tst.AssertEqual(t, "Stack.IsEmpty", s.IsEmpty(), false)
	tst.AssertEqual(t, "Stack.NotEmpty", s.NotEmpty(), true)
	tst.AssertListEqual(t, "Stack.Items", s.Items(), items)
	// Copy
	s2 := s.Copy()
	tst.AssertListEqual(t, "Stack.Copy.Items", s2.Items(), s.Items())
	// Clear
	s2.Clear()
	tst.AssertEqual(t, "Stack.Clear.IsEmpty", s2.IsEmpty(), true)
	// Check original stack is unchanged
	tst.AssertListEqual(t, "StackItems", s.Items(), items)
}

func TestStackPush(t *testing.T) {
	s := NewStack[int]()
	s.Push(1)
	tst.AssertEqual(t, "Push.MustTop", s.MustTop(), 1)
	s.Push(2)
	tst.AssertEqual(t, "Push.MustTop", s.MustTop(), 2)
	s.Push(3)
	tst.AssertEqual(t, "Push.MustTop", s.MustTop(), 3)
	tst.AssertListEqual(t, "Stack.Items", s.Items(), []int{1, 2, 3})
}

func TestStackTop(t *testing.T) {
	s := NewStack[int]()
	top := s.Top()
	tst.AssertEqual(t, "Stack.Top.IsNil", top.IsNil(), true)
	s.Push(1)
	top = s.Top()
	tst.AssertEqual2(t, "Stack.Top", top.IsNil(), false, top.Value(), 1)
	tst.AssertEqual(t, "Stack.MustTop", s.MustTop(), 1)

	s.Pop()
	defer tst.AssertPanic(t, "Stack.MustTop")
	s.MustTop() // should panic
}

func TestStackPop(t *testing.T) {
	s := NewStackFrom[int]([]int{1, 2})
	top := s.Pop()
	tst.AssertEqual2(t, "Stack.Pop", top.IsNil(), false, top.Value(), 2)
	topItem := s.MustPop()
	tst.AssertEqual(t, "Stack.MustPop", topItem, 1)
	top = s.Pop()
	tst.AssertEqual2(t, "Stack.Pop", top.IsNil(), true, top.Value(), 0)

	defer tst.AssertPanic(t, "Stack.MustPop")
	s.MustPop() // should panic
}

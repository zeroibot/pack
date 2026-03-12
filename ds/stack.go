package ds

import "fmt"

type Stack[T any] struct {
	items List[T]
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return new(Stack[T]{items: make(List[T], 0)})
}

// NewStackFrom creates a new stack from list of items (last item = stack top)
func NewStackFrom[T any](items []T) *Stack[T] {
	return new(Stack[T]{items: items})
}

// String returns the string representation of the stack
func (s *Stack[T]) String() string {
	return fmt.Sprintf("%v", s.items)
}

// Len returns the number of items in the stack
func (s *Stack[T]) Len() int {
	return s.items.Len()
}

// IsEmpty checks if stack is empty
func (s *Stack[T]) IsEmpty() bool {
	return s.items.IsEmpty()
}

// NotEmpty checks if stack is not empty
func (s *Stack[T]) NotEmpty() bool {
	return s.items.NotEmpty()
}

// Items returns the List of stack items
func (s *Stack[T]) Items() List[T] {
	return s.items
}

// Push adds an item to the top of the stack
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Top returns an Option that contains the top item.
// If the stack is empty, the Option contains nil
func (s *Stack[T]) Top() Option[T] {
	return s.items.Last(1)
}

// MustTop returns the top item, and panics if the stack is empty
func (s *Stack[T]) MustTop() T {
	option := s.Top()
	if option.IsNil() {
		panic("empty stack")
	}
	return option.Value()
}

// Pop returns an Option that contains the top item and removes it from the stack.
// If the stack is empty, the Option contains nil
func (s *Stack[T]) Pop() Option[T] {
	option := s.Top()
	if option.IsNil() {
		return option
	}
	s.items = s.items[:s.items.LastIndex()]
	return option
}

// MustPop returns the top item and removes it from the stack, and panics if the stack is empty
func (s *Stack[T]) MustPop() T {
	option := s.Pop()
	if option.IsNil() {
		panic("empty stack")
	}
	return option.Value()
}

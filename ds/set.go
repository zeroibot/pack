package ds

type Set[T comparable] struct {
	items Map[T, struct{}]
}

// NewEmptySet creates a new empty set
func NewEmptySet[T comparable]() *Set[T] {
	return &Set[T]{items: make(Map[T, struct{}])}
}

// NewSetFrom creates a new set from given items
func NewSetFrom[T comparable](items List[T]) *Set[T] {
	set := NewEmptySet[T]()
	set.AddItems(items)
	return set
}

// Len returns the Set size
func (s *Set[T]) Len() int {
	return s.items.Len()
}

// IsEmpty checks if the Set is empty
func (s *Set[T]) IsEmpty() bool {
	return s.items.IsEmpty()
}

// NotEmpty checks if the Set is not empty
func (s *Set[T]) NotEmpty() bool {
	return s.items.NotEmpty()
}

// Copy creates a new Set with copied items
func (s *Set[T]) Copy() *Set[T] {
	return &Set[T]{items: s.items.Copy()}
}

// Items returns the Set items, in arbitrary order
func (s *Set[T]) Items() List[T] {
	return s.items.Keys()
}

// Add adds an item to the set
func (s *Set[T]) Add(item T) {
	s.items[item] = struct{}{}
}

// AddItems adds items to the set
func (s *Set[T]) AddItems(items List[T]) {
	for _, item := range items {
		s.Add(item)
	}
}

// Has checks if set contains item
func (s *Set[T]) Has(item T) bool {
	_, hasItem := s.items[item]
	return hasItem
}

// HasNo checks if set does not contain item
func (s *Set[T]) HasNo(item T) bool {
	return !s.Has(item)
}

// Delete deletes an item from the set
func (s *Set[T]) Delete(item T) {
	delete(s.items, item)
}

// Union computes the union of two sets
func (s *Set[T]) Union(s2 *Set[T]) *Set[T] {
	s3 := NewEmptySet[T]()
	for _, items := range []Map[T, struct{}]{s.items, s2.items} {
		for item := range items {
			s3.Add(item)
		}
	}
	return s3
}

// Intersection computes the intersection of two sets
func (s *Set[T]) Intersection(s2 *Set[T]) *Set[T] {
	s3 := NewEmptySet[T]()
	for item := range s.items {
		if s2.Has(item) {
			s3.Add(item)
		}
	}
	return s3
}

// Difference computes the difference of two sets
func (s *Set[T]) Difference(s2 *Set[T]) *Set[T] {
	s3 := NewEmptySet[T]()
	for item := range s.items {
		if s2.HasNo(item) {
			s3.Add(item)
		}
	}
	return s3
}

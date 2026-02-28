package ds

// List extends the slice collection
type List[T any] []T

// Len returns the List length
func (l List[T]) Len() Int {
	return Int(len(l))
}

// Cap returns the List capacity
func (l List[T]) Cap() Int {
	return Int(cap(l))
}

// IsEmpty checks if the List is empty
func (l List[T]) IsEmpty() Boolean {
	return len(l) == 0
}

// NotEmpty checks if the List is not empty
func (l List[T]) NotEmpty() Boolean {
	return len(l) > 0
}

// Copy creates a new List with copied items
func (l List[T]) Copy() List[T] {
	items := append(List[T]{}, l...)
	return items
}

// Last returns the nth item from the back of the list (starts at 1)
func (l List[T]) Last(rank int) (T, Boolean) {
	numItems := len(l)
	if rank > numItems || rank <= 0 {
		var item T
		return item, false
	}
	return l[numItems-rank], true
}

package ds

// NumberList is a List of Numbers (Int, Uint, Float)
type NumberList[T Number] struct {
	List[T]
}

// NewNumberList creates a new NumberList, with the given items
func NewNumberList[T Number](items ...T) *NumberList[T] {
	return &NumberList[T]{List: items}
}

// NewRangeList creates a new NumberList, containing numbers from [start, end)
func NewRangeList[T Integer](start, end T) *NumberList[T] {
	items := make(List[T], 0, end-start)
	for i := start; i < end; i++ {
		items = append(items, i)
	}
	return &NumberList[T]{List: items}
}

// EmptyNumberList creates a new empty NumberList, with given capacity
func EmptyNumberList[T Number](capacity int) *NumberList[T] {
	return &NumberList[T]{List: make(List[T], 0, capacity)}
}

// Append adds a new Number to end of the List
func (n *NumberList[T]) Append(item T) {
	n.List = append(n.List, item)
}

// Sum computes the sum of number items
func (n *NumberList[T]) Sum() T {
	var total T = 0
	for _, number := range n.List {
		total += number
	}
	return total
}

// Product computes the product of number items
func (n *NumberList[T]) Product() T {
	var product T = 1
	for _, number := range n.List {
		product *= number
	}
	return product
}

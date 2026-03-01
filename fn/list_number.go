package fn

import "github.com/roidaradal/pack/ds"

// Sum computes the sum of number items
func Sum[T ds.Number](numbers ds.List[T]) T {
	var total T = 0
	for _, number := range numbers {
		total += number
	}
	return total
}

// Product computes the product of number items
func Product[T ds.Number](numbers ds.List[T]) T {
	var product T = 1
	for _, number := range numbers {
		product *= number
	}
	return product
}

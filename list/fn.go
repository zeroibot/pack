package list

import "github.com/roidaradal/pack/number"

// Sum computes the sum of number items
func Sum[T number.Type](numbers []T) T {
	var total T = 0
	for _, x := range numbers {
		total += x
	}
	return total
}

// Product computes the product of number items
func Product[T number.Type](numbers []T) T {
	var product T = 1
	for _, x := range numbers {
		product *= x
	}
	return product
}

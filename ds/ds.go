// Package ds contains data structures
package ds

// Boolean extends the bool type
type Boolean bool

// ToInt converts a Boolean to Int (true = 1, false = 0)
func (b Boolean) ToInt() Int {
	if b {
		return 1
	}
	return 0
}

// ToUint converts a Boolean to Uint (true = 1, false = 0)
func (b Boolean) ToUint() Uint {
	if b {
		return 1
	}
	return 0
}

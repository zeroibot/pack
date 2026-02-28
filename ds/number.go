package ds

import "fmt"

// Int extends the int type
type Int int

// Uint extends the uint type
type Uint uint

// Float extends the float64 type
type Float float64

// String representation of Int
func (i Int) String() String {
	return String(fmt.Sprintf("%d", i))
}

// String representation of Uint
func (u Uint) String() String {
	return String(fmt.Sprintf("%d", u))
}

// String representation of Float
func (f Float) String() String {
	return String(fmt.Sprintf("%f", f))
}

// StringDecimal returns the float string, with N decimal places
func (f Float) StringDecimal(decimalPlaces int) String {
	format := fmt.Sprintf("%%.%df", decimalPlaces)
	return String(fmt.Sprintf(format, f))
}

// Abs gets the absolute value of the Int
func (i Int) Abs() Int {
	if i < 0 {
		return -i
	}
	return i
}

// ToUint converts the Int to Uint
func (i Int) ToUint() Uint {
	return Uint(i)
}

// ToFloat converts the Int to Float
func (i Int) ToFloat() Float {
	return Float(i)
}

// ToBoolean converts the Int to Boolean (0 is false, else true)
func (i Int) ToBoolean() Boolean {
	return i != 0
}

// ToInt converts the Uint to Int
func (u Uint) ToInt() Int {
	return Int(u)
}

// ToFloat converts the Uint to Float
func (u Uint) ToFloat() Float {
	return Float(u)
}

// ToBoolean converts the Uint to Boolean (0 is false, else true)
func (u Uint) ToBoolean() Boolean {
	return u != 0
}

// ToInt converts the Float to Int
func (f Float) ToInt() Int {
	return Int(f)
}

// ToUint converts the Float to Uint
func (f Float) ToUint() Uint {
	return Uint(f)
}

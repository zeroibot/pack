package ds

// Tuple2 is a tuple with two types
type Tuple2[A, B any] struct {
	V1 A
	V2 B
}

// Tuple3 is a tuple with three types
type Tuple3[A, B, C any] struct {
	V1 A
	V2 B
	V3 C
}

// Tuple4 is a tuple with four types
type Tuple4[A, B, C, D any] struct {
	V1 A
	V2 B
	V3 C
	V4 D
}

// Pair is a tuple with two values of the same type
type Pair[T any] [2]T

// Triple is a tuple with three values of the same type
type Triple[T any] [3]T

// Quad is a tuple with four values of the same type
type Quad[T any] [4]T

// NewTuple2 creates a new Tuple2
func NewTuple2[A, B any](v1 A, v2 B) Tuple2[A, B] {
	return Tuple2[A, B]{v1, v2}
}

// NewTuple3 creates a new Tuple3
func NewTuple3[A, B, C any](v1 A, v2 B, v3 C) Tuple3[A, B, C] {
	return Tuple3[A, B, C]{v1, v2, v3}
}

// NewTuple4 creates a new Tuple4
func NewTuple4[A, B, C, D any](v1 A, v2 B, v3 C, v4 D) Tuple4[A, B, C, D] {
	return Tuple4[A, B, C, D]{v1, v2, v3, v4}
}

// Values returns the unpacked Tuple2 values
func (t Tuple2[A, B]) Values() (A, B) {
	return t.V1, t.V2
}

// Values returns the unpacked Tuple3 values
func (t Tuple3[A, B, C]) Values() (A, B, C) {
	return t.V1, t.V2, t.V3
}

// Values returns the unpacked Tuple4 values
func (t Tuple4[A, B, C, D]) Values() (A, B, C, D) {
	return t.V1, t.V2, t.V3, t.V4
}

// Values returns the unpacked Pair values
func (p Pair[T]) Values() (T, T) {
	return p[0], p[1]
}

// Values returns the unpacked Triple values
func (t Triple[T]) Values() (T, T, T) {
	return t[0], t[1], t[2]
}

// Values returns the unpacked Quad values
func (q Quad[T]) Values() (T, T, T, T) {
	return q[0], q[1], q[2], q[3]
}

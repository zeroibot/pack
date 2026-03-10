package ds

import "testing"

func TestTuple2(t *testing.T) {
	v1, v2 := "apple", 5
	t2 := NewTuple2(v1, v2)
	if t2.V1 != v1 || t2.V2 != v2 {
		t.Errorf("Tuple2.V1, V2 = %v, %v, want %v, %v", t2.V1, t2.V2, v1, v2)
	}
	a, b := t2.Values()
	if a != v1 || b != v2 {
		t.Errorf("Tuple2.Values() = %v, %v, want %v, %v", a, b, v1, v2)
	}
}

func TestTuple3(t *testing.T) {
	v1, v2, v3 := "apple", 5, 12.5
	t3 := NewTuple3(v1, v2, v3)
	if t3.V1 != v1 || t3.V2 != v2 || t3.V3 != v3 {
		t.Errorf("Tuple3.V1, V2, V3 = %v, %v, %v want %v, %v, %v", t3.V1, t3.V2, t3.V3, v1, v2, v3)
	}
	a, b, c := t3.Values()
	if a != v1 || b != v2 || c != v3 {
		t.Errorf("Tuple3.Values() = %v, %v, %v want %v, %v, %v", a, b, c, v1, v2, v3)
	}
}

func TestTuple4(t *testing.T) {
	v1, v2, v3, v4 := "apple", 5, 12.5, 'x'
	t4 := NewTuple4(v1, v2, v3, v4)
	if t4.V1 != v1 || t4.V2 != v2 || t4.V3 != v3 || t4.V4 != v4 {
		t.Errorf("Tuple4.V1, V2, V3, V4 = %v, %v, %v, %v want %v, %v, %v, %v", t4.V1, t4.V2, t4.V3, t4.V4, v1, v2, v3, v4)
	}
	a, b, c, d := t4.Values()
	if a != v1 || b != v2 || c != v3 || d != v4 {
		t.Errorf("Tuple4.Values() = %v, %v, %v, %v want %v, %v, %v, %v", a, b, c, d, v1, v2, v3, v4)
	}
}

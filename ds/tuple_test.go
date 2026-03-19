package ds

import (
	"testing"

	"github.com/roidaradal/tst"
)

func TestTuple2(t *testing.T) {
	v1, v2 := "apple", 5
	t2 := NewTuple2(v1, v2)
	a, b := t2.Unpack()
	tst.AssertEqual2(t, "Tuple2.V1,V2", t2.V1, v1, t2.V2, v2)
	tst.AssertEqual2(t, "Tuple2.Unpack", a, v1, b, v2)
}

func TestTuple3(t *testing.T) {
	v1, v2, v3 := "apple", 5, 12.5
	t3 := NewTuple3(v1, v2, v3)
	a, b, c := t3.Unpack()
	tst.AssertEqual3(t, "Tuple3.V1,V2,V3", t3.V1, v1, t3.V2, v2, t3.V3, v3)
	tst.AssertEqual3(t, "Tuple3.Unpack", a, v1, b, v2, c, v3)
}

func TestTuple4(t *testing.T) {
	v1, v2, v3, v4 := "apple", 5, 12.5, 'x'
	t4 := NewTuple4(v1, v2, v3, v4)
	a, b, c, d := t4.Unpack()
	tst.AssertEqual4(t, "Tuple4.V1,V2,V3,V4", t4.V1, v1, t4.V2, v2, t4.V3, v3, t4.V4, v4)
	tst.AssertEqual4(t, "Tuple4.Unpack", a, v1, b, v2, c, v3, d, v4)
}

func TestPair(t *testing.T) {
	v1, v2 := 6, 7
	p := Pair[int]{v1, v2}
	a, b := p.Unpack()
	tst.AssertEqual2(t, "Pair[0,1]", p[0], v1, p[1], v2)
	tst.AssertEqual2(t, "Pair.Unpack", a, v1, b, v2)
}

func TestTriple(t *testing.T) {
	v1, v2, v3 := 6, 7, 8
	r := Triple[int]{v1, v2, v3}
	a, b, c := r.Unpack()
	tst.AssertEqual3(t, "Triple[0,1,2]", r[0], v1, r[1], v2, r[2], v3)
	tst.AssertEqual3(t, "Triple.Unpack", a, v1, b, v2, c, v3)
}

func TestQuad(t *testing.T) {
	v1, v2, v3, v4 := 1, 2, 3, 4
	q := Quad[int]{v1, v2, v3, v4}
	a, b, c, d := q.Unpack()
	tst.AssertEqual4(t, "Quad[0,1,2,3]", q[0], v1, q[1], v2, q[2], v3, q[3], v4)
	tst.AssertEqual4(t, "Quad.Unpack", a, v1, b, v2, c, v3, d, v4)
}

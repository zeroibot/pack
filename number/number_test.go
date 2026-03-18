package number

import (
	"testing"

	"github.com/roidaradal/tst"
)

func TestAbs(t *testing.T) {
	testCases1 := []tst.P1W1[int, int]{
		{-5, 5},
		{0, 0},
		{3, 3},
	}
	testCases2 := []tst.P1W1[float64, float64]{
		{-5.5, 5.5},
		{0, 0},
		{3.25, 3.25},
	}
	tst.AllP1W1(t, testCases1, "Abs", Abs, tst.AssertEqual)
	tst.AllP1W1(t, testCases2, "Abs", Abs, tst.AssertEqual)
}

func TestCeilInt(t *testing.T) {
	testCases := []tst.P1W1[float64, int]{
		{6.95, 7},
		{3.14, 4},
		{1.0001, 2},
		{5, 5},
		{0, 0},
		{-1.92, -1},
		{-2.05, -2},
	}
	tst.AllP1W1(t, testCases, "CeilInt", CeilInt, tst.AssertEqual)
}

func TestFloorInt(t *testing.T) {
	testCases := []tst.P1W1[float64, int]{
		{6.95, 6},
		{3.14, 3},
		{1.0001, 1},
		{5, 5},
		{0, 0},
		{-1.92, -2},
		{-2.05, -3},
	}
	tst.AllP1W1(t, testCases, "FloorInt", FloorInt, tst.AssertEqual)
}

func TestRoundInt(t *testing.T) {
	testCases := []tst.P1W1[float64, int]{
		{6.95, 7},
		{3.14, 3},
		{1.0001, 1},
		{5, 5},
		{0, 0},
		{-1.92, -2},
		{-2.05, -2},
		{3.5, 4},
		{4.5, 5},
		{0.5, 1},
	}
	tst.AllP1W1(t, testCases, "RoundInt", RoundInt, tst.AssertEqual)
}

func TestRoundToEvenInt(t *testing.T) {
	testCases := []tst.P1W1[float64, int]{
		{6.95, 7},
		{3.14, 3},
		{1.0001, 1},
		{5, 5},
		{0, 0},
		{-1.92, -2},
		{-2.05, -2},
		{3.5, 4},
		{4.5, 4},
		{0.5, 0},
		{1.5, 2},
		{1.49, 1},
		{1.51, 2},
	}
	tst.AllP1W1(t, testCases, "RoundToEvenInt", RoundToEvenInt, tst.AssertEqual)
}

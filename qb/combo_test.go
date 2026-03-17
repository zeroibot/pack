package qb

import (
	"slices"
	"testing"

	"github.com/roidaradal/pack/list"
)

func TestCombos(t *testing.T) {
	type Person struct {
		Name    string
		Address string
		Age     int
		Job     string
		Score   int
	}
	type testCase struct {
		cond      DualCondition[Person]
		wantBools []bool
	}
	this := NewInstance(MySQL)
	p := &Person{}
	err := AddType(this, p)
	if err != nil {
		t.Errorf("AddType error: %v", err)
	}

	persons := []Person{
		{"John", "Astra Tower, Gotham City", 30, "manager", 95},
		{"Jane", "Star Village, India", 65, "assistant ceo", 85},
		{"Jill", "Tower of Doom, Japan", 20, "dev", 55},
		{"Jim", "Sub Division, Singapore City", 40, "qa", 67},
		{"Jack", "N/A", 18, "intern", 50},
	}

	condNone := NoCondition[Person]()
	condEqual := Equal[Person](this, &p.Name, "John")
	condNotEqual := NotEqual[Person](this, &p.Job, "manager")
	condPrefix := Prefix[Person](this, &p.Job, "assistant")
	condSuffix := Suffix[Person](this, &p.Address, "City")
	condSubstring := Substring[Person](this, &p.Address, "Tower")
	condGreater := Greater[Person](this, &p.Score, 75)
	condGreaterEqual := GreaterEqual[Person](this, &p.Age, 20)
	condLesser := Lesser[Person](this, &p.Age, 60)
	condLesserEqual := LesserEqual[Person](this, &p.Score, 50)
	condIn := In[Person](this, &p.Job, []string{"dev", "qa", "intern"})
	condNotIn := NotIn[Person](this, &p.Score, []int{67, 69})
	condAnd := And[Person](condEqual, condGreater)
	condOr := Or[Person](condLesserEqual, condIn)

	testCases := []testCase{
		{condNone, []bool{true, true, true, true, true}},
		{condEqual, []bool{true, false, false, false, false}},
		{condNotEqual, []bool{false, true, true, true, true}},
		{condPrefix, []bool{false, true, false, false, false}},
		{condSuffix, []bool{true, false, false, true, false}},
		{condSubstring, []bool{true, false, true, false, false}},
		{condGreater, []bool{true, true, false, false, false}},
		{condGreaterEqual, []bool{true, true, true, true, false}},
		{condLesser, []bool{true, false, true, true, true}},
		{condLesserEqual, []bool{false, false, false, false, true}},
		{condIn, []bool{false, false, true, true, true}},
		{condNotIn, []bool{true, true, true, false, true}},
		{condAnd, []bool{true, false, false, false, false}},
		{condOr, []bool{false, false, true, true, true}},
	}
	for _, x := range testCases {
		actualBools := list.Map(persons, x.cond.Test)
		if slices.Equal(actualBools, x.wantBools) == false {
			t.Errorf("Combo.Test() = %v, want %v", actualBools, x.wantBools)
		}
	}
}

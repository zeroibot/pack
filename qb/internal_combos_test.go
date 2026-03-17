package qb

import (
	"slices"
	"testing"

	"github.com/roidaradal/pack/ds"
	"github.com/roidaradal/pack/list"
)

func TestInternalCombos(t *testing.T) {
	type Person struct {
		Name    string
		Age     int `col:"age"`
		Job     string
		Details string `col:"-"`
	}
	this := NewInstance(MySQL)
	p := &Person{}
	err := AddType(this, p)
	if err != nil {
		t.Errorf("AddType error: %v", err)
	}
	p1 := Person{"John", 18, "dev", "regular"}
	p2 := Person{"James", 20, "qa", "prob"}
	p3 := Person{"Jill", 25, "admin", "regular"}
	p4 := Person{"Juno", 23, "dev", "prob"}
	p5 := Person{"Jack", 21, "sales", "intern"}
	people := ds.List[Person]{p1, p2, p3, p4, p5}
	// missingCombo.Test
	combo1 := missingCombo[Person]{}
	allFalse := people.All(func(person Person) bool {
		return combo1.Test(person) == false
	})
	if !allFalse {
		t.Errorf("Not all missingCombo.Test() returned false")
	}
	// matchAllCombo.Test
	combo2 := matchAllCombo[Person]{}
	allTrue := people.All(func(person Person) bool {
		return combo2.Test(person) == true
	})
	if !allTrue {
		t.Errorf("Not all matchAllCombo.Test() returned true")
	}
	// newValueCombo, valueCombo.Test
	test1 := createFieldValueTest[Person]("Name", func(name string) bool { return name == "Jill" })
	valueCond1 := newValueCondition(this, &p.Name, "Jill", opEqual)
	valueCombo1 := newValueCombo(valueCond1, test1)
	wantBools := []bool{false, false, true, false, false}
	actualBools := list.Map(people, valueCombo1.Test)
	if slices.Equal(wantBools, actualBools) == false {
		t.Errorf("valueCombo.Test() = %v, want %v", actualBools, wantBools)
	}
	test0 := createFieldValueTest[Person]("Unknown", func(unknown string) bool { return unknown != "A" })
	valueCombo2 := newValueCombo(valueCond1, test0)
	allFalse = people.All(func(person Person) bool {
		return valueCombo2.Test(person) == false
	})
	if !allFalse {
		t.Errorf("Not all valueCombo.Test() returned false")
	}
	// newListCombo, listCombo.Test
	jobs := []string{"dev", "qa"}
	test2 := createFieldValueTest[Person]("Job", func(job string) bool { return list.Has(jobs, job) })
	listCond1 := newListCondition(this, &p.Job, jobs, opIn, opEqual)
	listCombo1 := newListCombo(listCond1, test2)
	wantBools = []bool{true, true, false, true, false}
	actualBools = list.Map(people, listCombo1.Test)
	if slices.Equal(wantBools, actualBools) == false {
		t.Errorf("listCombo.Test() = %v, want %v", actualBools, wantBools)
	}
	// newMultiCombo, multiCombo.Test
	multiTest := func(person Person) bool { return test1(person) || test2(person) }
	multiCombo1 := newMultiCombo(ds.List[DualCondition[Person]]{valueCombo1, listCombo1}, opOr, multiTest)
	wantBools = []bool{true, true, true, true, false}
	actualBools = list.Map(people, multiCombo1.Test)
	if slices.Equal(wantBools, actualBools) == false {
		t.Errorf("multiCombo.Test() = %v, want %v", actualBools, wantBools)
	}
	// multiCombo.BuildCondition
	wantCond := "(`Name` = ? OR `Job` IN (?, ?))"
	wantValues := ds.List[any]{"Jill", "dev", "qa"}
	actualCond, actualValues := multiCombo1.BuildCondition()
	if actualCond != wantCond || slices.Equal(wantValues, actualValues) == false {
		t.Errorf("multiCombo.BuildCondition() = %q, %v, want %q, %v", actualCond, actualValues, wantCond, wantValues)
	}
}

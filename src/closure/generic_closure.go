package main

import (
	"fmt"
	"strings"
)

type Employee struct {
	Id   int
	Name string
}

type Person struct {
	Age  int
	Name string
}

type PetOwner struct {
	Person
	Pet string
}

func NameFilter(l int, name string,
	nameFinder func(int) string,
	populator func(int)) {
	for i := 0; i < l; i++ {
		if strings.Contains(nameFinder(i), name) {
			populator(i)
		}
	}
}

func main() {
	employees := []Employee{
		{Id: 1, Name: "Mike Lee"},
		{Id: 2, Name: "Johan Fred"},
		{Id: 3, Name: "Ginger"},
		{Id: 4, Name: "Riza"},
		{Id: 5, Name: "Tom Batsele"},
		{Id: 6, Name: "Fred Leonado"},
		{Id: 7, Name: "Jenny"},
		{Id: 8, Name: "Fred"},
		{Id: 9, Name: "red"},
		{Id: 10, Name: "FredMan"},
	}
	var allFredEmps []Employee
	NameFilter(len(employees), "Fred", func(i int) string {
		return employees[i].Name
	}, func(i int) {
		allFredEmps = append(allFredEmps, employees[i])
	})
	fmt.Println("All Employees which contain name of 'Fred':", allFredEmps)
}

//the following lines are duplicate code, should use closure feature above to achieve the generic
func FilterPerson(ps []Person, name string) []Person {
	var out []Person
	for _, v := range ps {
		if v.Name == name {
			out = append(out, v)
		}
	}
	return out
}

func FilterEmployee(es []Employee, name string) []Employee {

	var out []Employee
	for _, v := range es {
		if v.Name == name {
			out = append(out, v)
		}
	}
	return out
}

func FilterPetOwner(po []PetOwner, name string) []PetOwner {

	var out []PetOwner
	for _, v := range po {
		if v.Name == name {
			out = append(out, v)
		}
	}
	return out
}

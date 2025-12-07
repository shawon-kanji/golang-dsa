package main

import (
	"fmt"
	"sort"
)

type Person struct {
	name string
	age  int
}

type ByAge []Person

func (p ByAge) Len() int {
	return len(p)
}

func (p ByAge) Less(i, j int) bool {
	return p[i].age < p[j].age
}

func (p ByAge) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func main() {
	people := []Person{
		{name: "a", age: 10},
		{name: "b", age: 5},
		{name: "c", age: 15},
		{name: "d", age: 85},
		{name: "e", age: 50},
	}

	sort.Sort(ByAge(people))

	for _, v := range people {
		fmt.Println("person:: ", v)
	}
}

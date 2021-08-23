package main

import "fmt"

type person interface {
	getName() string
}

type human struct {
	name string
}

func (h human) getName() string {
	return h.name
}

type teacher struct {
	human
	school string
	degree string
}

type student struct {
	human
	school string
	grades map[string]int
}

type printer struct {
}

func (pr printer) printName(p person) {
	fmt.Printf("Name is %s\n", p.getName())
}

func main() {
	h := human{"Helen"}
	t := teacher{
		human{"Tom"},
		"No.1 School",
		"Master",
	}
	s := student{
		human{"Susan"},
		"No.1 School",
		map[string]int{
			"English": 3,
			"Math":    9,
		},
	}
	pr := printer{}
	pr.printName(h)
	pr.printName(t)
	pr.printName(s)
}

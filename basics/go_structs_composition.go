package main

import "fmt"

// https://medium.com/@simplyianm/why-gos-structs-are-superior-to-class-based-inheritance-b661ba897c67
type Animal interface {
	Name() string
}

type Dog struct {
}

func (d *Dog) Name() string {
	return "Donald"
}

func (d *Dog) Bark() string {
	return "Wang Wang!"
}

type PartyAnimal interface {
	Animal
	Party()
}

func (d *Dog) Party() {
	fmt.Printf("%s is a Party Animal", d.Name())
}

func main() {
	var animal PartyAnimal

	// TODO: reading pointer and *
	// works
	animal = &Dog{}
	// not work
	// animal = Dog{}
	fmt.Printf("%s barks %s\n", animal.Name(), animal.(*Dog).Bark())
	animal.Party()
}

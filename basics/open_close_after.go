package main

import (
	"fmt"
	"math"
)

type shape interface {
	area() float64
}

type circle struct {
	radius float64
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

type square struct {
	length float64
}

func (s square) area() float64 {
	return s.length * s.length
}

type calculator struct {
}

// signature passing shape in, so to add new shape, for example, rectangle, just adding the shape, no need to change sum method
func (c calculator) sum(shapes ...shape) float64 {
	var sum float64

	for _, shape := range shapes {
		sum += shape.area()
	}
	return sum
}

func main() {
	c := circle{5}
	s := square{10}
	calculator := calculator{}
	sum := calculator.sum(c, s)
	fmt.Printf("Total area is: %f\n", sum)
}

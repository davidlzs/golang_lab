package main

import (
	"fmt"
	"math"
)

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

// violate the open close principle; if you want to extend to add a new shape: hexgonal, you need edit the sum() method
func sum(shapes ...interface{}) float64 {
	var sum float64

	for _, shape := range shapes {
		switch shape.(type) {
		case circle:
			sum += shape.(circle).area()
		case square:
			sum += shape.(square).area()
		}
	}
	return sum
}

func main() {
	c := circle{5}
	s := square{10}
	sum := sum(c, s)
	fmt.Println(sum)
}

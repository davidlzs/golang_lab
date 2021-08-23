package main

import (
	"fmt"
	"math"
)

type shape interface {
	area() float64
}

// separate shape interface to object interface which has volume() method, because volume() method does not apply to 2D shapes
type object interface {
	shape
	volume() float64
}

type square struct {
	length float64
}

func (s square) area() float64 {
	return s.length * s.length
}

type cube struct {
	length float64
}

func (c cube) area() float64 {
	return math.Pow(c.length, 2)
}

func (c cube) volume() float64 {
	return math.Pow(c.length, 3)
}

func areaSum(shapes ...shape) float64 {
	var sum float64
	for _, s := range shapes {
		sum += s.area()
	}
	return sum
}

func areaVolumeSum(objects ...object) float64 {
	var sum float64
	for _, s := range objects {
		sum += s.area() + s.volume()
	}
	return sum
}

func main() {
	s1 := square{6}
	s2 := square{10}
	c1 := cube{5}
	c2 := cube{10}
	fmt.Printf("area sum; %f\n", areaSum(s1, s2))
	fmt.Printf("area volume sum: %f\n", areaVolumeSum(c1, c2))
}

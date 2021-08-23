package main

import (
	"encoding/json"
	"fmt"
	"math"
)

type shape interface {
	area() float64
	name() string
}

type circle struct {
	radius float64
}

func (c circle) name() string {
	return "circle"
}

// single responsiblity for calculating area
// only need to change area(), when need to change area()
func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

type square struct {
	length float64
}

func (s square) name() string {
	return "square"
}

func (s square) area() float64 {
	return s.length * s.length
}

// single responsibility for outputting
// only need to change outputter when want to change the output, for example, add yaml output
type outputter struct {
}

func (o outputter) Text(s shape) string {
	return fmt.Sprintf("%s area: %f", s.name(), s.area())
}

func (o outputter) JSON(s shape) string {
	res := struct {
		Name string  `json:"shape"`
		Area float64 `json:"area"`
	}{
		Name: s.name(),
		Area: s.area(),
	}

	bs, err := json.Marshal(res)

	if err != nil {
		panic(err)
	}
	return string(bs)
}

func main() {
	o := outputter{}

	c := circle{radius: 5}
	fmt.Println(o.Text(c))
	fmt.Println(o.JSON(c))

	s := square{length: 10}
	fmt.Println(o.Text(s))
	fmt.Println(o.JSON(s))
}

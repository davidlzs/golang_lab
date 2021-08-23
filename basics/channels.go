package main

import "fmt"

type intChan chan int

func sum(s []int, c intChan) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}

func main() {
	s := []int{7, 5, 3, 2, 0}

	c := make(intChan)
	go sum(s, c)
	x := <-c

	fmt.Println(x)
	fmt.Println("hello")

}

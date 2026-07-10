package main

import "fmt"

type Numbers interface {
	int | float32 | float64
}

func add[T Numbers](elems ...T) T {
	var total T

	for _, v := range elems {
		total += v
	}

	return total
}

func main() {

	a := add(1, 2, 3, 4)
	fmt.Println(a)

	b := add(1.2, 34.5, 56.78)
	fmt.Println(b)

}

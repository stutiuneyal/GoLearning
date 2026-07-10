package main

import (
	"fmt"
	"slices"
)

func main() {

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(len(numbers))

	numbers = append(numbers, 10)
	fmt.Println(len(numbers))

	n1 := numbers[2:5] // [2,5)
	fmt.Printf("%#v\n", n1)

	n2 := numbers[:4] // [0,4)
	fmt.Printf("%#v\n", n2)

	n3 := numbers[4:] // [4,len(numbers))
	fmt.Printf("%#v\n", n3)

	// slices package
	if slices.Contains(numbers, 11) {
		fmt.Println("slice contains 11")
	} else {
		fmt.Println("slice does not contain 11")
	}

	numbers = slices.Insert(numbers, 4, 100, 101)
	fmt.Printf("%#v\n", numbers)

	// appending a slice in a slice
	evenNumbers := []int{200, 300, 400, 500}
	numbers = append(numbers, evenNumbers...) // spread operator -> [1,2,3,4]... -> 1,2,3,4

}
